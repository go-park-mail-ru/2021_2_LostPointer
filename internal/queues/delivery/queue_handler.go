package delivery

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/queues"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
)

type QueueDelivery struct {
	queueLogic queues.QueueUsecase
	logger 	   *zap.SugaredLogger
}

func NewQueueDelivery(queueLogic queues.QueueUsecase, logger *zap.SugaredLogger) QueueDelivery {
	return QueueDelivery{
		queueLogic: queueLogic,
		logger: logger,
	}
}

func (queueD QueueDelivery) StoreQueue(ctx echo.Context) error {
	userID := ctx.Get("USER_ID").(int)
	requestID := ctx.Get("REQUEST_ID").(string)

	var queueData models.Queue
	err := ctx.Bind(&queueData)
	if err != nil {
		queueD.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)

		return ctx.JSON(http.StatusInternalServerError, &models.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	customError := queueD.queueLogic.StoreQueue(userID, &queueData)
	if customError != nil {
		if customError.ErrorType == http.StatusInternalServerError {
			queueD.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.OriginalError.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)

			return ctx.JSON(http.StatusInternalServerError, &models.Response{
				Status:  http.StatusInternalServerError,
				Message: customError.OriginalError.Error(),
			})
		}
	}

	return ctx.JSON(http.StatusOK, &models.Response{
		Status: http.StatusOK,
		Message: "Queue stored successfully",
	})
}

func (queueD QueueDelivery) GetQueue(ctx echo.Context) error {
	userID := ctx.Get("USER_ID").(int)
	requestID := ctx.Get("REQUEST_ID").(string)

	queue, customError := queueD.queueLogic.GetQueue(userID)
	if customError != nil {
		if customError.ErrorType == http.StatusInternalServerError {
			queueD.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.OriginalError.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)

			return ctx.JSON(http.StatusInternalServerError, &models.Response{
				Status:  http.StatusInternalServerError,
				Message: customError.OriginalError.Error(),
			})
		}
	}

	return ctx.JSON(http.StatusOK, queue)
}

func (queueD QueueDelivery) InitHandlers(server *echo.Echo) {
	server.POST("/api/v1/queue", queueD.StoreQueue)
	server.GET("/api/v1/queue", queueD.GetQueue)
}
