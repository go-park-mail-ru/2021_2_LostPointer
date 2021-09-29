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

const SessionTokenLength = 40

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

	sessionToken := utils.GetRandomString(SessionTokenLength)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{ // Setting up cookies
		Name: "Session_cookie",
		Value: sessionToken,
		Path: "/",
		Secure: true,
		HttpOnly: true,
		Expires: time.Now().Add(time.Hour * 24 * 30),
	}
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

	sessionToken := utils.GetRandomString(SessionTokenLength)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{ // Setting up cookies
		Name: "Session_cookie",
		Value: sessionToken,
		Path: "/",
		Secure: true,
		HttpOnly: true,
		Expires: time.Now().Add(time.Hour * 24 * 30),
	}
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

func LogoutHandler(redisConnection *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Session_cookie")
		log.Println("Cookies: ", cookie)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(cookie.Value)
		redisConnection.Del(cookie.Value)
		cookie.Expires = time.Now().AddDate(0, 0, -1)
		c.SetCookie(cookie)

		return c.NoContent(http.StatusOK)
	}
}

func GetHomePageHandler(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		selectionForHomePage, err := utils.GetSelectionForHomePage(db)
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
