package image

import (
	"2021_2_LostPointer/internal/constants"
	"bytes"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func TestCreateImages(t *testing.T) {
	const filename = "image.jpeg"

	img := image.NewRGBA(image.Rect(0, 0, 100, 50))
	img.Set(2, 3, color.RGBA{R: 255, A: 255})
	f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Error(err)
		}
	}(f)
	err := jpeg.Encode(f, img, nil)
	if err != nil {
		t.Error(err)
	}

	// Пришла непустая форма, все ОК
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, _ := writer.CreateFormFile("photo", filename)

	file, _ := os.Open("./" + filename)
	io.Copy(fw, file)
	writer.Close()

	req, _ := http.NewRequest("POST", "", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.ParseMultipartForm(0)

	r := NewImagesService()

	avatar, err := r.CreateImages(
		req.MultipartForm.File["photo"][0],
		os.Getenv("PLAYLIST_FULL_PREFIX"), map[int]string{
			100: constants.PlaylistArtworkExtension100px,
			384: constants.PlaylistArtworkExtension384px,
		})
	_ = os.Remove(filename)
	_ = os.Remove(avatar.Filename + constants.UserAvatarExtension150px)
	_ = os.Remove(avatar.Filename + constants.UserAvatarExtension500px)
	assert.NoError(t, err)

	// Пришла форма с пустой картинкой
	const brokenFilename = "image2.jpeg"
	f, _ = os.OpenFile(brokenFilename, os.O_WRONLY|os.O_CREATE, 0600)
	f.Write([]byte("a"))
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(f)
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)

	fw, _ = writer.CreateFormFile("photo", brokenFilename)

	file, _ = os.Open("./" + brokenFilename)
	io.Copy(fw, file)
	writer.Close()

	req, _ = http.NewRequest("POST", "", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.ParseMultipartForm(0)

	r = NewImagesService()

	avatar, err = r.CreateImages(
		req.MultipartForm.File["photo"][0],
		os.Getenv("PLAYLIST_FULL_PREFIX"), map[int]string{
			100: constants.PlaylistArtworkExtension100px,
			384: constants.PlaylistArtworkExtension384px,
		})
	_ = os.Remove(brokenFilename)
	assert.Error(t, err)
}
