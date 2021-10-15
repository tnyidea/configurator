package configurator

import (
	"errors"
	"github.com/spf13/viper"
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

	// Get configuration key names and validate config type
	keys, _, err := keyValueMap(v)
	if err != nil {
		return err
	}

	// At this point, we know that configType is a struct with string fields only

	// Parse configuration defaults
	defaults := parseTag(v, keys, "default")
	for _, key := range keys {
		if keyValue, ok := defaults[key]; ok {
			viperConfig.SetDefault(key, keyValue)
		}
	}

	// Bind environment variables
	env := parseTag(v, keys, "env")
	for _, key := range keys {
		if keyValue, ok := env[key]; ok {
			_ = viperConfig.BindEnv(key, keyValue)
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
	e := reflect.ValueOf(v).Elem()
	for _, key := range keys {
		e.FieldByName(key).SetString(viperConfig.GetString(key))
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
	keys, keyValues, err := keyValueMap(configType)
	if err != nil {
		return err
	}
	required := requiredKeys(configType, keys)

	var messages []string
	for key, keyRequired := range required {
		if keyRequired && keyValues[key] == "" {
			messages = append(messages, key+" not set")
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

func keyValueMap(configType interface{}) ([]string, map[string]string, error) {
	// assume v is a pointer to a struct
	// caller must first use checkKind

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
	// assume v is a pointer to a struct
	// caller must first use checkKind

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

func parseTag(v interface{}, keys []string, tag string) map[string]string {
	// assume v is a pointer to a struct
	// caller must first use checkKind

	e := reflect.ValueOf(v).Elem()

	tagValues := make(map[string]string)
	for _, k := range keys {
		f, _ := e.Type().FieldByName(k)
		tagValue := f.Tag.Get(tag)
		if tagValue != "" {
			tagValues[k] = tagValue
		}
	}

	return tagValues
}
