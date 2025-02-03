using System.Net.WebSockets;
using System.Runtime.CompilerServices;
using System.Text;
using System.Text.Json;

namespace WeatherFlow.Thingsboard;

public class ThingsBoardClient(string token, string deviceId)
{
    private const string WebSocketUrl = "wss://thingsboard.bda-itnovum.com/api/ws";

    public async Task<IAsyncEnumerable<DeviceData?>> StartDataStream(CancellationToken cancellationToken = default)
    {
        var webSocket = new ClientWebSocket();
        
        await webSocket.ConnectAsync(new Uri(WebSocketUrl), cancellationToken);

        await SendMessage(webSocket, ThingsBoardRequestBuilder.CreateWeatherDataRequest(deviceId, token));
        Console.WriteLine("sent message");

        return CreateDataStream(webSocket, cancellationToken);
    }

    private async IAsyncEnumerable<DeviceData?> CreateDataStream(
        ClientWebSocket webSocket, 
        [EnumeratorCancellation] CancellationToken cancellationToken)
    {
        while (webSocket.State == WebSocketState.Open && !cancellationToken.IsCancellationRequested)
        {
            DeviceData? data = null;
            try
            {
                var message = await ReceiveMessage(webSocket, cancellationToken);
                Console.WriteLine($"message: {message}");
                if (!string.IsNullOrEmpty(message))
                {
                    data = ParseDeviceData(message);   
                }
            }
            catch (Exception e)
            {
                Console.WriteLine($"Error reading from socket: {e.Message}");
            }

            if (data != null)
            {
                yield return data;   
            }
        }
    }

    private async Task SendMessage(ClientWebSocket webSocket, object message)
    {
        var options = new JsonSerializerOptions
        {
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase
        };
        var jsonMessage = JsonSerializer.Serialize(message, options);
        
        Console.WriteLine(jsonMessage);
        
        var buffer = Encoding.UTF8.GetBytes(jsonMessage);
        
        await webSocket.SendAsync(
            new ArraySegment<byte>(buffer), 
            WebSocketMessageType.Text, 
            true, 
            CancellationToken.None
        );
    }

    private async Task<string?> ReceiveMessage(ClientWebSocket webSocket, CancellationToken cancellationToken)
    {
        var buffer = new byte[1024 * 4];
        var result = await webSocket.ReceiveAsync(new ArraySegment<byte>(buffer), cancellationToken);
        Console.WriteLine($"result: {result.MessageType} / {result.Count}");
        return result.MessageType == WebSocketMessageType.Text 
            ? Encoding.UTF8.GetString(buffer, 0, result.Count) 
            : null;
    }

    private DeviceData? ParseDeviceData(string jsonMessage)
    {
        Console.WriteLine($"data: {jsonMessage}");
        // Implementieren Sie hier Ihre spezifische Parsing-Logik
        return new DeviceData();
    }
}

public class DeviceData
{
    public double Temperature { get; set; }
    public double RainGauge { get; set; }
    public double WindSpeed { get; set; }
}