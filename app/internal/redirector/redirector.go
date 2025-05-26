package redirector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RedirectResponse struct {
	RedirectTo string `json:"redirectTo"`
	Code       string `json:"code"`
	Nonce      string `json:"nonce"`
}

func FollowRedirect(jsonResponse string, serverUrl string) error {
	var redirectResp RedirectResponse

	err := json.Unmarshal([]byte(jsonResponse), &redirectResp)
	if err != nil {
		return fmt.Errorf("failed to parse redirect response: %v", err)
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
		return fmt.Errorf("failed to follow redirect: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read redirect response: %v", err)
	}

	fmt.Printf("Redirected response:\n%s\n", string(body))
	return nil
}
