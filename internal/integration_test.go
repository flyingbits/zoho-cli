//go:build integration

package internal_test

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func zoho(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.CommandContext(context.Background(), "./zoho", args...)
	cmd.Dir = ".."
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("zoho %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

func TestCRMModulesList(t *testing.T) {
	out := zoho(t, "crm", "modules", "list")
	if !strings.Contains(out, "[") {
		t.Error("expected JSON array output")
	}
}

func TestCRMModulesFields(t *testing.T) {
	out := zoho(t, "crm", "modules", "fields", "Contacts")
	if !strings.Contains(out, "[") {
		t.Error("expected JSON array output")
	}
}

func TestCRMRecordsList(t *testing.T) {
	out := zoho(t, "crm", "records", "list", "Contacts", "--fields", "id,Full_Name")
	if !strings.Contains(out, "[") && !strings.Contains(out, "data") {
		t.Error("expected JSON output")
	}
}

func TestCRMSearchGlobal(t *testing.T) {
	out := zoho(t, "crm", "search-global", "test")
	if len(out) == 0 {
		t.Error("expected some output")
	}
}

func TestCRMUsersList(t *testing.T) {
	out := zoho(t, "crm", "users", "list")
	if !strings.Contains(out, "[") {
		t.Error("expected JSON array output")
	}
}

func TestHelpAll(t *testing.T) {
	out := zoho(t, "--help-all")
	if !strings.Contains(out, "crm") {
		t.Error("expected crm in help-all output")
	}
	if !strings.Contains(out, "projects") {
		t.Error("expected projects in help-all output")
	}
	if !strings.Contains(out, "drive") {
		t.Error("expected drive in help-all output")
	}
}
