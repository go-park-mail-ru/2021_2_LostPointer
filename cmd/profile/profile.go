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

	"2021_2_LostPointer/internal/microservices/profile/proto"
	"2021_2_LostPointer/internal/microservices/profile/repository"
	"2021_2_LostPointer/internal/microservices/profile/usecase"
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

func main() {
	dbConnection := InitializeDatabase()
	storage := repository.NewUserSettingsStorage(dbConnection)
	defer func() {
		if dbConnection != nil {
			err := dbConnection.Close()
			if err != nil {
				log.Fatal("Error occurred during closing database connection")
			}
		}
	}()

	port := os.Getenv("PROFILE_PORT")
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("CANNOT LISTEN PORT: ", port, err.Error())
	}

	server := grpc.NewServer()
	proto.RegisterProfileServer(server, usecase.NewProfileService(*storage))
	log.Printf("STARTED PROFILE MICROSERVICE ON %s", port)
	err = server.Serve(listen)
	if err != nil {
		log.Println("CANNOT LISTEN PORT: ", port, err.Error())
	}
}
