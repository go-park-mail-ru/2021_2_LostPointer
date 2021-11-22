package main

import (
	"2021_2_LostPointer/internal/monitoring/delivery"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	api "2021_2_LostPointer/internal/api/delivery"
	authMicroservice "2021_2_LostPointer/internal/microservices/authorization/proto"
	musicMicroservice "2021_2_LostPointer/internal/microservices/music/proto"
	playlistsMicroservice "2021_2_LostPointer/internal/microservices/playlists/proto"
	profileMicroservice "2021_2_LostPointer/internal/microservices/profile/proto"
	"2021_2_LostPointer/internal/middleware"
	"2021_2_LostPointer/pkg/image"
)

func LoadMicroservices(server *echo.Echo) (authMicroservice.AuthorizationClient, profileMicroservice.ProfileClient,
	musicMicroservice.MusicClient, playlistsMicroservice.PlaylistsClient, []*grpc.ClientConn) {
	connections := make([]*grpc.ClientConn, 0)

	authPORT := os.Getenv("AUTH_PORT")
	authConn, err := grpc.Dial(
		os.Getenv("AUTH_HOST")+authPORT,
		grpc.WithInsecure(),
	)
	if err != nil {
		server.Logger.Fatal("cant connect to grpc")
	}
	connections = append(connections, authConn)

	profilePORT := os.Getenv("PROFILE_PORT")
	profileConn, err := grpc.Dial(
		os.Getenv("PROFILE_HOST")+profilePORT,
		grpc.WithInsecure(),
	)
	if err != nil {
		server.Logger.Fatal("cant connect to grpc")
	}
	connections = append(connections, profileConn)

	musicPORT := os.Getenv("MUSIC_PORT")
	musicConn, err := grpc.Dial(
		os.Getenv("MUSIC_HOST")+musicPORT,
		grpc.WithInsecure(),
	)
	if err != nil {
		server.Logger.Fatal("cant connect to grpc")
	}
	connections = append(connections, musicConn)

	playlistsPORT := os.Getenv("PLAYLISTS_PORT")
	playlistsConn, err := grpc.Dial(
		os.Getenv("PLAYLISTS_HOST")+playlistsPORT,
		grpc.WithInsecure(),
	)
	if err != nil {
		server.Logger.Fatal("cant connect to grpc")
	}
	connections = append(connections, playlistsConn)

	authorizationManager := authMicroservice.NewAuthorizationClient(authConn)
	profileManager := profileMicroservice.NewProfileClient(profileConn)
	musicManager := musicMicroservice.NewMusicClient(musicConn)
	playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)

	return authorizationManager, profileManager, musicManager, playlistsManager, connections
}

func main() {
	server := echo.New()
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		err := prLogger.Sync()
		if err != nil {
			log.Fatal("Error occurred in logger sync")
		}
	}(prLogger)

	auth, profile, music, playlists, conn := LoadMicroservices(server)
	defer func() {
		if len(conn) > 0 {
			for _, c := range conn {
				err := c.Close()
				if err != nil {
					log.Fatalf("Error occurred during closing connection")
				}
			}
		}
	}()
	imageServices := image.NewImagesService()
	appHandler := api.NewAPIMicroservices(logger, imageServices, auth, profile, music, playlists)

	monitor := delivery.RegisterMonitoring(server)
	middlewareHandler := middleware.NewMiddlewareHandler(auth, logger, monitor)

	appHandler.Init(server)
	middlewareHandler.InitMiddlewareHandlers(server)

	server.Static("/tracks", os.Getenv("TRACKS_PATH"))

	server.Logger.Fatal(server.Start(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))))
}
