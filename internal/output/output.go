package output

import (
	"encoding/json"
	"fmt"
	"os"
)

func JSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func JSONRaw(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		_, err2 := os.Stdout.Write(data)
		return err2
	}
	return JSON(v)
}

func Err(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
