package apiv1

import (
	"net/http"

	"github.com/aeilang/urlshortener/internal/model"
	"github.com/aeilang/urlshortener/internal/service"
	"github.com/labstack/echo/v4"
)

type URLHandler struct {
	urlService service.URLService
	baseURL    string
}

func NewURLHandler(urlService service.URLService, baseURL string) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		baseURL:    baseURL,
	}
}

func (h *URLHandler) CreateURL(c echo.Context) error {
	var req model.CreateURLRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	url, err := h.urlService.CreateURL(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, model.CreateURLResponse{
		ShortURL:  h.baseURL + "/" + url.ShortCode,
		ExpiresAt: url.ExpiresAt,
	})
}

func (h *URLHandler) RedirectURL(c echo.Context) error {
	shortCode := c.Param("code")
	url, err := h.urlService.GetURL(c.Request().Context(), shortCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if url == nil {
		return echo.NewHTTPError(http.StatusNotFound, "URL not found")
	}

	return c.Redirect(http.StatusMovedPermanently, url.OriginalUrl)
}
