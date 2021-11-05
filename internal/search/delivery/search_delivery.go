package delivery

import (
	"2021_2_LostPointer/internal/search"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	searcher search.Searcher
	logger      *zap.SugaredLogger
}

func NewHandler(searcher search.Searcher, logger *zap.SugaredLogger) Handler {
	return Handler{
		searcher:	searcher,
		logger:     logger,
	}
}

func (h *Handler) Search(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)

	searchText := ctx.FormValue("text")

	SearchPayload, customError := h.searcher.ByText(searchText)
	if customError != nil {
		if customError.ErrorType == http.StatusInternalServerError {
			h.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.OriginalError.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)
			return ctx.NoContent(http.StatusInternalServerError)
		}
	}

	return ctx.JSON(http.StatusOK, SearchPayload)
}

func (h Handler) InitHandlers(server *echo.Echo) {
	server.GET("/api/v1/search", h.Search)
}
