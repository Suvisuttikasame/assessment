package expense

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func SeedExpense(t *testing.T, body string) Expense {
	reqBody := bytes.NewBufferString(body)
	fmt.Println(reqBody)

	ex := Expense{}

	res := request(http.MethodPost, "http://localhost:2565/expenses", reqBody)
	err := res.Decode(&ex)
	fmt.Println(ex)

	if err != nil {
		t.Fatal("unaa\ble to seed demo data.", err.Error())
	}
	return ex
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func request(method, uri string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, uri, body)
	req.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4=")
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
