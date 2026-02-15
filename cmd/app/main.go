package main

import (
	"log"
)

func main() {
	app, cleanup, err := InitializeApp()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
