package usecase

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/queues"
	"net/http"
)

type QueueUseCase struct {
	queueDB queues.QueueRepository
}

func NewQueueUseCase(queueDB queues.QueueRepository) QueueUseCase {
	return QueueUseCase{queueDB: queueDB}
}

func (queueU QueueUseCase) StoreQueue(userID int, queueData *models.Queue) *models.CustomError {
	err := queueU.queueDB.StoreQueue(userID, queueData)
	if err != nil {
		return &models.CustomError{
			ErrorType: http.StatusInternalServerError,
			OriginalError: err,
		}
	}

	return nil
}

func (queueU QueueUseCase) GetQueue(userID int) (*models.Queue, *models.CustomError) {
	queueData, err := queueU.queueDB.GetQueue(userID)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType: http.StatusInternalServerError,
			OriginalError: err,
		}
	}

	return queueData, nil
}
