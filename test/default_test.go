package test

import (
	"encoding/json"
	"github.com/tnyidea/configurator"
	"log"
	"reflect"
	"testing"
)

type defaultTestType struct {
	Parameter1 string `json:"parameter1" env:"PARAMETER_1" default:"one" `
	Parameter2 string `json:"parameter2" default:"two" required:"true"`
	Parameter3 string `json:"parameter3" env:"PARAMETER_3" required:"true"`
}

func (p *defaultTestType) IsZero() bool {
	return reflect.ValueOf(*p).IsZero()
}

func (p *defaultTestType) Bytes() string {
	byteValue, _ := json.Marshal(p)
	return string(byteValue)
}

func (p *defaultTestType) String() string {
	byteValue, _ := json.MarshalIndent(p, "", "    ")
	return string(byteValue)
}

func TestSetAllDefaultValues(t *testing.T) {
	var testValue defaultTestType

	err := configurator.SetAllDefaultValues(&testValue)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	log.Println(&testValue)
}

func TestSetEmptyFieldDefaultValues(t *testing.T) {
	testValue := defaultTestType{
		Parameter1: "setAlready",
	}

	err := configurator.SetEmptyFieldDefaultValues(&testValue)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	log.Println(&testValue)
}
