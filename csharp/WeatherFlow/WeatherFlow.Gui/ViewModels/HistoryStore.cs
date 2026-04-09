namespace WeatherFlow.Gui.ViewModels;

public record DataPoint(DateTime Time, double Value);

public class HistoryStore
{
    private readonly Dictionary<string, List<DataPoint>> _series = new();
    private readonly object _lock = new();

    public void Record(string key, long timestampMs, string rawValue)
    {
        if (!double.TryParse(rawValue, System.Globalization.NumberStyles.Any,
                System.Globalization.CultureInfo.InvariantCulture, out var value))
            return;

        var time = DateTimeOffset.FromUnixTimeMilliseconds(timestampMs).LocalDateTime;

        lock (_lock)
        {
            if (!_series.TryGetValue(key, out var list))
            {
                list = new List<DataPoint>();
                _series[key] = list;
            }

            // Avoid duplicates (same timestamp)
            if (list.Count > 0 && list[^1].Time == time)
                return;

            list.Add(new DataPoint(time, value));
            PurgeOldEntries(list);
        }
    }

    public List<DataPoint> GetSeries(string key)
    {
        lock (_lock)
        {
            if (!_series.TryGetValue(key, out var list))
                return [];
            PurgeOldEntries(list);
            return [.. list];
        }
    }

    public bool HasData(string key)
    {
        lock (_lock)
            return _series.TryGetValue(key, out var list) && list.Count > 1;
    }

    private static void PurgeOldEntries(List<DataPoint> list)
    {
        var cutoff = DateTime.Today; // midnight today — removes all from yesterday and older
        list.RemoveAll(p => p.Time < cutoff);
    }
}
