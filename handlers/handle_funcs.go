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
		return c.JSON(http.StatusUnauthorized,
			&models.Response{Message: "Wrong username and/or password"})
	}

	sessionToken := utils.GetRandomString(SessionTokenLength)

	cookie := &http.Cookie{ // Setting up cookies
		Name:     "Session_cookie",
		Value:    sessionToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
	}
	c.SetCookie(cookie)

	err = utils.StoreSession(args.redisConnection,
		&models.Session{UserID: userID, Session: cookie.Value})
	if err != nil {
		log.Fatalln(err)
	}

	return c.JSON(http.StatusOK, &models.Response{Message: "User is authorized"})
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

	isCorrect, message, err := utils.ValidateSignUp(&user)
	if err != nil {
		return err
	}
	if !isCorrect {
		return c.JSON(http.StatusBadRequest, &models.Response{Message: message})
	}

	isEmailUnique, err := utils.IsUserEmailUnique(args.db, user.Email)
	if err != nil {
		return err
	}
	if !isEmailUnique {
		return c.JSON(http.StatusBadRequest, &models.Response{Message: "Email is already taken"})
	}
	isNicknameUnique, err := utils.IsUserNicknameUnique(args.db, user.Nickname)
	if err != nil {
		return err
	}
	if !isNicknameUnique {
		return c.JSON(http.StatusBadRequest, &models.Response{Message: "Nickname is already taken"})
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
		Name:     "Session_cookie",
		Value:    sessionToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
	}
	c.SetCookie(cookie)

	err = utils.StoreSession(args.redisConnection,
		&models.Session{UserID: userID, Session: cookie.Value})
	if err != nil {
		log.Fatalln(err)
	}

	return c.JSON(http.StatusCreated, &models.Response{Message: "User is created"})
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
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusUnauthorized, &models.Response{Message: "Unauthorized"})
		}
		redisConnection.Del(cookie.Value)
		cookie.Expires = time.Now().AddDate(0, 0, -1)
		cookie.Path = "/"
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, &models.Response{Message: "Logged out"})
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
			return c.JSON(http.StatusUnauthorized,
				&models.Response{Message: "User not authorized"})
		}
		id, err := utils.GetSessionUser(redisConnection, cookie.Value)
		if id == 0 {
			return c.JSON(http.StatusUnauthorized,
				&models.Response{Message: "User not authorized"})
		}
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK,
			&models.Response{Message: "User is authorized"})
	}
}
