using WeatherFlow.Thingsboard;

namespace WeatherFlow;

public class Dashboard
{
    private readonly ThingsBoardClient _client;
    private readonly Dictionary<string, string> _weather = new();
    private readonly Dictionary<string, string> _attributes = new();
    private readonly Dictionary<string, (string Value, long Ts)> _waterLevels = new();
    private int _alarmCount;
    private bool _initialized;

    private static readonly (string Key, string Label, string Unit)[] WeatherFields =
    [
        ("airTemperature", "Temperatur", "°C"),
        ("airHumidity", "Luftfeuchte", "%"),
        ("barometricPressure", "Luftdruck", "hPa"),
        ("windSpeed", "Wind", "m/s"),
        ("windDirectionSensor", "Windrichtung", "°"),
        ("rainGauge", "Regen", "mm"),
        ("uvIndex", "UV-Index", ""),
        ("lightIntensity", "Licht", "lux"),
        ("battery", "Batterie", "%")
    ];

    private static readonly string[] WindDirs =
        ["N", "NNO", "NO", "ONO", "O", "OSO", "SO", "SSO", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"];

    public Dashboard(ThingsBoardClient client) => _client = client;

    public void HandleMessage(WsMessage msg)
    {
        switch (msg.CmdId)
        {
            case 1: ProcessWeather(msg); break;
            case 2: ProcessAttributes(msg); break;
            case 10: ProcessAlarms(msg); break;
            case 11: ProcessWaterLevels(msg); break;
        }
    }

    private void ProcessWeather(WsMessage msg)
    {
        var entries = msg.Update ?? msg.Data?.Data;
        if (entries == null) return;

        foreach (var e in entries)
        {
            if (!e.Latest.TryGetValue("TIME_SERIES", out var ts)) continue;
            foreach (var (key, val) in ts)
                _weather[key] = val.Value;
        }
        Redraw();
    }

    private void ProcessAttributes(WsMessage msg)
    {
        var entries = msg.Data?.Data;
        if (entries == null) return;

        foreach (var e in entries)
        {
            if (!e.Latest.TryGetValue("ATTRIBUTE", out var attrs)) continue;
            foreach (var (key, val) in attrs)
                _attributes[key] = val.Value;
        }
        Redraw();
    }

    private void ProcessAlarms(WsMessage msg)
    {
        _alarmCount = msg.Data?.TotalElements ?? 0;
    }

    private void ProcessWaterLevels(WsMessage msg)
    {
        var entries = msg.Update ?? msg.Data?.Data;
        if (entries == null) return;

        foreach (var e in entries)
        {
            var name = _client.GetEntityName(e.EntityId.Id)
                       ?? e.EntityId.Id[..12] + "...";
            if (e.Latest.TryGetValue("TIME_SERIES", out var ts) &&
                ts.TryGetValue("waterLevel", out var wl))
            {
                _waterLevels[name] = (wl.Value, wl.Ts);
            }
        }
        Redraw();
    }

    private void Redraw()
    {
        if (!_initialized)
        {
            Console.Clear();
            _initialized = true;
        }

        Console.SetCursorPosition(0, 0);
        Console.CursorVisible = false;

        var now = DateTime.Now.ToString("HH:mm:ss");
        WriteColor($"  ╔══════════════════════════════════════════════════════════╗", ConsoleColor.DarkCyan);
        WriteColor($"  ║  WEATHERFLOW DASHBOARD                     {now,8}  ║", ConsoleColor.DarkCyan);
        WriteColor($"  ╚══════════════════════════════════════════════════════════╝", ConsoleColor.DarkCyan);
        Console.WriteLine();

        // Weather
        WriteColor("  ☀ WETTERSTATION", ConsoleColor.Cyan);
        if (_attributes.TryGetValue("latitude", out var lat) && _attributes.TryGetValue("longitude", out var lon))
            WriteColor($"    {lat}°N, {lon}°E, {_attributes.GetValueOrDefault("altitude", "?")}m", ConsoleColor.DarkGray);
        WriteSeparator();

        foreach (var (key, label, unit) in WeatherFields)
        {
            if (!_weather.TryGetValue(key, out var val)) continue;
            var display = FormatWeatherValue(key, val);
            Write($"  {label,-16}", ConsoleColor.Gray);
            Write($"{display,10}", ConsoleColor.White);
            WriteColor($" {unit}", ConsoleColor.DarkGray);
        }
        WriteSeparator();

        // Water levels
        Console.WriteLine();
        WriteColor("  ≋ KREBSBACH-PEGEL", ConsoleColor.Blue);
        WriteSeparator();

        Write($"  {"Standort",-36}", ConsoleColor.DarkGray);
        Write($"{"mm",8}", ConsoleColor.DarkGray);
        Write($"{"Zeit",10}", ConsoleColor.DarkGray);
        Console.WriteLine();

        foreach (var (name, (value, ts)) in _waterLevels.OrderBy(x => x.Key))
        {
            var shortName = name.Replace("Krebsbach_", "").Replace("Krebsbach ", "");
            var time = DateTimeOffset.FromUnixTimeMilliseconds(ts).LocalDateTime.ToString("HH:mm:ss");
            var bar = WaterBar(value);

            Write($"  {shortName,-36}", ConsoleColor.Gray);
            Write($"{value,8}", ConsoleColor.White);
            Write($"{time,10}  ", ConsoleColor.DarkGray);
            Console.WriteLine(bar);
        }

        WriteSeparator();

        // Alarms
        if (_alarmCount == 0)
            WriteColor("  ✓ Keine aktiven Alarme", ConsoleColor.Green);
        else
            WriteColor($"  ⚠ {_alarmCount} aktive Alarme!", ConsoleColor.Red);

        Console.WriteLine();
        WriteColor("  Ctrl+C zum Beenden", ConsoleColor.DarkGray);

        // Clear any leftover lines
        var clearLine = new string(' ', Console.WindowWidth > 0 ? Console.WindowWidth : 80);
        for (var i = 0; i < 3; i++)
            Console.WriteLine(clearLine);
    }

    private static string FormatWeatherValue(string key, string value)
    {
        if (key == "barometricPressure" && double.TryParse(value, out var pa))
            return $"{pa / 100:F1}";
        if (key == "windDirectionSensor" && double.TryParse(value, out var deg))
        {
            var idx = (int)((deg + 11.25) / 22.5) % 16;
            return $"{value} ({WindDirs[idx]})";
        }
        return value;
    }

    private static string WaterBar(string value)
    {
        if (!double.TryParse(value, out var mm)) return "";
        const int width = 20;
        const double max = 600;
        var filled = Math.Min((int)(mm / max * width), width);

        var color = mm switch
        {
            > 500 => ConsoleColor.Red,
            > 400 => ConsoleColor.Yellow,
            _ => ConsoleColor.Green
        };

        var result = new string('█', filled) + new string('░', width - filled);
        Console.ForegroundColor = color;
        Console.Write(result[..filled]);
        Console.ForegroundColor = ConsoleColor.DarkGray;
        Console.Write(result[filled..]);
        Console.ResetColor();
        return "";
    }

    private static void WriteSeparator()
    {
        WriteColor($"  {"",58}", ConsoleColor.DarkGray, '─');
    }

    private static void Write(string text, ConsoleColor color)
    {
        Console.ForegroundColor = color;
        Console.Write(text);
        Console.ResetColor();
    }

    private static void WriteColor(string text, ConsoleColor color, char? fill = null)
    {
        Console.ForegroundColor = color;
        if (fill.HasValue)
            Console.WriteLine(new string(fill.Value, text.Length > 2 ? text.Length : 58));
        else
            Console.WriteLine(text);
        Console.ResetColor();
    }
}
