using WeatherFlow.Auth;
using WeatherFlow.Thingsboard;

namespace WeatherFlow;

class Program
{
    private static async Task Main()
    {
        using (var cts = new CancellationTokenSource())
        {
            var cancellationTask = Task.Run(() =>
            {
                Console.WriteLine("Press 'q' to terminate program");
                while (true)
                {
                    if (Console.ReadKey(true).Key == ConsoleKey.Q)
                    {
                        cts.Cancel();
                        break;
                    }
                }
            }, cts.Token);

            try
            {
                var authService = new AuthService();
                var response = await authService.Login();

                var permissions = await authService.Permissions(response.Token);

                var client = new ThingsBoardClient(response.Token, "26945210-05ec-11ef-ac80-dde635ebcdb2");
                var dataStream = await client.StartDataStream(cts.Token);

                await foreach (var data in dataStream.WithCancellation(cts.Token))
                {
                    if (data == null)
                    {
                        continue;
                    }
                    Console.WriteLine($"Temperature: {data.Temperature}");
                }
            }
            catch (OperationCanceledException)
            {
                Console.WriteLine("program terminated");
            }
            catch (Exception e)
            {
                Console.WriteLine($"Exception: {e.Message}");
            }
        }
    }
}