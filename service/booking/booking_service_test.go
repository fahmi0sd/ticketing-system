package booking_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/fahmi0sd/ticketing-system/service/booking"
	routepkg "github.com/fahmi0sd/ticketing-system/service/route"
	"github.com/stretchr/testify/assert"
)

type mockBookingRepo struct {
	createFn func(b booking.Booking) (booking.Booking, error)
}

func (m *mockBookingRepo) Create(b booking.Booking) (booking.Booking, error) {
	return m.createFn(b)
}
func (m *mockBookingRepo) GetByID(id int) (booking.Booking, error) {
	return booking.Booking{}, nil
}
func (m *mockBookingRepo) GetByUserID(userID int) ([]booking.Booking, error) {
	return nil, nil
}
func (m *mockBookingRepo) GetByExternalID(externalID string) (booking.Booking, error) {
	return booking.Booking{}, nil
}
func (m *mockBookingRepo) UpdateStatus(id int, status booking.Status, paidAt *time.Time) error {
	return nil
}
func (m *mockBookingRepo) FindWithRelations(id int) (booking.Booking, error) {
	return booking.Booking{}, nil
}

type mockRouteRepo struct {
	getByIDFn    func(id int) (routepkg.Route, error)
	updateSoldFn func(id, newSold int) error
}

func (m *mockRouteRepo) GetAll(filter routepkg.SearchFilter) ([]routepkg.Route, error) {
	return nil, nil
}
func (m *mockRouteRepo) GetByID(id int) (routepkg.Route, error) {
	return m.getByIDFn(id)
}
func (m *mockRouteRepo) UpdateSold(id, newSold int) error {
	return m.updateSoldFn(id, newSold)
}

type mockMidtrans struct {
	createPaymentFn func(externalID string, amount float64) (string, error)
}

func (m *mockMidtrans) CreatePayment(externalID string, amount float64) (string, error) {
	return m.createPaymentFn(externalID, amount)
}

func newTestService(
	bookingRepo booking.Repository,
	routeRepo routepkg.Repository,
	midtrans booking.MidtransClient,
) booking.Service {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return booking.NewService(logger, bookingRepo, routeRepo, midtrans, "test-server-key")
}

func TestCreate_Success(t *testing.T) {
	bookingRepo := &mockBookingRepo{
		createFn: func(b booking.Booking) (booking.Booking, error) {
			b.ID = 1
			return b, nil
		},
	}
	routeRepo := &mockRouteRepo{
		getByIDFn: func(id int) (routepkg.Route, error) {
			return routepkg.Route{ID: 1, Price: 150000, Quota: 40, Sold: 10}, nil
		},
		updateSoldFn: func(id, newSold int) error { return nil },
	}
	midtrans := &mockMidtrans{
		createPaymentFn: func(externalID string, amount float64) (string, error) {
			return "https://sandbox.midtrans.com/pay/xxx", nil
		},
	}

	svc := newTestService(bookingRepo, routeRepo, midtrans)
	result, err := svc.Create(1, booking.CreateRequest{RouteID: 1, Quantity: 2})

	assert.NoError(t, err)
	assert.Equal(t, booking.StatusPending, result.Status)
	assert.Equal(t, 300000.0, result.TotalPrice)
	assert.NotEmpty(t, result.PaymentURL)
}

func TestCreate_QuotaExceeded(t *testing.T) {
	routeRepo := &mockRouteRepo{
		getByIDFn: func(id int) (routepkg.Route, error) {
			return routepkg.Route{ID: 1, Price: 150000, Quota: 5, Sold: 4}, nil
		},
	}

	svc := newTestService(&mockBookingRepo{}, routeRepo, &mockMidtrans{})
	_, err := svc.Create(1, booking.CreateRequest{RouteID: 1, Quantity: 2})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient seats")
}

func TestCreate_RouteNotFound(t *testing.T) {
	routeRepo := &mockRouteRepo{
		getByIDFn: func(id int) (routepkg.Route, error) {
			return routepkg.Route{}, errors.New("record not found")
		},
	}

	svc := newTestService(&mockBookingRepo{}, routeRepo, &mockMidtrans{})
	_, err := svc.Create(1, booking.CreateRequest{RouteID: 999, Quantity: 1})

	assert.Error(t, err)
	assert.Equal(t, "route not found", err.Error())
}

func TestCreate_MidtransFailure(t *testing.T) {
	routeRepo := &mockRouteRepo{
		getByIDFn: func(id int) (routepkg.Route, error) {
			return routepkg.Route{ID: 1, Price: 150000, Quota: 40, Sold: 10}, nil
		},
	}
	midtrans := &mockMidtrans{
		createPaymentFn: func(externalID string, amount float64) (string, error) {
			return "", errors.New("midtrans timeout")
		},
	}

	svc := newTestService(&mockBookingRepo{}, routeRepo, midtrans)
	_, err := svc.Create(1, booking.CreateRequest{RouteID: 1, Quantity: 1})

	assert.Error(t, err)
	assert.Equal(t, "failed to create payment link", err.Error())
}
