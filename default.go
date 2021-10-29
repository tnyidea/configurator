package configurator

import (
	"reflect"
)

func SetDefaultValues(v interface{}) error {
	err := checkKindStructPtr(v)
	if err != nil {
		return err
	}

	// TODO this is temporary for initial version; future versions should allow other types in a struct
	err = checkStructFieldsKindString(reflect.Indirect(reflect.ValueOf(v)).Interface())
	if err != nil {
		return err
	}

	dm := parseDefaultMetadata(v)
	setValuesFromFieldNameValueMap(v, dm.fieldNameDefaultValueMap)

	return nil
}

type defaultMetadata struct {
	fieldNames               []string
	fieldNameDefaultValueMap map[string]interface{}
}

func parseDefaultMetadata(v interface{}) defaultMetadata {
	// assume v is a pointer to a struct

	rve := reflect.ValueOf(v).Elem()

	var dm defaultMetadata
	for i := 0; i < rve.NumField(); i++ {
		field := rve.Type().Field(i)
		fieldName := field.Name
		tagValue := field.Tag.Get("default")
		if tagValue != "" {
			if dm.fieldNameDefaultValueMap == nil {
				dm.fieldNameDefaultValueMap = make(map[string]interface{})
			}
			dm.fieldNames = append(dm.fieldNames, fieldName)
			dm.fieldNameDefaultValueMap[fieldName] = tagValue
		}
	}

	return dm
}
