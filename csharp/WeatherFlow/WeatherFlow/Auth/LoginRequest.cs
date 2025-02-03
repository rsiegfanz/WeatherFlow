using System.Text.Json.Serialization;

namespace WeatherFlow.Auth;

public class LoginRequest
{
    [JsonPropertyName("publicId")]
    public string PublicId { get; set; }
}