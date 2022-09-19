package main

import (
	"fmt"
	"keyvault-demo/config"
	"keyvault-demo/domain/entity"
	"keyvault-demo/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type controller struct {
	echo       *echo.Echo
	studentSrv service.StudentService
}

func NewController(e *echo.Echo) controller {
	return controller{
		echo:       e,
		studentSrv: service.NewStudentService(),
	}
}

func (ctr *controller) createStudent(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var student = new(entity.Student)

	c.Bind(student)

	id, err := ctr.studentSrv.Create(ctx, *student)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error creating student", err)
	}

	stdResult, err := ctr.studentSrv.FindByID(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error get student", err)
	}

	return c.JSON(http.StatusCreated, stdResult)
}

func (ctr *controller) findAll(c echo.Context) (err error) {
	ctx := c.Request().Context()

	students, err := ctr.studentSrv.FindAll(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, students)
}

func (ctr *controller) rotate(c echo.Context) (err error) {
	ctx := c.Request().Context()

	success, fails, err := ctr.studentSrv.Rotate(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"succcess": success,
		"fails":    fails,
	})
}

func (ctr *controller) Start() {
	s := ctr.echo.Group("student")
	s.GET("", ctr.findAll)
	s.POST("", ctr.createStudent)
	s.GET("/rotate", ctr.rotate)
}

func main() {
	config.ReadConfig(".env")
	e := echo.New()
	e.Use(middleware.Logger())

	sctr := NewController(e)
	sctr.Start()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Configuration().Port)))
}
