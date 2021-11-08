package main

import (
	"2021_2_LostPointer/internal/microservices/music/proto"
	"2021_2_LostPointer/internal/microservices/music/repository"
	"2021_2_LostPointer/internal/microservices/music/usecase"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
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
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalln("NO CONNECTION TO DATABASE", err.Error())
	}
	db.SetConnMaxLifetime(time.Second * 300)

	return db
}

func InitializeStorages() *repository.MusicStorage {
	dbConnection := InitializeDatabase()

	musicStorage := repository.NewMusicStorage(dbConnection)
	return musicStorage
}

func main() {
	storage := InitializeStorages()
	port := os.Getenv("MUSIC_PORT")
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("CANNOT LISTEN PORT: ", port, err.Error())
	}

	server := grpc.NewServer()
	proto.RegisterMusicServer(server, usecase.NewMusicService(*storage))
	log.Printf("STARTED MUSIC MICROSERVICE ON %s", port)
	err = server.Serve(listen)
	if err != nil {
		log.Fatal("CANNOT LISTEN PORT: ", port, err.Error())
	}
}
