package utils

import (
	"2021_2_LostPointer/models"
	"github.com/go-redis/redis"
	"time"
)


func StoreSession(redisConnection *redis.Client, session *models.Session) error {
	err := redisConnection.Set(session.Session, session.UserID, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

