package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fahmi0sd/go-utils/postgres"
	bookingCtrl "github.com/fahmi0sd/ticketing-system/app/echo-server/controller/booking"
	routeCtrl "github.com/fahmi0sd/ticketing-system/app/echo-server/controller/route"
	userCtrl "github.com/fahmi0sd/ticketing-system/app/echo-server/controller/user"
	"github.com/fahmi0sd/ticketing-system/app/echo-server/router"
	midtrans "github.com/fahmi0sd/ticketing-system/pkg"
	bookingRepo "github.com/fahmi0sd/ticketing-system/repository/booking"
	routeRepo "github.com/fahmi0sd/ticketing-system/repository/route"
	userRepo "github.com/fahmi0sd/ticketing-system/repository/user"
	bookingSvc "github.com/fahmi0sd/ticketing-system/service/booking"
	routeSvc "github.com/fahmi0sd/ticketing-system/service/route"
	userSvc "github.com/fahmi0sd/ticketing-system/service/user"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	// Database Connection
	database := postgres.GetPostgresConnection()
	logger.Info("database connected")

	// Config
	jwtSecret := os.Getenv("JWT_SECRET")
	midtransKey := os.Getenv("MIDTRANS_SERVER_KEY")
	midtransSnapURL := os.Getenv("MIDTRANS_SNAP_URL")

	midtransClient := midtrans.NewClient(midtransKey, midtransSnapURL)
	logger.Info("midtrans client initialized", slog.String("snap_url", midtransSnapURL))

	// Repositories
	usrRepo := userRepo.NewGormRepository(database)
	rtRepo := routeRepo.NewGormRepository(database)
	bkgRepo := bookingRepo.NewGormRepository(database)

	// Services
	usrSvc := userSvc.NewService(logger, usrRepo, jwtSecret)
	rtSvc := routeSvc.NewService(logger, rtRepo)
	bkgSvc := bookingSvc.NewService(logger, bkgRepo, rtRepo, midtransClient, midtransKey)

	// Controllers
	usrCtrl := userCtrl.NewController(logger, usrSvc)
	rtCtrl := routeCtrl.NewController(logger, rtSvc)
	bkgCtrl := bookingCtrl.NewController(logger, bkgSvc)

	// Echo
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","level":"INFO","method":"${method}","uri":"${uri}","status":${status},"latency_human":"${latency_human}"}` + "\n",
	}))
	e.Pre(middleware.RemoveTrailingSlash())

	router.RegisterPath(e, jwtSecret, usrCtrl, rtCtrl, bkgCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	go func() {
		if err := e.Start(addr); err != http.ErrServerClosed {
			log.Fatal("server error: " + err.Error())
		}
	}()

	logger.Info("Ticketing API running", slog.String("addr", addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("failed to shutdown server:", err)
	}

	logger.Info("server shutdown gracefully")
}
