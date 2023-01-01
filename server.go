package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Suvisuttikasame/assessment/expense"
	"github.com/labstack/echo/v4"
)

func main() {
	fmt.Println("Initiating database ...")
	expense.InitDb()
	fmt.Println("Successfully initiate database")

	e := echo.New()

	e.POST("/expenses", expense.CreateExpenses)
	// fmt.Println("Please use server.go for main file")
	// fmt.Println("start at port:", os.Getenv("PORT"))
	fmt.Println("server is running on port:", os.Getenv("PORT"))
	log.Fatal(e.Start(":" + os.Getenv("PORT")))
}
