package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/ednesic/coursemanagement/cache"
	"github.com/ednesic/coursemanagement/handlers"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.elastic.co/apm/module/apmechov4"
)

func main() {
	var err error
	// var serviceName = "course"
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if os.Getenv("ENV") == "prod" {
		e.Logger.SetLevel(log.INFO)
	}

	cache.GetInstance().Initialize(map[string]string{"server1": os.Getenv("REDIS_HOST")})
	err = storage.GetInstance().Initialize(
		ctx,
		os.Getenv("DB_HOST"),
		os.Getenv("DB"),
	)

	if err != nil {
		e.Logger.Fatal("Could not resolve Data access layer: ", err)
	}

	// apm.DefaultTracer.Service.Name = serviceName
	// opentracing.InitGlobalTracer(apmot.New())

	// e.Use(opentracingMid.OpenTracing(serviceName))
	e.Use(apmechov4.Middleware())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("2M"))
	//e.Use(metrics.NewMetric())
	e.Use(middleware.Logger())

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	gCourse := e.Group("/courses")
	gCourse.DELETE("/:name", handlers.DelCourse)
	gCourse.GET("/:name", handlers.GetCourse)
	gCourse.GET("", handlers.GetCourses)
	gCourse.POST("", handlers.SetCourse)
	gCourse.PUT("", handlers.PutCourse)

	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	cache.GetInstance().Disconnect()
	storage.GetInstance().Disconnect()
}
