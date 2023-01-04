//go:build unit

package expense

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
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

func TestCreateExpenseUnit(t *testing.T) {
	//arrange
	e := echo.New()
	reqBody := bytes.NewBufferString(`{
		"title": "buy a new phone",
		"amount": 39000,
		"note": "buy a new phone",
		"tags": ["gadget", "shopping"]
	}`)
	req := httptest.NewRequest(http.MethodPost, "/", reqBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var mock sqlmock.Sqlmock
	var err error
	md := Expense{
		Title:  "buy a new phone",
		Amount: 39000,
		Note:   "buy a new phone",
		Tags:   []string{"gadget", "shopping"},
	}
	rt := Expense{}

	Db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal("unable to create mock db", err)
	}
	defer Db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id`)).
		WithArgs(md.Title, md.Amount, md.Note, pq.Array(&md.Tags)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	//action
	err = CreateExpenses(c)
	assert.Nil(t, err)
	err = json.NewDecoder(rec.Body).Decode(&rt)

	//assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, 1, rt.Id)
	assert.Equal(t, "buy a new phone", rt.Title)
	assert.Equal(t, float32(39000), rt.Amount)
	assert.Equal(t, "buy a new phone", rt.Note)
	assert.Equal(t, []string{"gadget", "shopping"}, rt.Tags)

}

func TestGetExpenseByIdUnit(t *testing.T) {
	//arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	var mock sqlmock.Sqlmock
	var err error
	md := Expense{
		Title:  "buy a new phone",
		Amount: 39000,
		Note:   "buy a new phone",
		Tags:   []string{"gadget", "shopping"},
	}
	rt := Expense{}

	Db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal("unable to create mock db", err)
	}
	defer Db.Close()

	prep := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT * FROM expenses WHERE id = $1`))

	prep.ExpectQuery().
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(1, md.Title, md.Amount, md.Note, pq.Array(&md.Tags)))

	//action
	err = GetExpensesById(c)
	assert.Nil(t, err)
	err = json.NewDecoder(rec.Body).Decode(&rt)

	//assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, rt.Id)
	assert.Equal(t, "buy a new phone", rt.Title)
	assert.Equal(t, float32(39000), rt.Amount)
	assert.Equal(t, "buy a new phone", rt.Note)
	assert.Equal(t, []string{"gadget", "shopping"}, rt.Tags)
}

func TestUpdateExpensesByIdUnit(t *testing.T) {
	//arrange
	e := echo.New()
	reqBody := bytes.NewBufferString(`{
		"title": "buy a new phone",
		"amount": 39000,
		"note": "buy a new phone",
		"tags": ["gadget", "shopping"]
	}`)
	req := httptest.NewRequest(http.MethodPut, "/", reqBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	var mock sqlmock.Sqlmock
	var err error
	md := Expense{
		Title:  "buy a new phone",
		Amount: 39000,
		Note:   "buy a new phone",
		Tags:   []string{"gadget", "shopping"},
	}
	rt := Expense{}

	Db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal("unable to create mock db", err)
	}
	defer Db.Close()

	prep := mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING *`))
	prep.ExpectQuery().
		WithArgs(md.Title, md.Amount, md.Note, pq.Array(&md.Tags), "1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(1, md.Title, md.Amount, md.Note, pq.Array(&md.Tags)))

	//action
	err = UpdateExpensesById(c)
	assert.Nil(t, err)
	err = json.NewDecoder(rec.Body).Decode(&rt)

	//assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, rt.Id)
	assert.Equal(t, "buy a new phone", rt.Title)
	assert.Equal(t, float32(39000), rt.Amount)
	assert.Equal(t, "buy a new phone", rt.Note)
	assert.Equal(t, []string{"gadget", "shopping"}, rt.Tags)
}

func TestGetExpensesUnit(t *testing.T) {
	//arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var mock sqlmock.Sqlmock
	var err error
	md := Expense{
		Title:  "buy a new phone",
		Amount: 39000,
		Note:   "buy a new phone",
		Tags:   []string{"gadget", "shopping"},
	}
	rt := []Expense{}

	Db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatal("unable to create mock db", err)
	}
	defer Db.Close()

	prep := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT * FROM expenses`))

	prep.ExpectQuery().
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(1, md.Title, md.Amount, md.Note, pq.Array(&md.Tags)))

	//action
	err = GetExpenses(c)
	assert.Nil(t, err)
	err = json.NewDecoder(rec.Body).Decode(&rt)

	//assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEqual(t, 0, len(rt))
}
