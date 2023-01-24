package main

import (
	"firebird/internal/app"
	"os"
	"log"
)

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
    if err != nil {
        log.Fatalln(err)
    }
    log.SetOutput(logFile)

	app := app.NewApp()
	err = app.Config.ConfigLoad(".env")
	if err != nil {
		panic(err)
	}
	err = app.Queries.QueriesLoad("queries.yml")
	if err != nil {
		panic(err)
	}

	err = app.Setup()
	if err != nil {
		panic(err)
	}

	err = app.Run()
	if err != nil {
		panic(err)
	}

}
