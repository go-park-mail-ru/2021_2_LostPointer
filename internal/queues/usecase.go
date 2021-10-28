package queues

import "2021_2_LostPointer/internal/models"

type QueueUsecase interface {
	StoreQueue(int, *models.Queue) *models.CustomError
	GetQueue(int) (*models.Queue, *models.CustomError)
}
