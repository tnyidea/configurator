package configurator

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"reflect"
	"sort"
	"strings"
)

func ParseEnvConfig(v interface{}, configFile ...string) error {
	// configType must be a pointer to a struct and not nil
	err := checkKind(v)
	if err != nil {
		return err
	}

	viperConfig := viper.New()

	// Get configuration field names and validate config type
	fieldNames, _, err := fieldValueMap(v)
	if err != nil {
		return err
	}

	// At this point, we know that configType is a struct with string fields only

	// Parse configuration defaults
	defaults := parseTagValues(v, fieldNames, "default")
	for _, fieldName := range fieldNames {
		if fieldValue, ok := defaults[fieldName]; ok {
			viperConfig.SetDefault(fieldName, fieldValue)
		}
	}

	// Bind environment variables
	env := parseTagValues(v, fieldNames, "env")
	for _, fieldName := range fieldNames {
		if fieldValue, ok := env[fieldName]; ok {
			_ = viperConfig.BindEnv(fieldName, fieldValue)
		}
	}

	// Parse configuration file (if provided)
	if configFile != nil {
		filename := configFile[0]
		tokens := strings.Split(filename, ".")

		var ok bool
		if len(tokens) != 2 && tokens[1] != "env" {
			log.Println("ignoring file: " + filename + ": must have .env extension")
			ok = false
		}
		if ok {
			viperConfig.SetConfigFile(tokens[0])
			err := viperConfig.ReadInConfig()
			log.Println(viperConfig)
			if err != nil {
				return errors.New("invalid configuration: " + err.Error())
			}
		}
	}

	// Parse environment variables
	viperConfig.AutomaticEnv()

	// Validate Configuration
	//required := requiredKeys(configType, keys)
	//if err != nil {
	//	return err
	//}
	//var flag bool
	//var messages []string
	//for k, v := range required {
	//	if v && !viperConfig.IsSet(k) {
	//		messages = append(messages, k+" not set")
	//		flag = true
	//	}
	//}
	//if flag {
	//	return errors.New("invalid configuration: missing required values:\n" + strings.Join(messages, "\n"))
	//}
	//
	e := reflect.ValueOf(v).Elem()
	for _, fieldName := range fieldNames {
		e.FieldByName(fieldName).SetString(viperConfig.GetString(fieldName))
	}

	return nil
}

func ValidateConfig(configType interface{}) error {
	// configType must be a pointer to a struct and not nil
	err := checkKind(configType)
	if err != nil {
		return err
	}

	// Get configuration key names, values, and required keys
	fieldNames, fieldValues, err := fieldValueMap(configType)
	if err != nil {
		return err
	}
	requiredFields := requiredFieldMap(configType, fieldNames)

	var messages []string
	for fieldName, isRequired := range requiredFields {
		if isRequired && fieldValues[fieldName] == "" {
			messages = append(messages, fieldName+" not set")
		}
	}
	if messages != nil {
		return errors.New("invalid configuration: missing required values:\n" + strings.Join(messages, "\n"))
	}

	return nil
}

func checkKind(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("invalid type: must be a pointer to a struct")
	}
	if reflect.Indirect(rv).Kind() != reflect.Struct {
		return errors.New("invalid type: must be a pointer to a struct")
	}

	return nil
}

func fieldValueMap(configType interface{}) ([]string, map[string]string, error) {
	// assume v is a pointer to a struct
	// caller must first use checkKind

	rve := reflect.ValueOf(configType).Elem()

	var fieldNames []string
	fieldValues := make(map[string]string)
	for i := 0; i < rve.NumField(); i++ {
		field := rve.Type().Field(i)
		if field.Type.Kind() != reflect.String {
			return nil, nil, errors.New("invalid configType: struct fields must be of type string")
		}
		fieldNames = append(fieldNames, field.Name)
		fieldValues[field.Name] = rve.Field(i).String()
	}
	sort.Strings(fieldNames)

	return fieldNames, fieldValues, nil
}

func requiredFieldMap(configType interface{}, keys []string) map[string]bool {
	// assume v is a pointer to a struct
	// caller must first use checkKind

	requiredTags := parseTagValues(configType, keys, "config")

	result := make(map[string]bool)
	for _, k := range keys {
		if v, ok := requiredTags[k]; !ok {
			result[k] = false
		} else {
			result[k] = strings.Contains(v, "required")
		}
	}

	return result
}

func parseTagValues(v interface{}, fieldNames []string, tag string) map[string]string {
	// assume v is a pointer to a struct
	// caller must first use checkKind

	rve := reflect.ValueOf(v).Elem()

	tagValues := make(map[string]string)
	for _, fieldName := range fieldNames {
		field, _ := rve.Type().FieldByName(fieldName)
		tagValue := field.Tag.Get(tag)
		if tagValue != "" {
			tagValues[fieldName] = tagValue
		}
	}

	return tagValues
}
