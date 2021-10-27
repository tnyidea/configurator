package configurator

import (
	"bufio"
	"errors"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"sort"
	"strings"
)

func SetEnvFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return err
		}

		// Only handle items of form VARIABLE=value
		tokens := strings.Split(scanner.Text(), "=")
		if len(tokens) == 2 {
			_ = os.Setenv(tokens[0], tokens[1])
		}
	}

	return nil
}

func ParseEnvConfig(v interface{}) error {
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

	// Parse environment
	viperConfig.AutomaticEnv()

	// Set the values of v
	rve := reflect.ValueOf(v).Elem()
	for _, fieldName := range fieldNames {
		rve.FieldByName(fieldName).SetString(viperConfig.GetString(fieldName))
	}

	return nil
}

func ValidateConfig(v interface{}) error {
	// configType must be a pointer to a struct and not nil
	err := checkKind(v)
	if err != nil {
		return err
	}

	// Get configuration key names, values, and required keys
	fieldNames, fieldValues, err := fieldValueMap(v)
	if err != nil {
		return err
	}
	requiredFields := requiredFieldMap(v, fieldNames)

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

// Helpers
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

func fieldValueMap(configType interface{}) ([]string, map[string]string, error) {
	// assume v is a pointer to a struct
	// caller must first use checkKind

	rve := reflect.ValueOf(configType).Elem()

	var fieldNames []string
	fieldValues := make(map[string]string)
	for i := 0; i < rve.NumField(); i++ {
		field := rve.Type().Field(i)
		if field.Tag.Get("config") != "" && field.Type.Kind() != reflect.String {
			return nil, nil, errors.New("invalid configType: struct fields with config tag must be of type string")
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
