package mapping

import "reflect"

// Dereference If t is a pointer type, returns its element type that t point to.
// Otherwise, return t
func Dereference(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}

	return t
}
