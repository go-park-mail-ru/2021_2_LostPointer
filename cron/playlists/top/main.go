package main

import (
	"2021_2_LostPointer/internal/constants"
	"context"
	"database/sql"
	"fmt"
	"github.com/robfig/cron"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

const (
	redisKey     = "top10"
	userID       = 261
	title        = "LostPointer top 10"
	artwork      = "top_10"
	artworkColor = "#e60f5a"
	isPublic     = true
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
		DB:       3,
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
	err = rows.Err()
	if err != nil {
		return tracksID, err
	}

	return tracksID, nil
}

func initializePlaylist(db *sql.DB, tracksID []int) (int, error) {
	query := `INSERT INTO playlists(title, user_id, artwork, artwork_color, is_public) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id int
	err := db.QueryRow(query, title, userID, artwork, artworkColor, isPublic).Scan(&id)
	if err != nil {
		return id, err
	}


	query = `INSERT INTO playlist_tracks(playlist, track) VALUES ($1, $2)`
	for _, trackID := range tracksID {
		_, err := db.Exec(query, id, trackID)
		if err != nil {
			return id, err
		}
	}

	return id, nil
}

func updatePlaylist(db *sql.DB, playlistID string, tracksID []int) error {
	query := `SELECT id FROM playlist_tracks WHERE playlist = $1`

	rows, err := db.Query(query, playlistID)
	if err != nil {
		return err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	var currentID int
	id := make([]int, 0)
	for rows.Next() {
		if err := rows.Scan(&currentID); err != nil {
			return err
		}

		id = append(id, currentID)
	}

	query = `UPDATE playlist_tracks SET track=$1 WHERE id=$2`
	for i := 0; i < len(tracksID); i++ {
		_, err := db.Exec(query, tracksID[i], id[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func CreatePlaylist() {
	log.Println("Cron called function")

	dbConn := ConnectDatabase()
	redisConn := ConnectRedis()

	tracksID, err := getTopTracks(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	playlistID, err := redisConn.Get(context.Background(), redisKey).Result()
	if err == redis.Nil {
		createdPlaylistID, err := initializePlaylist(dbConn, tracksID)
		playlistID = strconv.Itoa(createdPlaylistID)
		if err != nil {
			log.Fatal(err)
		}
		err = redisConn.Set(context.Background(), redisKey, playlistID, constants.CookieLifetime).Err()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = updatePlaylist(dbConn, playlistID, tracksID)
		if err != nil {
			log.Fatal(err)
		}
	}
	
}

func main() {
	c := cron.New()

	err := c.AddFunc("1 * * * *", CreatePlaylist)
	if err != nil {
		log.Fatal(err)
	}

	go c.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
