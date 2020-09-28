package main

import (
	"github.com/go-http-utils/logger"
	"log"
	"net/http"
	"os"
)

func main() {

	app := newServer()
	err := http.ListenAndServe(":8080", logger.Handler(app, os.Stdout, logger.DevLoggerType))
	if err != nil {
		log.Panic(err)
	}
}
