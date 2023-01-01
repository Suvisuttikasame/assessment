package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateExpenses(c echo.Context) error {
	ex := Expense{}
	err := c.Bind(&ex)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, ex)
}
