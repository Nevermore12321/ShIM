package jsonx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Unmarshal unmarshals data bytes into v.
func Unmarshal(data []byte, v any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(v); err != nil {
		return formatJsonError(string(data), err)
	}
	return nil
}

func formatJsonError(v string, err error) error {
	return fmt.Errorf("jsonx string: `%s`, error: `%w`", v, err)
}
