package main

import (
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"runtime"
)

type SpecArgument struct {
	DefaultValue string `json:default_value`
	Label        string `json:label`
	Regexp       string `json:regexp`
}

type Spec struct {
	CommandLine string                  `json:command_line`
	Arguments   map[string]SpecArgument `json:arguments`
}

func make_filepath(home string, name string) (string) {
	return filepath.Join(home, "spec", runtime.GOOS, name) + ".json"
}

func do_load(neubot_home string, nettest_name string) (*Spec, error) {
	spec_filepath := make_filepath(neubot_home, nettest_name)
	content, err := ioutil.ReadFile(spec_filepath)
	if err != nil {
		return nil, err
	}
	var spec *Spec = nil
	err = json.Unmarshal(content, &spec)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

func SpecLoad(neubot_home string, nettest_name string) (*Spec, error) {
	spec, err := do_load(neubot_home, nettest_name)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

func SpecRunSync(spec *Spec) error {
	return nil
}
