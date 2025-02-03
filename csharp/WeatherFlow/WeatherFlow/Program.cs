using WeatherFlow.Auth;
using WeatherFlow.Thingsboard;

Console.WriteLine("Starting");

try
{
    var authService = new AuthService();
    var response = await authService.Login();
// Console.WriteLine(response.Token);

    var permissions = await authService.Permissions(response.Token);

    var client = new ThingsBoardClient(response.Token, "26945210-05ec-11ef-ac80-dde635ebcdb2");
    var dataStream = await client.StartDataStream();

    await foreach (var data in dataStream)
    {
        Console.WriteLine($"Temperature: {data.Temperature}");
    }
}
catch (Exception e)
{
    Console.WriteLine($"Exception: {e.Message}");
}

Console.ReadKey();

// var uri = new Uri("");
//
// using (var ws = new ClientWebSocket())
// {
//     try
//     {
//         Console.WriteLine("Connecting...");
//         await ws.ConnectAsync(uri, CancellationToken.None);
//         Console.WriteLine("Connection established");
//         
//         
//         
//     }
//     catch (Exception e)
//     {
//         Console.WriteLine($"Error: {e.Message}");
//     }
// }