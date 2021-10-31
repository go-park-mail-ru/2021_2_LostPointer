package main

import (
	deliveryArtist "2021_2_LostPointer/internal/artist/delivery"
	repositoryArtist "2021_2_LostPointer/internal/artist/repository"
	usecaseArtist "2021_2_LostPointer/internal/artist/usecase"
	"2021_2_LostPointer/internal/middleware"

	deliveryTrack "2021_2_LostPointer/internal/track/delivery"
	repositoryTrack "2021_2_LostPointer/internal/track/repository"
	usecaseTrack "2021_2_LostPointer/internal/track/usecase"

	deliveryAlbum "2021_2_LostPointer/internal/album/delivery"
	repositoryAlbum "2021_2_LostPointer/internal/album/repository"
	usecaseAlbum "2021_2_LostPointer/internal/album/usecase"

	deliveryPlaylist "2021_2_LostPointer/internal/playlist/delivery"
	repositoryPlaylist "2021_2_LostPointer/internal/playlist/repository"
	usecasePlaylist "2021_2_LostPointer/internal/playlist/usecase"

	deliveryQueue "2021_2_LostPointer/internal/queues/delivery"
	repositoryQueue "2021_2_LostPointer/internal/queues/repository"
	usecaseQueue "2021_2_LostPointer/internal/queues/usecase"

	deliveryUser "2021_2_LostPointer/internal/users/delivery"
	repositoryUser "2021_2_LostPointer/internal/users/repository"
	usecaseUser "2021_2_LostPointer/internal/users/usecase"

	authorizationMicro "2021_2_LostPointer/internal/microservices/authorization/delivery"

	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

type RequestHandlers struct {
	userHandlers       deliveryUser.UserDelivery
	artistHandlers     deliveryArtist.ArtistDelivery
	trackHandlers      deliveryTrack.TrackDelivery
	albumHandlers      deliveryAlbum.AlbumDelivery
	playlistHandlers   deliveryPlaylist.PlaylistDelivery
	middlewareHandlers middleware.Middleware
	queueHandlers      deliveryQueue.QueueDelivery
}

func NewRequestHandler(db *sql.DB, redisConnQueue *redis.Client, logger *zap.SugaredLogger,	sessionChecker authorizationMicro.SessionCheckerClient) *RequestHandlers {

	userDB := repositoryUser.NewUserRepository(db)
	fileSystem := repositoryUser.NewFileSystem()
	userUseCase := usecaseUser.NewUserUserCase(userDB, fileSystem, sessionChecker)
	userHandlers := deliveryUser.NewUserDelivery(logger, userUseCase)

	artistRepo := repositoryArtist.NewArtistRepository(db)
	artistUseCase := usecaseArtist.NewArtistUseCase(artistRepo)
	artistHandlers := deliveryArtist.NewArtistDelivery(artistUseCase, logger)

	trackRepo := repositoryTrack.NewTrackRepository(db)
	trackUseCase := usecaseTrack.NewTrackUseCase(trackRepo)
	trackHandlers := deliveryTrack.NewTrackDelivery(trackUseCase, logger)

	albumRepo := repositoryAlbum.NewAlbumRepository(db)
	albumUseCase := usecaseAlbum.NewAlbumUseCase(albumRepo)
	albumHandlers := deliveryAlbum.NewAlbumDelivery(albumUseCase, logger)

	playlistRepo := repositoryPlaylist.NewPlaylistRepository(db)
	playlistUseCase := usecasePlaylist.NewPlaylistUseCase(playlistRepo)
	playlistHandlers := deliveryPlaylist.NewPlaylistDelivery(playlistUseCase, logger)

	queueRepo := repositoryQueue.NewQueueRepository(db, redisConnQueue)
	queueUseCase := usecaseQueue.NewQueueUseCase(queueRepo)
	queueHandlers := deliveryQueue.NewQueueDelivery(queueUseCase, logger)

	middlewareHandlers := middleware.NewMiddlewareHandler(logger, userUseCase, sessionChecker)

	api := &(RequestHandlers{
		userHandlers:       userHandlers,
		artistHandlers:     artistHandlers,
		trackHandlers:      trackHandlers,
		albumHandlers:      albumHandlers,
		playlistHandlers:   playlistHandlers,
		middlewareHandlers: middlewareHandlers,
		queueHandlers: queueHandlers,
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
	db.SetConnMaxLifetime(time.Second * 300)

	return db
}

func InitializeRedisUsers() *redis.Client {
	var AddrConfig string
	if len(os.Getenv("REDIS_PORT")) == 0 {
		AddrConfig = os.Getenv("REDIS_HOST")
	} else {
		AddrConfig = fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	}

	redisConnUsers := redis.NewClient(&redis.Options{
		Addr:     AddrConfig,
		Password: os.Getenv("REDIS_PASS"),
		DB:       1,
	})

	return redisConnUsers
}

func InitializeRedisQueues() *redis.Client {
	var AddrConfig string
	if len(os.Getenv("REDIS_PORT")) == 0 {
		AddrConfig = os.Getenv("REDIS_HOST")
	} else {
		AddrConfig = fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	}

	redisConnUsers := redis.NewClient(&redis.Options{
		Addr:     AddrConfig,
		Password: os.Getenv("REDIS_PASS"),
		DB:       2,
	})

	return redisConnUsers
}

func LoadMicroservices(server *echo.Echo) (authorizationMicro.SessionCheckerClient, []*grpc.ClientConn) {
	connections := make([]*grpc.ClientConn, 0)

	authPORT := os.Getenv("AUTH_PORT")
	log.Println("AUTH_PORT", authPORT)

	authConn, err := grpc.Dial(
		"127.0.0.1" + authPORT,
		grpc.WithInsecure(),
	)
	connections = append(connections, authConn)

	if err != nil {
		server.Logger.Fatal("cant connect to grpc")
	}

	authorizationManager := authorizationMicro.NewSessionCheckerClient(authConn)

	return authorizationManager, connections
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
	redisConnUsers := InitializeRedisUsers()
	defer func() {
		if redisConnUsers != nil {
			redisConnUsers.Close()
		}
	}()
	redisConnQueues := InitializeRedisQueues()
	defer func() {
		if redisConnQueues != nil {
			redisConnQueues.Close()
		}
	}()

	auth, conn := LoadMicroservices(server)
	defer func() {
		if len(conn) > 0 {
			for i, _ := range conn {
				conn[i].Close()
			}
		}
	}()

	api := NewRequestHandler(db, redisConnQueues, logger, auth)

	api.userHandlers.InitHandlers(server)
	api.artistHandlers.InitHandlers(server)
	api.trackHandlers.InitHandlers(server)
	api.albumHandlers.InitHandlers(server)
	api.playlistHandlers.InitHandlers(server)
	api.queueHandlers.InitHandlers(server)
	api.middlewareHandlers.InitMiddlewareHandlers(server)

	server.Static("/tracks", os.Getenv("TRACKS_PATH"))

	server.Logger.Fatal(server.Start(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))))
}
