package delivery

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/playlist"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
)

const NoPlaylists = "No playlists"
const SelectionLimit = 4

type PlaylistDelivery struct {
	PlaylistUseCase playlist.PlaylistUseCase
	Logger          *zap.SugaredLogger
}

func NewPlaylistDelivery(playlistUseCae playlist.PlaylistUseCase, logger *zap.SugaredLogger) PlaylistDelivery {
	return PlaylistDelivery{PlaylistUseCase: playlistUseCae, Logger: logger}
}

func (playlistDelivery PlaylistDelivery) Home(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)

	artists, err := playlistDelivery.PlaylistUseCase.GetHome(SelectionLimit)
	if err != nil {
		playlistDelivery.Logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.OriginalError.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: NoPlaylists},
		)
	}

	playlistDelivery.Logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, artists)
}

func (playlistDelivery PlaylistDelivery) InitHandlers(server *echo.Echo) {
	server.GET("api/v1/home/playlists", playlistDelivery.Home)
}
