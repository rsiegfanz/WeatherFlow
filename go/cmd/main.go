package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/rsiegfanz/WeatherFlow/pkg/auth"
	"github.com/rsiegfanz/WeatherFlow/pkg/client"
)

const (
	exitOK        = 0
	exitUsage     = 1
	exitAuth      = 2
	exitConnect   = 3
	exitLogSetup  = 4
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: weatherflow -public-id <ID> [-output <file>] [-error-log <file>]\n\n")
		fmt.Fprintf(os.Stderr, "Streams real-time weather and water level data from ThingsBoard.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  weatherflow -public-id d58b18a0-1440-11ef-aef4-af283e5094d9\n")
		fmt.Fprintf(os.Stderr, "  weatherflow -public-id d58b18a0-1440-11ef-aef4-af283e5094d9 -output data.log -error-log errors.log\n")
	}

	publicID := flag.String("public-id", "", "ThingsBoard public customer ID (required)")
	outputFile := flag.String("output", "weatherflow.log", "output file for messages (use - for stdout only)")
	errorLogFile := flag.String("error-log", "weatherflow-errors.log", "error log file")
	flag.Parse()

	if *publicID == "" {
		flag.Usage()
		os.Exit(exitUsage)
	}

	errorLog, err := setupErrorLog(*errorLogFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open error log file: %v\n", err)
		os.Exit(exitLogSetup)
	}
	defer errorLog.Close()

	msgLogger, outputCloser, err := setupOutputLog(*outputFile)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	if outputCloser != nil {
		defer outputCloser.Close()
	}

	token, err := auth.Authenticate(*publicID)
	if err != nil {
		log.Fatalf("Authentication error (exit %d): %v", exitAuth, err)
	}

	c := client.New(token, msgLogger)
	defer c.Close()

	if err := c.Connect(); err != nil {
		log.Fatalf("Connection failed (exit %d): %v", exitConnect, err)
	}
}

func setupErrorLog(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	// Errors go to both stderr and the error log file
	log.SetOutput(io.MultiWriter(os.Stderr, f))
	log.SetFlags(log.LstdFlags)
	return f, nil
}

func setupOutputLog(path string) (*log.Logger, *os.File, error) {
	if path == "-" {
		return log.New(os.Stdout, "", log.LstdFlags), nil, nil
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}

	// Messages go to both stdout and the output file
	w := io.MultiWriter(os.Stdout, f)
	return log.New(w, "", log.LstdFlags), f, nil
}
