package jsonx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
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

// Unmarshal unmarshals data from io.Reader into v.
func UnmarshalFromReader(reader io.Reader, v any) error {
	var buffer strings.Builder
	teeReader := io.TeeReader(reader, &buffer)

	decoder := json.NewDecoder(teeReader)
	decoder.UseNumber()
	if err := decoder.Decode(v); err != nil {
		return formatJsonError(buffer.String(), err)
	}
	return nil
}

func formatJsonError(v string, err error) error {
	return fmt.Errorf("jsonx string: `%s`, error: `%w`", v, err)
}
