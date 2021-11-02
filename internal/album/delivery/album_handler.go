package delivery

import (
	"2021_2_LostPointer/internal/album"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/utils/constants"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type AlbumDelivery struct {
	AlbumUseCase album.AlbumUseCase
	Logger       *zap.SugaredLogger
}

func NewAlbumDelivery(albumUseCase album.AlbumUseCase, logger *zap.SugaredLogger) AlbumDelivery {
	return AlbumDelivery{AlbumUseCase: albumUseCase, Logger: logger}
}

func (albumDelivery AlbumDelivery) Home(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)

	artists, err := albumDelivery.AlbumUseCase.GetHome(constants.AlbumCollectionLimit)
	if err != nil {
		albumDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: constants.NoArtists},
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