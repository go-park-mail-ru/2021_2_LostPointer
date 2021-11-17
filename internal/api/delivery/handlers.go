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
	imageService image.ImageService

	authMicroservice      authorization.AuthorizationClient
	profileMicroservice   profile.ProfileClient
	musicMicroservice     music.MusicClient
	playlistsMicroservice playlists.PlaylistsClient
}

func NewAPIMicroservices(logger *zap.SugaredLogger, imageService image.ImageService, auth authorization.AuthorizationClient,
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

func (api *APIMicroservices) ParseErrorByCode(ctx echo.Context, requestID string, err error) error {
	if e, temp := status.FromError(err); temp {
		if e.Code() == codes.Internal {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
		if e.Code() == codes.InvalidArgument || e.Code() == codes.NotFound {
			api.logger.Info(
				zap.String("ID", requestID),
				zap.String("MESSAGE", e.Message()),
				zap.Int("ANSWER STATUS", http.StatusBadRequest))
			return ctx.JSON(http.StatusOK, &models.Response{
				Status:  http.StatusBadRequest,
				Message: e.Message(),
			})
		}
		if e.Code() == codes.PermissionDenied {
			api.logger.Info(
				zap.String("ID", requestID),
				zap.String("MESSAGE", e.Message()),
				zap.Int("ANSWER STATUS", http.StatusForbidden))
			return ctx.JSON(http.StatusOK, &models.Response{
				Status:  http.StatusForbidden,
				Message: e.Message(),
			})
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
	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: constants.UserAuthorizedMessage,
	})
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
	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusCreated),
	)

	return ctx.JSON(http.StatusCreated, &models.Response{
		Status:  http.StatusCreated,
		Message: constants.UserCreatedMessage,
	})
}

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
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
	}

	avatar, err := api.authMicroservice.GetAvatar(context.Background(), &authorization.UserID{ID: int64(userID)})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK,
		struct {
			Status int    `json:"status"`
			Avatar string `json:"avatar"`
		}{http.StatusOK, avatar.Filename})
}

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
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
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

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: constants.LoggedOutMessage,
	})
}

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
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
	}

	settings, err := api.profileMicroservice.GetSettings(context.Background(), &profile.ProfileUserID{ID: int64(userID)})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}
	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, &models.UserSettings{
		Email:       settings.Email,
		Nickname:    settings.Nickname,
		SmallAvatar: settings.SmallAvatar,
		BigAvatar:   settings.BigAvatar,
	})
}

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
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
	}

	oldSettings, err := api.profileMicroservice.GetSettings(context.Background(), &profile.ProfileUserID{ID: int64(userID)})
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
		newAvatarFilename, err = api.imageService.CreateAvatar(file)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
		oldAvatarFilename := oldSettings.BigAvatar[len(os.Getenv("USERS_ROOT_PREFIX")) : len(oldSettings.BigAvatar)-len(constants.LittleAvatarPostfix)]
		err = api.imageService.DeleteAvatar(oldAvatarFilename)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
	}

	_, err = api.profileMicroservice.UpdateSettings(context.Background(), &profile.UploadSettings{
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

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: constants.SettingsUploadedMessage,
	})
}

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
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
	}

	cookie, _ := ctx.Cookie("Session_cookie")
	token, _ := csrf.Tokens.Create(cookie.Value, time.Now().Unix()+constants.CSRFTokenLifetime)
	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: token,
	})
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
		&music.RandomTracksOptions{Amount: constants.TracksCollectionLimit, IsAuthorized: isAuthorized})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	tracks := make([]models.Track, 0, constants.TracksCollectionLimit)
	for _, current := range tracksListProto.Tracks {
		var track models.Track
		track.BindProtoTrack(current)
		tracks = append(tracks, track)
	}
	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, tracks)
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

	albumsListProto, err := api.musicMicroservice.RandomAlbums(context.Background(), &music.RandomAlbumsOptions{Amount: constants.AlbumCollectionLimit})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	albums := make([]models.Album, 0, constants.AlbumCollectionLimit)
	for _, current := range albumsListProto.Albums {
		var album models.Album
		album.BindProtoAlbum(current)
		albums = append(albums, album)
	}
	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, albums)
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

	artistsListProto, err := api.musicMicroservice.RandomArtists(context.Background(), &music.RandomArtistsOptions{Amount: constants.ArtistsCollectionLimit})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	artists := make([]models.Artist, 0, constants.ArtistsCollectionLimit)
	for _, current := range artistsListProto.Artists {
		var artist models.Artist
		artist.BindProtoArtist(current)
		artists = append(artists, artist)
	}
	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, artists)
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
		IsAuthorized: isAuthorized,
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var artistData models.Artist
	artistData.BindProtoArtist(artistDataProto)
	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, artistData)
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

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: "Incremented track listen count",
	})
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
		IsAuthorized: isAuthorized,
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var albumData models.AlbumPage
	albumData.BindProtoAlbumPage(albumDataProto)
	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, albumData)
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

	searchResultProto, err := api.musicMicroservice.Find(context.Background(), &music.FindOptions{Text: text, IsAuthorized: isAuthorized})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var searchResult models.SearchResult
	searchResult.BindProtoSearchResponse(searchResultProto)

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, searchResult)
}

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
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
	}

	title := ctx.FormValue("title")
	artwork, err := ctx.FormFile("artwork")
	var artworkFilename string
	if err != nil {
		artworkFilename = ""
	} else {
		artworkFilename = artwork.Filename
	}
	artworkColor := constants.DefaultPlaylistArtworkColor
	if len(artworkFilename) != 0 {
		artworkFilename, artworkColor, err = api.imageService.CreatePlaylistArtwork(artwork)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
	}

	playlistIDProto, err := api.playlistsMicroservice.CreatePlaylist(context.Background(), &playlists.CreatePlaylistOptions{
		UserID:       int64(userID),
		Title:        title,
		Artwork:      artworkFilename,
		ArtworkColor: artworkColor,
	})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var playlistID models.PlaylistID
	playlistID.BindProto(playlistIDProto)

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusCreated),
	)
	return ctx.JSON(http.StatusCreated, playlistID)
}

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
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
	}

	playlistID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	title := ctx.FormValue("title")
	artwork, err := ctx.FormFile("artwork")
	var artworkFilename string
	if err != nil {
		artworkFilename = ""
	} else {
		artworkFilename = artwork.Filename
	}
	var artworkColor string
	if len(artworkFilename) != 0 {
		artworkFilename, artworkColor, err = api.imageService.CreatePlaylistArtwork(artwork)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
	}

	oldArtworkProto, err := api.playlistsMicroservice.UpdatePlaylist(context.Background(), &playlists.UpdatePlaylistOptions{
		PlaylistID:   int64(playlistID),
		Title:        title,
		UserID:       int64(userID),
		Artwork:      artworkFilename,
		ArtworkColor: artworkColor,
	})
	if err != nil {
		_ = api.imageService.DeletePlaylistArtwork(artworkFilename)
		return api.ParseErrorByCode(ctx, requestID, err)
	}
	_ = api.imageService.DeletePlaylistArtwork(oldArtworkProto.OldArtworkFilename)

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: constants.PlaylistTitleUpdatedMessage,
	})
}

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
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
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
	_ = api.imageService.DeletePlaylistArtwork(oldArtworkProto.OldArtworkFilename)

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: constants.PlaylistDeletedMessage,
	})
}

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
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
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

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusCreated),
	)
	return ctx.JSON(http.StatusCreated, &models.Response{
		Status:  http.StatusCreated,
		Message: constants.TrackAddedToPlaylistMessage,
	})
}

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
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
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

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: constants.TrackDeletedFromPlaylistMessage,
	})
}

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

	playlistsProto, err := api.musicMicroservice.UserPlaylists(context.Background(), &music.UserPlaylistsOptions{UserID: int64(userID)})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	var userPlaylists models.UserPlaylists
	userPlaylists.BindProto(playlistsProto)

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, userPlaylists)
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

	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
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

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK, playlistPage)
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

	// Playlists
	server.POST("/api/v1/playlists", api.CreatePlaylist)
	server.PATCH("/api/v1/playlists/:id", api.UpdatePlaylist)
	server.DELETE("/api/v1/playlists/:id", api.DeletePlaylist)
	server.POST("/api/v1/playlist/track", api.AddTrack)
	server.DELETE("/api/v1/playlist/track", api.DeleteTrack)

	// CSRF
	server.GET("/api/v1/csrf", api.GenerateCSRF)
}
