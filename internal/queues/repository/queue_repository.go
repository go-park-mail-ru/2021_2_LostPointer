package repository

import (
	"2021_2_LostPointer/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type QueueRepository struct {
	DB *sql.DB
	Redis *redis.Client
}

func NewQueueRepository(dbConn *sql.DB, redisConn *redis.Client) QueueRepository {
	return QueueRepository{
		DB: dbConn,
		Redis: redisConn,
	}
}

func (queueR QueueRepository) StoreQueue(userID int, queueData *models.Queue) error {
	stringData, err := json.Marshal(queueData)
	if err != nil {
		return err
	}
	err = queueR.Redis.Set(context.Background(), strconv.Itoa(userID), stringData, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (queueR QueueRepository) GetQueue(userID int) (*models.Queue, error) {
	var queueData models.Queue
	stringData, err := queueR.Redis.Get(context.Background(), strconv.Itoa(userID)).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(stringData), &queueData)
	if err != nil {
		return nil, err
	}
	return &queueData, nil
}
