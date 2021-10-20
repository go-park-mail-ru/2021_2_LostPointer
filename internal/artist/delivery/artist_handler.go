package delivery

import (
	"2021_2_LostPointer/internal/artist/usecase"
	"2021_2_LostPointer/internal/models"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const InvalidParameter = "Invalid parameter"
const DatabaseNotResponding = "Database not responding"

type ArtistDelivery struct {
	ArtistUseCase usecase.ArtistUseCase
	Logger        *zap.SugaredLogger
}

func NewArtistDelivery(artistUseCase usecase.ArtistUseCase, logger *zap.SugaredLogger) ArtistDelivery {
	return ArtistDelivery{ArtistUseCase: artistUseCase, Logger: logger}
}

func (artistDelivery ArtistDelivery) GetProfile(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)
	isAuthorized := ctx.Get("IS_AUTHORIZED").(bool)
	artistID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		artistDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  400,
			Message: InvalidParameter},
		)
	}

	artist, customErr := artistDelivery.ArtistUseCase.GetProfile(artistID, isAuthorized)
	if customErr != nil {
		artistDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", customErr.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: DatabaseNotResponding},
		)
	}

	artistDelivery.Logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, artist)
}

func (artistDelivery ArtistDelivery) InitHandlers(server *echo.Echo) {
	server.GET("api/v1/artist/:id", artistDelivery.GetProfile)
}
