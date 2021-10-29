package test

import (
	"encoding/json"
	"reflect"
)

type testType struct {
	Parameter1 string `json:"parameter1" env:"PARAMETER_1" default:"one" `
	Parameter2 string `json:"parameter2" default:"two" required:"true"`
	Parameter3 string `json:"parameter3" env:"PARAMETER_3" required:"true"`
}

func (p *testType) IsZero() bool {
	return reflect.ValueOf(*p).IsZero()
}

func (p *testType) Bytes() string {
	byteValue, _ := json.Marshal(p)
	return string(byteValue)
}

func (p *testType) String() string {
	byteValue, _ := json.MarshalIndent(p, "", "    ")
	return string(byteValue)
}
