package images

import (
	"io/ioutil"
	"log"
	"mime/multipart"
	"testing"
)




func TestAvatarRepository_CreateImage(t *testing.T) {
	//file := multipart.FileHeader{Filename: "/Users/vitaly/Desktop/Dev/go_all/2021_2_LostPointer/internal/utils/images/image.jpeg"}

	r := NewAvatarRepository()

	_, err := r.CreateImage(&file)
	if err != nil {
		log.Fatalln(err)
	}
}