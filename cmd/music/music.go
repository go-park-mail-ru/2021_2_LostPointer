package main

import (
	"database/sql"
	"fmt"
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

func main() {
	dbConnection := InitializeDatabase()
	storage := repository.NewMusicStorage(dbConnection)
	defer func() {
		if dbConnection != nil {
			err := dbConnection.Close()
			if err != nil {
				log.Fatal("Error occurred during closing database connection")
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
