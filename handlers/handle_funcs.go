package handlers

import (
	"2021_2_LostPointer/models"
	"database/sql"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

func LoginUserHandler(db *sql.DB) echo.HandlerFunc{
	return func(c echo.Context) error {
		var user models.User
		if err := c.Bind(&user); err != nil {
			return err
		}
		userExists, err := models.UserExistsLogin(db, user)
		if err != nil {
			return err
		}

		if !userExists {
			return c.JSON(http.StatusNotFound, "ERROR: User not found")
		}
		sessionToken, err := uuid.NewV4()
		if err != nil {
			return err
		}
		cookie := new(http.Cookie)
		cookie.Name = "Session_cookie"
		cookie.Value = sessionToken.String()
		cookie.HttpOnly = true
		cookie.Expires = time.Now().Add(365 * 24 * time.Hour)
		c.SetCookie(cookie)
		return c.JSON(http.StatusOK, "OK: We can authorize user")
	}
}

func SignUpHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		if err := c.Bind(&user); err != nil {
			return err
		}
		isUnique, err := models.IsUserUnique(db, user)
		if err != nil {
			return err
		}

		if !isUnique {
			return c.JSON(http.StatusBadRequest, "ERROR: User is not unique")
		}
		err = models.CreateUser(db, user)
		if err != nil {
			return err
		}
		sessionToken, err := uuid.NewV4()
		if err != nil {
			return err
		}
		cookie := new(http.Cookie)
		cookie.Name = "Session_cookie"
		cookie.Value = sessionToken.String()
		cookie.HttpOnly = true
		cookie.Expires = time.Now().Add(365 * 24 * time.Hour)
		c.SetCookie(cookie)
		return c.JSON(http.StatusCreated, "OK: User created")
	}
}

func GetHomePageHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "DEVELOPING")
	}
}