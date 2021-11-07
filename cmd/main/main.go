package main

import (
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	api "2021_2_LostPointer/internal/api/delivery"
	authMicroservice "2021_2_LostPointer/internal/microservices/authorization/proto"
	profileMicroservice "2021_2_LostPointer/internal/microservices/profile/proto"
	"2021_2_LostPointer/internal/middleware"
	"2021_2_LostPointer/pkg/image"
)

func LoadMicroservices(server *echo.Echo) (authMicroservice.AuthorizationClient, profileMicroservice.ProfileClient,
	[]*grpc.ClientConn) {
	connections := make([]*grpc.ClientConn, 0)

	authPORT := os.Getenv("AUTH_PORT")
	authConn, err := grpc.Dial(
		"127.0.0.1"+authPORT,
		grpc.WithInsecure(),
	)
	if err != nil {
		server.Logger.Fatal("cant connect to grpc")
	}
	connections = append(connections, authConn)

	profilePORT := os.Getenv("PROFILE_PORT")
	profileConn, err := grpc.Dial(
		"127.0.0.1"+profilePORT,
		grpc.WithInsecure(),
	)
	if err != nil {
		server.Logger.Fatal("cant connect to grpc")
	}
	connections = append(connections, profileConn)

	authorizationManager := authMicroservice.NewAuthorizationClient(authConn)
	profileManager := profileMicroservice.NewProfileClient(profileConn)

	return authorizationManager, profileManager, connections
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

	auth, profile, conn := LoadMicroservices(server)
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

	avatarsServices := image.NewAvatarsService()

	appHandler := api.NewAPIMicroservices(logger, avatarsServices, auth, profile)
	middlewareHandler := middleware.NewMiddlewareHandler(auth, logger)

	appHandler.Init(server)
	middlewareHandler.InitMiddlewareHandlers(server)

	server.Logger.Fatal(server.Start(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))))
}
