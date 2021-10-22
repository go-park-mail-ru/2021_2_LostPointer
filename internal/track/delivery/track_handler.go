package delivery

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/track"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
)

const NoArtists = "No tracks"
const SelectionLimit = 10

type TrackDelivery struct {
	TrackUseCase track.TrackUseCase
	Logger       *zap.SugaredLogger
}

func NewTrackDelivery(trackUseCase track.TrackUseCase, logger *zap.SugaredLogger) TrackDelivery {
	return TrackDelivery{TrackUseCase: trackUseCase, Logger: logger}
}

func (trackDelivery TrackDelivery) Home(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)
	isAuthorized := ctx.Get("IS_AUTHORIZED").(bool)

	artists, err := trackDelivery.TrackUseCase.GetHome(SelectionLimit, isAuthorized)
	if err != nil {
		trackDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: NoArtists},
		)
	}

	trackDelivery.Logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, artists)
}

func (trackDelivery TrackDelivery) InitHandlers(server *echo.Echo) {
	server.GET("api/v1/home/tracks", trackDelivery.Home)
}
