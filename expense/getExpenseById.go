package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GetExpensesById(c echo.Context) error {
	id := c.Param("id")

	stmt, err := Db.Prepare(`SELECT * FROM expenses WHERE id = $1`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "unable to setup query statement" + err.Error()})
	}

	row := stmt.QueryRow(id)
	ex := Expense{}

	err = row.Scan(&ex.Id, &ex.Title, &ex.Amount, &ex.Note, pq.Array(&ex.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense's not found"})
	case nil:
		return c.JSON(http.StatusOK, ex)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expense:" + err.Error()})
	}

}
