package configurator

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
)

// TODO if something wants to be config'd it must have a config tag..
//   of the form config:"env=ENV_VAR,default=two,required"
type TestConfigType struct {
	Parameter1 string `json:"parameter1" env:"PARAMETER_1" default:"one" `
	Parameter2 string `json:"parameter2" default:"two" config:"required"`
	Parameter3 string `json:"parameter3" env:"PARAMETER_3" config:"required"`
}

func (p *TestConfigType) IsZero() bool {
	return reflect.DeepEqual(*p, TestConfigType{})
}

func (p *TestConfigType) Bytes() string {
	byteValue, _ := json.Marshal(p)
	return string(byteValue)
}

func (p *TestConfigType) String() string {
	byteValue, _ := json.MarshalIndent(p, "", "    ")
	return string(byteValue)
}

func TestConfigTypeIsStructCase1(t *testing.T) {
	var configType TestConfigType

	ptr := &configType

	// configType must be a pointer to a struct
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		log.Println("invalid configType: must be a pointer to a struct")
		t.FailNow()
	}
	if reflect.Indirect(reflect.ValueOf(ptr)).Kind() != reflect.Struct {
		log.Println("invalid configType: not a struct")
		t.FailNow()
	}
}

func TestConfigTypeIsStructCase2(t *testing.T) {
	configType := TestConfigType{}

	ptr := &configType

	// configType must be a pointer to a struct
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		log.Println("invalid configType: must be a pointer to a struct")
		t.FailNow()
	}
	if reflect.Indirect(reflect.ValueOf(ptr)).Kind() != reflect.Struct {
		log.Println("invalid configType: not a struct")
		t.FailNow()
	}
}

func TestConfigTypeIsStructCase3(t *testing.T) {
	configType := TestConfigType{
		Parameter1: "One",
		Parameter2: "Two",
		Parameter3: "Three",
	}

	ptr := &configType

	// configType must be a pointer to a struct
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		log.Println("invalid configType: must be a pointer to a struct")
		t.FailNow()
	}
	if reflect.Indirect(reflect.ValueOf(ptr)).Kind() != reflect.Struct {
		log.Println("invalid configType: not a struct")
		t.FailNow()
	}
}

func TestConfigTypeIsStructCase4(t *testing.T) {
	var configType TestConfigType

	ptr := &configType

	// configType must be a pointer to a struct
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		log.Println("invalid configType: must be a pointer to a struct")
		t.FailNow()
	}
	if reflect.Indirect(reflect.ValueOf(ptr)).Kind() != reflect.Struct {
		log.Println("invalid configType: not a struct")
		t.FailNow()
	}
}

func TestConfigTypeIsStructCase5(t *testing.T) {
	configType := TestConfigType{
		Parameter1: "One",
		Parameter2: "Two",
		Parameter3: "Three",
	}

	ptr := &configType

	// configType must be a pointer to a struct
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		log.Println("invalid configType: must be a pointer to a struct")
		t.FailNow()
	}
	if reflect.Indirect(reflect.ValueOf(ptr)).Kind() != reflect.Struct {
		log.Println("invalid configType: not a struct")
		t.FailNow()
	}

	//if reflect.ValueOf(ptr).Elem().Kind() != reflect.Struct {
	//	t.FailNow()
	//}
}

func TestKeyNames(t *testing.T) {
	configType := TestConfigType{
		Parameter1: "One",
		Parameter2: "Two",
		Parameter3: "Three",
	}

	keys, _, err := fieldValueMap(&configType)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	expectedKeys := []string{"Parameter1", "Parameter2", "Parameter3"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		log.Println("unexpected keys:", keys)
		log.Println("expected:", expectedKeys)
		t.FailNow()
	}
}

func TestKeyNamesEmptyConfigType(t *testing.T) {
	var configType TestConfigType

	keys, _, err := fieldValueMap(&configType)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	expectedKeys := []string{"Parameter1", "Parameter2", "Parameter3"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		log.Println("unexpected keys:", keys)
		log.Println("expected:", expectedKeys)
		t.FailNow()
	}
}

func TestFieldValueMap(t *testing.T) {
	configType := TestConfigType{
		Parameter1: "One",
		Parameter2: "Two",
		Parameter3: "Three",
	}

	expectedMap := map[string]string{
		"Parameter1": "One",
		"Parameter2": "Two",
		"Parameter3": "Three",
	}

	_, configMap, err := fieldValueMap(&configType)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	log.Println(configMap)
	log.Println(expectedMap)
}

func TestFieldValueMapEmptyConfigType(t *testing.T) {
	var configType TestConfigType

	expectedMap := map[string]string{
		"Parameter1": "One",
		"Parameter2": "Two",
		"Parameter3": "Three",
	}

	_, configMap, err := fieldValueMap(&configType)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	log.Println(configMap)
	log.Println(expectedMap)
}

func TestDefaultTags(t *testing.T) {
	var configType TestConfigType

	_, defaults := parseTagValues(&configType, "default")

	expectedDefaults := map[string]string{
		"Parameter1": "one",
		"Parameter2": "two",
	}

	if !reflect.DeepEqual(defaults, expectedDefaults) {
		log.Println("unexpected defaults:", defaults)
		log.Println("expected:", expectedDefaults)
		t.FailNow()
	}
}

func TestEnvTags(t *testing.T) {
	var configType TestConfigType

	_, env := parseTagValues(&configType, "env")

	expectedEnv := map[string]string{
		"Parameter1": "PARAMETER_1",
		"Parameter3": "PARAMETER_3",
	}

	if !reflect.DeepEqual(env, expectedEnv) {
		log.Println("unexpected env tags:", env)
		log.Println("expected:", expectedEnv)
		t.FailNow()
	}
}

func TestRequiredKeys(t *testing.T) {
	var configType TestConfigType

	requiredFields := requiredFieldMap(&configType)

	expectedRequiredFields := map[string]bool{
		"Parameter2": true,
		"Parameter3": true,
	}

	if !reflect.DeepEqual(requiredFields, expectedRequiredFields) {
		log.Println("unexpected required fields:", requiredFields)
		log.Println("expected:", expectedRequiredFields)
		t.FailNow()
	}
}

func TestParseEnvConfig(t *testing.T) {
	_ = os.Setenv("PARAMETER_3", "three")

	var configType TestConfigType
	err := ParseEnvConfig(&configType)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	log.Println(&configType)
}

func TestParseEnvConfigFile(t *testing.T) {
	_ = os.Unsetenv("PARAMETER_3")

	err := SetEnvFromFile("config_test.env")
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	var configType TestConfigType
	err = ParseEnvConfig(&configType)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	log.Println(&configType)
}

func TestValidateConfig(t *testing.T) {
	_ = os.Setenv("PARAMETER_3", "three")

	var configType TestConfigType

	err := ParseEnvConfig(&configType)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	err = ValidateConfig(&configType)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	log.Println(&configType)
}
