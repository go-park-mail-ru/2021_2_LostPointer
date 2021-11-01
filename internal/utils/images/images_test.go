package images

import (
	"github.com/chai2010/webp"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"os"
	"testing"
)

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