//nolint:typecheck
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"2021_2_LostPointer/internal/microservices/authorization/proto"
	"2021_2_LostPointer/internal/microservices/authorization/repository"
	"2021_2_LostPointer/internal/microservices/authorization/usecase"
)

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
		DB:       1,
	})

	return redisConnection
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
	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalln("NO CONNECTION TO DATABASE", err.Error())
	}
	database.SetConnMaxLifetime(time.Second * 300)

	return database
}

func main() {
	redisConnection := InitializeRedis()
	dbConnection := InitializeDatabase()
	defer func() {
		if redisConnection != nil {
			err := redisConnection.Close()
			if err != nil {
				log.Fatal("Error occurred during closing redis connection")
			}
		}
	}()
	defer func() {
		if dbConnection != nil {
			err := dbConnection.Close()
			if err != nil {
				log.Fatal("Error occurred during closing database connection")
			}
		}
	}()
	storage := repository.NewAuthStorage(dbConnection, redisConnection)

	port := os.Getenv("AUTH_PORT")
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("CANNOT LISTEN PORT: ", port, err.Error())
	}

	server := grpc.NewServer()
	proto.RegisterAuthorizationServer(server, usecase.NewAuthService(storage))
	log.Printf("STARTED AUTHORIZATION MICROSERVICE ON %s", port)
	err = server.Serve(listen)
	if err != nil {
		log.Println("CANNOT LISTEN PORT: ", port, err.Error())
	}
}
