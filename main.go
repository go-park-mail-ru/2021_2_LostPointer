package main

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"log"
	"os"

	deliveryUser "2021_2_LostPointer/internal/users/delivery"
	repositoryUser "2021_2_LostPointer/internal/users/repository"
	usecaseUser "2021_2_LostPointer/internal/users/usecase"
)

const redisDB = 1

type RequestHandlers struct {
	userHandler deliveryUser.UserDelivery
}

func NewRequestHandler(db *sql.DB, redisConnection *redis.Client) *RequestHandlers {
	userDB := repositoryUser.NewUserRepository(db)

	userUseCase := usecaseUser.NewUserUserCase(userDB, redisConnection)

	userH := deliveryUser.NewUserDelivery(userUseCase)

	api := &(RequestHandlers{
		userHandler: userH,
	})

	return api
}

func InitializeDatabase() *sql.DB {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBNAME"),
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalln("NO CONNECTION TO DATABASE", err.Error())
	}

	return db
}

func InitializeRedis() *redis.Client {
	redisConnection := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASS"),
		DB:       redisDB,
	})

	return redisConnection
}

func main() {
	server := echo.New()

	db := InitializeDatabase()
	defer func() {
		if db != nil {
			db.Close()
		}
	}()
	redisConnection := InitializeRedis()
	defer func() {
		if redisConnection != nil {
			redisConnection.Close()
		}
	}()

	api := NewRequestHandler(db, redisConnection)

	api.userHandler.InitHandlers(server)

	server.Logger.Fatal(server.Start(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))))
}
