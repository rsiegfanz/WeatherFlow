using System.Collections.ObjectModel;
using System.ComponentModel;
using System.Runtime.CompilerServices;
using System.Windows.Media;
using WeatherFlow.Thingsboard;

namespace WeatherFlow.Gui.ViewModels;

public class DashboardViewModel : INotifyPropertyChanged
{
    public event PropertyChangedEventHandler? PropertyChanged;

    // Weather
    private string _temperature = "--";
    private string _humidity = "--";
    private string _pressure = "--";
    private string _windSpeed = "--";
    private string _windDirection = "--";
    private string _windDirLabel = "";
    private string _rain = "--";
    private string _uvIndex = "--";
    private string _light = "--";
    private string _battery = "--";
    private string _lastUpdate = "--:--:--";

    // Location
    private string _location = "";
    private string _firmware = "";

    // Status
    private string _status = "Verbinde...";
    private Brush _statusColor = Brushes.Gray;
    private int _alarmCount;

    public string Temperature { get => _temperature; set => Set(ref _temperature, value); }
    public string Humidity { get => _humidity; set => Set(ref _humidity, value); }
    public string Pressure { get => _pressure; set => Set(ref _pressure, value); }
    public string WindSpeed { get => _windSpeed; set => Set(ref _windSpeed, value); }
    public string WindDirection { get => _windDirection; set => Set(ref _windDirection, value); }
    public string WindDirLabel { get => _windDirLabel; set => Set(ref _windDirLabel, value); }
    public string Rain { get => _rain; set => Set(ref _rain, value); }
    public string UvIndex { get => _uvIndex; set => Set(ref _uvIndex, value); }
    public string Light { get => _light; set => Set(ref _light, value); }
    public string Battery { get => _battery; set => Set(ref _battery, value); }
    public string LastUpdate { get => _lastUpdate; set => Set(ref _lastUpdate, value); }
    public string Location { get => _location; set => Set(ref _location, value); }
    public string Firmware { get => _firmware; set => Set(ref _firmware, value); }
    public string Status { get => _status; set => Set(ref _status, value); }
    public Brush StatusColor { get => _statusColor; set => Set(ref _statusColor, value); }
    public int AlarmCount { get => _alarmCount; set { Set(ref _alarmCount, value); OnPropertyChanged(nameof(AlarmText)); OnPropertyChanged(nameof(AlarmBrush)); } }
    public string AlarmText => _alarmCount == 0 ? "Keine Alarme" : $"{_alarmCount} aktive Alarme!";
    public Brush AlarmBrush => _alarmCount == 0 ? new SolidColorBrush(Color.FromRgb(166, 227, 161)) : new SolidColorBrush(Color.FromRgb(243, 139, 168));

    public ObservableCollection<WaterLevelItem> WaterLevels { get; } = new();
    public HistoryStore History { get; } = new();

    private static readonly string[] WindDirs =
        ["N", "NNO", "NO", "ONO", "O", "OSO", "SO", "SSO", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"];

    public static readonly Dictionary<string, (string Label, string Unit)> SensorInfo = new()
    {
        ["airTemperature"] = ("Temperatur", "°C"),
        ["airHumidity"] = ("Luftfeuchte", "%"),
        ["barometricPressure"] = ("Luftdruck", "hPa"),
        ["windSpeed"] = ("Windgeschwindigkeit", "m/s"),
        ["windDirectionSensor"] = ("Windrichtung", "°"),
        ["rainGauge"] = ("Regenmenge", "mm"),
        ["uvIndex"] = ("UV-Index", ""),
        ["lightIntensity"] = ("Lichtintensität", "lux"),
        ["battery"] = ("Batterie", "%"),
    };

    public void ProcessMessage(WsMessage msg, ThingsBoardClient client)
    {
        switch (msg.CmdId)
        {
            case 1: UpdateWeather(msg); break;
            case 2: UpdateAttributes(msg); break;
            case 10: UpdateAlarms(msg); break;
            case 11: UpdateWaterLevels(msg, client); break;
        }
    }

    private void UpdateWeather(WsMessage msg)
    {
        var entries = msg.Update ?? msg.Data?.Data;
        if (entries == null) return;

        foreach (var e in entries)
        {
            if (!e.Latest.TryGetValue("TIME_SERIES", out var ts)) continue;

            if (ts.TryGetValue("airTemperature", out var t)) Temperature = t.Value;
            if (ts.TryGetValue("airHumidity", out var h)) Humidity = h.Value;
            if (ts.TryGetValue("barometricPressure", out var p) && double.TryParse(p.Value, out var pa))
                Pressure = $"{pa / 100:F1}";
            if (ts.TryGetValue("windSpeed", out var ws)) WindSpeed = ws.Value;
            if (ts.TryGetValue("windDirectionSensor", out var wd))
            {
                WindDirection = wd.Value;
                if (double.TryParse(wd.Value, out var deg))
                    WindDirLabel = WindDirs[(int)((deg + 11.25) / 22.5) % 16];
            }
            if (ts.TryGetValue("rainGauge", out var r)) Rain = r.Value;
            if (ts.TryGetValue("uvIndex", out var uv)) UvIndex = uv.Value;
            if (ts.TryGetValue("lightIntensity", out var li)) Light = FormatNumber(li.Value);
            if (ts.TryGetValue("battery", out var b)) Battery = b.Value;

            // Record history for all sensors
            foreach (var (key, val) in ts)
            {
                if (key == "barometricPressure" && double.TryParse(val.Value, out var paVal))
                    History.Record(key, val.Ts, $"{paVal / 100:F1}");
                else
                    History.Record(key, val.Ts, val.Value);
            }

            LastUpdate = DateTimeOffset.FromUnixTimeMilliseconds(t?.Ts ?? ts.Values.First().Ts)
                .LocalDateTime.ToString("HH:mm:ss");
        }

        Status = "Verbunden";
        StatusColor = new SolidColorBrush(Color.FromRgb(166, 227, 161));
    }

    private void UpdateAttributes(WsMessage msg)
    {
        var entries = msg.Data?.Data;
        if (entries == null) return;

        foreach (var e in entries)
        {
            if (!e.Latest.TryGetValue("ATTRIBUTE", out var attrs)) continue;
            var lat = attrs.GetValueOrDefault("latitude")?.Value ?? "";
            var lon = attrs.GetValueOrDefault("longitude")?.Value ?? "";
            var alt = attrs.GetValueOrDefault("altitude")?.Value ?? "";
            if (!string.IsNullOrEmpty(lat))
                Location = $"{lat}°N, {lon}°E, {alt}m";
            var fw = attrs.GetValueOrDefault("firmwareVersion")?.Value ?? "";
            var hw = attrs.GetValueOrDefault("hardwareVersion")?.Value ?? "";
            if (!string.IsNullOrEmpty(fw))
                Firmware = $"FW {fw} / HW {hw}";
        }
    }

    private void UpdateAlarms(WsMessage msg)
    {
        AlarmCount = msg.Data?.TotalElements ?? 0;
    }

    private void UpdateWaterLevels(WsMessage msg, ThingsBoardClient client)
    {
        var entries = msg.Update ?? msg.Data?.Data;
        if (entries == null) return;

        foreach (var e in entries)
        {
            if (!e.Latest.TryGetValue("TIME_SERIES", out var ts)) continue;
            if (!ts.TryGetValue("waterLevel", out var wl)) continue;

            var name = client.GetEntityName(e.EntityId.Id) ?? e.EntityId.Id[..12];
            var shortName = name.Replace("Krebsbach_", "").Replace("Krebsbach ", "");

            if (!double.TryParse(wl.Value, out var mm)) continue;
            var time = DateTimeOffset.FromUnixTimeMilliseconds(wl.Ts).LocalDateTime.ToString("HH:mm:ss");

            History.Record($"water:{shortName}", wl.Ts, wl.Value);

            var existing = WaterLevels.FirstOrDefault(w => w.Name == shortName);
            if (existing != null)
            {
                existing.Value = mm;
                existing.Time = time;
            }
            else
            {
                WaterLevels.Add(new WaterLevelItem { Name = shortName, Value = mm, Time = time });
            }
        }
    }

    private static string FormatNumber(string value)
    {
        if (double.TryParse(value, out var n) && n >= 1000)
            return $"{n:N0}";
        return value;
    }

    private void Set<T>(ref T field, T value, [CallerMemberName] string? name = null)
    {
        if (EqualityComparer<T>.Default.Equals(field, value)) return;
        field = value;
        OnPropertyChanged(name);
    }

    private void OnPropertyChanged([CallerMemberName] string? name = null) =>
        PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(name));
}

public class WaterLevelItem : INotifyPropertyChanged
{
    public event PropertyChangedEventHandler? PropertyChanged;

    private string _name = "";
    private double _value;
    private string _time = "";

    public string Name { get => _name; set { _name = value; Notify(); } }
    public double Value { get => _value; set { _value = value; Notify(); Notify(nameof(BarWidth)); Notify(nameof(BarColor)); Notify(nameof(ValueText)); } }
    public string Time { get => _time; set { _time = value; Notify(); } }

    public double BarWidth => Math.Min(_value / 600.0 * 200, 200);
    public string ValueText => $"{_value:F0} mm";

    public Brush BarColor => _value switch
    {
        > 500 => new SolidColorBrush(Color.FromRgb(243, 139, 168)),
        > 400 => new SolidColorBrush(Color.FromRgb(249, 226, 175)),
        _ => new SolidColorBrush(Color.FromRgb(166, 227, 161))
    };

    private void Notify([CallerMemberName] string? name = null) =>
        PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(name));
}
