package main

import (
	"2021_2_LostPointer/handlers"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func CORSMiddlewareWrapper(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		dynamicCORSConfig := middleware.CORSConfig{
			AllowOrigins:     []string{req.Header.Get("Origin")}, // Разобраться и поменять!!!
			AllowHeaders:     []string{"Accept", "Cache-Control", "Content-Type", "X-Requested-With"},
			AllowCredentials: true,
		}
		CORSMiddleware := middleware.CORSWithConfig(dynamicCORSConfig)
		CORSHandler := CORSMiddleware(next)
		return CORSHandler(ctx)
	}
}

func main() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBHOST"), os.Getenv("DBPORT"),
		os.Getenv("DBNAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	redisConnection := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       1,
	})

	e := echo.New()
	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"http://lostpointer.site", "http://localhost"},
	//	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	//}))

	e.Use(CORSMiddlewareWrapper)

	//e.Use(middleware.CORS())

	e.GET("/api/v1/home", handlers.GetHomePageHandler(db))
	e.GET("/api/v1/auth", handlers.AuthHandler(redisConnection))
	e.POST("/api/v1/user/signup", handlers.SignUpHandler(db, redisConnection))
	e.POST("/api/v1/user/signin", handlers.LoginUserHandler(db, redisConnection))
	e.DELETE("/api/v1/user/signin", handlers.LogoutHandler())

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))))
}
