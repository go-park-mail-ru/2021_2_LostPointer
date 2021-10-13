package delivery

import (
	"2021_2_LostPointer/pkg/music"
	"github.com/labstack/echo"
	"net/http"
)

type MusicHandlers struct {
	MusicUseCase music.MusicUseCaseIFace
}

func NewMusicDelivery(musicUseCase music.MusicUseCaseIFace) MusicHandlers {
	return MusicHandlers{MusicUseCase: musicUseCase}
}

func (musicHandlers MusicHandlers) Home(ctx echo.Context) error {
	collection, err := musicHandlers.MusicUseCase.GetMusicCollection()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, collection)
}

func (musicHandlers MusicHandlers) InitHandlers(server *echo.Echo) {
	server.GET("api/v1/home", musicHandlers.Home)
}
