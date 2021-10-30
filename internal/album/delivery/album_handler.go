package delivery

import (
	"2021_2_LostPointer/internal/album"
	"2021_2_LostPointer/internal/models"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
)

const NoArtists = "No artists"
const SelectionLimit = 4

type AlbumDelivery struct {
	AlbumUseCase album.AlbumUseCase
	Logger       *zap.SugaredLogger
}

func NewAlbumDelivery(albumUseCase album.AlbumUseCase, logger *zap.SugaredLogger) AlbumDelivery {
	return AlbumDelivery{AlbumUseCase: albumUseCase, Logger: logger}
}

func (albumDelivery AlbumDelivery) Home(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)

	artists, err := albumDelivery.AlbumUseCase.GetHome(SelectionLimit)
	if err != nil {
		albumDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: NoArtists},
		)
	}

	albumDelivery.Logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, artists)
}

func (albumDelivery AlbumDelivery) InitHandlers(server *echo.Echo) {
	server.GET("api/v1/home/albums", albumDelivery.Home)
}