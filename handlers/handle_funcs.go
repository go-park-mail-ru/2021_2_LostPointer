package handlers

import (
	"2021_2_LostPointer/models"
	"2021_2_LostPointer/utils"
	"database/sql"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

func LoginUserHandler(db *sql.DB, redisConnection *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		if err := c.Bind(&user); err != nil {
			return err
		}
		userID, err := utils.UserExistsLogin(db, user)
		if err != nil {
			return err
		}
		if userID == 0 {
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
		cookie.Expires = time.Now().Add(time.Hour)
		c.SetCookie(cookie)

		err = utils.StoreSession(redisConnection,
			&models.Session{UserID: userID, Session: cookie.Value})
		if err != nil {
			log.Fatalln(err)
		}

		return c.JSON(http.StatusOK, "OK: We can authorize user")
	}
}

func SignUpHandler(db *sql.DB, redisConnection *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		if err := c.Bind(&user); err != nil {
			return err
		}
		isUnique, err := utils.IsUserUnique(db, user)
		if err != nil {
			return err
		}

		if !isUnique {
			return c.JSON(http.StatusBadRequest, "ERROR: User is not unique")
		}
		userID, err := utils.CreateUser(db, user)
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
		cookie.Expires = time.Now().Add(time.Hour)
		c.SetCookie(cookie)

		err = utils.StoreSession(redisConnection,
			&models.Session{UserID: userID, Session: cookie.Value})
		if err != nil {
			log.Fatalln(err)
		}

		return c.JSON(http.StatusCreated, "OK: User created")
	}
}

func GetHomePageHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		selectionForHomePage, err := models.GetSelectionForHomePage(db)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, selectionForHomePage)
	}
}
