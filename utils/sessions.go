package utils

import (
	"2021_2_LostPointer/models"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"time"
)

func StoreSession(redisConnection *redis.Client, session *models.Session) error {
	err := redisConnection.Set(session.Session, session.UserID, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetSessionUser(redisConnection *redis.Client, session string) (int, error) {
	log.Println(session)
	res, err := redisConnection.Get(session).Result()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	id, err := strconv.Atoi(res)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return id, err
}
