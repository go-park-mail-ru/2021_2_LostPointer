package delivery

//func TestArtistDelivery_GetProfile(t *testing.T) {
//	config := zap.NewDevelopmentConfig()
//	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
//	prLogger, _ := config.Build()
//	logger := prLogger.Sugar()
//	defer prLogger.Sync()
//
//	track := models.Track{
//		Id:       1,
//		Title:    "awa",
//		Explicit: true,
//		File:     "awa",
//		Duration: 1,
//		Lossless: true,
//		Cover:    "awa",
//	}
//	tracksUnAuth := track
//	tracksUnAuth.File = ""
//
//	album := models.Album{
//		Id:             1,
//		Title:          "awa",
//		Year:           1,
//		Artwork:        "awa",
//		TracksDuration: 1,
//	}
//
//	artist := models.Artist{
//		Id:     1,
//		Name:   "awa",
//		Avatar: "awa",
//		Tracks: []models.Track{track},
//		Albums: []models.Album{album},
//	}
//	artistUnAuth := artist
//	artistUnAuth.Tracks = []models.Track{tracksUnAuth}
//
//	tests := []struct {
//		name          string
//		param         string
//		useCaseMock   *mock.MockArtistUseCase
//		expected      *models.Artist
//		expectedError bool
//	}{
//		{
//			name:  "get profile",
//			param: "1",
//			useCaseMock: &mock.MockArtistUseCase{
//				GetProfileFunc: func(id int, isAuthorized bool) (*models.Artist, *models.CustomError) {
//					return &artist, nil
//				}},
//			expected:      &artist,
//			expectedError: false,
//		},
//		{
//			name:  "wrong param",
//			param: "str",
//			useCaseMock: &mock.MockArtistUseCase{
//				GetProfileFunc: func(id int, isAuthorized bool) (*models.Artist, *models.CustomError) {
//					return &artist, nil
//				}},
//			expected:      &artist,
//			expectedError: true,
//		},
//		{
//			name:  "GetProfile() error",
//			param: "1",
//			useCaseMock: &mock.MockArtistUseCase{
//				GetProfileFunc: func(id int, isAuthorized bool) (*models.Artist, *models.CustomError) {
//					return nil, &models.CustomError{
//						ErrorType:     http.StatusInternalServerError,
//						OriginalError: errors.New("error"),
//						Message:       "error",
//					}
//				}},
//			expected:      &artist,
//			expectedError: true,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			server := echo.New()
//			request := httptest.NewRequest(echo.GET, "/", nil)
//			recorder := httptest.NewRecorder()
//			ctx := server.NewContext(request, recorder)
//			ctx.SetPath("api/v1/artist/:id")
//			ctx.SetParamNames("id")
//			ctx.SetParamValues(test.param)
//			ctx.Set("REQUEST_ID", "1")
//			ctx.Set("IS_AUTHORIZED", true)
//			delivery := NewArtistDelivery(test.useCaseMock, logger)
//			_ = delivery.GetProfile(ctx)
//			body := recorder.Body
//			status := recorder.Result().Status
//			var result models.Artist
//			_ = json.Unmarshal(body.Bytes(), &result)
//			if test.expectedError {
//				assert.Equal(t, "500 Internal Server Error", status)
//			} else {
//				assert.Equal(t, *test.expected, result)
//				assert.Equal(t, "200 OK", status)
//			}
//		})
//	}
//}
//
//func TestArtistDelivery_Home(t *testing.T) {
//	config := zap.NewDevelopmentConfig()
//	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
//	prLogger, _ := config.Build()
//	logger := prLogger.Sugar()
//	defer prLogger.Sync()
//
//	artist := models.Artist{
//		Id:     1,
//		Name:   "awa",
//		Avatar: "awa",
//	}
//	tests := []struct {
//		name          string
//		param         string
//		useCaseMock   *mock.MockArtistUseCase
//		expected      []models.Artist
//		expectedError bool
//	}{
//		{
//			name:  "get home",
//			param: "1",
//			useCaseMock: &mock.MockArtistUseCase{
//				GetHomeFunc: func(amount int) ([]models.Artist, *models.CustomError) {
//					return []models.Artist{artist}, nil
//				}},
//			expected:      []models.Artist{artist},
//			expectedError: false,
//		},
//		{
//			name:  "GetHome() error",
//			param: "1",
//			useCaseMock: &mock.MockArtistUseCase{
//				GetHomeFunc: func(amount int) ([]models.Artist, *models.CustomError) {
//					return nil, &models.CustomError{
//						ErrorType:     http.StatusInternalServerError,
//						OriginalError: errors.New("error"),
//						Message:       "error",
//					}
//				}},
//			expected:      []models.Artist{artist},
//			expectedError: true,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			server := echo.New()
//			request := httptest.NewRequest(echo.GET, "/", nil)
//			recorder := httptest.NewRecorder()
//			ctx := server.NewContext(request, recorder)
//			ctx.SetPath("api/v1/home/artists")
//			ctx.SetParamNames("id")
//			ctx.SetParamValues(test.param)
//			ctx.Set("REQUEST_ID", "1")
//			ctx.Set("IS_AUTHORIZED", true)
//			delivery := NewArtistDelivery(test.useCaseMock, logger)
//			_ = delivery.Home(ctx)
//			body := recorder.Body
//			status := recorder.Result().Status
//			var result []models.Artist
//			_ = json.Unmarshal(body.Bytes(), &result)
//			if test.expectedError {
//				assert.Equal(t, "500 Internal Server Error", status)
//			} else {
//				assert.Equal(t, test.expected, result)
//				assert.Equal(t, "200 OK", status)
//			}
//		})
//	}
//}
