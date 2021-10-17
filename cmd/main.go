package main

import (
	"2021_2_LostPointer/pkg/middleware"
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
	userHandlers       deliveryUser.UserDelivery
	musicHandlers      handlersMusic.MusicHandlers
	middlewareHandlers middleware.Middleware
}

func NewRequestHandler(db *sql.DB, redisConnection *redis.Client) *RequestHandlers {
	userDB := repositoryUser.NewUserRepository(db)
	redisStore := repositoryUser.NewRedisStore(redisConnection)
	fileSystem := repositoryUser.NewFileSystem()
	userUseCase := usecaseUser.NewUserUserCase(userDB, redisStore, fileSystem)
	userHandlers := deliveryUser.NewUserDelivery(userUseCase)

	musicRepo := repositoryMusic.NewMusicRepository(db)
	musicUseCase := usecaseMusic.NewMusicUseCase(musicRepo, userUseCase)
	musicHandlers := handlersMusic.NewMusicDelivery(musicUseCase)

	middlewareHandlers := middleware.NewMiddlewareHandler(userUseCase)

	api := &(RequestHandlers{
		userHandlers:       userHandlers,
		musicHandlers:      musicHandlers,
		middlewareHandlers: middlewareHandlers,
	})

	return api
}

func InitializeDatabase() *sql.DB {
	connectionString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
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
	var AddrConfig string
	if len(os.Getenv("REDIS_PORT")) == 0 {
		AddrConfig = os.Getenv("REDIS_HOST")
	} else {
		AddrConfig = fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	}

	redisConnection := redis.NewClient(&redis.Options{
		Addr:     AddrConfig,
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
	api.middlewareHandlers.InitMiddlewareHandlers(server)

	server.Logger.Fatal(server.Start(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))))
}
