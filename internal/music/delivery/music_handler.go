package delivery

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/music"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
)

const NoMusic = "No music today"

type MusicHandlers struct {
	MusicUseCase music.MusicUseCase
	Logger       *zap.SugaredLogger
}

func NewMusicDelivery(musicUseCase music.MusicUseCase, logger *zap.SugaredLogger) MusicHandlers {
	return MusicHandlers{MusicUseCase: musicUseCase, Logger: logger}
}

func (musicHandlers MusicHandlers) Home(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)
	isAuthorized := ctx.Get("IS_AUTHORIZED").(bool)

	collection, err := musicHandlers.MusicUseCase.GetMusicCollection(isAuthorized)
	if err != nil {
		musicHandlers.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: NoMusic},
			)
	}

	musicHandlers.Logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, collection)
}

func (musicHandlers MusicHandlers) InitHandlers(server *echo.Echo) {
	server.GET("api/v1/home", musicHandlers.Home)
}
