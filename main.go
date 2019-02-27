package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"

	"github.com/ednesic/coursemanagement/handlers"
	"github.com/ednesic/coursemanagement/services"
	"github.com/ednesic/coursemanagement/servivemanager"
)

func main() {
	var err error
	e := echo.New()

	servivemanager.CourseService, err = services.NewCourseService(os.Getenv("COURSE_DB_HOST"), os.Getenv("COURSE_DB"))

	if err != nil {
		e.Logger.Fatal("Could not resolve course service", err)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("2M"))

	//e.Server.ReadTimeout = time.Duration(1 * time.Second)
	//e.Server.WriteTimeout= time.Duration(1 * time.Second)

	e.GET("/courses/:name", handlers.GetCourse)
	e.GET("/courses", handlers.GetCourses)
	e.PUT("/courses", handlers.PutCourse)
	e.POST("/courses", handlers.SetCourse)
	e.DELETE("/courses/:name", handlers.DelCourse)

	e.Logger.Fatal(e.Start(":9090"))
}