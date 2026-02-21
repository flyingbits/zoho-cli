package output

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestJSON(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	data := map[string]string{"key": "value"}
	if err := JSON(data); err != nil {
		t.Fatal(err)
	}

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = old

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["key"] != "value" {
		t.Errorf("expected 'value', got %s", result["key"])
	}
}

func TestJSONRaw(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	if err := JSONRaw([]byte(`{"foo":"bar"}`)); err != nil {
		t.Fatal(err)
	}

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = old

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result["foo"] != "bar" {
		t.Errorf("expected 'bar', got %s", result["foo"])
	}
}

func TestJSONNull(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	if err := JSON(nil); err != nil {
		t.Fatal(err)
	}

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = old

	if got := bytes.TrimSpace(buf.Bytes()); string(got) != "null" {
		t.Errorf("expected 'null', got %s", got)
	}
}
