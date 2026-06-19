package booking

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/fahmi0sd/ticketing-system/service/booking"
	"github.com/fahmi0sd/ticketing-system/util/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger     *slog.Logger
	bookingSvc booking.Service
}

func NewController(logger *slog.Logger, svc booking.Service) *Controller {
	return &Controller{logger: logger, bookingSvc: svc}
}

func (ctrl *Controller) Create(c echo.Context) error {
	userID := c.Get("id").(int)

	req := new(booking.CreateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}
	if err := validator.New().Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("validation failed: "+err.Error()))
	}

	created, err := ctrl.bookingSvc.Create(userID, *req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, response.Error(err.Error()))
		}
		if strings.Contains(err.Error(), "insufficient") {
			return c.JSON(http.StatusConflict, response.Error(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success("booking created, please complete payment", created))
}

func (ctrl *Controller) GetMyBookings(c echo.Context) error {
	userID := c.Get("id").(int)

	bookings, err := ctrl.bookingSvc.GetMyBookings(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("success", bookings))
}

func (ctrl *Controller) GetByID(c echo.Context) error {
	userID := c.Get("id").(int)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid booking id"))
	}

	b, err := ctrl.bookingSvc.GetByID(userID, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, response.Error(err.Error()))
		}
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, response.Error(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("success", b))
}

func (ctrl *Controller) Cancel(c echo.Context) error {
	userID := c.Get("id").(int)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid booking id"))
	}

	if err := ctrl.bookingSvc.Cancel(userID, id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, response.Error(err.Error()))
		}
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, response.Error(err.Error()))
		}
		if strings.Contains(err.Error(), "cannot cancel") || strings.Contains(err.Error(), "already cancelled") {
			return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("booking cancelled successfully", nil))
}

func (ctrl *Controller) HandleWebhook(c echo.Context) error {
	var payload booking.WebhookPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid webhook payload"))
	}

	if err := ctrl.bookingSvc.HandleWebhook(payload); err != nil {
		ctrl.logger.Error("webhook processing error", slog.Any("err", err))
		if strings.Contains(err.Error(), "invalid webhook signature") {
			return c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("webhook processed", nil))
}

func (ctrl *Controller) GetPaymentStatus(c echo.Context) error {
	userID := c.Get("id").(int)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid booking id"))
	}

	b, err := ctrl.bookingSvc.GetPaymentStatus(userID, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, response.Error(err.Error()))
		}
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, response.Error(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("success", b))
}
