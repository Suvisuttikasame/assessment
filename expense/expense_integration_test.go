//go:build integration

package expense

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
