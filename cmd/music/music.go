package main

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"2021_2_LostPointer/internal/microservices/music/proto"
	"2021_2_LostPointer/internal/microservices/music/repository"
	"2021_2_LostPointer/internal/microservices/music/usecase"
)

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
		DB:       2,
	})

	return redisConnection
}

func main() {
	dbConnection := InitializeDatabase()
	redisConnection := InitializeRedis()
	storage := repository.NewMusicStorage(dbConnection, redisConnection)
	defer func() {
		if dbConnection != nil {
			err := dbConnection.Close()
			if err != nil {
				log.Fatalf("Could not close database connection: %v", err)
			}
		}
		if redisConnection != nil {
			err := redisConnection.Close()
			if err != nil {
				log.Fatalf("Could not close redis connection: %v", err)
			}
		}
	}()

	port := os.Getenv("MUSIC_PORT")
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("CANNOT LISTEN PORT: %s error: %s", port, err.Error())
	}

	server := grpc.NewServer()
	proto.RegisterMusicServer(server, usecase.NewMusicService(storage))
	log.Printf("STARTED MUSIC MICROSERVICE ON %s", port)
	err = server.Serve(listen)
	if err != nil {
		log.Printf("CANNOT LISTEN PORT: %s error: %s", port, err.Error())
	}
}
