package user

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/fahmi0sd/ticketing-system/service/user"
	"github.com/fahmi0sd/ticketing-system/util/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger  *slog.Logger
	userSvc user.Service
}

func NewController(logger *slog.Logger, svc user.Service) *Controller {
	return &Controller{logger: logger, userSvc: svc}
}

type registerRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required,min=6"`
	Phone    string `json:"phone"     validate:"omitempty"`
}

type loginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (ctrl *Controller) Register(c echo.Context) error {
	req := new(registerRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}
	if err := validator.New().Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("validation failed: "+err.Error()))
	}

	created, err := ctrl.userSvc.Register(user.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
	})
	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			return c.JSON(http.StatusConflict, response.Error(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success("user registered successfully", created))
}

func (ctrl *Controller) Login(c echo.Context) error {
	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}
	if err := validator.New().Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("email and password are required"))
	}

	token, err := ctrl.userSvc.Login(req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, response.Error(err.Error()))
		}
		if strings.Contains(err.Error(), "wrong") {
			return c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("login successful", map[string]string{
		"token": token,
	}))
}

func (ctrl *Controller) GetMe(c echo.Context) error {
	userID := c.Get("id").(int)

	u, err := ctrl.userSvc.GetMe(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("success", u))
}
