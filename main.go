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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://lostpointer.site", "http://localhost:3000"},
		AllowHeaders:     []string{"Accept", "Cache-Control", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
	}))
	e.GET("/api/v1/home", handlers.GetHomePageHandler(db))
	e.GET("/api/v1/auth", handlers.AuthHandler(redisConnection))
	e.POST("/api/v1/user/signup", handlers.SignUpHandler(db, redisConnection))
	e.POST("/api/v1/user/signin", handlers.LoginUserHandler(db, redisConnection))
	e.POST("/api/v1/user/logout", handlers.LogoutHandler(redisConnection))

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))))
}
