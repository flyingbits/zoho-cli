package auth

import (
	"os"
	"testing"
	"time"
)

func TestTokenValid(t *testing.T) {
	c := &AuthConfig{}
	if c.TokenValid() {
		t.Error("empty config should not be valid")
	}

	c.AccessToken = "test"
	c.ExpiresAt = time.Now().UTC().Add(10 * time.Minute)
	if !c.TokenValid() {
		t.Error("token with future expiry should be valid")
	}

	c.ExpiresAt = time.Now().UTC().Add(2 * time.Minute)
	if c.TokenValid() {
		t.Error("token expiring within 5 minute buffer should not be valid")
	}

	c.ExpiresAt = time.Now().UTC().Add(-1 * time.Minute)
	if c.TokenValid() {
		t.Error("expired token should not be valid")
	}
}

func TestConfigDir(t *testing.T) {
	original := os.Getenv("ZOHO_CLI_CONFIG_DIR")
	defer os.Setenv("ZOHO_CLI_CONFIG_DIR", original)

	os.Setenv("ZOHO_CLI_CONFIG_DIR", "/tmp/zoho-test")
	if ConfigDir() != "/tmp/zoho-test" {
		t.Errorf("expected /tmp/zoho-test, got %s", ConfigDir())
	}

	os.Unsetenv("ZOHO_CLI_CONFIG_DIR")
	dir := ConfigDir()
	if dir == "" {
		t.Error("ConfigDir should return a path")
	}
}

func TestCacheKey(t *testing.T) {
	k1 := cacheKey("token1")
	k2 := cacheKey("token2")
	if k1 == k2 {
		t.Error("different tokens should produce different cache keys")
	}
	if len(k1) != 16 {
		t.Errorf("cache key should be 16 chars, got %d", len(k1))
	}
}

func TestResolveAuthNoConfig(t *testing.T) {
	original := os.Getenv("ZOHO_CLIENT_ID")
	defer os.Setenv("ZOHO_CLIENT_ID", original)
	os.Unsetenv("ZOHO_CLIENT_ID")
	os.Unsetenv("ZOHO_CLIENT_SECRET")
	os.Unsetenv("ZOHO_REFRESH_TOKEN")

	origDir := os.Getenv("ZOHO_CLI_CONFIG_DIR")
	defer os.Setenv("ZOHO_CLI_CONFIG_DIR", origDir)
	os.Setenv("ZOHO_CLI_CONFIG_DIR", "/tmp/zoho-cli-test-nonexistent")

	_, err := ResolveAuth()
	if err == nil {
		t.Error("expected error when no auth configured")
	}
}
