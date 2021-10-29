package configurator

import (
	"errors"
	"reflect"
)

func checkKindStructPtr(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("invalid type: must be a pointer to a struct")
	}
	if reflect.Indirect(rv).Kind() != reflect.Struct {
		return errors.New("invalid type: must be a pointer to a struct")
	}

	return nil
}

func checkStructFieldsKindString(v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Struct {
		return errors.New("invalid type: must be a struct with string fields only")
	}

	for i := 0; i < rv.NumField(); i++ {
		if rv.Field(i).Kind() != reflect.String {
			return errors.New("invalid type: must be a struct with string fields only")
		}
	}

	return nil
}
