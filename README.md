# WeatherFlow

WebSocket client for streaming real-time weather and water level data from a [ThingsBoard](https://thingsboard.io/) IoT platform.

## Data Sources

**Weather station** (Device `2CF7F1C044300113`, located at 50.179°N, 8.925°E):
- Air temperature, humidity, barometric pressure
- Wind speed and direction
- Rain gauge, UV index, light intensity
- Battery status

**Water level sensors** (10x Dragino LDDS along the Krebsbach):
- Real-time water level readings
- Active alarms

## Build

```bash
cd go

# Linux (amd64) — default on most PCs/servers
GOOS=linux GOARCH=amd64 go build -o weatherflow-linux-amd64 ./cmd/main.go

# Raspberry Pi 3/4/5 (64-bit OS)
GOOS=linux GOARCH=arm64 go build -o weatherflow-linux-arm64 ./cmd/main.go

# Raspberry Pi Zero / Pi 1 / Pi 2 (32-bit OS)
GOOS=linux GOARCH=arm GOARM=6 go build -o weatherflow-linux-arm ./cmd/main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o weatherflow-darwin-arm64 ./cmd/main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o weatherflow-darwin-amd64 ./cmd/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o weatherflow.exe ./cmd/main.go

# Current platform (quick build)
go build -o weatherflow ./cmd/main.go
```

Copy the matching binary to the target machine — no Go installation needed there.

## Usage

```bash
./weatherflow -public-id <THINGSBOARD_PUBLIC_ID>
```

The current public ID for the BDA/itnovum ThingsBoard instance:

```
d58b18a0-1440-11ef-aef4-af283e5094d9
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-public-id` | (required) | ThingsBoard public customer ID |
| `-pretty` | off | Formatted, colored terminal output instead of raw JSON |
| `-output` | `weatherflow.log` | Output file for messages (use `-` for stdout only) |
| `-error-log` | `weatherflow-errors.log` | Error log file |
| `-help` | | Show usage |

### Examples

```bash
# Basic: output to terminal + weatherflow.log, errors to weatherflow-errors.log
./weatherflow -public-id d58b18a0-1440-11ef-aef4-af283e5094d9

# Custom file paths
./weatherflow -public-id d58b18a0-1440-11ef-aef4-af283e5094d9 -output data.log -error-log errors.log

# Pretty terminal output (colored, formatted)
./weatherflow -public-id d58b18a0-1440-11ef-aef4-af283e5094d9 -pretty

# Terminal only (no output file)
./weatherflow -public-id d58b18a0-1440-11ef-aef4-af283e5094d9 -output -
```

Stop with `Ctrl+C`.

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | OK |
| 1 | Invalid usage (missing flags) |
| 2 | Authentication failed |
| 3 | WebSocket connection failed |
| 4 | Log file setup failed |

## Project Structure

```
go/
  cmd/main.go              Entry point
  pkg/auth/auth.go         ThingsBoard public auth (HTTP POST)
  pkg/client/client.go     WebSocket client, message reader
  pkg/payload/             WebSocket subscription payload structs
```

## Requirements

- Go 1.23+
- Network access to `thingsboard.bda-itnovum.com`
