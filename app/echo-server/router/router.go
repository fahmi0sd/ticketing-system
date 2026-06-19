package router

import (
	"net/http"

	"github.com/fahmi0sd/go-utils/middleware"
	bookingCtrl "github.com/fahmi0sd/ticketing-system/app/echo-server/controller/booking"
	routeCtrl "github.com/fahmi0sd/ticketing-system/app/echo-server/controller/route"
	userCtrl "github.com/fahmi0sd/ticketing-system/app/echo-server/controller/user"
	"github.com/labstack/echo/v4"
)

func RegisterPath(
	e *echo.Echo,
	jwtSecret string,
	ctrlUser *userCtrl.Controller,
	ctrlRoute *routeCtrl.Controller,
	ctrlBooking *bookingCtrl.Controller,
) {
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})

	jwtMiddleware := middleware.JWTMiddleware2(jwtSecret)

	api := e.Group("/api")

	users := api.Group("/users")
	users.POST("/register", ctrlUser.Register)
	users.POST("/login", ctrlUser.Login)
	users.GET("/me", ctrlUser.GetMe, jwtMiddleware)

	routes := api.Group("/routes")
	routes.GET("", ctrlRoute.GetAll)
	routes.GET("/:id", ctrlRoute.GetByID)

	bookings := api.Group("/bookings")
	bookings.POST("", ctrlBooking.Create, jwtMiddleware)
	bookings.GET("", ctrlBooking.GetMyBookings, jwtMiddleware)
	bookings.POST("/webhook", ctrlBooking.HandleWebhook) // public: dipanggil Midtrans
	bookings.GET("/:id", ctrlBooking.GetByID, jwtMiddleware)
	bookings.PUT("/:id/cancel", ctrlBooking.Cancel, jwtMiddleware)
	bookings.GET("/:id/payment", ctrlBooking.GetPaymentStatus, jwtMiddleware)
}
