package delivery

import (
	"2021_2_LostPointer/internal/search"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type SearchDelivery struct {
	searchLogic search.SearchUsecase
	logger      *zap.SugaredLogger
}

func NewSearchDelivery(searchLogic search.SearchUsecase, logger *zap.SugaredLogger) SearchDelivery {
	return SearchDelivery{
		searchLogic: searchLogic,
		logger:      logger,
	}
}

func (searchD SearchDelivery) Search(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)

	searchText := ctx.FormValue("text")

	searchResponseData, customError := searchD.searchLogic.SearchByText(searchText)
	if customError != nil {
		if customError.ErrorType == http.StatusInternalServerError {
			searchD.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.OriginalError.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)
			return ctx.NoContent(http.StatusInternalServerError)
		}
	}

	return ctx.JSON(http.StatusOK, searchResponseData)
}

func (searchD SearchDelivery) InitHandlers(server *echo.Echo) {
	server.GET("/api/v1/search", searchD.Search)
}
