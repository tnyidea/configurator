package configurator

import (
	"errors"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

func ParseEnvConfig(configType interface{}, configFile ...string) error {
	//if !reflect.ValueOf(configType).CanAddr() {
	//	return errors.New("invalid configType: must pass address to configType")
	//}

	viperConfig := viper.New()

	// Get configuration key names and validate config type
	keys, err := keyNames(configType)
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
	required := requiredKeys(configType, keys)
	var flag bool
	var messages []string
	for k, v := range required {
		if v && !viperConfig.IsSet(k) {
			messages = append(messages, k+" not set")
			flag = true
		}
	}
	if flag {
		return errors.New("invalid configuration: missing required values:\n" + strings.Join(messages, "\n"))
	}

	s := reflect.ValueOf(configType).Elem()
	for _, k := range keys {
		s.FieldByName(k).SetString(viperConfig.GetString(k))
	}

	return nil
}

func keyNames(configType interface{}) ([]string, error) {
	s := reflect.ValueOf(configType).Elem()

	if s.Kind() != reflect.Struct {
		return nil, errors.New("invalid configType: not a struct")
	}

	var result []string
	for i := 0; i < s.NumField(); i++ {
		f := s.Type().Field(i)
		if f.Type.Kind() != reflect.String {
			return nil, errors.New("invalid configType: struct fields must be of type string")
		}
		result = append(result, f.Name)
	}

	return result, nil
}

func requiredKeys(configType interface{}, keys []string) map[string]bool {
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
	s := reflect.ValueOf(configType).Elem()

	result := make(map[string]string)
	for _, k := range keys {
		f, _ := s.Type().FieldByName(k)
		v := f.Tag.Get(tagName)
		if v != "" {
			result[k] = v
		}
	}

	return result
}
