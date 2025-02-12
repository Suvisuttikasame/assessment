package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Suvisuttikasame/assessment/customMiddleware"
	"github.com/Suvisuttikasame/assessment/expense"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("Initiating database ...")
	//init db connection & create table
	Port := os.Getenv("PORT")
	Url := os.Getenv("DATABASE_URL")
	expense.InitDb(Url)
	fmt.Println("Successfully initiate database")

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.BasicAuth(customMiddleware.Authentication))

	e.POST("/expenses", expense.CreateExpenses)
	e.GET("/expenses", expense.GetExpenses)
	e.GET("/expenses/:id", expense.GetExpensesById)
	e.PUT("/expenses/:id", expense.UpdateExpensesById)

	// fmt.Println("Please use server.go for main file")
	// fmt.Println("start at port:", os.Getenv("PORT"))
	fmt.Println("server is running on port:", Port)
	go func() {
		fmt.Println(e.Start(":" + Port))
	}()

	//create buffer to listen to os signal
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	//wait til receive signal
	<-gracefulStop
	//define context
	//clear all req & shut down server with in 10 sec
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("server is shutting down...")
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("server is shut down gracefully")

}
