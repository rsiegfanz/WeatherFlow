package main

import (
	"flag"
	"log"

	"github.com/rsiegfanz/WeatherFlow/pkg/auth"
	"github.com/rsiegfanz/WeatherFlow/pkg/client"
)

func main() {
	publicID := flag.String("public-id", "", "ThingsBoard public customer ID (required)")
	flag.Parse()

	if *publicID == "" {
		log.Fatal("Missing required flag: -public-id")
	}

	token, err := auth.Authenticate(*publicID)
	if err != nil {
		log.Fatalf("Authentication error: %v", err)
	}

	c := client.New(token)
	defer c.Close()

	if err := c.Connect(); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
}
