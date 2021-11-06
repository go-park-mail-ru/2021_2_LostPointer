package main

import (
	api "2021_2_LostPointer/internal/api/delivery"
	authMicroservice "2021_2_LostPointer/internal/microservices/authorization/proto"
	"fmt"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"log"
	"os"
)

func LoadMicroservices(server *echo.Echo) (authMicroservice.AuthorizationClient, []*grpc.ClientConn) {
	connections := make([]*grpc.ClientConn, 0)

	authPORT := os.Getenv("AUTH_PORT")
	log.Println("AUTH_PORT", authPORT)

	authConn, err := grpc.Dial(
		"127.0.0.1"+authPORT,
		grpc.WithInsecure(),
	)
	connections = append(connections, authConn)

	if err != nil {
		server.Logger.Fatal("cant connect to grpc")
	}

	authorizationManager := authMicroservice.NewAuthorizationClient(authConn)

	return authorizationManager, connections
}

func main() {
	server := echo.New()
	auth, conn := LoadMicroservices(server)
	defer func() {
		if len(conn) > 0 {
			for i, _ := range conn {
				conn[i].Close()
			}
		}
	}()

	app := api.NewApiMicroservices(auth)
	app.Init(server)

	server.Logger.Fatal(server.Start(fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))))
}