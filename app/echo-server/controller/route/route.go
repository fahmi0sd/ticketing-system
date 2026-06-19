package route

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/fahmi0sd/ticketing-system/service/route"
	"github.com/fahmi0sd/ticketing-system/util/response"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger   *slog.Logger
	routeSvc route.Service
}

func NewController(logger *slog.Logger, svc route.Service) *Controller {
	return &Controller{logger: logger, routeSvc: svc}
}

func (ctrl *Controller) GetAll(c echo.Context) error {
	filter := route.SearchFilter{
		Origin:      c.QueryParam("origin"),
		Destination: c.QueryParam("destination"),
		Type:        c.QueryParam("type"),
		Date:        c.QueryParam("date"),
	}

	routes, err := ctrl.routeSvc.GetAll(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("success", routes))
}

func (ctrl *Controller) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid route id"))
	}

	rt, err := ctrl.routeSvc.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("success", rt))
}
