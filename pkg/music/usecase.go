package music

import (
	"2021_2_LostPointer/pkg/models"
	"github.com/labstack/echo"
)

//go:generate moq -out ../mock/music_usecase_mock.go -pkg mock . MusicUseCaseIFace:MockMusicUseCaseIFace
type MusicUseCaseIFace interface {
	GetMusicCollection(ctx echo.Context) (*models.MusicCollection, error)
}
