package delivery

import (
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlaylistDelivery_Home(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

	playlist := models.Playlist{
		Id:   1,
		Name: "awa",
		User: 1,
	}
	tests := []struct {
		name          string
		param         string
		useCaseMock   *mock.MockPlaylistUseCase
		expected      []models.Playlist
		expectedError bool
	}{
		{
			name:  "get home",
			param: "1",
			useCaseMock: &mock.MockPlaylistUseCase{
				GetHomeFunc: func(amount int) ([]models.Playlist, *models.CustomError) {
					return []models.Playlist{playlist}, nil
				}},
			expected:      []models.Playlist{playlist},
			expectedError: false,
		},
		{
			name:  "GetHome() error",
			param: "1",
			useCaseMock: &mock.MockPlaylistUseCase{
				GetHomeFunc: func(amount int) ([]models.Playlist, *models.CustomError) {
					return nil, &models.CustomError{
						ErrorType:     http.StatusInternalServerError,
						OriginalError: errors.New("error"),
						Message:       "error",
					}
				}},
			expected:      []models.Playlist{playlist},
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
			ctx.Set("IS_AUTHORIZED", true)
			delivery := NewPlaylistDelivery(test.useCaseMock, logger)
			_ = delivery.Home(ctx)
			body := recorder.Body
			status := recorder.Result().Status
			var result []models.Playlist
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
