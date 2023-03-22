// Package appconfig implements a main routine that is the starting point of
// location-mgmt microservice.
package appconfig

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	middleKit "github.com/oboadagd/kit-go/middleware/echo"
	pgKit "github.com/oboadagd/kit-go/postgresql"
	"github.com/oboadagd/location-common/recordtype"
	"github.com/oboadagd/location-history-mgmt/controller"
	"github.com/oboadagd/location-history-mgmt/migration"
	"github.com/oboadagd/location-history-mgmt/repository"
	"github.com/oboadagd/location-history-mgmt/router"
	"github.com/oboadagd/location-history-mgmt/service"
	grpcserver "github.com/oboadagd/location-history-mgmt/userlocation/server"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// StartApp implements configuration and start-up of microservice.
// It also runs a goroutine that starts up of grpc server.
func StartApp() {

	echoInstance := echo.New()

	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(echoInstance)

	echoInstance.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency_human=${latency_human}",
	}))
	echoInstance.Use(middleware.Recover())

	if err := envconfig.Process("LIST", &recordtype.Cfg); err != nil {
		err = errors.Wrap(err, "parse environment variables")
		return
	}

	db := pgKit.NewPgDB(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", recordtype.Cfg.DBHost, recordtype.Cfg.DBPort),
		User:     recordtype.Cfg.DBUser,
		Password: recordtype.Cfg.DBPass,
		Database: recordtype.Cfg.DBName,
	})
	migration.Init(db)

	locationRepository := repository.NewLocationRepository(db)
	locationHistoryRepository := repository.NewLocationHistoryRepository(db)
	locationService := service.NewLocationService(locationRepository, locationHistoryRepository)
	locationController := controller.NewLocationController(locationService)

	errorHandlerMiddle := middleKit.NewErrorHandlerMiddleware()

	r := router.NewRouter(echoInstance, locationController, errorHandlerMiddle)
	r.Init()

	go func() {
		log.Infof(grpcserver.GrpcServe(locationService, echoInstance.AcquireContext()).Error())
	}()

	// Start server
	go func() {
		if err := echoInstance.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGSTOP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := echoInstance.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
