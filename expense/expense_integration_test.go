//go:build integration

package expense

import (
	"bytes"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateExpenses(t *testing.T) {
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
}

func TestGetExpensesById(t *testing.T) {

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

}

func TestUpdateExpensesById(t *testing.T) {
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
}

func TestGetExpenses(t *testing.T) {
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
}
