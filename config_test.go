package configurator

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
)

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

func TestConfigTypeIsStruct(t *testing.T) {
	configType := &TestConfigType{}

	if reflect.ValueOf(configType).Elem().Kind() != reflect.Struct {
		t.FailNow()
	}
}

func TestKeyNames(t *testing.T) {
	var configType TestConfigType

	keys, err := keyNames(&configType)
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

func TestDefaultTags(t *testing.T) {
	var configType TestConfigType

	keys, err := keyNames(&configType)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	defaults := parseTag(&configType, keys, "default")

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

	keys, err := keyNames(&configType)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	env := parseTag(&configType, keys, "env")

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

	keys, err := keyNames(&configType)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	required := requiredKeys(&configType, keys)

	expectedRequired := map[string]bool{
		"Parameter1": false,
		"Parameter2": true,
		"Parameter3": true,
	}

	if !reflect.DeepEqual(required, expectedRequired) {
		log.Println("unexpected required tags:", required)
		log.Println("expected:", expectedRequired)
		t.FailNow()
	}
}

func TestParseEnvConfig(t *testing.T) {
	_ = os.Setenv("PARAMETER_3", "three")

	var config TestConfigType
	err := ParseEnvConfig(&config)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	log.Println(&config)
}

func TestAWSService(t *testing.T) {
	serviceKey := os.Getenv("AWS_SERVICE_KEY")
	log.Println("SERVICE KEY EMPTY:", serviceKey == "")
}
