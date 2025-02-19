package context

import (
	"context"
	"reflect"
)

// Get - Get the value associated to the key
func Get(ctx context.Context, key interface{}) (value interface{}) {

	value = ctx.Value(key)

	if value == nil {
		value = reflect.New(reflect.TypeOf(key).Elem()).Interface()
	}

	return value
}
