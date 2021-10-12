package main

import (
	handlersMusic "2021_2_LostPointer/pkg/music/delivery"
	repositoryMusic "2021_2_LostPointer/pkg/music/repository"
	usecaseMusic "2021_2_LostPointer/pkg/music/usecase"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"log"
	"os"

	deliveryUser "2021_2_LostPointer/pkg/users/delivery"
	repositoryUser "2021_2_LostPointer/pkg/users/repository"
	usecaseUser "2021_2_LostPointer/pkg/users/usecase"
)

const redisDB = 1

type RequestHandlers struct {
	userHandlers  deliveryUser.UserDelivery
	musicHandlers handlersMusic.MusicHandlers
}

func NewRequestHandler(db *sql.DB, redisConnection *redis.Client) *RequestHandlers {
	userDB := repositoryUser.NewUserRepository(db)
	redisStore := repositoryUser.NewRedisStore(redisConnection)
	userUseCase := usecaseUser.NewUserUserCase(userDB, redisStore)
	userHandlers := deliveryUser.NewUserDelivery(userUseCase)

	musicHandlers := handlersMusic.NewMusicHandlers(usecaseMusic.NewMusicUseCase(repositoryMusic.NewMusicRepository(db)))

	api := &(RequestHandlers{
		userHandlers:  userHandlers,
		musicHandlers: musicHandlers,
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

	api.userHandlers.InitHandlers(server)
	api.musicHandlers.InitHandlers(server)

	server.Logger.Fatal(server.Start(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))))
}
