package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/omin8tor/zoho-cli/internal"
	zohodc "github.com/omin8tor/zoho-cli/internal/dc"
)

func normalizeExpiresIn(raw any) int {
	switch v := raw.(type) {
	case float64:
		n := int(v)
		if n > 86400 {
			n = n / 1000
		}
		if n < 1 {
			return 1
		}
		return n
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			return 3600
		}
		if n > 86400 {
			n = n / 1000
		}
		if n < 1 {
			return 1
		}
		return n
	default:
		return 3600
	}
}

func RefreshAccessToken(config *AuthConfig) (string, error) {
	accountsURL := config.AccountsURL
	if accountsURL == "" {
		accountsURL = zohodc.AccountsURL(config.DC)
	}

	params := url.Values{
		"client_id":     {config.ClientID},
		"client_secret": {config.ClientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {config.RefreshToken},
	}

	resp, err := http.PostForm(accountsURL+"/oauth/v2/token", params)
	if err != nil {
		return "", internal.NewAuthError(fmt.Sprintf("Token refresh request failed: %v", err))
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return "", internal.NewAuthError(fmt.Sprintf("Token refresh failed: unexpected response: %s", body))
	}

	if errMsg, ok := data["error"]; ok {
		return "", internal.NewAuthError(fmt.Sprintf("Token refresh failed: %v", errMsg))
	}

	accessToken, ok := data["access_token"].(string)
	if !ok || accessToken == "" {
		return "", internal.NewAuthError(fmt.Sprintf("Token refresh failed: no access_token in response: %s", body))
	}

	expiresIn := normalizeExpiresIn(data["expires_in"])
	apiDomain, _ := data["api_domain"].(string)
	if apiDomain == "" {
		apiDomain = config.APIDomain
	}

	config.AccessToken = accessToken
	config.APIDomain = apiDomain
	config.ExpiresAt = time.Now().UTC().Add(time.Duration(expiresIn) * time.Second)

	_ = SaveCachedAccessToken(config.RefreshToken, accessToken, expiresIn)

	if config.Source == "config" {
		_ = SaveTokens(config.RefreshToken, accessToken, expiresIn, config.DC, config.AccountsURL, apiDomain, config.Scopes)
	}

	return accessToken, nil
}

func EnsureAccessToken(config *AuthConfig, forceRefresh bool) (string, error) {
	if !forceRefresh && config.TokenValid() {
		return config.AccessToken, nil
	}
	return RefreshAccessToken(config)
}
