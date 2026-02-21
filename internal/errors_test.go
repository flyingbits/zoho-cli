package internal

import "testing"

func TestNewError(t *testing.T) {
	e := NewError("test")
	if e.Message != "test" {
		t.Errorf("expected 'test', got %s", e.Message)
	}
	if e.ExitCode != ExitError {
		t.Errorf("expected %d, got %d", ExitError, e.ExitCode)
	}
	if e.Error() != "test" {
		t.Errorf("Error() should return message")
	}
}

func TestNewAuthError(t *testing.T) {
	e := NewAuthError("auth fail")
	if e.ExitCode != ExitAuth {
		t.Errorf("expected %d, got %d", ExitAuth, e.ExitCode)
	}
}

func TestNewNotFoundError(t *testing.T) {
	e := NewNotFoundError("not found")
	if e.ExitCode != ExitNotFound {
		t.Errorf("expected %d, got %d", ExitNotFound, e.ExitCode)
	}
}

func TestNewValidationError(t *testing.T) {
	e := NewValidationError("bad input")
	if e.ExitCode != ExitValidation {
		t.Errorf("expected %d, got %d", ExitValidation, e.ExitCode)
	}
}

func TestNewAPIError(t *testing.T) {
	tests := []struct {
		status   int
		wantExit int
	}{
		{401, ExitAuth},
		{404, ExitNotFound},
		{500, ExitError},
		{400, ExitError},
	}
	for _, tt := range tests {
		e := NewAPIError(tt.status, "body")
		if e.ExitCode != tt.wantExit {
			t.Errorf("status %d: expected exit %d, got %d", tt.status, tt.wantExit, e.ExitCode)
		}
		if e.StatusCode != tt.status {
			t.Errorf("expected status %d, got %d", tt.status, e.StatusCode)
		}
	}
}

func TestExitCodes(t *testing.T) {
	if ExitSuccess != 0 {
		t.Error("ExitSuccess should be 0")
	}
	if ExitError != 1 {
		t.Error("ExitError should be 1")
	}
	if ExitAuth != 2 {
		t.Error("ExitAuth should be 2")
	}
	if ExitNotFound != 3 {
		t.Error("ExitNotFound should be 3")
	}
	if ExitValidation != 4 {
		t.Error("ExitValidation should be 4")
	}
}
