package delivery

import (
	"2021_2_LostPointer/internal/csrf"
	"2021_2_LostPointer/pkg/image"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"2021_2_LostPointer/internal/constants"
	authorization "2021_2_LostPointer/internal/microservices/authorization/proto"
	music "2021_2_LostPointer/internal/microservices/music/proto"
	profile "2021_2_LostPointer/internal/microservices/profile/proto"
	"2021_2_LostPointer/internal/models"
)

type APIMicroservices struct {
	logger         *zap.SugaredLogger
	avatarsService image.AvatarsService

	authMicroservice    authorization.AuthorizationClient
	profileMicroservice profile.ProfileClient
	musicMicroservice   music.MusicClient
}

func NewAPIMicroservices(logger *zap.SugaredLogger, avatarsService image.AvatarsService, auth authorization.AuthorizationClient,
	profile profile.ProfileClient, music music.MusicClient) APIMicroservices {
	return APIMicroservices{
		logger:              logger,
		avatarsService:      avatarsService,
		authMicroservice:    auth,
		profileMicroservice: profile,
		musicMicroservice:   music,
	}
}

type MyStringType string
type MyIntType int

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
		&authorization.AuthData{Login: authData.Email, Password: authData.Password})
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
		&authorization.RegisterData{Login: registerData.Email, Password: registerData.Password, Nickname: registerData.Nickname})
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
	cookie.Expires = time.Now().AddDate(0, 0, -1)

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
		newAvatarFilename, err = api.avatarsService.CreateImage(file)
		if err != nil {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
		oldAvatarFilename := oldSettings.BigAvatar[0 : len(oldSettings.BigAvatar)-11]
		err = api.avatarsService.DeleteImage(oldAvatarFilename)
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
		Status: http.StatusOK,
		Message: "Incremented track listen count",
	})
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
	server.POST("/api/v1/inc_listencount", api.IncrementListenCount)

	// CSRF
	server.GET("/api/v1/csrf", api.GenerateCSRF)
}
