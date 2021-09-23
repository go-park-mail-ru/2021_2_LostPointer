package main

import (
	"2021_2_LostPointer/handlers"
	"database/sql"
	"fmt"
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

	e := echo.New()
	e.GET("/", handlers.GetHomePageHandler())
	e.POST("/signup", handlers.CreateUserHandler(db))
	e.POST("/signin", handlers.LoginUserHandler(db)) // This one works
	e.Logger.Fatal(e.Start(":3030"))
}