package main

import (
	"log"

	"idc-finance/internal/platform/app"
)

func main() {
	application := app.New()
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
