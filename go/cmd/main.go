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
		fmt.Fprintf(os.Stderr, "Usage: weatherflow -public-id <ID> [-pretty] [-output <file>] [-error-log <file>]\n\n")
		fmt.Fprintf(os.Stderr, "Streams real-time weather and water level data from ThingsBoard.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  weatherflow -public-id d58b18a0-1440-11ef-aef4-af283e5094d9\n")
		fmt.Fprintf(os.Stderr, "  weatherflow -public-id d58b18a0-1440-11ef-aef4-af283e5094d9 -output data.log -error-log errors.log\n")
	}

	publicID := flag.String("public-id", "", "ThingsBoard public customer ID (required)")
	pretty := flag.Bool("pretty", false, "formatted terminal output instead of raw JSON")
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

	msgLogger, outputCloser, err := setupOutputLog(*outputFile, *pretty)
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

	c := client.New(token, msgLogger, *pretty)
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

func setupOutputLog(path string, pretty bool) (*log.Logger, *os.File, error) {
	if path == "-" {
		if pretty {
			return log.New(io.Discard, "", 0), nil, nil
		}
		return log.New(os.Stdout, "", log.LstdFlags), nil, nil
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}

	if pretty {
		// Pretty mode: raw JSON only to file, formatted output to terminal separately
		return log.New(f, "", log.LstdFlags), f, nil
	}

	// Normal mode: raw JSON to both stdout and file
	w := io.MultiWriter(os.Stdout, f)
	return log.New(w, "", log.LstdFlags), f, nil
}
