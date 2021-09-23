package handlers

import (
	"2021_2_LostPointer/models"
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func CreateUserHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Test: createUserHandler")
	}
}

func LoginUserHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		c.Bind(&user)
		isExists, err := models.IsUserExists(db, &user)
		if err == nil {
			if !isExists {
				return c.JSON(http.StatusNotFound, "User not registered")
			}
			return c.JSON(http.StatusOK, "User exists")
		} else {
			fmt.Println(err)
			return err
		}
	}
}

func GetHomePageHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Test: homePageHandler")
	}
}
