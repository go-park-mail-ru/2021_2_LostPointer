package images

import (
	"2021_2_LostPointer/internal/constants"
	"bytes"
	"github.com/chai2010/webp"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func TestAvatarRepository_CreateImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))
	img.Set(2, 3, color.RGBA{R: 255, A: 255})
	f, _ := os.OpenFile("image.jpeg", os.O_WRONLY|os.O_CREATE, 0600)
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

	const filename = "image.jpeg"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Создали поле
	fw, _ := writer.CreateFormFile("photo", filename)

	file, _ := os.Open("./" + filename)
	io.Copy(fw, file)
	writer.Close()

	req, _ := http.NewRequest("POST", "", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = req.ParseMultipartForm(0)

	r := NewAvatarRepository()

	avatar, err := r.CreateImage(req.MultipartForm.File["photo"][0])
	_ = os.Remove("image.jpeg")
	_ = os.Remove(avatar + constants.LittleAvatarPostfix)
	_ = os.Remove(avatar + constants.BigAvatarPostfix)
	assert.NoError(t, err)
}

func TestAvatarRepository_DeleteImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))
	img.Set(2, 3, color.RGBA{R: 255, A: 255})
	f, _ := os.OpenFile("out_150px.webp", os.O_WRONLY|os.O_CREATE, 0600)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Error(err)
		}
	}(f)
	err := webp.Encode(f, img, nil)
	if err != nil {
		t.Error(err)
	}

	f, _ = os.OpenFile("out_500px.webp", os.O_WRONLY|os.O_CREATE, 0600)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Error(err)
		}
	}(f)
	err = webp.Encode(f, img, nil)
	if err != nil {
		t.Error(err)
	}

	r := NewAvatarRepository()

	err = r.DeleteImage("out")

	assert.NoError(t, err)
}
