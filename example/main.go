package main

import (
	"fmt"
	"log"
	"os"

	"github.com/geetarista/gosumi"
)

func main() {
	email := os.Getenv("APPLE_EMAIL")
	password := os.Getenv("APPLE_PASSWORD")

	icloud, err := gosumi.New(email, password)
	if err != nil {
		log.Printf("Unable to initialize iCloud client: %s", err)
		return
	}

	devices := icloud.Devices
	for _, device := range devices {
		fmt.Printf("%+v\n", device)
	}
}
