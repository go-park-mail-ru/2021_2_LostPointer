package delivery

import (
	"2021_2_LostPointer/internal/artist"
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/models"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ArtistDelivery struct {
	ArtistUseCase artist.ArtistUseCase
	Logger        *zap.SugaredLogger
}

func NewArtistDelivery(artistUseCase artist.ArtistUseCase, logger *zap.SugaredLogger) ArtistDelivery {
	return ArtistDelivery{ArtistUseCase: artistUseCase, Logger: logger}
}

func (artistDelivery ArtistDelivery) GetProfile(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)
	userID := ctx.Get("USER_ID").(int)
	var isAuthorized bool
	if userID != -1 {
		isAuthorized = true
	}
	artistID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		artistDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusBadRequest,
			Message: constants.InvalidParameter},
		)
	}

	art, customErr := artistDelivery.ArtistUseCase.GetProfile(artistID, isAuthorized)
	if customErr != nil {
		artistDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", customErr.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: constants.DatabaseNotResponding},
		)
	}

	artistDelivery.Logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, art)
}

func (artistDelivery ArtistDelivery) Home(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)

	artists, err := artistDelivery.ArtistUseCase.GetHome(constants.ArtistsCollectionLimit)
	if err != nil {
		artistDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: constants.NoArtists},
		)
	}

	artistDelivery.Logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, artists)
}

func (artistDelivery ArtistDelivery) InitHandlers(server *echo.Echo) {
	server.GET("api/v1/artist/:id", artistDelivery.GetProfile)
	server.GET("api/v1/home/artists", artistDelivery.Home)
}
