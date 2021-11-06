package delivery

import (
	"2021_2_LostPointer/internal/album"
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/models"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type AlbumDelivery struct {
	AlbumUseCase album.AlbumUseCase
	Logger       *zap.SugaredLogger
}

func NewAlbumDelivery(albumUseCase album.AlbumUseCase, logger *zap.SugaredLogger) AlbumDelivery {
	return AlbumDelivery{AlbumUseCase: albumUseCase, Logger: logger}
}

func (albumDelivery *AlbumDelivery) Home(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)

	albums, err := albumDelivery.AlbumUseCase.GetHome(constants.AlbumCollectionLimit)
	if err != nil {
		albumDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: constants.NoAlbums},
		)
	}

	albumDelivery.Logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, albums)
}

func (albumDelivery *AlbumDelivery) GetByArtist(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)

	artistID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		albumDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusBadRequest,
			Message: constants.InvalidParameter},
		)
	}
	albums, customErr := albumDelivery.AlbumUseCase.GetByArtist(artistID, constants.AlbumsDefaultAmountForArtist)
	if customErr != nil {
		albumDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", customErr.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: constants.NoAlbums},
		)
	}

	albumDelivery.Logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, albums)
}

func (albumDelivery *AlbumDelivery) InitHandlers(server *echo.Echo) {
	server.GET("api/v1/home/albums", albumDelivery.Home)
	server.GET("api/v1/artist/:id/albums", albumDelivery.GetByArtist)
}
