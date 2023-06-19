package mapping

import (
	"errors"
	"reflect"
)

var (
	errValueNotSettable = errors.New("value is not settable")
	errUnsupportedType  = errors.New("unsupported type on setting field value")
	errTypeMismatch     = errors.New("type mismatch")
)

type (
	// unmarshalOptions customize options for unmarshal with given tag key
	unmarshalOptions struct {
		fillDefault  bool
		fromString   bool
		canonicalKey func(key string) string
	}

	// Unmarshaler is used to unmarshal with given tag key.
	Unmarshaler struct {
		key  string
		opts unmarshalOptions
	}

	// UnmarshalOpt defines the method to customize an Unmarshaler.
	UnmarshalOpt func(*unmarshalOptions)
)

// NewUnmarshaler initial the Unmarshaler object
func NewUnmarshaler(key string, opts ...UnmarshalOpt) *Unmarshaler {
	unmarshaler := &Unmarshaler{
		key: key,
	}

	for _, opt := range opts {
		opt(&unmarshaler.opts)
	}

	return unmarshaler
}

// Unmarshal unmarshal m into v.
func (u *Unmarshaler) Unmarshal(m any, v any) error {
	valueType := reflect.TypeOf(v)

	// if converted type v is not a pointer, error! not settable
	// because must modify object v, so need v is a pointer
	if valueType.Kind() != reflect.Ptr {
		return errValueNotSettable
	}

	// obtain pointed element data type
	elementType := Dereference(valueType)

	// Determine the type of m
	switch mv := m.(type) {
	case map[string]any: // key: value , v must be struct
		if elementType.Kind() != reflect.Struct {
			return errTypeMismatch
		}
		return u.UnmarshalValuer(mapValuer(iv), v)
	case []any: // key: - a -b , v must be slice
		if elementType.Kind() != reflect.Slice {
			return errTypeMismatch
		}

	default:
		return errUnsupportedType
	}
}
