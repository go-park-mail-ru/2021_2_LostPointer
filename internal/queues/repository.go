package queues

import "2021_2_LostPointer/internal/models"

type QueueRepository interface {
	StoreQueue(int, *models.Queue) error
	GetQueue(int) (*models.Queue, error)
}
