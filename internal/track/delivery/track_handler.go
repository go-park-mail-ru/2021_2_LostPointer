package delivery

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/track"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type TrackDelivery struct {
	TrackUseCase track.TrackUseCase
	Logger       *zap.SugaredLogger
}

func NewTrackDelivery(trackUseCase track.TrackUseCase, logger *zap.SugaredLogger) TrackDelivery {
	return TrackDelivery{TrackUseCase: trackUseCase, Logger: logger}
}

func (trackDelivery *TrackDelivery) Home(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)
	userID := ctx.Get("USER_ID").(int)
	var isAuthorized bool
	if userID != -1 {
		isAuthorized = true
	}

	tracks, err := trackDelivery.TrackUseCase.GetHome(constants.TracksCollectionLimit, isAuthorized)
	if err != nil {
		trackDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: constants.NoTracks},
		)
	}

	trackDelivery.Logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, tracks)
}

func (trackDelivery *TrackDelivery) GetByArtist(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)
	userID := ctx.Get("USER_ID").(int)
	var isAuthorized bool
	if userID != -1 {
		isAuthorized = true
	}

	artistID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		trackDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusBadRequest,
			Message: constants.InvalidParameter},
		)
	}
	tracks, customErr := trackDelivery.TrackUseCase.GetByArtist(artistID, constants.TracksDefaultAmountForArtist, isAuthorized)
	if customErr != nil {
		trackDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", customErr.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: constants.NoTracks},
		)
	}

	trackDelivery.Logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, tracks)
}

func (trackDelivery *TrackDelivery) GetByAlbum(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)
	userID := ctx.Get("USER_ID").(int)
	var isAuthorized bool
	if userID != -1 {
		isAuthorized = true
	}

	artistID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		trackDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusBadRequest,
			Message: constants.InvalidParameter},
		)
	}
	tracks, customErr := trackDelivery.TrackUseCase.GetByAlbum(artistID, isAuthorized)
	if customErr != nil {
		trackDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", customErr.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: constants.NoTracks},
		)
	}

	trackDelivery.Logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, tracks)
}

func (trackDelivery *TrackDelivery) IncrementListenCount(ctx echo.Context) error {
	var trackID models.TrackID
	requestID := ctx.Get("REQUEST_ID").(string)

	err := ctx.Bind(&trackID)
	if err != nil {
		trackDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	customError := trackDelivery.TrackUseCase.IncrementListenCount(trackID.Id)
	if customError != nil {
		trackDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", customError.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: "Incremented track listen count",
	})
}

func (trackDelivery TrackDelivery) InitHandlers(server *echo.Echo) {
	server.GET("api/v1/home/tracks", trackDelivery.Home)
	server.POST("/api/v1/inc_listencount", trackDelivery.IncrementListenCount)
	server.GET("/api/v1/artist/:id/tracks", trackDelivery.GetByArtist)
	server.GET("/api/v1/album/:id/tracks", trackDelivery.GetByAlbum)
}
