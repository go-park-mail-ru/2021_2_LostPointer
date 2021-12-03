package image

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"2021_2_LostPointer/internal/constants"
)

func TestCreateImages(t *testing.T) {
	const filename = "image.jpeg"

	img := image.NewRGBA(image.Rect(0, 0, 100, 50))
	img.Set(2, 3, color.RGBA{R: 255, A: 255})
	filePtr, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Error(err)
		}
	}(filePtr)
	err := jpeg.Encode(filePtr, img, nil)
	if err != nil {
		t.Error(err)
	}

	// Пришла непустая форма, все ОК
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fwriter, _ := writer.CreateFormFile("photo", filename)

	file, _ := os.Open("./" + filename)
	_, err = io.Copy(fwriter, file)
	if err != nil {
		t.Error(err)
	}
	writer.Close()

	req, _ := http.NewRequestWithContext(context.Background(), "POST", "", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = req.ParseMultipartForm(0)
	if err != nil {
		t.Error(err)
	}

	service := NewImagesService()

	avatar, err := service.CreateImages(
		req.MultipartForm.File["photo"][0],
		os.Getenv("PLAYLIST_FULL_PREFIX"), map[int]string{
			100: constants.PlaylistArtworkExtension100px,
			384: constants.PlaylistArtworkExtension384px,
		})
	_ = os.Remove(filename)
	_ = os.Remove(avatar.Filename + constants.PlaylistArtworkExtension100px)
	_ = os.Remove(avatar.Filename + constants.PlaylistArtworkExtension384px)
	assert.NoError(t, err)

	// Пришла форма с пустой картинкой
	const brokenFilename = "image2.jpeg"
	filePtr, _ = os.OpenFile(brokenFilename, os.O_WRONLY|os.O_CREATE, 0600)
	_, err = filePtr.Write([]byte("a"))
	if err != nil {
		t.Error(err)
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(filePtr)
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)

	fwriter, _ = writer.CreateFormFile("photo", brokenFilename)

	file, _ = os.Open("./" + brokenFilename)
	_, err = io.Copy(fwriter, file)
	if err != nil {
		t.Error(err)
	}

	writer.Close()

	req, _ = http.NewRequestWithContext(context.Background(), "POST", "", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = req.ParseMultipartForm(0)
	if err != nil {
		t.Error(err)
	}

	service = NewImagesService()

	_, err = service.CreateImages(
		req.MultipartForm.File["photo"][0],
		os.Getenv("PLAYLIST_FULL_PREFIX"), map[int]string{
			100: constants.PlaylistArtworkExtension100px,
			384: constants.PlaylistArtworkExtension384px,
		})
	_ = os.Remove(brokenFilename)
	assert.Error(t, err)
}
