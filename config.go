package configurator

import (
	"bufio"
	"errors"
	"github.com/spf13/viper"
	"log"
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
	// TODO this is not quite the best idea -- we are taking apart the field map based on an env tag
	//  if that env tag has a non-string value then this craps out (which is what we want... but maybe just do that in
	//  check kind?
	_, _, err = fieldValueMap(v)
	if err != nil {
		return err
	}

	// At this point, we know that configType is a struct with string fields only

	// Parse configuration defaults
	defaultTagFieldNames, defaultTagValues := parseTagValues(v, "default")
	log.Println("THE DEFAULTS ARE:", defaultTagFieldNames)
	log.Println("THE DEFAULT VALUES ARE:", defaultTagValues)
	for _, fieldName := range defaultTagFieldNames {
		viperConfig.SetDefault(fieldName, defaultTagValues[fieldName])
	}

	// Bind environment variables and parse the environment
	envTagFieldNames, envTagValues := parseTagValues(v, "env")
	for _, fieldName := range envTagFieldNames {
		_ = viperConfig.BindEnv(fieldName, envTagValues[fieldName])
	}
	viperConfig.AutomaticEnv()

	// Finally, configure v: Set the values of v
	rve := reflect.ValueOf(v).Elem()
	for _, fieldName := range envTagFieldNames {
		log.Println("SETTING FIELD:", fieldName)
		log.Println("WITH VALUE:", viperConfig.GetString(fieldName))
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

	// Get required field names and current field values
	requiredFields := requiredFieldMap(v)
	_, fieldValues, err := fieldValueMap(v)
	if err != nil {
		return err
	}

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

func parseTagValues(v interface{}, tag string) ([]string, map[string]string) {
	// assume v is a pointer to a struct

	rve := reflect.ValueOf(v).Elem()

	var fieldNames []string
	tagValues := make(map[string]string)
	for i := 0; i < rve.NumField(); i++ {
		field := rve.Type().Field(i)
		fieldName := field.Name
		tagValue := field.Tag.Get(tag)
		if tagValue != "" {
			fieldNames = append(fieldNames, fieldName)
			tagValues[fieldName] = tagValue
		}
	}

	return fieldNames, tagValues
}

func fieldValueMap(configType interface{}) ([]string, map[string]string, error) {
	// assume v is a pointer to a struct
	// caller must first use checkKind

	rve := reflect.ValueOf(configType).Elem()

	var fieldNames []string
	fieldValues := make(map[string]string)
	for i := 0; i < rve.NumField(); i++ {
		field := rve.Type().Field(i)
		// TODO this is not ideal, as it checks the env tag only... move this to checkKind
		if field.Tag.Get("env") != "" && field.Type.Kind() != reflect.String {
			return nil, nil, errors.New("invalid configType: struct fields with 'env' tag must be of type string")
		}
		fieldNames = append(fieldNames, field.Name)
		fieldValues[field.Name] = rve.Field(i).String()
	}
	sort.Strings(fieldNames)

	return fieldNames, fieldValues, nil
}

func requiredFieldMap(configType interface{}) map[string]bool {
	// assume v is a pointer to a struct
	// caller must first use checkKind

	requiredFields, requiredTags := parseTagValues(configType, "config")

	result := make(map[string]bool)
	for _, k := range requiredFields {
		if v, ok := requiredTags[k]; !ok {
			result[k] = false
		} else {
			result[k] = strings.Contains(v, "required")
		}
	}

	return result
}
