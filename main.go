package main

import (
	"os"

	"./src/system/app"

	"flag"

	"github.com/joho/godotenv"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8000", "Assigning the port that the server should listen on.")

	flag.Parse()

	if err := godotenv.Load("config.ini"); err != nil {
		panic(err)
	}

	envPort := os.Getenv("PORT")
	if len(envPort) > 0 {
		port = envPort
	}
}

func main() {
	s := app.NewServer()

	s.Init(port)
	s.Start()
}
