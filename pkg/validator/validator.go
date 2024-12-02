package validator

import (
	"net/http"

	valid "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *valid.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func NewCustomeValidator() *CustomValidator {
	return &CustomValidator{
		validator: valid.New(),
	}
}
