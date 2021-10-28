package repository
//
//import (
//	"context"
//	"github.com/go-redis/redis/v8"
//	"strconv"
//	"time"
//)
//
//type SessionRepository struct {
//	sessionDB *redis.Client
//}
//
//func NewSessionRepository(redisConnection *redis.Client) SessionRepository {
//	return SessionRepository{sessionDB: redisConnection}
//}
//
//func (sessionR SessionRepository) CreateSession(id int, cookieValue string) error {
//	err := sessionR.sessionDB.Set(context.Background(), cookieValue, id, time.Hour).Err()
//	return err
//}
//
//func (sessionR SessionRepository) GetUserIdByCookie(cookieValue string) (int, error) {
//	res, err := sessionR.sessionDB.Get(context.Background(), cookieValue).Result()
//	id, _ := strconv.Atoi(res)
//	return id, err
//}
//
//func (sessionR SessionRepository) DeleteSession(cookieValue string) error {
//	err := sessionR.sessionDB.Del(context.Background(), cookieValue).Err()
//	return err
//}
