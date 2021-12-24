package main

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	redisKey = "top10"
	userID = 261
	title = "LostPointer top 10"
	artwork = "top_10_100px.webp"
	artworkColor = "#e60f5a"
	isPublic = true
)

func ConnectRedis() *redis.Client {
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

func ConnectDatabase() *sql.DB {
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

func getTopTracks(db *sql.DB) ([]int, error) {
	query := `SELECT id FROM tracks ORDER BY listen_count DESC LIMIT 10`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	var id int
	tracksID := make([]int, 0)
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return tracksID, err
		}

		tracksID = append(tracksID, id)
	}

	return tracksID, nil
}

func UpdatePlaylist() {
	log.Println("Cron called function")

	dbConn := ConnectDatabase()
	// redisConn := ConnectRedis()

	tracksID, err := getTopTracks(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(tracksID)
}

func main() {
	c := cron.New()

	err := c.AddFunc("1 * * * * *", UpdatePlaylist)
	if err != nil {
		log.Fatal(err)
	}

	go c.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}