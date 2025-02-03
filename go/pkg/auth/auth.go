package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const url = "https://thingsboard.bda-itnovum.com/api/auth/login/public"
const publicID = "d58b18a0-1440-11ef-aef4-af283e5094d9"

type loginRequest struct {
	PublicID string `json:"publicId"`
}

type authResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func Authenticate() (string, error) {
	loginRequest := loginRequest{
		PublicID: publicID,
	}

	jsonPayload, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResp authResponse
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse authentication response: %v", err)
	}

	return authResp.Token, nil
}
