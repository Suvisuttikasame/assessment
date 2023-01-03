//go:build unit

package expense

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpensesValidation(t *testing.T) {
	t.Run("should return title error : this field should not empty. when title input is empty", func(t *testing.T) {
		//arrange
		e := echo.New()
		reqBody := bytes.NewBufferString(`{
			"title": "",
			"amount": 39000,
			"note": "buy a new phone",
			"tags": ["gadget", "shopping"]
		}`)
		req := httptest.NewRequest(http.MethodPost, "/expenses", reqBody)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		var r Err

		//action
		err := CreateExpenses(c)
		assert.Nil(t, err)
		err = json.NewDecoder(rec.Body).Decode(&r)

		//assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "title error : this field should not empty.", r.Message)
	})

	t.Run("should return amount error : this field should not less than 0. when amount input is minus value", func(t *testing.T) {
		//arrange
		e := echo.New()
		reqBody := bytes.NewBufferString(`{
			"title": "buy a new phone",
			"amount": -199,
			"note": "buy a new phone",
			"tags": ["gadget", "shopping"]
		}`)
		req := httptest.NewRequest(http.MethodPost, "/expenses", reqBody)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		var r Err

		//action
		err := CreateExpenses(c)
		assert.Nil(t, err)
		err = json.NewDecoder(rec.Body).Decode(&r)

		//assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "amount error : this field should not less than 0.", r.Message)
	})

	t.Run("should return tags error : this field should have at least 1. when tags input is empty", func(t *testing.T) {
		//arrange
		e := echo.New()
		reqBody := bytes.NewBufferString(`{
			"title": "buy a new phone",
			"amount": 39000,
			"note": "buy a new phone",
			"tags": []
		}`)
		req := httptest.NewRequest(http.MethodPost, "/expenses", reqBody)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		var r Err

		//action
		err := CreateExpenses(c)
		assert.Nil(t, err)
		err = json.NewDecoder(rec.Body).Decode(&r)

		//assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "tags error : this field should have at least 1.", r.Message)
	})

}

func TestUpdateExpensesByIdValidation(t *testing.T) {
	t.Run("should return title error : this field should not empty. when title input is empty", func(t *testing.T) {
		//arrange
		e := echo.New()
		id := "1"
		reqBody := bytes.NewBufferString(`{
			"title": "",
			"amount": 39000,
			"note": "buy a new phone",
			"tags": ["gadget", "shopping"]
		}`)
		req := httptest.NewRequest(http.MethodPut, "/expenses/"+id, reqBody)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		var r Err

		//action
		err := UpdateExpensesById(c)
		assert.Nil(t, err)
		err = json.NewDecoder(rec.Body).Decode(&r)

		//assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "title error : this field should not empty.", r.Message)
	})

	t.Run("should return amount error : this field should not less than 0. when amount input is minus value", func(t *testing.T) {
		//arrange
		e := echo.New()
		id := "1"
		reqBody := bytes.NewBufferString(`{
			"title": "buy a new phone",
			"amount": -199,
			"note": "buy a new phone",
			"tags": ["gadget", "shopping"]
		}`)
		req := httptest.NewRequest(http.MethodPut, "/expenses/"+id, reqBody)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		var r Err

		//action
		err := UpdateExpensesById(c)
		assert.Nil(t, err)
		err = json.NewDecoder(rec.Body).Decode(&r)

		//assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "amount error : this field should not less than 0.", r.Message)
	})

	t.Run("should return tags error : this field should have at least 1. when tags input is empty", func(t *testing.T) {
		//arrange
		e := echo.New()
		id := "1"
		reqBody := bytes.NewBufferString(`{
			"title": "buy a new phone",
			"amount": 39000,
			"note": "buy a new phone",
			"tags": []
		}`)
		req := httptest.NewRequest(http.MethodPut, "/expenses/"+id, reqBody)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		var r Err

		//action
		err := UpdateExpensesById(c)
		assert.Nil(t, err)
		err = json.NewDecoder(rec.Body).Decode(&r)

		//assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "tags error : this field should have at least 1.", r.Message)
	})
}
