package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetExpenses(c echo.Context) error {
	stmt, err := Db.Prepare(`SELECT * FROM expenses`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "unable to setup query statement" + err.Error()})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "unable to query statement" + err.Error()})
	}

	exs := []Expense{}

	for rows.Next() {
		ex := Expense{}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "unable to scan expense" + err.Error()})
		}
		exs = append(exs, ex)
	}
	return c.JSON(http.StatusOK, exs)

}
