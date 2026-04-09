using WeatherFlow.Auth;
using WeatherFlow.Thingsboard;

namespace WeatherFlow;

class Program
{
    private const string DefaultPublicId = "d58b18a0-1440-11ef-aef4-af283e5094d9";

    private static async Task<int> Main(string[] args)
    {
        var publicId = GetArg(args, "--public-id") ?? PromptForPublicId();

        using var cts = new CancellationTokenSource();
        Console.CancelKeyPress += (_, e) => { e.Cancel = true; cts.Cancel(); };

        try
        {
            Console.WriteLine("Verbinde...");
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

    private static string PromptForPublicId()
    {
        Console.WriteLine("╔══════════════════════════════════════════════════════════╗");
        Console.WriteLine("║  WEATHERFLOW                                           ║");
        Console.WriteLine("╚══════════════════════════════════════════════════════════╝");
        Console.WriteLine();
        Console.WriteLine($"  ThingsBoard Public ID [{DefaultPublicId}]:");
        Console.Write("  > ");

        var input = Console.ReadLine()?.Trim();
        return string.IsNullOrEmpty(input) ? DefaultPublicId : input;
    }

    private static string? GetArg(string[] args, string name)
    {
        var idx = Array.IndexOf(args, name);
        return idx >= 0 && idx + 1 < args.Length ? args[idx + 1] : null;
    }
}
