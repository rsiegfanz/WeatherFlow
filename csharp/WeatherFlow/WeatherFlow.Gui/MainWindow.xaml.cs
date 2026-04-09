using System.Windows;
using System.Windows.Controls;
using WeatherFlow.Auth;
using WeatherFlow.Gui.ViewModels;
using WeatherFlow.Thingsboard;

namespace WeatherFlow.Gui;

public partial class MainWindow : Window
{
    private const string DefaultPublicId = "d58b18a0-1440-11ef-aef4-af283e5094d9";

    private readonly DashboardViewModel _vm = new();
    private CancellationTokenSource? _cts;

    public MainWindow()
    {
        InitializeComponent();
        DataContext = _vm;
        Loaded += OnLoaded;
        Closed += (_, _) => _cts?.Cancel();
    }

    private async void OnLoaded(object sender, RoutedEventArgs e)
    {
        var publicId = PromptForPublicId();
        if (publicId == null)
        {
            Close();
            return;
        }

        _cts = new CancellationTokenSource();

        try
        {
            _vm.Status = "Authentifiziere...";
            var token = await AuthService.Authenticate(publicId);

            _vm.Status = "Verbinde WebSocket...";
            using var client = new ThingsBoardClient(token);

            await client.RunAsync(msg =>
            {
                Dispatcher.Invoke(() => _vm.ProcessMessage(msg, client));
            }, _cts.Token);
        }
        catch (OperationCanceledException) { }
        catch (Exception ex)
        {
            _vm.Status = "Fehler";
            _vm.StatusColor = new System.Windows.Media.SolidColorBrush(
                System.Windows.Media.Color.FromRgb(243, 139, 168));
            MessageBox.Show(ex.Message, "Fehler", MessageBoxButton.OK, MessageBoxImage.Error);
        }
    }

    private static string? PromptForPublicId()
    {
        var dialog = new Window
        {
            Title = "WeatherFlow - Verbinden",
            Width = 500, Height = 200,
            WindowStartupLocation = WindowStartupLocation.CenterScreen,
            ResizeMode = ResizeMode.NoResize,
            Background = new System.Windows.Media.SolidColorBrush(
                System.Windows.Media.Color.FromRgb(30, 30, 46))
        };

        var panel = new StackPanel { Margin = new Thickness(24) };

        panel.Children.Add(new TextBlock
        {
            Text = "ThingsBoard Public ID:",
            Foreground = new System.Windows.Media.SolidColorBrush(
                System.Windows.Media.Color.FromRgb(205, 214, 244)),
            FontSize = 14, Margin = new Thickness(0, 0, 0, 8)
        });

        var textBox = new TextBox
        {
            Text = DefaultPublicId,
            FontSize = 14,
            Padding = new Thickness(8, 6, 8, 6),
            Background = new System.Windows.Media.SolidColorBrush(
                System.Windows.Media.Color.FromRgb(42, 42, 60)),
            Foreground = new System.Windows.Media.SolidColorBrush(
                System.Windows.Media.Color.FromRgb(205, 214, 244)),
            BorderBrush = new System.Windows.Media.SolidColorBrush(
                System.Windows.Media.Color.FromRgb(137, 180, 250)),
            SelectionStart = 0, SelectionLength = DefaultPublicId.Length
        };
        panel.Children.Add(textBox);

        var button = new Button
        {
            Content = "Verbinden",
            FontSize = 14,
            Padding = new Thickness(20, 8, 20, 8),
            Margin = new Thickness(0, 16, 0, 0),
            HorizontalAlignment = HorizontalAlignment.Right,
            Background = new System.Windows.Media.SolidColorBrush(
                System.Windows.Media.Color.FromRgb(137, 180, 250)),
            Foreground = new System.Windows.Media.SolidColorBrush(
                System.Windows.Media.Color.FromRgb(30, 30, 46)),
            BorderThickness = new Thickness(0)
        };

        string? result = null;
        button.Click += (_, _) => { result = textBox.Text.Trim(); dialog.Close(); };
        textBox.KeyDown += (_, args) =>
        {
            if (args.Key == System.Windows.Input.Key.Enter) { result = textBox.Text.Trim(); dialog.Close(); }
        };

        panel.Children.Add(button);
        dialog.Content = panel;
        dialog.ShowDialog();

        return string.IsNullOrEmpty(result) ? null : result;
    }
}
