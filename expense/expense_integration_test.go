//go:build integration

package expense

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	host     = "pq_test"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbName   = "go-postest-db"
)

func TestCreateExpenses(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		InitTestDb(t)

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
		InitTestDb(t)

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
		ex := Expense{}

		//action
		res := request(http.MethodGet, "http://localhost:2565/expenses/"+strconv.Itoa(id), nil)
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
		InitTestDb(t)

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
		InitTestDb(t)

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

func InitTestDb(t *testing.T) {
	//connect to postgres db
	var err error
	connStr := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable", host, user, password, port, dbName)
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	createTb := `CREATE TABLE IF NOT EXISTS expenses(
		id SERIAL PRIMARY KEY,
		title TEXT,
		AMOUNT FLOAT,
		NOTE TEXT,
		TAGS TEXT[])`

	_, err = Db.Exec(createTb)

	if err != nil {
		t.Fatal("can not create table expense", err)
	}
	t.Log("successfully connect to db")
}

func SeedExpense(t *testing.T, body string) Expense {
	reqBody := bytes.NewBufferString(body)

	ex := Expense{}

	res := request(http.MethodPost, "http://localhost:2565/expenses", reqBody)
	err := res.Decode(&ex)

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
