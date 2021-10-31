package delivery

import (
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrackDelivery_Home(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

	track := models.Track{
		Id:          1,
		Title:       "awa",
		Artist:      "awa",
		Explicit:    true,
		Genre:       "awa",
		Number:      1,
		File:        "awa",
		ListenCount: 1,
		Duration:    1,
		Lossless:    true,
		Album: models.Album{
			Id:      1,
			Title:   "awa",
			Artwork: "awa",
		},
	}
	tests := []struct {
		name          string
		param         string
		useCaseMock   *mock.MockTrackUseCase
		expected      []models.Track
		expectedError bool
	}{
		{
			name:  "get home",
			param: "1",
			useCaseMock: &mock.MockTrackUseCase{
				GetHomeFunc: func(amount int, isAuthorized bool) ([]models.Track, *models.CustomError) {
					return []models.Track{track}, nil
				}},
			expected:      []models.Track{track},
			expectedError: false,
		},
		{
			name:  "GetHome() error",
			param: "1",
			useCaseMock: &mock.MockTrackUseCase{
				GetHomeFunc: func(amount int, isAuthorized bool) ([]models.Track, *models.CustomError) {
					return nil, &models.CustomError{
						ErrorType:     http.StatusInternalServerError,
						OriginalError: errors.New("error"),
						Message:       "error",
					}
				}},
			expected:      []models.Track{track},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := echo.New()
			request := httptest.NewRequest(echo.GET, "/", nil)
			recorder := httptest.NewRecorder()
			ctx := server.NewContext(request, recorder)
			ctx.SetPath("api/v1/home/artists")
			ctx.SetParamNames("id")
			ctx.SetParamValues(test.param)
			ctx.Set("REQUEST_ID", "1")
			ctx.Set("USER_ID", 1)
			delivery := NewTrackDelivery(test.useCaseMock, logger)
			_ = delivery.Home(ctx)
			body := recorder.Body
			status := recorder.Result().Status
			var result []models.Track
			_ = json.Unmarshal(body.Bytes(), &result)
			if test.expectedError {
				assert.Equal(t, "500 Internal Server Error", status)
			} else {
				assert.Equal(t, test.expected, result)
				assert.Equal(t, "200 OK", status)
			}
		})
	}
}
