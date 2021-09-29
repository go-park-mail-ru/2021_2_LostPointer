package handlers

import (
	"2021_2_LostPointer/models"
	"2021_2_LostPointer/utils"
	"database/sql"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"time"
)

type Arguments struct {
	db              *sql.DB
	redisConnection *redis.Client
}

func LoginUser(c echo.Context, args *Arguments) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return err
	}
	userID, err := utils.UserExistsLogin(args.db, user)
	if err != nil {
		return err
	}
	if userID == 0 {
		return c.JSON(http.StatusNotFound, "ERROR: User not found")
	}

	sessionToken := utils.GetRandomString(40)
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "Session_cookie"
	cookie.Path = "/"
	cookie.Domain = "http://localhost:3000"
	cookie.SameSite = http.SameSiteNoneMode
	cookie.Secure = true
	cookie.Value = sessionToken
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(time.Hour * 24 * 30)
	c.SetCookie(cookie)

	err = utils.StoreSession(args.redisConnection,
		&models.Session{UserID: userID, Session: cookie.Value})
	if err != nil {
		log.Fatalln(err)
	}

	return c.JSON(http.StatusOK, "OK: User is authorized")
}

func LoginUserHandler(db *sql.DB, redisConnection *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := LoginUser(c, &Arguments{db: db, redisConnection: redisConnection})
		if err != nil {
			return err
		}
		return nil
	}
}

func SignUp(c echo.Context, args *Arguments) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return err
	}
	isUnique, err := utils.IsUserUnique(args.db, user)
	if err != nil {
		return err
	}

	if !isUnique {
		return c.JSON(http.StatusBadRequest, "ERROR: User is not unique")
	}
	userID, err := utils.CreateUser(args.db, user)
	if err != nil {
		return err
	}

	sessionToken := utils.GetRandomString(40)
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "Session_cookie"
	cookie.Value = sessionToken
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Domain = "http://localhost:3000"
	cookie.SameSite = http.SameSiteNoneMode
	cookie.Secure = true
	cookie.Domain = "http://localhost:3000"
	cookie.Expires = time.Now().Add(time.Hour * 24 * 30)
	c.SetCookie(cookie)

	err = utils.StoreSession(args.redisConnection,
		&models.Session{UserID: userID, Session: cookie.Value})
	if err != nil {
		log.Fatalln(err)
	}

	return c.JSON(http.StatusCreated, "OK: User created")
}

func SignUpHandler(db *sql.DB, redisConnection *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := SignUp(c, &Arguments{db: db, redisConnection: redisConnection})
		if err != nil {
			return err
		}
		return nil
	}
}

func LogoutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Session_cookie")
		if err != nil {
			return err
		}
		cookie.Expires = time.Now().AddDate(0, 0, -1)
		c.SetCookie(cookie)

		return c.NoContent(http.StatusOK)
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

func AuthHandler(redisConnection *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Session_cookie")
		if err != nil {
			log.Println(err)
			return err
		}
		id, err := utils.GetSessionUser(redisConnection, cookie.Value)
		if err != nil {
			log.Println(err)
			return err
		}
		return c.JSON(http.StatusOK, id)
	}
}
