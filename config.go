package configurator

import (
	"errors"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

func ParseEnvConfig(configType interface{}, filename string) (interface{}, error) {
	viperConfig := viper.New()

	// Get configuration key names and validate config type
	keys, err := keyNames(configType)
	if err != nil {
		return nil, err
	}

	// At this point, we know that configType is a struct with string fields only

	// Parse configuration defaults
	defaults := parseTags(configType, keys, "default")
	for _, k := range keys {
		if v, ok := defaults[k]; ok {
			viperConfig.SetDefault(k, v)
		}
	}

	// Bind environment variables
	env := parseTags(configType, keys, "env")
	for _, k := range keys {
		if v, ok := env[k]; ok {
			_ = viperConfig.BindEnv(k, v)
		}
	}

	// Parse configuration file (if provided) filename := filenameBinding(configType)
	if filename != "" {
		viperConfig.SetConfigFile(filename)
		err := viperConfig.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	// Parse environment variables
	viperConfig.AutomaticEnv()

	// Validate Configuration
	required := parseTags(configType, keys, "required")
	var flag bool
	var messages []string
	for _, k := range keys {
		v, ok := required[k]
		if ok && v == "true" && !viperConfig.IsSet(k) {
			messages = append(messages, k+" not set")
			flag = true
		}
	}
	if flag {
		return nil, errors.New("invalid configuration: missing required values:\n" + strings.Join(messages, "\n"))
	}

	ps := reflect.New(reflect.TypeOf(configType))
	for _, k := range keys {
		s := ps.Elem()
		f := s.FieldByName(k)
		f.SetString(viperConfig.GetString(k))
	}

	return ps, nil
}

func keyNames(configType interface{}) ([]string, error) {
	s := reflect.ValueOf(&configType).Elem()

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

func parseTags(configType interface{}, keys []string, tagName string) map[string]string {
	s := reflect.ValueOf(&configType).Elem()

	var result map[string]string
	for _, k := range keys {
		f, _ := s.Type().FieldByName(k)
		v := f.Tag.Get(tagName)
		if v != "" {
			result[k] = v
		}
	}

	return result
}
