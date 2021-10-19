package delivery

import (
	"2021_2_LostPointer/pkg/models"
	"2021_2_LostPointer/pkg/music"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
)

const NoMusic = "No music today"

type MusicHandlers struct {
	MusicUseCase music.MusicUseCaseIFace
	logger       *zap.SugaredLogger
}

func NewMusicDelivery(musicUseCase music.MusicUseCaseIFace) MusicHandlers {
	return MusicHandlers{MusicUseCase: musicUseCase}
}

func (musicHandlers MusicHandlers) Home(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)

	collection, err := musicHandlers.MusicUseCase.GetMusicCollection(ctx)
	if err != nil {
		musicHandlers.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{Message: NoMusic})
	}

	musicHandlers.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, collection)
}

func (musicHandlers MusicHandlers) InitHandlers(server *echo.Echo) {
	server.GET("api/v1/home", musicHandlers.Home)
}
