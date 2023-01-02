//go:build integration

package expense

import (
	"bytes"
	"net/http"
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
