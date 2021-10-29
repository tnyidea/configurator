package configurator

import (
	"bufio"
	"os"
	"reflect"
	"strings"
)

func SetValuesFromOsEnv(v interface{}) error {
	err := checkKindStructPtr(v)
	if err != nil {
		return err
	}

	em := parseEnvMetadata(v)

	var fieldNameValueMap map[string]interface{}
	for _, envVariable := range em.envVariables {
		if envValue := os.Getenv(envVariable); envValue != "" {
			if fieldNameValueMap == nil {
				fieldNameValueMap = make(map[string]interface{})
			}
			fieldNameValueMap[em.envVariableFieldNameMap[envVariable]] = envValue
		}
	}

	if fieldNameValueMap != nil {
		setValuesFromFieldNameValueMap(v, fieldNameValueMap)
	}

	return nil
}

func SetValuesFromEnvFile(v interface{}, filename string) error {
	err := checkKindStructPtr(v)
	if err != nil {
		return err
	}

	// Parse the env file
	// Expected format is ENV_VAR=value, one per line
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	var envVariableValueMap map[string]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return err
		}

		// Only handle items of form VARIABLE=value
		tokens := strings.Split(scanner.Text(), "=")
		if len(tokens) == 2 {
			if envVariableValueMap == nil {
				envVariableValueMap = make(map[string]string)
			}
			envVariableValueMap[tokens[0]] = tokens[1]
		}
	}

	if envVariableValueMap != nil {
		em := parseEnvMetadata(v)

		var fieldNameValueMap map[string]interface{}
		for _, envVariable := range em.envVariables {
			if envValue, ok := envVariableValueMap[envVariable]; ok {
				if fieldNameValueMap == nil {
					fieldNameValueMap = make(map[string]interface{})
				}
				fieldNameValueMap[em.envVariableFieldNameMap[envVariable]] = envValue
			}
		}

		if fieldNameValueMap != nil {
			setValuesFromFieldNameValueMap(v, fieldNameValueMap)
		}
	}

	return nil
}

type envMetadata struct {
	fieldNames              []string
	envVariables            []string
	fieldNameEnvVariableMap map[string]string
	envVariableFieldNameMap map[string]string
}

func parseEnvMetadata(v interface{}) envMetadata {
	// assume v is a pointer to a struct

	rve := reflect.ValueOf(v).Elem()

	var em envMetadata
	for i := 0; i < rve.NumField(); i++ {
		field := rve.Type().Field(i)
		fieldName := field.Name
		envVariable := field.Tag.Get("env")
		if envVariable != "" {
			if em.fieldNameEnvVariableMap == nil {
				em.fieldNameEnvVariableMap = make(map[string]string)
			}
			if em.envVariableFieldNameMap == nil {
				em.envVariableFieldNameMap = make(map[string]string)
			}
			em.fieldNames = append(em.fieldNames, fieldName)
			em.envVariables = append(em.envVariables, envVariable)
			em.fieldNameEnvVariableMap[fieldName] = envVariable
			em.envVariableFieldNameMap[envVariable] = fieldName
		}
	}

	return em
}
