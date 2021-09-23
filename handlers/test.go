package handlers

import (
	"github.com/labstack/echo"
	"net/http"
)

func CreateUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Test: createUserHandler")
	}
}

func LoginUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Test: loginUserHandler")
	}
}

func GetHomePageHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Test: homePageHandler")
	}
}
