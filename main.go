package main

import (
	"log"

	"github.com/izaakdale/service-ids/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
