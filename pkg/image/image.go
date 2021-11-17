package image

import (
	"bufio"
	"io"
	"mime/multipart"
	"os"

	"github.com/CapsLock-Studio/go-webpbin"
	"github.com/oliamb/cutter"
	uuid "github.com/satori/go.uuid"
	"github.com/sunshineplan/imgconv"

	"2021_2_LostPointer/internal/constants"
)

type Service interface {
	CreateAvatar(*multipart.FileHeader) (string, error)
	DeleteAvatar(string) error
}

type ImageService struct{}

func NewImageService() ImageService {
	return ImageService{}
}

//nolint:cyclop
func (service *ImageService) CreateAvatar(file *multipart.FileHeader) (string, error) {
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
		src = imgconv.Resize(src, imgconv.ResizeOption{Height: constants.BigAvatarHeight})
	} else {
		src = imgconv.Resize(src, imgconv.ResizeOption{Width: constants.BigAvatarHeight})
	}
	// Cropping image
	avatarLarge, err := cutter.Crop(src, cutter.Config{Width: constants.BigAvatarHeight, Height: constants.BigAvatarHeight, Mode: cutter.Centered})
	if err != nil {
		return "", err
	}
	// Create image and encode it into WEBP
	out, err := os.Create(os.Getenv("USERS_FULL_PREFIX") + fileName + constants.BigAvatarPostfix)
	if err != nil {
		return "", err
	}
	writer := io.Writer(out)
	if err = webpbin.Encode(writer, avatarLarge); err != nil {
		return "", err
	}

	if height < width {
		src = imgconv.Resize(src, imgconv.ResizeOption{Height: constants.LittleAvatarHeight})
	} else {
		src = imgconv.Resize(src, imgconv.ResizeOption{Width: constants.LittleAvatarHeight})
	}
	avatarSmall, err := cutter.Crop(src, cutter.Config{Width: constants.LittleAvatarHeight, Height: constants.LittleAvatarHeight, Mode: cutter.Centered})
	if err != nil {
		return "", err
	}
	out, err = os.Create(os.Getenv("USERS_FULL_PREFIX") + fileName + constants.LittleAvatarPostfix)
	if err != nil {
		return "", err
	}
	writer = io.Writer(out)
	if err = webpbin.Encode(writer, avatarSmall); err != nil {
		return "", err
	}

	return fileName, nil
}

func (service *ImageService) CreatePlaylistArtwork(file *multipart.FileHeader) (string, error) {
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
		src = imgconv.Resize(src, imgconv.ResizeOption{Height: constants.BigPlaylistArtworkHeight})
	} else {
		src = imgconv.Resize(src, imgconv.ResizeOption{Width: constants.BigPlaylistArtworkHeight})
	}
	// Cropping image
	avatarLarge, err := cutter.Crop(src, cutter.Config{Width: constants.BigPlaylistArtworkHeight, Height: constants.BigPlaylistArtworkHeight, Mode: cutter.Centered})
	if err != nil {
		return "", err
	}
	// Create image and encode it into WEBP
	out, err := os.Create(os.Getenv("PLAYLIST_FULL_PREFIX") + fileName + constants.BigPlaylistArtworkPostfix)
	if err != nil {
		return "", err
	}
	writer := io.Writer(out)
	if err = webpbin.Encode(writer, avatarLarge); err != nil {
		return "", err
	}

	if height < width {
		src = imgconv.Resize(src, imgconv.ResizeOption{Height: constants.LittlePlaylistArtworkHeight})
	} else {
		src = imgconv.Resize(src, imgconv.ResizeOption{Width: constants.LittlePlaylistArtworkHeight})
	}
	avatarSmall, err := cutter.Crop(src, cutter.Config{Width: constants.LittlePlaylistArtworkHeight, Height: constants.LittlePlaylistArtworkHeight, Mode: cutter.Centered})
	if err != nil {
		return "", err
	}
	out, err = os.Create(os.Getenv("PLAYLIST_FULL_PREFIX") + fileName + constants.LittlePlaylistArtworkPostfix)
	if err != nil {
		return "", err
	}
	writer = io.Writer(out)
	if err = webpbin.Encode(writer, avatarSmall); err != nil {
		return "", err
	}

	return fileName, nil
}

func (service *ImageService) DeletePlaylistArtwork(filename string) error {
	doesFileExist := true
	if _, err := os.Stat(os.Getenv("PLAYLIST_FULL_PREFIX") + filename + constants.LittlePlaylistArtworkPostfix); os.IsNotExist(err) {
		doesFileExist = false
	}

	if filename != constants.PlaylistArtworkDefaultFilename && doesFileExist {
		err := os.Remove(os.Getenv("PLAYLIST_FULL_PREFIX") + filename + constants.LittlePlaylistArtworkPostfix)
		if err != nil {
			return err
		}
		err = os.Remove(os.Getenv("PLAYLIST_FULL_PREFIX") + filename + constants.BigPlaylistArtworkPostfix)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *ImageService) DeleteAvatar(filename string) error {
	doesFileExist := true
	if _, err := os.Stat(os.Getenv("USERS_FULL_PREFIX") + filename + constants.LittleAvatarPostfix); os.IsNotExist(err) {
		doesFileExist = false
	}

	if filename != constants.AvatarDefaultFileName && doesFileExist {
		err := os.Remove(os.Getenv("USERS_FULL_PREFIX") + filename + constants.LittleAvatarPostfix)
		if err != nil {
			return err
		}
		err = os.Remove(os.Getenv("USERS_FULL_PREFIX") + filename + constants.BigAvatarPostfix)
		if err != nil {
			return err
		}
	}

	return nil
}
