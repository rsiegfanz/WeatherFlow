package main

import (
	"log"

	"github.com/rsiegfanz/WeatherFlow/pkg/auth"
	"github.com/rsiegfanz/WeatherFlow/pkg/client"
)

func main() {
	token, err := auth.Authenticate()
	if err != nil {
		log.Fatalf("Authentication error: %v", err)
	}

	c := client.New(token)
	defer c.Close()

	if err := c.Connect(); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
}
