using System.Globalization;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Shapes;
using WeatherFlow.Gui.ViewModels;

namespace WeatherFlow.Gui;

public partial class ChartWindow : Window
{
    private readonly List<DataPoint> _data;
    private readonly string _unit;

    private static readonly Brush LineBrush = new SolidColorBrush(Color.FromRgb(137, 180, 250));
    private static readonly Brush FillBrush = new SolidColorBrush(Color.FromArgb(40, 137, 180, 250));
    private static readonly Brush GridBrush = new SolidColorBrush(Color.FromArgb(30, 205, 214, 244));
    private static readonly Brush DotBrush = new SolidColorBrush(Color.FromRgb(148, 226, 213));
    private static readonly Brush LabelBrush = new SolidColorBrush(Color.FromRgb(108, 112, 134));

    public ChartWindow(string title, string unit, List<DataPoint> data)
    {
        InitializeComponent();
        Title = $"WeatherFlow — {title}";
        TitleText.Text = $"{title} ({unit})";
        _data = data;
        _unit = unit;

        if (_data.Count == 0)
        {
            SummaryText.Text = "Keine Daten vorhanden.";
            return;
        }

        var min = _data.Min(d => d.Value);
        var max = _data.Max(d => d.Value);
        var last = _data[^1];
        SummaryText.Text = $"Aktuell: {last.Value:F1}{_unit}  |  Min: {min:F1}{_unit}  |  Max: {max:F1}{_unit}  |  {_data.Count} Messpunkte";

        Loaded += (_, _) => DrawChart();
        SizeChanged += (_, _) => DrawChart();
        ChartCanvas.MouseMove += OnMouseMove;
        ChartCanvas.MouseLeave += (_, _) => HoverText.Text = "";
    }

    private void DrawChart()
    {
        ChartCanvas.Children.Clear();
        YAxisCanvas.Children.Clear();
        XAxisCanvas.Children.Clear();

        if (_data.Count < 2) return;

        var w = ChartCanvas.ActualWidth;
        var h = ChartCanvas.ActualHeight;
        if (w <= 0 || h <= 0) return;

        var minVal = _data.Min(d => d.Value);
        var maxVal = _data.Max(d => d.Value);
        var range = maxVal - minVal;
        if (range < 0.001) range = 1;

        // Add 10% padding
        minVal -= range * 0.1;
        maxVal += range * 0.1;
        range = maxVal - minVal;

        var minTime = _data[0].Time;
        var maxTime = _data[^1].Time;
        var timeRange = (maxTime - minTime).TotalSeconds;
        if (timeRange < 1) timeRange = 1;

        // Grid lines (horizontal)
        var gridLines = 5;
        for (var i = 0; i <= gridLines; i++)
        {
            var y = h - (h * i / gridLines);
            var val = minVal + range * i / gridLines;

            var line = new Line
            {
                X1 = 0, X2 = w, Y1 = y, Y2 = y,
                Stroke = GridBrush, StrokeThickness = 1
            };
            ChartCanvas.Children.Add(line);

            var label = new TextBlock
            {
                Text = val.ToString("F1"), FontSize = 10, Foreground = LabelBrush
            };
            Canvas.SetLeft(label, 2);
            Canvas.SetTop(label, y - 8);
            YAxisCanvas.Children.Add(label);
        }

        // Build points
        var points = new PointCollection();
        foreach (var dp in _data)
        {
            var x = (dp.Time - minTime).TotalSeconds / timeRange * w;
            var y = h - (dp.Value - minVal) / range * h;
            points.Add(new Point(x, y));
        }

        // Fill area
        var fillPoints = new PointCollection(points) { new(w, h), new(0, h) };
        var fill = new Polygon { Points = fillPoints, Fill = FillBrush };
        ChartCanvas.Children.Add(fill);

        // Line
        var polyline = new Polyline
        {
            Points = points,
            Stroke = LineBrush,
            StrokeThickness = 2,
            StrokeLineJoin = PenLineJoin.Round
        };
        ChartCanvas.Children.Add(polyline);

        // Dots
        foreach (var pt in points)
        {
            var dot = new Ellipse { Width = 5, Height = 5, Fill = DotBrush };
            Canvas.SetLeft(dot, pt.X - 2.5);
            Canvas.SetTop(dot, pt.Y - 2.5);
            ChartCanvas.Children.Add(dot);
        }

        // X-Axis labels
        var labelCount = Math.Min(8, _data.Count);
        var step = Math.Max(1, _data.Count / labelCount);
        for (var i = 0; i < _data.Count; i += step)
        {
            var x = (_data[i].Time - minTime).TotalSeconds / timeRange * w;
            var label = new TextBlock
            {
                Text = _data[i].Time.ToString("HH:mm"), FontSize = 10, Foreground = LabelBrush
            };
            Canvas.SetLeft(label, x - 15);
            Canvas.SetTop(label, 5);
            XAxisCanvas.Children.Add(label);
        }
    }

    private void OnMouseMove(object sender, MouseEventArgs e)
    {
        if (_data.Count < 2) return;

        var pos = e.GetPosition(ChartCanvas);
        var w = ChartCanvas.ActualWidth;
        if (w <= 0) return;

        var minTime = _data[0].Time;
        var maxTime = _data[^1].Time;
        var timeRange = (maxTime - minTime).TotalSeconds;
        if (timeRange < 1) return;

        var hoverTime = minTime.AddSeconds(pos.X / w * timeRange);

        // Find nearest point
        var nearest = _data.MinBy(d => Math.Abs((d.Time - hoverTime).TotalSeconds));
        if (nearest != null)
            HoverText.Text = $"{nearest.Time:HH:mm:ss}  —  {nearest.Value:F1} {_unit}";
    }
}
