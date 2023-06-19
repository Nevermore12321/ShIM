package mapping

import (
	"github.com/Nevermore12321/ShIM/pkg/jsonx"
	"io"
)

// json tag key. e.g. `json:"test"`
const jsonTagKey = "json"

var jsonUnmarshaler = NewUnmarshaler(jsonTagKey)

// getJsonUnmarshaler initial the Unmarshaler for json tag key
func getJsonUnmarshaler(opts ...UnmarshalOpt) *Unmarshaler {
	if len(opts) > 0 {
		return NewUnmarshaler(jsonTagKey, opts...)
	}

	return jsonUnmarshaler
}

// UnmarshalJsonMap unmarshal m(map) into v.
func UnmarshalJsonMap(m map[string]any, v any, otps ...UnmarshalOpt) error {
	return getJsonUnmarshaler(otps...).Unmarshal(m, v)
}

// UnmarshalJsonBytes unmarshal content(content bytes) into v.
func UnmarshalJsonBytes(content []byte, v any, opts ...UnmarshalOpt) error {
	var m any
	if err := jsonx.Unmarshal(content, v); err != nil {
		return err
	}
	return getJsonUnmarshaler(opts...).Unmarshal(m, v)
}

// UnmarshalJsonReader unmarshal content from reader into v.
func UnmarshalJsonReader(reader io.Reader, v any, opts ...UnmarshalOpt) error {
	var m any
	if err := jsonx.UnmarshalFromReader(reader, v); err != nil {
		return err
	}
	return getJsonUnmarshaler(opts...).Unmarshal(m, v)
}
