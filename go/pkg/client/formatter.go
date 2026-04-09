package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	blue    = "\033[34m"
	cyan    = "\033[36m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	magenta = "\033[35m"
	white   = "\033[97m"
)

var weatherLabels = map[string]struct{ Label, Unit string }{
	"airTemperature":      {"Temperatur", "°C"},
	"airHumidity":         {"Luftfeuchte", "%"},
	"barometricPressure":  {"Luftdruck", "Pa"},
	"windSpeed":           {"Wind", "m/s"},
	"windDirectionSensor": {"Windrichtung", "°"},
	"rainGauge":           {"Regen", "mm"},
	"uvIndex":             {"UV-Index", ""},
	"lightIntensity":      {"Licht", "lux"},
	"battery":             {"Batterie", "%"},
}

var weatherOrder = []string{
	"airTemperature", "airHumidity", "barometricPressure",
	"windSpeed", "windDirectionSensor", "rainGauge",
	"uvIndex", "lightIntensity", "battery",
}

var attrLabels = map[string]string{
	"latitude":        "Breitengrad",
	"longitude":       "Längengrad",
	"altitude":        "Höhe (m)",
	"firmwareVersion": "Firmware",
	"hardwareVersion": "Hardware",
	"active":          "Aktiv",
}

func (c *Client) prettyPrint(raw []byte) {
	var msg wsMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}

	switch {
	case msg.CmdID == 1 && msg.Data != nil:
		c.printWeatherInitial(msg.Data.Data)
	case msg.CmdID == 1 && msg.Update != nil:
		var updates []entityData
		if json.Unmarshal(msg.Update, &updates) == nil {
			c.printWeatherUpdate(updates)
		}
	case msg.CmdID == 2 && msg.Data != nil:
		c.printAttributes(msg.Data.Data)
	case msg.CmdID == 11 && msg.Data != nil:
		c.printWaterLevelsInitial(msg.Data.Data)
	case msg.CmdID == 11 && msg.Update != nil:
		var updates []entityData
		if json.Unmarshal(msg.Update, &updates) == nil {
			c.printWaterLevelUpdate(updates)
		}
	case msg.CmdID == 10 && msg.Data != nil:
		c.printAlarms(msg.Data.Data)
	}
}

func (c *Client) printWeatherInitial(entities []entityData) {
	if len(entities) == 0 {
		return
	}

	ts := entities[0].Latest["TIME_SERIES"]
	t := tsTime(ts)

	fmt.Println()
	printHeader("WETTERSTATION", t)
	printSeparator()

	for _, key := range weatherOrder {
		if v, ok := ts[key]; ok {
			info := weatherLabels[key]
			val := formatWeatherValue(key, v.Value)
			fmt.Printf("  %s%-14s%s  %s%s%s %s\n", cyan, info.Label, reset, bold+white, val, reset, dim+info.Unit+reset)
		}
	}
	printSeparator()
}

func (c *Client) printWeatherUpdate(updates []entityData) {
	if len(updates) == 0 {
		return
	}

	ts := updates[0].Latest["TIME_SERIES"]
	t := tsTime(ts)

	fmt.Printf("\n  %s⟳ Wetter-Update%s  %s%s%s\n", green+bold, reset, dim, t, reset)

	changed := []string{}
	for _, key := range weatherOrder {
		if v, ok := ts[key]; ok {
			info := weatherLabels[key]
			val := formatWeatherValue(key, v.Value)
			changed = append(changed, fmt.Sprintf("%s%s%s=%s%s%s%s", cyan, info.Label, reset, white, val, info.Unit, reset))
		}
	}
	fmt.Printf("  %s\n", strings.Join(changed, "  "))
}

func (c *Client) printAttributes(entities []entityData) {
	if len(entities) == 0 {
		return
	}

	attrs := entities[0].Latest["ATTRIBUTE"]
	fmt.Println()
	fmt.Printf("  %s%sGeräteinformationen%s\n", bold, blue, reset)
	printSeparator()

	keys := []string{"latitude", "longitude", "altitude", "firmwareVersion", "hardwareVersion", "active"}
	for _, key := range keys {
		if v, ok := attrs[key]; ok {
			label := attrLabels[key]
			fmt.Printf("  %s%-14s%s  %s%s%s\n", cyan, label, reset, white, v.Value, reset)
		}
	}
	printSeparator()
}

func (c *Client) printWaterLevelsInitial(entities []entityData) {
	if len(entities) == 0 {
		return
	}

	fmt.Println()
	printHeader("KREBSBACH-PEGEL", "")
	fmt.Printf("  %s%-40s  %6s  %8s%s\n", dim, "Standort", "mm", "Zeit", reset)
	printSeparator()

	for _, e := range entities {
		name := c.getEntityName(e)
		wl := e.Latest["TIME_SERIES"]["waterLevel"]
		t := time.UnixMilli(wl.Ts).Format("15:04:05")
		bar := waterBar(wl.Value)
		fmt.Printf("  %-40s  %s%6s%s  %s%s%s  %s\n", name, bold, wl.Value, reset, dim, t, reset, bar)
	}
	printSeparator()
	fmt.Printf("  %s%d Sensoren aktiv%s\n", dim, len(entities), reset)
}

func (c *Client) printWaterLevelUpdate(updates []entityData) {
	for _, e := range updates {
		name := c.getEntityName(e)
		if wl, ok := e.Latest["TIME_SERIES"]["waterLevel"]; ok {
			t := time.UnixMilli(wl.Ts).Format("15:04:05")
			bar := waterBar(wl.Value)
			fmt.Printf("  %s⟳ Pegel%s  %s%s%s  %-34s  %s%smm%s  %s\n",
				blue+bold, reset, dim, t, reset, name, white+bold, wl.Value, reset, bar)
		}
	}
}

func (c *Client) printAlarms(entities []entityData) {
	if len(entities) == 0 {
		fmt.Printf("  %s✓ Keine aktiven Alarme%s\n", green, reset)
		return
	}
	fmt.Printf("  %s⚠ %d aktive Alarme%s\n", yellow+bold, len(entities), reset)
}

func (c *Client) getEntityName(e entityData) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if name, ok := c.entityNames[e.EntityID.ID]; ok {
		return strings.Replace(name, "Krebsbach_", "", 1)
	}
	return e.EntityID.ID[:12] + "..."
}

func printHeader(title string, ts string) {
	if ts != "" {
		fmt.Printf("  %s%s%s%s  %s%s%s\n", bold, blue, title, reset, dim, ts, reset)
	} else {
		fmt.Printf("  %s%s%s%s\n", bold, blue, title, reset)
	}
}

func printSeparator() {
	fmt.Printf("  %s%s%s\n", dim, strings.Repeat("─", 58), reset)
}

func formatWeatherValue(key, value string) string {
	if key == "barometricPressure" {
		// Pa -> hPa
		var pa float64
		if _, err := fmt.Sscanf(value, "%f", &pa); err == nil {
			return fmt.Sprintf("%.1f", pa/100)
		}
	}
	if key == "windDirectionSensor" {
		var deg float64
		if _, err := fmt.Sscanf(value, "%f", &deg); err == nil {
			dirs := []string{"N", "NNO", "NO", "ONO", "O", "OSO", "SO", "SSO",
				"S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
			idx := int((deg+11.25)/22.5) % 16
			return fmt.Sprintf("%s (%s)", value, dirs[idx])
		}
	}
	return value
}

func waterBar(value string) string {
	var mm float64
	if _, err := fmt.Sscanf(value, "%f", &mm); err != nil {
		return ""
	}

	maxMM := 600.0
	width := 20
	filled := int(mm / maxMM * float64(width))
	if filled > width {
		filled = width
	}

	color := green
	if mm > 400 {
		color = yellow
	}
	if mm > 500 {
		color = "\033[31m" // red
	}

	return fmt.Sprintf("%s%s%s%s", color, strings.Repeat("█", filled), dim+strings.Repeat("░", width-filled), reset)
}

func tsTime(ts map[string]tsValue) string {
	for _, v := range ts {
		return time.UnixMilli(v.Ts).Format("15:04:05")
	}
	return ""
}
