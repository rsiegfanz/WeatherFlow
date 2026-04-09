using System.Text.Json.Serialization;

namespace WeatherFlow.Auth;

public class AuthResponse
{
    [JsonPropertyName("token")]
    public required string Token { get; set; }

    [JsonPropertyName("refreshToken")]
    public string? RefreshToken { get; set; }
}