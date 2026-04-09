using WeatherFlow.Auth;
using WeatherFlow.Thingsboard;

namespace WeatherFlow;

class Program
{
    private static async Task<int> Main(string[] args)
    {
        var publicId = GetArg(args, "--public-id");
        if (string.IsNullOrEmpty(publicId))
        {
            Console.Error.WriteLine("Usage: weatherflow --public-id <ID>");
            Console.Error.WriteLine();
            Console.Error.WriteLine("Example:");
            Console.Error.WriteLine("  weatherflow --public-id d58b18a0-1440-11ef-aef4-af283e5094d9");
            return 1;
        }

        using var cts = new CancellationTokenSource();
        Console.CancelKeyPress += (_, e) => { e.Cancel = true; cts.Cancel(); };

        try
        {
            var token = await AuthService.Authenticate(publicId);

            using var client = new ThingsBoardClient(token);
            var dashboard = new Dashboard(client);

            Console.Clear();
            Console.CursorVisible = false;

            await client.RunAsync(dashboard.HandleMessage, cts.Token);
        }
        catch (OperationCanceledException)
        {
            Console.CursorVisible = true;
            Console.WriteLine("\nBeendet.");
        }
        catch (AuthException ex)
        {
            Console.Error.WriteLine($"Authentication error: {ex.Message}");
            return 2;
        }
        catch (Exception ex)
        {
            Console.Error.WriteLine($"Error: {ex.Message}");
            return 3;
        }

        return 0;
    }

    private static string? GetArg(string[] args, string name)
    {
        var idx = Array.IndexOf(args, name);
        return idx >= 0 && idx + 1 < args.Length ? args[idx + 1] : null;
    }
}
