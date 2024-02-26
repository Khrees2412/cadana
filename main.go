package main

import (
	"github.com/joho/godotenv"
	"github.com/khrees2412/cadana/api"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}

func main() {
	api.Start()
}
