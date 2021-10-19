package delivery

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/music"
	"github.com/labstack/echo"
	"net/http"
)

type MusicHandlers struct {
	MusicUseCase music.MusicUseCase
}

func NewMusicDelivery(musicUseCase music.MusicUseCase) MusicHandlers {
	return MusicHandlers{MusicUseCase: musicUseCase}
}

func (musicHandlers MusicHandlers) Home(ctx echo.Context) error {
	collection, err := musicHandlers.MusicUseCase.GetMusicCollection(ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.Response{Message: "No music today"})
	}

	return ctx.JSON(http.StatusOK, collection)
}

func (musicHandlers MusicHandlers) InitHandlers(server *echo.Echo) {
	server.GET("api/v1/home", musicHandlers.Home)
}
