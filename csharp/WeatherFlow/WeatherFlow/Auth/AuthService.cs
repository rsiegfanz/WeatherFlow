using System.Net.Http;
using System.Text;
using System.Text.Json;

namespace WeatherFlow.Auth;

public class AuthService
{
    private const string LoginUrl = "https://thingsboard.bda-itnovum.com/api/auth/login/public";
    private static readonly HttpClient Client = new();

    public static async Task<string> Authenticate(string publicId)
    {
        var payload = JsonSerializer.Serialize(new LoginRequest { PublicId = publicId });
        using var content = new StringContent(payload, Encoding.UTF8, "application/json");

        var response = await Client.PostAsync(LoginUrl, content);
        var body = await response.Content.ReadAsStringAsync();

        if (!response.IsSuccessStatusCode)
            throw new AuthException($"Auth failed ({(int)response.StatusCode}): {body}");

        var authResponse = JsonSerializer.Deserialize<AuthResponse>(body)
            ?? throw new AuthException($"Failed to parse auth response: {body}");

        if (string.IsNullOrEmpty(authResponse.Token))
            throw new AuthException("No token in auth response");

        return authResponse.Token;
    }
}

public class AuthException(string message) : Exception(message);
