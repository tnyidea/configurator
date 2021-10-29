package configurator

import (
	"reflect"
)

func SetAllDefaultValues(v interface{}) error {
	return setDefaultValues(v, false)
}

func SetEmptyFieldDefaultValues(v interface{}) error {
	return setDefaultValues(v, true)
}

func setDefaultValues(v interface{}, setZeroOnly bool) error {
	err := checkKindStructPtr(v)
	if err != nil {
		return err
	}

	dm := parseDefaultMetadata(v)

	fieldNameValueMap := dm.fieldNameDefaultValueMap

	if setZeroOnly { // if setZeroOnly, delete the non-zero struct fields from the map
		sm := parseStructMetadata(v)
		for _, fieldName := range dm.fieldNames {
			if !reflect.ValueOf(sm.fieldNameValueMap[fieldName]).IsZero() {
				delete(fieldNameValueMap, fieldName)
			}
		}
	}
	setValuesFromFieldNameValueMap(v, fieldNameValueMap)

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
