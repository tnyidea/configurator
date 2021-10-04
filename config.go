package configurator

import (
	"errors"
	"github.com/spf13/viper"
	"reflect"
	"sort"
	"strings"
)

func ParseEnvConfig(configType interface{}, configFile ...string) error {
	// configType must be a pointer to a struct
	err := checkConfigTypeKind(configType)
	if err != nil {
		return err
	}

	viperConfig := viper.New()

	// Get configuration key names and validate config type
	keys, _, err := keyValueMap(configType)
	if err != nil {
		return err
	}

	// At this point, we know that configType is a struct with string fields only

	// Parse configuration defaults
	defaults := parseTag(configType, keys, "default")
	for _, k := range keys {
		if v, ok := defaults[k]; ok {
			viperConfig.SetDefault(k, v)
		}
	}

	// Bind environment variables
	env := parseTag(configType, keys, "env")
	for _, k := range keys {
		if v, ok := env[k]; ok {
			_ = viperConfig.BindEnv(k, v)
		}
	}

	// Parse configuration file (if provided) filename := filenameBinding(configType)
	var filename string
	if configFile != nil {
		filename = configFile[0]
	}
	if filename != "" {
		viperConfig.SetConfigFile(filename)
		err := viperConfig.ReadInConfig()
		if err != nil {
			return errors.New("invalid configuration: " + err.Error())
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
	e := reflect.ValueOf(configType).Elem()
	for _, k := range keys {
		e.FieldByName(k).SetString(viperConfig.GetString(k))
	}

	return nil
}

func ValidateConfig(configType interface{}) error {
	// configType must be a pointer to a struct
	err := checkConfigTypeKind(configType)
	if err != nil {
		return err
	}

	// Get configuration key names, values, and required keys
	keys, keyValues, err := keyValueMap(configType)
	if err != nil {
		return err
	}
	required := requiredKeys(configType, keys)

	var flag bool
	var messages []string
	for k, v := range required {
		if v && keyValues[k] == "" {
			messages = append(messages, k+" not set")
			flag = true
		}
	}
	if flag {
		return errors.New("invalid configuration: missing required values:\n" + strings.Join(messages, "\n"))
	}

	return nil
}

func checkConfigTypeKind(configType interface{}) error {
	// configType must be a pointer to a struct
	if reflect.ValueOf(configType).Kind() != reflect.Ptr {
		return errors.New("invalid configType: must be a pointer to a struct")
	}
	if reflect.Indirect(reflect.ValueOf(configType)).Kind() != reflect.Struct {
		return errors.New("invalid configType: not a struct")
	}

	return nil
}

func keyValueMap(configType interface{}) ([]string, map[string]string, error) {
	// assume configType is a pointer to a struct, caller must first use checkConfigTypeKind

	e := reflect.ValueOf(configType).Elem()

	var keys []string
	keyValues := make(map[string]string)
	for i := 0; i < e.NumField(); i++ {
		f := e.Type().Field(i)
		if f.Type.Kind() != reflect.String {
			return nil, nil, errors.New("invalid configType: struct fields must be of type string")
		}
		keys = append(keys, f.Name)
		keyValues[f.Name] = e.Field(i).String()
	}
	sort.Strings(keys)

	return keys, keyValues, nil
}

func requiredKeys(configType interface{}, keys []string) map[string]bool {
	// assume configType is a pointer to a struct, caller must first use checkConfigTypeKind

	requiredTags := parseTag(configType, keys, "config")

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

func parseTag(configType interface{}, keys []string, tagName string) map[string]string {
	// assume configType is a pointer to a struct, caller must first use checkConfigTypeKind

	e := reflect.ValueOf(configType).Elem()

	tagValues := make(map[string]string)
	for _, k := range keys {
		f, _ := e.Type().FieldByName(k)
		v := f.Tag.Get(tagName)
		if v != "" {
			tagValues[k] = v
		}
	}

	return tagValues
}
