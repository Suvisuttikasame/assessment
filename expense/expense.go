package expense

import (
	"fmt"
)

type Expense struct {
	Id     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float32  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

func (e *Expense) validation() error {
	if e.Title == "" {
		return fmt.Errorf("title error : this field should not empty.")
	}
	if e.Amount < 0 {
		return fmt.Errorf("amount error : this field should not less than 0.")
	}
	if len(e.Tags) == 0 {
		return fmt.Errorf("tags error : this field should have at least 1.")
	}
	return nil
}
