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

## Usage

```bash
cd go
go build -o weatherflow ./cmd/main.go

./weatherflow -public-id <THINGSBOARD_PUBLIC_ID>
```

The current public ID for the BDA/itnovum ThingsBoard instance:

```
d58b18a0-1440-11ef-aef4-af283e5094d9
```

Example:

```bash
./weatherflow -public-id d58b18a0-1440-11ef-aef4-af283e5094d9
```

Stop with `Ctrl+C`.

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
