package configurator

import (
	"fmt"
	"reflect"
)

type structMetadata struct {
	fieldNames        []string
	fieldNameValueMap map[string]interface{}
}

func parseStructMetadata(v interface{}) structMetadata {
	// assume v is a pointer to a struct

	rve := reflect.ValueOf(v).Elem()

	var sm structMetadata
	sm.fieldNameValueMap = make(map[string]interface{})
	for i := 0; i < rve.NumField(); i++ {
		field := rve.Type().Field(i)
		sm.fieldNames = append(sm.fieldNames, field.Name)
		sm.fieldNameValueMap[field.Name] = rve.Field(i).Interface()
	}

	return sm
}

func setValuesFromFieldNameValueMap(v interface{}, m map[string]interface{}) {
	// Assume that v is a non-nil pointer to a struct with string fields only
	// TODO this is temporary as future types are planned for support
	rv := reflect.Indirect(reflect.ValueOf(v))

	for fieldName, fieldValue := range m {
		rv.FieldByName(fieldName).SetString(fmt.Sprintf("%v", fieldValue))
	}
}
