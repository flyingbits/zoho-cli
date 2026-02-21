package dc

import "testing"

func TestGetDC(t *testing.T) {
	cfg := GetDC("com")
	if cfg.Accounts != "https://accounts.zoho.com" {
		t.Errorf("expected accounts.zoho.com, got %s", cfg.Accounts)
	}
	if cfg.CRM != "https://zohoapis.com" {
		t.Errorf("expected zohoapis.com, got %s", cfg.CRM)
	}
}

func TestGetDCFallback(t *testing.T) {
	cfg := GetDC("nonexistent")
	if cfg.Accounts != "https://accounts.zoho.com" {
		t.Error("expected fallback to com DC")
	}
}

func TestAllDCs(t *testing.T) {
	for _, dc := range ValidDCs {
		cfg := GetDC(dc)
		if cfg.Accounts == "" {
			t.Errorf("DC %s has empty accounts URL", dc)
		}
		if cfg.CRM == "" {
			t.Errorf("DC %s has empty CRM URL", dc)
		}
		if cfg.Projects == "" {
			t.Errorf("DC %s has empty Projects URL", dc)
		}
		if cfg.WorkDrive == "" {
			t.Errorf("DC %s has empty WorkDrive URL", dc)
		}
		if cfg.Writer == "" {
			t.Errorf("DC %s has empty Writer URL", dc)
		}
		if cfg.Cliq == "" {
			t.Errorf("DC %s has empty Cliq URL", dc)
		}
		if cfg.Download == "" {
			t.Errorf("DC %s has empty Download URL", dc)
		}
	}
}

func TestHelperFunctions(t *testing.T) {
	tests := []struct {
		fn   func(string) string
		dc   string
		want string
	}{
		{AccountsURL, "com", "https://accounts.zoho.com"},
		{AccountsURL, "eu", "https://accounts.zoho.eu"},
		{CliqURL, "com", "https://cliq.zoho.com"},
		{CRMURL, "com", "https://zohoapis.com"},
		{ProjectsURL, "com", "https://projectsapi.zoho.com"},
		{WorkDriveURL, "com", "https://workdrive.zoho.com"},
		{WriterURL, "com", "https://www.zohoapis.com/writer"},
		{DownloadURL, "com", "https://download.zoho.com"},
	}
	for _, tt := range tests {
		got := tt.fn(tt.dc)
		if got != tt.want {
			t.Errorf("expected %s, got %s", tt.want, got)
		}
	}
}

func TestValidDCsCount(t *testing.T) {
	if len(ValidDCs) != 9 {
		t.Errorf("expected 9 valid DCs, got %d", len(ValidDCs))
	}
}
