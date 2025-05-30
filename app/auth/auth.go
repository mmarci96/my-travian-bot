package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Optin    bool   `json:"mobileOptimizations"`
	Width    string `json:"w"`
}

func CreateLoginRequest(email, password, serverUrl string) (string, error) {
	url := serverUrl + "/api/v1/auth/login"

	payload := LoginRequest{
		Name:     email,
		Password: password,
		Optin:    true,
		Width:    "1440:900",
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return string(body), nil
}

type RedirectResponse struct {
	RedirectTo string `json:"redirectTo"`
	Code       string `json:"code"`
	Nonce      string `json:"nonce"`
}

func FollowRedirect(jsonResponse string, serverUrl string) (string, error) {
	var redirectResp RedirectResponse

	err := json.Unmarshal([]byte(jsonResponse), &redirectResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse redirect response: %v", err)
	}

	// Unescape redirect path
	redirectPath := strings.ReplaceAll(redirectResp.RedirectTo, `\/`, `/`)

	// Handle relative vs absolute URL
	var fullURL string
	if strings.HasPrefix(redirectPath, "/") {
		fullURL = strings.TrimRight(serverUrl, "/") + redirectPath
	} else {
		fullURL = redirectPath
	}

	// Make the GET request to the redirect URL
	resp, err := http.Get(fullURL)
	if err != nil {
		return "", fmt.Errorf("failed to follow redirect: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read redirect response: %v", err)
	}

	// fmt.Printf("Redirected response:\n%s\n", string(body))
	return string(body), nil
}
