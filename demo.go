package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Demo struct {
	Title     string
	Variables DemoVariables
	Setup     DemoStepList
	Steps     DemoStepList
	Teardown  DemoStepList
}

type DemoStepList []DemoStep

type DemoVariables []DemoVariable

type DemoVariable struct {
	Name  string
	Value string
	Cmd   string
	Echo  bool
}

type DemoStep struct {
	Comment    string
	Cmd        string
	Background bool
	Echo       bool
	Capture    string
}

func (s *DemoStep) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawStep DemoStep
	raw := rawStep{
		Echo: true,
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*s = DemoStep(raw)
	return nil
}

func (v *DemoVariable) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawVar DemoVariable
	raw := rawVar{
		Echo: true,
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*v = DemoVariable(raw)
	return nil
}

func LoadDemoFile(filename string) (Demo, error) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return Demo{}, err
	}

	var demo Demo
	err = yaml.Unmarshal(data, &demo)
	return demo, err
}

func (v DemoVariables) mapping(key string) string {
	for _, entry := range v {
		if entry.Name == key {
			return entry.Value
		}
	}

	return ""
}

func (v DemoVariables) sanitize(key string) string {
	for _, entry := range v {
		if entry.Name == key {
			if entry.Echo == false {
				return fmt.Sprintf("${%s}", key)
			} else {
				return entry.Value
			}
		}
	}

	return ""
}

func RunDemo(steps DemoStepList, vars DemoVariables) error {
	action := DisplayStepN(0)
	var ok bool
	var err error

	for {
		action, ok, err = action(steps, vars)

		if err != nil {
			return err
		}

		if !ok {
			break
		}
	}

	return nil
}
