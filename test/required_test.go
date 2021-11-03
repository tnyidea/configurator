package test

import (
	"github.com/tnyidea/configurator"
	"log"
	"testing"
)

type requiredTestType struct {
	Parameter1 string `json:"parameter1" env:"PARAMETER_1" default:"one" `
	Parameter2 string `json:"parameter2" default:"two" required:"true"`
	Parameter3 string `json:"parameter3" env:"PARAMETER_3" required:"true"`
}

func TestCheckRequiredValuesNotSet(t *testing.T) {
	var rtv requiredTestType

	err := configurator.CheckRequiredValues(&rtv)
	if err == nil {
		log.Println("test error: empty value ")
		t.FailNow()
	}

	log.Println(rtv)
}

func TestCheckRequiredValuesSet(t *testing.T) {
	rtv := requiredTestType{
		Parameter1: "",
		Parameter2: "two",
		Parameter3: "three",
	}

	err := configurator.CheckRequiredValues(&rtv)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	log.Println(rtv)
}
