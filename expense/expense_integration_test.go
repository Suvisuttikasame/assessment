//go:build integration

package expense

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpenses(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		InitDb("postgres://pcfxcojb:klhxSm6mVTqPkH1OZB9GWSASMEkMguZG@tiny.db.elephantsql.com/pcfxcojb")

		e.POST("/expenses", CreateExpenses)
		e.Start(":2565")
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", "localhost:2565", 30*time.Second)
		if err != nil {
			// log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	//arrange
	reqBody := bytes.NewBufferString(`{
		"title": "buy a new phone",
		"amount": 39000,
		"note": "buy a new phone",
		"tags": ["gadget", "shopping"]
	}`)
	var ex Expense

	//act
	res := request(http.MethodPost, "http://localhost:2565/expenses", reqBody)
	err := res.Decode(&ex)

	//assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, ex.Id)
	assert.Equal(t, "buy a new phone", ex.Title)
	assert.Equal(t, float32(39000), ex.Amount)
	assert.Equal(t, "buy a new phone", ex.Note)
	assert.Equal(t, []string{"gadget", "shopping"}, ex.Tags)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestGetExpensesById(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		InitDb("postgres://pcfxcojb:klhxSm6mVTqPkH1OZB9GWSASMEkMguZG@tiny.db.elephantsql.com/pcfxcojb")

		e.GET("/expenses/:id", GetExpensesById)
		e.POST("/expenses", CreateExpenses)
		e.Start(":2565")
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", "localhost:2565", 30*time.Second)
		if err != nil {
			// log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	t.Run("should return the same as input when input is valid", func(t *testing.T) {
		//arrange
		m := SeedExpense(t, `{
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		}`)
		id := m.Id
		t.Log(id)
		ex := Expense{}

		//action
		res := request(http.MethodGet, "http://localhost:2565/expenses/"+strconv.Itoa(id), nil)
		t.Log(res)
		err := res.Decode(&ex)

		//assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, id, ex.Id)
		assert.Equal(t, "apple smoothie", ex.Title)
		assert.Equal(t, float32(89), ex.Amount)
		assert.Equal(t, "no discount", ex.Note)
		assert.Equal(t, []string{"beverage"}, ex.Tags)

	})

	t.Run("should return expense's not found when input param is 999999999", func(t *testing.T) {
		//arrange
		id := 999999999
		ex := Err{}

		//action
		res := request(http.MethodGet, "http://localhost:2565/expenses/"+strconv.Itoa(id), nil)
		err := res.Decode(&ex)

		//assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, "expense's not found", ex.Message)

	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestUpdateExpensesById(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		InitDb("postgres://pcfxcojb:klhxSm6mVTqPkH1OZB9GWSASMEkMguZG@tiny.db.elephantsql.com/pcfxcojb")

		e.PUT("/expenses/:id", UpdateExpensesById)
		e.POST("/expenses", CreateExpenses)
		e.Start(":2565")
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", "localhost:2565", 30*time.Second)
		if err != nil {
			// log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	m := SeedExpense(t, `{
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`)
	id := m.Id
	reqBody := bytes.NewBufferString(`{
		"title": "buy a new phone",
		"amount": 39000,
		"note": "buy a new phone",
		"tags": ["gadget", "shopping"]
	}`)
	ex := Expense{}

	res := request(http.MethodPut, "http://localhost:2565/expenses/"+strconv.Itoa(id), reqBody)
	err := res.Decode(&ex)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, ex.Id, id)
	assert.Equal(t, "buy a new phone", ex.Title)
	assert.Equal(t, float32(39000), ex.Amount)
	assert.Equal(t, "buy a new phone", ex.Note)
	assert.Equal(t, []string{"gadget", "shopping"}, ex.Tags)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestGetExpenses(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		InitDb("postgres://pcfxcojb:klhxSm6mVTqPkH1OZB9GWSASMEkMguZG@tiny.db.elephantsql.com/pcfxcojb")

		e.GET("/expenses", GetExpenses)
		e.POST("/expenses", CreateExpenses)
		e.Start(":2565")
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", "localhost:2565", 30*time.Second)
		if err != nil {
			// log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	_ = SeedExpense(t, `{
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`)
	exs := []Expense{}

	res := request(http.MethodGet, "http://localhost:2565/expenses", nil)
	err := res.Decode(&exs)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEqual(t, 0, len(exs))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}
