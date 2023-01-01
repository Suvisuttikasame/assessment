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

func TestCreateExpensesValidation(t *testing.T) {
	t.Run("should return title error : this field should not empty. when title input is empty", func(t *testing.T) {
		//arrange
		reqBody := bytes.NewBufferString(`{
			"title": "",
			"amount": 39000,
			"note": "buy a new phone",
			"tags": ["gadget", "shopping"]
		}`)
		var e Err

		//action
		res := request(http.MethodPost, "http://localhost:2565/expenses", reqBody)
		err := res.Decode(&e)

		//assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "title error : this field should not empty.", e.Message)
	})

	t.Run("should return amount error : this field should not less than 0. when amount input is minus value", func(t *testing.T) {
		//arrange
		reqBody := bytes.NewBufferString(`{
			"title": "buy a new phone",
			"amount": -199,
			"note": "buy a new phone",
			"tags": ["gadget", "shopping"]
		}`)
		var e Err

		//action
		res := request(http.MethodPost, "http://localhost:2565/expenses", reqBody)
		err := res.Decode(&e)

		//assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "amount error : this field should not less than 0.", e.Message)
	})

	t.Run("should return tags error : this field should have at least 1. when tags input is empty", func(t *testing.T) {
		//arrange
		reqBody := bytes.NewBufferString(`{
			"title": "buy a new phone",
			"amount": 39000,
			"note": "buy a new phone",
			"tags": []
		}`)
		var e Err

		//action
		res := request(http.MethodPost, "http://localhost:2565/expenses", reqBody)
		err := res.Decode(&e)

		//assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "tags error : this field should have at least 1.", e.Message)
	})

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
