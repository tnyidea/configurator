package configurator

import (
	"errors"
	"reflect"
)

func CheckRequiredValues(v interface{}) []error {
	err := checkKindStructPtr(v)
	if err != nil {
		return []error{err}
	}

	rm := parseRequiredMetadata(v)

	sm := parseStructMetadata(v)

	// Check that all required field values are !reflect.ValueOf().IsZero()
	var errs []error
	for fieldName, fieldIsRequired := range rm.fieldNameIsRequiredMap {
		if fieldIsRequired && reflect.ValueOf(sm.fieldNameValueMap[fieldName]).IsZero() {
			errs = append(errs, errors.New("invalid value: "+fieldName+" not set"))
		}
	}

	return errs
}

type requiredMetadata struct {
	fieldNames             []string
	fieldNameIsRequiredMap map[string]bool
}

func parseRequiredMetadata(v interface{}) requiredMetadata {
	// assume v is a pointer to a struct

	rve := reflect.ValueOf(v).Elem()

	var rm requiredMetadata
	for i := 0; i < rve.NumField(); i++ {
		field := rve.Type().Field(i)

		if field.Tag.Get("required") == "true" {
			if rm.fieldNameIsRequiredMap == nil {
				rm.fieldNameIsRequiredMap = make(map[string]bool)
			}
			rm.fieldNames = append(rm.fieldNames, field.Name)
			rm.fieldNameIsRequiredMap[field.Name] = true
		}
	}

	return rm
}
