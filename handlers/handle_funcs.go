package handlers

import (
	"2021_2_LostPointer/models"
	"database/sql"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

func CreateUserHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Test: createUserHandler")
	}
}

func LoginUserHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		err := c.Bind(&user)
		if err != nil {
			return err
		}
		isExists, err := models.IsUserExists(db, &user)
		if err != nil {
			return err
		}
		if !isExists {
			return c.JSON(http.StatusNotFound, "ERR: User is not registered")
		}
		sessionToken, err := uuid.NewV4()
		if err != nil {
			return err
		}
		cookie := &http.Cookie{
			Name: "Session_token",
			Value: sessionToken.String(),
			HttpOnly: true,
			Expires: time.Now().Add(365 * 24 * time.Hour),
		}
		c.SetCookie(cookie)
		return c.JSON(http.StatusOK, "OK: User is registered")
	}
}

func GetHomePageHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Test: homePageHandler")
	}
}
