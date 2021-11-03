package main

import (
	session "2021_2_LostPointer/internal/microservices/authorization/delivery"
	"2021_2_LostPointer/internal/microservices/authorization/usecase"
	sessionsRepository "2021_2_LostPointer/internal/sessions/repository"
	repositoryUser "2021_2_LostPointer/internal/users/repository"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

func InitializeDataBases() (repositoryUser.UserRepository, sessionsRepository.SessionRepository) {
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

	sessionDB := sessionsRepository.NewSessionRepository(redisConnUsers)
	userDB := repositoryUser.NewUserRepository(db)

	return userDB, sessionDB
}

func main() {
	users, sessions := InitializeDataBases()
	port := os.Getenv("AUTH_PORT")
	log.Println(port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("CANNOT LISTEN PORT : ", port, err.Error())
	}

	server := grpc.NewServer()

	session.RegisterSessionCheckerServer(server, usecase.NewAuthorizationUseCase(users, sessions))

	fmt.Println("starting server at " + port)
	err = server.Serve(lis)
	if err != nil {
		log.Fatal("CANNOT LISTEN PORT : ", port, err.Error())
	}
}
