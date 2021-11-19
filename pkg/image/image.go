package image

import (
	"bufio"
	"io"
	"mime/multipart"
	"os"

	"github.com/CapsLock-Studio/go-webpbin"
	"github.com/cenkalti/dominantcolor"
	"github.com/oliamb/cutter"
	uuid "github.com/satori/go.uuid"
	"github.com/sunshineplan/imgconv"

	"2021_2_LostPointer/internal/constants"
)

type Service interface {
	CreateAvatar(*multipart.FileHeader) (string, error)
	DeleteAvatar(string) error
}

type ImagesService struct{}

func NewImagesService() ImagesService {
	return ImagesService{}
}

//nolint:cyclop
func (service *ImagesService) CreateAvatar(file *multipart.FileHeader) (string, error) {
	// Open image and decode it into image.Image type
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(f multipart.File) {
		_ = f.Close()
	}(f)
	reader := bufio.NewReader(f)
	src, err := imgconv.Decode(reader)
	if err != nil {
		return "", err
	}
	// Get image width and height
	_, err = f.Seek(0, 0)
	if err != nil {
		return "", err
	}
	reader = bufio.NewReader(f)
	tmp, err := imgconv.DecodeConfig(reader)
	if err != nil {
		return "", err
	}
	width := tmp.Width
	height := tmp.Height
	// Generate filename for image
	fileName := uuid.NewV4().String()
	// Resizing image
	if height < width {
		src = imgconv.Resize(src, imgconv.ResizeOption{Height: constants.UserAvatarHeight500px})
	} else {
		src = imgconv.Resize(src, imgconv.ResizeOption{Width: constants.UserAvatarHeight500px})
	}
	// Cropping image
	avatarLarge, err := cutter.Crop(src, cutter.Config{Width: constants.UserAvatarHeight500px, Height: constants.UserAvatarHeight500px, Mode: cutter.Centered})
	if err != nil {
		return "", err
	}
	// Create image and encode it into WEBP
	out, err := os.Create(os.Getenv("USERS_FULL_PREFIX") + fileName + constants.UserAvatarExtension500px)
	if err != nil {
		return "", err
	}
	writer := io.Writer(out)
	if err = webpbin.Encode(writer, avatarLarge); err != nil {
		return "", err
	}

	if height < width {
		src = imgconv.Resize(src, imgconv.ResizeOption{Height: constants.UserAvatarHeight150px})
	} else {
		src = imgconv.Resize(src, imgconv.ResizeOption{Width: constants.UserAvatarHeight150px})
	}
	avatarSmall, err := cutter.Crop(src, cutter.Config{Width: constants.UserAvatarHeight150px, Height: constants.UserAvatarHeight150px, Mode: cutter.Centered})
	if err != nil {
		return "", err
	}
	out, err = os.Create(os.Getenv("USERS_FULL_PREFIX") + fileName + constants.UserAvatarExtension150px)
	if err != nil {
		return "", err
	}
	writer = io.Writer(out)
	if err = webpbin.Encode(writer, avatarSmall); err != nil {
		return "", err
	}

	return fileName, nil
}

//nolint:cyclop
func (service *ImagesService) CreatePlaylistArtwork(file *multipart.FileHeader) (string, string, error) {
	// Open image and decode it into image.Image type
	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer func(f multipart.File) {
		_ = f.Close()
	}(f)
	reader := bufio.NewReader(f)
	src, err := imgconv.Decode(reader)
	if err != nil {
		return "", "", err
	}
	// Get image width and height
	_, err = f.Seek(0, 0)
	if err != nil {
		return "", "", err
	}
	reader = bufio.NewReader(f)
	tmp, err := imgconv.DecodeConfig(reader)
	if err != nil {
		return "", "", err
	}
	width := tmp.Width
	height := tmp.Height
	// Generate filename for image
	fileName := uuid.NewV4().String()
	// Resizing image
	if height < width {
		src = imgconv.Resize(src, imgconv.ResizeOption{Height: constants.PlaylistArtworkHeight384px})
	} else {
		src = imgconv.Resize(src, imgconv.ResizeOption{Width: constants.PlaylistArtworkHeight384px})
	}
	// Cropping image
	avatarLarge, err := cutter.Crop(src, cutter.Config{Width: constants.PlaylistArtworkHeight384px, Height: constants.PlaylistArtworkHeight384px, Mode: cutter.Centered})
	if err != nil {
		return "", "", err
	}
	// Create image and encode it into WEBP
	out, err := os.Create(os.Getenv("PLAYLIST_FULL_PREFIX") + fileName + constants.PlaylistArtworkExtension384px)
	if err != nil {
		return "", "", err
	}
	writer := io.Writer(out)
	if err = webpbin.Encode(writer, avatarLarge); err != nil {
		return "", "", err
	}

	if height < width {
		src = imgconv.Resize(src, imgconv.ResizeOption{Height: constants.PlaylistArtworkHeight100px})
	} else {
		src = imgconv.Resize(src, imgconv.ResizeOption{Width: constants.PlaylistArtworkHeight100px})
	}
	avatarSmall, err := cutter.Crop(src, cutter.Config{Width: constants.PlaylistArtworkHeight100px, Height: constants.PlaylistArtworkHeight100px, Mode: cutter.Centered})
	if err != nil {
		return "", "", err
	}
	out, err = os.Create(os.Getenv("PLAYLIST_FULL_PREFIX") + fileName + constants.PlaylistArtworkExtension100px)
	if err != nil {
		return "", "", err
	}
	writer = io.Writer(out)
	if err = webpbin.Encode(writer, avatarSmall); err != nil {
		return "", "", err
	}

	artworkColor := dominantcolor.Hex(dominantcolor.Find(avatarSmall))

	return fileName, artworkColor, nil
}

//nolint:dupl
func (service *ImagesService) DeletePlaylistArtwork(filename string) error {
	doesFileExist := true
	if _, err := os.Stat(os.Getenv("PLAYLIST_FULL_PREFIX") + filename + constants.PlaylistArtworkExtension100px); os.IsNotExist(err) {
		doesFileExist = false
	}

	if filename != constants.PlaylistArtworkDefaultFilename && doesFileExist {
		err := os.Remove(os.Getenv("PLAYLIST_FULL_PREFIX") + filename + constants.PlaylistArtworkExtension100px)
		if err != nil {
			return err
		}
		err = os.Remove(os.Getenv("PLAYLIST_FULL_PREFIX") + filename + constants.PlaylistArtworkExtension384px)
		if err != nil {
			return err
		}
	}

	return nil
}

//nolint:dupl
func (service *ImagesService) DeleteAvatar(filename string) error {
	doesFileExist := true
	if _, err := os.Stat(os.Getenv("USERS_FULL_PREFIX") + filename + constants.UserAvatarExtension150px); os.IsNotExist(err) {
		doesFileExist = false
	}

	if filename != constants.AvatarDefaultFileName && doesFileExist {
		err := os.Remove(os.Getenv("USERS_FULL_PREFIX") + filename + constants.UserAvatarExtension150px)
		if err != nil {
			return err
		}
		err = os.Remove(os.Getenv("USERS_FULL_PREFIX") + filename + constants.UserAvatarExtension500px)
		if err != nil {
			return err
		}
	}

	return nil
}
