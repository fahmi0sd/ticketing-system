package booking

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	routepkg "github.com/fahmi0sd/ticketing-system/service/route"
)

type MidtransClient interface {
	CreatePayment(externalID string, amount float64) (paymentURL string, err error)
}

type service struct {
	logger      *slog.Logger
	repo        Repository
	routeRepo   routepkg.Repository
	midtrans    MidtransClient
	midtransKey string
}

type Service interface {
	Create(userID int, req CreateRequest) (Booking, error)
	GetMyBookings(userID int) ([]Booking, error)
	GetByID(userID, bookingID int) (Booking, error)
	Cancel(userID, bookingID int) error
	HandleWebhook(payload WebhookPayload) error
	GetPaymentStatus(userID, bookingID int) (Booking, error)
}

func NewService(
	logger *slog.Logger,
	repo Repository,
	routeRepo routepkg.Repository,
	midtrans MidtransClient,
	midtransKey string,
) Service {
	return &service{
		logger:      logger,
		repo:        repo,
		routeRepo:   routeRepo,
		midtrans:    midtrans,
		midtransKey: midtransKey,
	}
}

func (s *service) Create(userID int, req CreateRequest) (Booking, error) {
	route, err := s.routeRepo.GetByID(req.RouteID)
	if err != nil {
		return Booking{}, errors.New("route not found")
	}

	available := route.Quota - route.Sold
	if req.Quantity > available {
		return Booking{}, fmt.Errorf("insufficient seats: only %d available", available)
	}

	totalPrice := route.Price * float64(req.Quantity)
	externalID := fmt.Sprintf("TKT-%d-%d", userID, time.Now().UnixNano())
	expiredAt := time.Now().Add(15 * time.Minute)

	paymentURL, err := s.midtrans.CreatePayment(externalID, totalPrice)
	if err != nil {
		s.logger.Error("midtrans create payment error", slog.Any("err", err))
		return Booking{}, errors.New("failed to create payment link")
	}

	b := Booking{
		UserID:     userID,
		RouteID:    req.RouteID,
		Quantity:   req.Quantity,
		TotalPrice: totalPrice,
		Status:     StatusPending,
		PaymentURL: paymentURL,
		ExternalID: externalID,
		ExpiredAt:  &expiredAt,
	}

	created, err := s.repo.Create(b)
	if err != nil {
		s.logger.Error("create booking error", slog.Any("err", err))
		return Booking{}, errors.New("failed to create booking")
	}

	if err := s.routeRepo.UpdateSold(req.RouteID, route.Sold+req.Quantity); err != nil {
		s.logger.Error("update sold error", slog.Any("err", err))
	}

	return created, nil
}

func (s *service) GetMyBookings(userID int) ([]Booking, error) {
	bookings, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to fetch bookings")
	}
	return bookings, nil
}

func (s *service) GetByID(userID, bookingID int) (Booking, error) {
	b, err := s.repo.FindWithRelations(bookingID)
	if err != nil {
		return Booking{}, errors.New("booking not found")
	}
	if b.UserID != userID {
		return Booking{}, errors.New("forbidden: booking does not belong to you")
	}
	return b, nil
}

func (s *service) Cancel(userID, bookingID int) error {
	b, err := s.repo.GetByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}
	if b.UserID != userID {
		return errors.New("forbidden: booking does not belong to you")
	}
	if b.Status == StatusPaid {
		return errors.New("cannot cancel a paid booking")
	}
	if b.Status == StatusCancelled {
		return errors.New("booking already cancelled")
	}
	return s.repo.UpdateStatus(bookingID, StatusCancelled, nil)
}

func (s *service) HandleWebhook(payload WebhookPayload) error {
	if !s.verifySignature(payload) {
		return errors.New("invalid webhook signature")
	}

	b, err := s.repo.GetByExternalID(payload.OrderID)
	if err != nil {
		return errors.New("booking not found for order_id: " + payload.OrderID)
	}

	switch payload.TransactionStatus {
	case "settlement", "capture":
		now := time.Now()
		return s.repo.UpdateStatus(b.ID, StatusPaid, &now)
	case "expire":
		return s.repo.UpdateStatus(b.ID, StatusExpired, nil)
	case "cancel", "deny":
		return s.repo.UpdateStatus(b.ID, StatusCancelled, nil)
	}

	return nil
}

func (s *service) GetPaymentStatus(userID, bookingID int) (Booking, error) {
	b, err := s.repo.GetByID(bookingID)
	if err != nil {
		return Booking{}, errors.New("booking not found")
	}
	if b.UserID != userID {
		return Booking{}, errors.New("forbidden: booking does not belong to you")
	}
	return b, nil
}

func (s *service) verifySignature(p WebhookPayload) bool {
	raw := p.OrderID + p.StatusCode + p.GrossAmount + s.midtransKey
	h := sha512.New()
	h.Write([]byte(raw))
	expected := hex.EncodeToString(h.Sum(nil))
	return strings.EqualFold(expected, p.SignatureKey)
}
