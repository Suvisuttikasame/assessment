package expense

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitDb(url string) {
	var err error
	Db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Db connection error", err)
	}

	createTb := `CREATE TABLE IF NOT EXISTS expenses(
					id SERIAL PRIMARY KEY,
					title TEXT,
					AMOUNT FLOAT,
					NOTE TEXT,
					TAGS TEXT[])`

	_, err = Db.Exec(createTb)

	if err != nil {
		log.Fatal("can not create table expense", err)
	}
}
