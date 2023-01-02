package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func UpdateExpensesById(c echo.Context) error {
	id := c.Param("id")
	b := Expense{}
	err := c.Bind(&b)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	err = b.validation()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := Db.Prepare(`UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING *`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "unable to setup query statement" + err.Error()})
	}

	ex := Expense{}

	update := stmt.QueryRow(b.Title, b.Amount, b.Note, pq.Array(b.Tags), id)
	err = update.Scan(&ex.Id, &ex.Title, &ex.Amount, &ex.Note, pq.Array(&ex.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "updated expense's not found"})
	case nil:
		return c.JSON(http.StatusOK, ex)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan updated expense:" + err.Error()})
	}

}
