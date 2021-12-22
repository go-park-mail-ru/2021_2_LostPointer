//nolint:varnamelen
package delivery

import (
	"2021_2_LostPointer/internal/csrf"
	"2021_2_LostPointer/pkg/image"
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"2021_2_LostPointer/internal/constants"
	authorization "2021_2_LostPointer/internal/microservices/authorization/proto"
	music "2021_2_LostPointer/internal/microservices/music/proto"
	playlists "2021_2_LostPointer/internal/microservices/playlists/proto"
	profile "2021_2_LostPointer/internal/microservices/profile/proto"
	"2021_2_LostPointer/internal/models"
)

type APIMicroservices struct {
	logger       *zap.SugaredLogger
	imageService image.ImagesService

	authMicroservice      authorization.AuthorizationClient
	profileMicroservice   profile.ProfileClient
	musicMicroservice     music.MusicClient
	playlistsMicroservice playlists.PlaylistsClient
}

func NewAPIMicroservices(logger *zap.SugaredLogger, imageService image.ImagesService, auth authorization.AuthorizationClient,
	profile profile.ProfileClient, music music.MusicClient, playlists playlists.PlaylistsClient) APIMicroservices {
	return APIMicroservices{
		logger:                logger,
		imageService:          imageService,
		authMicroservice:      auth,
		profileMicroservice:   profile,
		musicMicroservice:     music,
		playlistsMicroservice: playlists,
	}
}

//nolint:dupl
func (api *APIMicroservices) ParseErrorByCode(ctx echo.Context, requestID string, err error) error {
	if currentError, temp := status.FromError(err); temp {
		if currentError.Code() == codes.Internal {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
		if currentError.Code() == codes.InvalidArgument {
			api.logger.Info(
				zap.String("ID", requestID),
				zap.String("MESSAGE", currentError.Message()),
				zap.Int("ANSWER STATUS", http.StatusBadRequest))

			response := &models.Response{
				Status:  http.StatusBadRequest,
				Message: currentError.Message(),
			}
			jsonResponse, err := easyjson.Marshal(response)
			if err != nil {
				api.logger.Error(
					zap.String("ID", requestID),
					zap.String("ERROR", err.Error()),
					zap.Int("ANSWER STATUS", http.StatusInternalServerError))
				return ctx.NoContent(http.StatusInternalServerError)
			}

			return ctx.JSONBlob(http.StatusOK, jsonResponse)
		}
		if currentError.Code() == codes.PermissionDenied {
			api.logger.Info(
				zap.String("ID", requestID),
				zap.String("MESSAGE", currentError.Message()),
				zap.Int("ANSWER STATUS", http.StatusForbidden))

			response := &models.Response{
				Status:  http.StatusForbidden,
				Message: currentError.Message(),
			}
			jsonResponse, err := easyjson.Marshal(response)
			if err != nil {
				api.logger.Error(
					zap.String("ID", requestID),
					zap.String("ERROR", err.Error()),
					zap.Int("ANSWER STATUS", http.StatusInternalServerError))
				return ctx.NoContent(http.StatusInternalServerError)
			}

			return ctx.JSONBlob(http.StatusOK, jsonResponse)
		}
		if currentError.Code() == codes.NotFound {
			api.logger.Info(
				zap.String("ID", requestID),
				zap.String("MESSAGE", currentError.Message()),
				zap.Int("ANSWER STATUS", http.StatusNotFound))

			response := &models.Response{
				Status:  http.StatusNotFound,
				Message: currentError.Message(),
			}
			jsonResponse, err := easyjson.Marshal(response)
			if err != nil {
				api.logger.Error(
					zap.String("ID", requestID),
					zap.String("ERROR", err.Error()),
					zap.Int("ANSWER STATUS", http.StatusInternalServerError))
				return ctx.NoContent(http.StatusInternalServerError)
			}

			return ctx.JSONBlob(http.StatusOK, jsonResponse)
		}
	}
	return nil
}

func (api *APIMicroservices) Login(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	var authData models.AuthData
	if err := ctx.Bind(&authData); err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	cookies, err := api.authMicroservice.Login(context.Background(),
		&authorization.AuthData{Email: authData.Email, Password: authData.Password})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}
	cookie := &http.Cookie{
		Name:     "Session_cookie",
		Value:    cookies.Cookies,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(constants.CookieLifetime),
	}
	ctx.SetCookie(cookie)

	response := &models.Response{
		Status:  http.StatusOK,
		Message: constants.UserAuthorizedMessage,
	}
	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonResponse)
}

func (api *APIMicroservices) Register(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	var registerData models.RegisterData

	if err := ctx.Bind(&registerData); err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	cookies, err := api.authMicroservice.Register(context.Background(),
		&authorization.RegisterData{Email: registerData.Email, Password: registerData.Password, Nickname: registerData.Nickname})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}
	cookie := &http.Cookie{
		Name:     "Session_cookie",
		Value:    cookies.Cookies,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(constants.CookieLifetime),
	}
	ctx.SetCookie(cookie)

	response := &models.Response{
		Status:  http.StatusCreated,
		Message: constants.UserCreatedMessage,
	}
	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusCreated),
	)
	return ctx.JSONBlob(http.StatusCreated, jsonResponse)
}

//nolint:dupl
func (api *APIMicroservices) GetUserAvatar(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	avatar, err := api.authMicroservice.GetAvatar(context.Background(), &authorization.UserID{ID: int64(userID)})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	response := &models.AvatarResponse{
		Status: http.StatusOK,
		Avatar: avatar.Filename,
	}
	jsonAvatarResponse, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonAvatarResponse)
}

//nolint:dupl
func (api *APIMicroservices) Logout(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	cookie, err := ctx.Cookie("Session_cookie")
	if err != nil {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, marshalErr := easyjson.Marshal(response)
		if marshalErr != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", marshalErr.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	_, err = api.authMicroservice.Logout(context.Background(), &authorization.Cookie{Cookies: cookie.Value})
	if err != nil {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusConflict))
		return ctx.NoContent(http.StatusConflict)
	}
	cookie = &http.Cookie{
		Name:     "Session_cookie",
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().AddDate(0, 0, -1),
	}
	ctx.SetCookie(cookie)

	response := &models.Response{
		Status:  http.StatusOK,
		Message: constants.LoggedOutMessage,
	}
	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonResponse)
}

//nolint:dupl
func (api *APIMicroservices) GetSettings(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	settingsProto, err := api.profileMicroservice.GetSettings(context.Background(), &profile.GetSettingsOptions{ID: int64(userID)})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var settings models.UserSettings
	settings.BindProto(settingsProto)

	jsonResponse, err := easyjson.Marshal(settings)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonResponse)
}

//nolint:dupl,cyclop
func (api *APIMicroservices) UpdateSettings(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	oldSettings, err := api.profileMicroservice.GetSettings(context.Background(), &profile.GetSettingsOptions{ID: int64(userID)})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var newAvatarFilename string
	email := ctx.FormValue("email")
	nickname := ctx.FormValue("nickname")
	oldPassword := ctx.FormValue("old_password")
	newPassword := ctx.FormValue("new_password")
	file, err := ctx.FormFile("avatar")
	if err != nil {
		newAvatarFilename = ""
	} else {
		newAvatarFilename = file.Filename
	}

	if len(newAvatarFilename) != 0 {
		var createdImageData *models.ImageData
		createdImageData, err = api.imageService.CreateImages(
			file,
			os.Getenv("USERS_FULL_PREFIX"),
			map[int]string{
				150: constants.UserAvatarExtension150px,
				500: constants.UserAvatarExtension500px,
			})
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
		newAvatarFilename = createdImageData.Filename
	}

	_, err = api.profileMicroservice.UpdateSettings(context.Background(), &profile.UpdateSettingsOptions{
		UserID:         int64(userID),
		Email:          email,
		Nickname:       nickname,
		OldPassword:    oldPassword,
		NewPassword:    newPassword,
		AvatarFilename: newAvatarFilename,
		OldSettings: &profile.UserSettings{
			Email:       oldSettings.Email,
			Nickname:    oldSettings.Nickname,
			SmallAvatar: oldSettings.SmallAvatar,
			BigAvatar:   oldSettings.BigAvatar,
		},
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	oldAvatarFilename := oldSettings.BigAvatar[len(os.Getenv("USERS_ROOT_PREFIX")) : len(oldSettings.BigAvatar)-len(constants.UserAvatarExtension150px)]
	err = api.imageService.DeleteImages(
		os.Getenv("USERS_FULL_PREFIX"),
		oldAvatarFilename,
		[]string{constants.UserAvatarExtension150px, constants.UserAvatarExtension500px},
		constants.AvatarDefaultFileName,
	)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	response := &models.Response{
		Status:  http.StatusOK,
		Message: constants.SettingsUploadedMessage,
	}
	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonResponse)
}

//nolint:dupl
func (api *APIMicroservices) GenerateCSRF(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	cookie, _ := ctx.Cookie("Session_cookie")
	token, _ := csrf.Tokens.Create(cookie.Value, time.Now().Unix()+constants.CSRFTokenLifetime)

	response := &models.Response{
		Status:  http.StatusOK,
		Message: token,
	}
	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSONBlob(http.StatusOK, jsonResponse)
}

//nolint:dupl
func (api *APIMicroservices) GetHomeTracks(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	var isAuthorized bool
	if userID != -1 {
		isAuthorized = true
	}

	tracksListProto, err := api.musicMicroservice.RandomTracks(context.Background(),
		&music.RandomTracksOptions{Amount: constants.HomePageTracksAmount, UserID: int64(userID), IsAuthorized: isAuthorized})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	tracks := models.Tracks{}
	for _, current := range tracksListProto.Tracks {
		var track models.Track
		track.BindProto(current)
		tracks = append(tracks, track)
	}

	jsonTracks, err := easyjson.Marshal(tracks)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonTracks)
}

//nolint:dupl
func (api *APIMicroservices) GetHomeAlbums(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	albumsListProto, err := api.musicMicroservice.RandomAlbums(context.Background(), &music.RandomAlbumsOptions{Amount: constants.HomePageAlbumsAmount})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	albums := models.Albums{}
	for _, current := range albumsListProto.Albums {
		var album models.Album
		album.BindProto(current)
		albums = append(albums, album)
	}

	jsonAlbums, err := easyjson.Marshal(albums)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonAlbums)
}

//nolint:dupl
func (api *APIMicroservices) GetHomeArtists(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	artistsListProto, err := api.musicMicroservice.RandomArtists(context.Background(), &music.RandomArtistsOptions{Amount: constants.HomePageArtistsAmount})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	artists := models.Artists{}
	for _, current := range artistsListProto.Artists {
		var artist models.Artist
		artist.BindProto(current)
		artists = append(artists, artist)
	}

	jsonArtists, err := easyjson.Marshal(artists)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonArtists)
}

//nolint:dupl
func (api *APIMicroservices) GetArtistProfile(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	var isAuthorized bool
	if userID != -1 {
		isAuthorized = true
	}
	artistID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	artistDataProto, err := api.musicMicroservice.ArtistProfile(context.Background(), &music.ArtistProfileOptions{
		ArtistID:     int64(artistID),
		UserID:       int64(userID),
		IsAuthorized: isAuthorized,
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var artistData models.Artist
	artistData.BindProto(artistDataProto)

	jsonArtistData, err := easyjson.Marshal(artistData)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonArtistData)
}

func (api *APIMicroservices) IncrementListenCount(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	var trackID models.TrackID
	err := ctx.Bind(&trackID)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	_, err = api.musicMicroservice.IncrementListenCount(context.Background(), &music.IncrementListenCountOptions{ID: trackID.ID})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	response := &models.Response{
		Status:  http.StatusOK,
		Message: "Incremented track listen count",
	}
	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonResponse)
}

//nolint:dupl
func (api *APIMicroservices) GetAlbumPage(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	var isAuthorized bool
	if userID != -1 {
		isAuthorized = true
	}
	albumID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	albumDataProto, err := api.musicMicroservice.AlbumPage(context.Background(), &music.AlbumPageOptions{
		AlbumID:      int64(albumID),
		UserID:       int64(userID),
		IsAuthorized: isAuthorized,
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var albumData models.AlbumPage
	albumData.BindProto(albumDataProto)

	jsonAlbumData, err := easyjson.Marshal(albumData)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonAlbumData)
}

func (api *APIMicroservices) SearchMusic(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	var isAuthorized bool
	if userID != -1 {
		isAuthorized = true
	}

	text := ctx.FormValue("text")

	searchResultProto, err := api.musicMicroservice.Find(context.Background(), &music.FindOptions{Text: text, UserID: int64(userID), IsAuthorized: isAuthorized})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var searchResult models.SearchResult
	searchResult.BindProto(searchResultProto)

	jsonSearchResult, err := easyjson.Marshal(searchResult)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonSearchResult)
}

//nolint:dupl,cyclop
func (api *APIMicroservices) CreatePlaylist(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	var artworkFilename, artworkColor string
	title := ctx.FormValue("title")
	isPublic, _ := strconv.ParseBool(ctx.FormValue("is_public"))
	artwork, err := ctx.FormFile("artwork")
	if err != nil {
		artworkFilename = ""
	} else {
		artworkFilename = artwork.Filename
	}

	if len(artworkFilename) != 0 {
		var createdImageData *models.ImageData
		createdImageData, err = api.imageService.CreateImages(
			artwork,
			os.Getenv("PLAYLIST_FULL_PREFIX"),
			map[int]string{
				100: constants.PlaylistArtworkExtension100px,
				384: constants.PlaylistArtworkExtension384px,
			})
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
		artworkFilename = createdImageData.Filename
		artworkColor = createdImageData.ArtworkColor
	}

	playlistIDProto, err := api.playlistsMicroservice.CreatePlaylist(context.Background(), &playlists.CreatePlaylistOptions{
		UserID:       int64(userID),
		Title:        title,
		IsPublic:     isPublic,
		Artwork:      artworkFilename,
		ArtworkColor: artworkColor,
	})
	if err != nil {
		deleteErr := api.imageService.DeleteImages(
			os.Getenv("PLAYLIST_FULL_PREFIX"),
			artworkFilename,
			[]string{constants.PlaylistArtworkExtension100px, constants.PlaylistArtworkExtension384px},
			constants.PlaylistArtworkDefaultFilename,
		)
		if deleteErr != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", deleteErr.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var playlistID models.PlaylistID
	playlistID.BindProto(playlistIDProto)

	jsonPlaylistID, err := easyjson.Marshal(playlistID)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusCreated),
	)
	return ctx.JSONBlob(http.StatusCreated, jsonPlaylistID)
}

//nolint:cyclop,dupl
func (api *APIMicroservices) UpdatePlaylist(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	playlistID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	var artworkFilename, artworkColor string
	var isPublic bool
	title := ctx.FormValue("title")
	artwork, err := ctx.FormFile("artwork")
	isPublic, _ = strconv.ParseBool(ctx.FormValue("is_public"))
	if err != nil {
		artworkFilename = ""
	} else {
		artworkFilename = artwork.Filename
	}
	if len(artworkFilename) != 0 {
		var createdImageData *models.ImageData
		createdImageData, err = api.imageService.CreateImages(
			artwork,
			os.Getenv("PLAYLIST_FULL_PREFIX"),
			map[int]string{
				100: constants.PlaylistArtworkExtension100px,
				384: constants.PlaylistArtworkExtension384px,
			})
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
		artworkFilename = createdImageData.Filename
		artworkColor = createdImageData.ArtworkColor
	}

	artworkProto, err := api.playlistsMicroservice.UpdatePlaylist(context.Background(), &playlists.UpdatePlaylistOptions{
		PlaylistID:   int64(playlistID),
		Title:        title,
		UserID:       int64(userID),
		Artwork:      artworkFilename,
		ArtworkColor: artworkColor,
		IsPublic:     isPublic,
	})
	if err != nil {
		deleteErr := api.imageService.DeleteImages(
			os.Getenv("PLAYLIST_FULL_PREFIX"),
			artworkFilename,
			[]string{constants.PlaylistArtworkExtension100px, constants.PlaylistArtworkExtension384px},
			constants.PlaylistArtworkDefaultFilename,
		)
		if deleteErr != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", deleteErr.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return api.ParseErrorByCode(ctx, requestID, err)
	}
	err = api.imageService.DeleteImages(
		os.Getenv("PLAYLIST_FULL_PREFIX"),
		artworkProto.OldArtworkFilename,
		[]string{constants.PlaylistArtworkExtension100px, constants.PlaylistArtworkExtension384px},
		constants.PlaylistArtworkDefaultFilename,
	)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	response := &models.PlaylistArtworkColor{ArtworkColor: artworkProto.ArtworkColor}
	jsonPlaylistArtworkColor, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonPlaylistArtworkColor)
}

//nolint:dupl
func (api *APIMicroservices) DeletePlaylist(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	playlistID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	oldArtworkProto, err := api.playlistsMicroservice.DeletePlaylist(context.Background(), &playlists.DeletePlaylistOptions{
		PlaylistID: int64(playlistID),
		UserID:     int64(userID),
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}
	err = api.imageService.DeleteImages(
		os.Getenv("PLAYLIST_FULL_PREFIX"),
		oldArtworkProto.OldArtworkFilename,
		[]string{constants.PlaylistArtworkExtension100px, constants.PlaylistArtworkExtension384px},
		constants.PlaylistArtworkDefaultFilename,
	)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	response := &models.Response{
		Status:  http.StatusOK,
		Message: constants.PlaylistDeletedMessage,
	}
	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonResponse)
}

//nolint:dupl
func (api *APIMicroservices) AddTrack(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	var requestData models.PlaylistTrack
	if err := ctx.Bind(&requestData); err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	_, err := api.playlistsMicroservice.AddTrack(context.Background(), &playlists.AddTrackOptions{
		TrackID:    requestData.TrackID,
		PlaylistID: requestData.PlaylistID,
		UserID:     int64(userID),
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	response := &models.Response{
		Status:  http.StatusCreated,
		Message: constants.TrackAddedToPlaylistMessage,
	}
	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusCreated),
	)
	return ctx.JSONBlob(http.StatusCreated, jsonResponse)
}

//nolint:dupl
func (api *APIMicroservices) DeleteTrack(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	var requestData models.PlaylistTrack
	if err := ctx.Bind(&requestData); err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	_, err := api.playlistsMicroservice.DeleteTrack(context.Background(), &playlists.DeleteTrackOptions{
		PlaylistID: requestData.PlaylistID,
		TrackID:    requestData.TrackID,
		UserID:     int64(userID),
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	response := &models.Response{
		Status:  http.StatusOK,
		Message: constants.TrackDeletedFromPlaylistMessage,
	}
	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonResponse)
}

//nolint:dupl
func (api *APIMicroservices) GetUserPlaylists(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	playlistsProto, err := api.musicMicroservice.UserPlaylists(context.Background(), &music.UserPlaylistsOptions{UserID: int64(userID)})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var userPlaylists models.UserPlaylists
	userPlaylists.BindProto(playlistsProto)

	jsonUserPlaylists, err := easyjson.Marshal(userPlaylists)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonUserPlaylists)
}

func (api *APIMicroservices) GetPlaylistPage(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	playlistID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	playlistPageDataProto, err := api.musicMicroservice.PlaylistPage(context.Background(), &music.PlaylistPageOptions{
		PlaylistID: int64(playlistID),
		UserID:     int64(userID),
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var playlistPage models.PlaylistPage
	playlistPage.BindProto(playlistPageDataProto)

	jsonPlaylistPage, err := easyjson.Marshal(playlistPage)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonPlaylistPage)
}

//nolint:dupl
func (api *APIMicroservices) AddTrackToFavorites(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	trackID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	_, err = api.musicMicroservice.AddTrackToFavorites(context.Background(), &music.AddTrackToFavoritesOptions{
		UserID:  int64(userID),
		TrackID: int64(trackID),
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	response := &models.Response{
		Status:  http.StatusCreated,
		Message: constants.TrackAddedToFavoritesMessage,
	}
	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusCreated),
	)
	return ctx.JSONBlob(http.StatusCreated, jsonResponse)
}

//nolint:dupl
func (api *APIMicroservices) DeleteTrackFromFavorites(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	trackID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	_, err = api.musicMicroservice.DeleteTrackFromFavorites(context.Background(), &music.DeleteTrackFromFavoritesOptions{
		UserID:  int64(userID),
		TrackID: int64(trackID),
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	response := &models.Response{
		Status:  http.StatusCreated,
		Message: constants.TrackDeletedFromFavoritesMessage,
	}
	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusCreated),
	)
	return ctx.JSONBlob(http.StatusCreated, jsonResponse)
}

//nolint:dupl
func (api *APIMicroservices) GetUserFavorites(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))

		response := &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		}
		jsonResponse, err := easyjson.Marshal(response)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.JSONBlob(http.StatusOK, jsonResponse)
	}

	tracksListProto, err := api.musicMicroservice.GetFavoriteTracks(context.Background(),
		&music.UserFavoritesOptions{UserID: int64(userID)})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	tracks := models.Tracks{}
	for _, current := range tracksListProto.Tracks {
		var track models.Track
		track.BindProto(current)
		tracks = append(tracks, track)
	}

	jsonTracks, err := easyjson.Marshal(tracks)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonTracks)
}

func (api *APIMicroservices) DeletePlaylistArtwork(ctx echo.Context) error {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		api.logger.Error(
			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	userID, ok := ctx.Get("USER_ID").(int)
	if !ok {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	playlistID, err := strconv.Atoi(ctx.FormValue("id"))
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	oldArtworkProto, err := api.playlistsMicroservice.DeletePlaylistArtwork(context.Background(), &playlists.DeletePlaylistArtworkOptions{
		PlaylistID: int64(playlistID),
		UserID:     int64(userID),
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	err = api.imageService.DeleteImages(
		os.Getenv("PLAYLIST_FULL_PREFIX"),
		oldArtworkProto.OldArtworkFilename,
		[]string{constants.PlaylistArtworkExtension100px, constants.PlaylistArtworkExtension384px},
		constants.PlaylistArtworkDefaultFilename,
	)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	response := &models.PlaylistArtworkColor{ArtworkColor: constants.PlaylistArtworkDefaultColor}
	jsonPlaylistArtworkColor, err := easyjson.Marshal(response)
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSONBlob(http.StatusOK, jsonPlaylistArtworkColor)
}

func (api *APIMicroservices) Init(server *echo.Echo) {
	// Authorization
	server.POST("/api/v1/user/signin", api.Login)
	server.POST("/api/v1/user/signup", api.Register)
	server.GET("/api/v1/auth", api.GetUserAvatar)
	server.POST("/api/v1/user/logout", api.Logout)

	// Profile
	server.GET("/api/v1/user/settings", api.GetSettings)
	server.PATCH("/api/v1/user/settings", api.UpdateSettings)

	// Music
	server.GET("/api/v1/home/tracks", api.GetHomeTracks)
	server.GET("/api/v1/home/albums", api.GetHomeAlbums)
	server.GET("/api/v1/home/artists", api.GetHomeArtists)
	server.GET("/api/v1/artist/:id", api.GetArtistProfile)
	server.GET("/api/v1/album/:id", api.GetAlbumPage)
	server.POST("/api/v1/inc_listencount", api.IncrementListenCount)
	server.GET("/api/v1/music/search", api.SearchMusic)
	server.GET("/api/v1/playlists", api.GetUserPlaylists)
	server.GET("/api/v1/playlists/:id", api.GetPlaylistPage)
	server.POST("api/v1/track/like/:id", api.AddTrackToFavorites)
	server.DELETE("api/v1/track/like/:id", api.DeleteTrackFromFavorites)
	server.GET("api/v1/track/favorites", api.GetUserFavorites)

	// Playlists
	server.POST("/api/v1/playlists", api.CreatePlaylist)
	server.PATCH("/api/v1/playlists/:id", api.UpdatePlaylist)
	server.DELETE("/api/v1/playlists/:id", api.DeletePlaylist)
	server.POST("/api/v1/playlist/track", api.AddTrack)
	server.DELETE("/api/v1/playlist/track", api.DeleteTrack)
	server.DELETE("/api/v1/playlist/artwork", api.DeletePlaylistArtwork)

	// CSRF
	server.GET("/api/v1/csrf", api.GenerateCSRF)
}
