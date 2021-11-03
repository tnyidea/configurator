package test

import (
	"github.com/tnyidea/configurator"
	"log"
	"os"
	"testing"
)

type envTestType struct {
	Parameter1 string `json:"parameter1" env:"PARAMETER_1" default:"one" `
	Parameter2 string `json:"parameter2" default:"two" required:"true"`
	Parameter3 string `json:"parameter3" env:"PARAMETER_3" required:"true"`
}

func TestSetValuesFromOsEnv(t *testing.T) {
	_ = os.Setenv("PARAMETER_3", "three")

	var testValue envTestType
	err := configurator.SetValuesFromOsEnv(&testValue)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	log.Println(&testValue)
}

func TestSetValuesFromEnvFile(t *testing.T) {
	_ = os.Unsetenv("PARAMETER_3")

	var testValue envTestType

	err := configurator.SetValuesFromEnvFile(&testValue, "env_test.env")
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	log.Println(&testValue)
}
