package main

import (
	"github.com/joho/godotenv"
	"github.com/khrees2412/cadana/taskone"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}

func main() {
	taskone.StartServer()
	//tasktwo.Start()
}
