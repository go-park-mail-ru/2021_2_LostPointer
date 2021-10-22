package main

import (
	deliveryArtist "2021_2_LostPointer/internal/artist/delivery"
	repositoryArtist "2021_2_LostPointer/internal/artist/repository"
	usecaseArtist "2021_2_LostPointer/internal/artist/usecase"
	middleware "2021_2_LostPointer/internal/middleware"

	deliveryTrack "2021_2_LostPointer/internal/track/delivery"
	repositoryTrack "2021_2_LostPointer/internal/track/repository"
	usecaseTrack "2021_2_LostPointer/internal/track/usecase"

	deliveryAlbum "2021_2_LostPointer/internal/album/delivery"
	repositoryAlbum "2021_2_LostPointer/internal/album/repository"
	usecaseAlbum "2021_2_LostPointer/internal/album/usecase"

	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"

	deliveryUser "2021_2_LostPointer/internal/users/delivery"
	repositoryUser "2021_2_LostPointer/internal/users/repository"
	usecaseUser "2021_2_LostPointer/internal/users/usecase"
)

const redisDB = 1

type RequestHandlers struct {
	userHandlers deliveryUser.UserDelivery
	//musicHandlers      handlersMusic.MusicHandlers
	artistHandlers     deliveryArtist.ArtistDelivery
	trackHandlers      deliveryTrack.TrackDelivery
	albumHandlers      deliveryAlbum.AlbumDelivery
	middlewareHandlers middleware.Middleware
}

func NewRequestHandler(db *sql.DB, redisConnection *redis.Client, logger *zap.SugaredLogger) *RequestHandlers {
	userDB := repositoryUser.NewUserRepository(db)
	redisStore := repositoryUser.NewRedisStore(redisConnection)
	fileSystem := repositoryUser.NewFileSystem()
	userUseCase := usecaseUser.NewUserUserCase(userDB, redisStore, fileSystem)
	userHandlers := deliveryUser.NewUserDelivery(logger, userUseCase)

	//musicRepo := repositoryMusic.NewMusicRepository(db)
	//musicUseCase := usecaseMusic.NewMusicUseCase(musicRepo)
	//musicHandlers := handlersMusic.NewMusicDelivery(musicUseCase, logger)

	artistRepo := repositoryArtist.NewArtistRepository(db)
	artistUseCase := usecaseArtist.NewArtistUseCase(artistRepo)
	artistHandlers := deliveryArtist.NewArtistDelivery(artistUseCase, logger)

	trackRepo := repositoryTrack.NewTrackRepository(db)
	trackUseCase := usecaseTrack.NewTrackUseCase(trackRepo)
	trackHandlers := deliveryTrack.NewTrackDelivery(trackUseCase, logger)

	albumRepo := repositoryAlbum.NewAlbumRepository(db)
	albumUseCase := usecaseAlbum.NewAlbumUseCase(albumRepo)
	albumHandlers := deliveryAlbum.NewAlbumDelivery(albumUseCase, logger)

	middlewareHandlers := middleware.NewMiddlewareHandler(logger, userUseCase)

	api := &(RequestHandlers{
		userHandlers: userHandlers,
		//musicHandlers:      musicHandlers,
		artistHandlers:     artistHandlers,
		trackHandlers:      trackHandlers,
		albumHandlers:      albumHandlers,
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
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

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

	api := NewRequestHandler(db, redisConnection, logger)

	api.userHandlers.InitHandlers(server)
	//api.musicHandlers.InitHandlers(server)
	api.artistHandlers.InitHandlers(server)
	api.trackHandlers.InitHandlers(server)
	api.albumHandlers.InitHandlers(server)
	api.middlewareHandlers.InitMiddlewareHandlers(server)

	server.Logger.Fatal(server.Start(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))))
}
