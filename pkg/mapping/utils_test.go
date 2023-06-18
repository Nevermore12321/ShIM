package mapping

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestDereference(t *testing.T) {
	// test data
	i := 8
	s := "shim"
	num := struct {
		f float64
	}{
		f: 8.8,
	}
	testCase := []struct {
		t      reflect.Type
		expect reflect.Kind
	}{
		{
			t:      reflect.TypeOf(i),
			expect: reflect.Int,
		},
		{
			t:      reflect.TypeOf(&i),
			expect: reflect.Int,
		},
		{
			t:      reflect.TypeOf(s),
			expect: reflect.String,
		},
		{
			t:      reflect.TypeOf(&s),
			expect: reflect.String,
		},
		{
			t:      reflect.TypeOf(num.f),
			expect: reflect.Float64,
		},
		{
			t:      reflect.TypeOf(&num.f),
			expect: reflect.Float64,
		},
	}

	// start testing
	for _, item := range testCase {
		t.Run(item.t.String(), func(t *testing.T) {
			assert.Equal(t, item.expect, Dereference(item.t).Kind(), "case now: "+item.t.String())
		})
	}
}
