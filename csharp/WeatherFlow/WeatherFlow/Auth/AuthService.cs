using System.Text;
using System.Text.Json;

namespace WeatherFlow.Auth;

public class AuthService
{
    private const string UrlLogin = "https://thingsboard.bda-itnovum.com/api/auth/login/public";

    private const string UrlAllowedPermissions = "https://thingsboard.bda-itnovum.com/api/permissions/allowedPermissions";
    
    private const string PublicId = "d58b18a0-1440-11ef-aef4-af283e5094d9";
    // private readonly string _publicId = Guid.NewGuid().ToString();
    
    private static readonly HttpClient Client = new();
    
    public async Task<AuthResponse> Login()
    {
        var loginRequest = new LoginRequest
        {
            PublicId = PublicId
        };
        var payload = JsonSerializer.Serialize(loginRequest);

        using var content = new StringContent(payload, Encoding.UTF8, "application/json");
        
        HttpResponseMessage response;
        try
        {
            response = await Client.PostAsync(UrlLogin, content);
        }
        catch (Exception ex)
        {
            throw new Exception($"Error sending login request: {ex.Message}");
        }
        
        var responseBody = await response.Content.ReadAsStringAsync();

        if (!response.IsSuccessStatusCode)
        {
            throw new Exception($"Auth failed {(int)response.StatusCode}: {responseBody}");
        }

        AuthResponse? authResponse;
        try
        {
            authResponse = JsonSerializer.Deserialize<AuthResponse>(responseBody);
        }
        catch (Exception ex)
        {
            throw new Exception($"Parsing error: {ex.Message}");
        }

        if (authResponse == null || string.IsNullOrEmpty(authResponse.Token))
        {
            throw new Exception($"No token received. Body: {responseBody}");
        }

        Console.WriteLine("login successful");
        return authResponse;
    }

    private void Refresh()
    {
        
    }

    public async Task<string> Permissions(string token)
    {
        Client.DefaultRequestHeaders.Add("X-Authorization", $"Bearer {token}");

        var response = await Client.GetAsync(UrlAllowedPermissions);
        var responseBody = await response.Content.ReadAsStringAsync();
        return responseBody;
    }
}