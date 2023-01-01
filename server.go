package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Suvisuttikasame/assessment/expense"
	"github.com/labstack/echo/v4"
)

func main() {
	fmt.Println("Initiating database ...")
	//init db connection & create table
	expense.InitDb()
	fmt.Println("Successfully initiate database")

	e := echo.New()

	e.POST("/expenses", expense.CreateExpenses)

	// fmt.Println("Please use server.go for main file")
	// fmt.Println("start at port:", os.Getenv("PORT"))
	fmt.Println("server is running on port:", os.Getenv("PORT"))
	go func() {
		fmt.Println(e.Start(":" + os.Getenv("PORT")))
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
