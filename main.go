package main

import (
	"2021_2_LostPointer/handlers"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
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
	e.GET("/", handlers.GetHomePageHandler(db))
	e.POST("/signup", handlers.SignUpHandler(db, redisConnection))
	e.POST("/signin", handlers.LoginUserHandler(db, redisConnection))
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))))
}
