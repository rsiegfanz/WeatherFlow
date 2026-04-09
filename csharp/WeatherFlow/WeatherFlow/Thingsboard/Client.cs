using System.Collections.Concurrent;
using System.IO;
using System.Net.WebSockets;
using System.Text;
using System.Text.Json;

namespace WeatherFlow.Thingsboard;

public class ThingsBoardClient : IDisposable
{
    private const string WsUrl = "wss://thingsboard.bda-itnovum.com/api/ws";

    private readonly string _token;
    private readonly ClientWebSocket _ws = new();
    private readonly ConcurrentDictionary<string, string> _entityNames = new();

    public ThingsBoardClient(string token) => _token = token;

    public async Task RunAsync(Action<WsMessage> onMessage, CancellationToken ct)
    {
        await _ws.ConnectAsync(new Uri(WsUrl), ct);
        await SendPayload(ct);

        var buffer = new byte[64 * 1024];
        using var ms = new MemoryStream();

        while (_ws.State == WebSocketState.Open && !ct.IsCancellationRequested)
        {
            ms.SetLength(0);
            WebSocketReceiveResult result;
            do
            {
                result = await _ws.ReceiveAsync(new ArraySegment<byte>(buffer), ct);
                if (result.MessageType == WebSocketMessageType.Close) return;
                ms.Write(buffer, 0, result.Count);
            } while (!result.EndOfMessage);

            var json = Encoding.UTF8.GetString(ms.GetBuffer(), 0, (int)ms.Length);
            var msg = JsonSerializer.Deserialize<WsMessage>(json);
            if (msg == null) continue;

            LearnEntityNames(msg);
            onMessage(msg);
        }
    }

    public string? GetEntityName(string entityId) =>
        _entityNames.TryGetValue(entityId, out var name) ? name : null;

    private async Task SendPayload(CancellationToken ct)
    {
        var request = RequestBuilder.Build(_token);
        var json = JsonSerializer.Serialize(request);
        var bytes = Encoding.UTF8.GetBytes(json);
        await _ws.SendAsync(new ArraySegment<byte>(bytes), WebSocketMessageType.Text, true, ct);
    }

    private void LearnEntityNames(WsMessage msg)
    {
        var entries = msg.Data?.Data ?? msg.Update;
        if (entries == null) return;

        foreach (var e in entries)
        {
            if (e.Latest.TryGetValue("ATTRIBUTE", out var attrs) &&
                attrs.TryGetValue("displayName", out var dn) && !string.IsNullOrEmpty(dn.Value))
            {
                _entityNames[e.EntityId.Id] = dn.Value;
            }
            else if (!_entityNames.ContainsKey(e.EntityId.Id) &&
                     e.Latest.TryGetValue("ENTITY_FIELD", out var fields) &&
                     fields.TryGetValue("label", out var label) && !string.IsNullOrEmpty(label.Value))
            {
                _entityNames[e.EntityId.Id] = label.Value;
            }
        }
    }

    public void Dispose()
    {
        if (_ws.State == WebSocketState.Open)
        {
            try { _ws.CloseAsync(WebSocketCloseStatus.NormalClosure, "", CancellationToken.None).Wait(2000); }
            catch { /* ignore */ }
        }
        _ws.Dispose();
    }
}
