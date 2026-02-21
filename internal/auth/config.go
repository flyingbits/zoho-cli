package auth

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/omin8tor/zoho-cli/internal"
	zohodc "github.com/omin8tor/zoho-cli/internal/dc"
)

type AuthConfig struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	DC           string
	AccountsURL  string
	APIDomain    string
	AccessToken  string
	ExpiresAt    time.Time
	Scopes       string
	Source       string
}

func (c *AuthConfig) TokenValid() bool {
	if c.AccessToken == "" || c.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().UTC().Add(5 * time.Minute).Before(c.ExpiresAt)
}

func ConfigDir() string {
	if d := os.Getenv("ZOHO_CLI_CONFIG_DIR"); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "zoho-cli")
}

func ConfigFile() string { return filepath.Join(ConfigDir(), "config.toml") }
func TokensFile() string { return filepath.Join(ConfigDir(), "tokens.json") }
func CacheDir() string   { return filepath.Join(ConfigDir(), "cache") }

func cacheKey(refreshToken string) string {
	h := sha256.Sum256([]byte(refreshToken))
	return fmt.Sprintf("%x", h)[:16]
}

type cachedToken struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   string `json:"expires_at"`
}

func LoadCachedAccessToken(refreshToken string) (string, time.Time, bool) {
	path := filepath.Join(CacheDir(), cacheKey(refreshToken)+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", time.Time{}, false
	}
	var ct cachedToken
	if err := json.Unmarshal(data, &ct); err != nil || ct.AccessToken == "" || ct.ExpiresAt == "" {
		return "", time.Time{}, false
	}
	expiresAt, err := time.Parse(time.RFC3339, ct.ExpiresAt)
	if err != nil {
		return "", time.Time{}, false
	}
	if time.Now().UTC().Add(5 * time.Minute).After(expiresAt) {
		return "", time.Time{}, false
	}
	return ct.AccessToken, expiresAt, true
}

func SaveCachedAccessToken(refreshToken, accessToken string, expiresIn int) error {
	dir := CacheDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	path := filepath.Join(dir, cacheKey(refreshToken)+".json")
	expiresAt := time.Now().UTC().Add(time.Duration(expiresIn) * time.Second)
	ct := cachedToken{AccessToken: accessToken, ExpiresAt: expiresAt.Format(time.RFC3339)}
	data, _ := json.Marshal(ct)
	return os.WriteFile(path, data, 0600)
}

type tokensFile struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    string `json:"access_token_expires_at"`
	DC           string `json:"dc"`
	AccountsURL  string `json:"accounts_url"`
	APIDomain    string `json:"api_domain"`
	Scopes       string `json:"scopes"`
}

func loadTokens() *tokensFile {
	data, err := os.ReadFile(TokensFile())
	if err != nil {
		return nil
	}
	var tf tokensFile
	if err := json.Unmarshal(data, &tf); err != nil {
		return nil
	}
	return &tf
}

func SaveTokens(refreshToken, accessToken string, expiresIn int, dc, accountsURL, apiDomain, scopes string) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	expiresAt := time.Now().UTC().Add(time.Duration(expiresIn) * time.Second)
	tf := tokensFile{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		ExpiresAt:    expiresAt.Format(time.RFC3339),
		DC:           dc,
		AccountsURL:  accountsURL,
		APIDomain:    apiDomain,
		Scopes:       scopes,
	}
	data, _ := json.MarshalIndent(tf, "", "  ")
	if err := os.WriteFile(TokensFile(), data, 0600); err != nil {
		return err
	}
	return SaveCachedAccessToken(refreshToken, accessToken, expiresIn)
}

type tomlConfig struct {
	Auth struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"auth"`
}

func fromEnv() *AuthConfig {
	clientID := os.Getenv("ZOHO_CLIENT_ID")
	clientSecret := os.Getenv("ZOHO_CLIENT_SECRET")
	refreshToken := os.Getenv("ZOHO_REFRESH_TOKEN")

	if clientID == "" || clientSecret == "" || refreshToken == "" {
		return nil
	}

	dc := os.Getenv("ZOHO_DC")
	if dc == "" {
		dc = "com"
	}
	accountsURL := os.Getenv("ZOHO_ACCOUNTS_URL")
	if accountsURL == "" {
		accountsURL = zohodc.AccountsURL(dc)
	}

	var accessToken string
	var expiresAt time.Time
	if at, exp, ok := LoadCachedAccessToken(refreshToken); ok {
		accessToken = at
		expiresAt = exp
	}

	return &AuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
		DC:           dc,
		AccountsURL:  accountsURL,
		AccessToken:  accessToken,
		ExpiresAt:    expiresAt,
		Source:       "env",
	}
}

func fromConfigFile() *AuthConfig {
	data, err := os.ReadFile(ConfigFile())
	if err != nil {
		return nil
	}

	var cfg tomlConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	if cfg.Auth.ClientID == "" || cfg.Auth.ClientSecret == "" {
		return nil
	}

	tokens := loadTokens()
	if tokens == nil || tokens.RefreshToken == "" {
		return nil
	}

	var expiresAt time.Time
	if tokens.ExpiresAt != "" {
		expiresAt, _ = time.Parse(time.RFC3339, tokens.ExpiresAt)
	}

	dc := tokens.DC
	if dc == "" {
		dc = "com"
	}

	return &AuthConfig{
		ClientID:     cfg.Auth.ClientID,
		ClientSecret: cfg.Auth.ClientSecret,
		RefreshToken: tokens.RefreshToken,
		DC:           dc,
		AccountsURL:  tokens.AccountsURL,
		APIDomain:    tokens.APIDomain,
		AccessToken:  tokens.AccessToken,
		ExpiresAt:    expiresAt,
		Scopes:       tokens.Scopes,
		Source:       "config",
	}
}

func ResolveAuth() (*AuthConfig, error) {
	if cfg := fromEnv(); cfg != nil {
		return cfg, nil
	}
	if cfg := fromConfigFile(); cfg != nil {
		return cfg, nil
	}
	return nil, internal.NewAuthError("Not authenticated. Run `zoho auth login` to authenticate.")
}
