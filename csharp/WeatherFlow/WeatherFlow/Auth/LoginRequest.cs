using System.Text.Json.Serialization;

namespace WeatherFlow.Auth;

public class LoginRequest
{
    [JsonPropertyName("publicId")]
    public required string PublicId { get; set; }
}