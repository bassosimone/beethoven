package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"text/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

type SpecArgument struct {
	DefaultValue string
	Label        string
	Regexp       string
}

type Spec struct {
	CommandLine string
	Arguments   map[string]SpecArgument
}

func make_filepath(home string, name string) (string) {
	return filepath.Join(home, "spec", runtime.GOOS, name) + ".json"
}

func SpecLoad(neubot_home string, nettest_name string) (Spec, error) {
	var spec Spec
	spec_filepath := make_filepath(neubot_home, nettest_name)
	content, err := ioutil.ReadFile(spec_filepath)
	if err != nil {
		return spec, err
	}
	err = json.Unmarshal(content, &spec)
	if err != nil {
		return spec, err
	}
	return spec, nil
}

func SpecCmdline(spec Spec, arguments map[string]string) (string, error) {

	for skey, svalue := range spec.Arguments {
		if _, found := arguments[skey]; !found {
			log.Printf("Set %s to '%s' (default)\n", skey, svalue.DefaultValue)
			arguments[skey] = svalue.DefaultValue
		}
	}

	for akey, avalue := range arguments {
		log.Printf("checking key %s, value %s\n", akey, avalue)
		svalue, found := spec.Arguments[akey]
		if !found {
			return "", errors.New("unmapped argument: " + akey)
		}
		log.Printf("validating regexp: '%s'", svalue.Regexp)
		if svalue.Regexp == "" {
			return "", errors.New("missing regexp for: " + akey)
		}
		matched, err := regexp.MatchString(svalue.Regexp, avalue)
		if err != nil {
			return "", err
		}
		if !matched {
			return "", errors.New("regexp does not match for: " + akey)
		}
		log.Printf("Argument okay")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	arguments["cwd"] = cwd
	arguments["python"] = DefaultPython()

	log.Printf("cmdline: %s\n", spec.CommandLine)
	tmpl, err := template.New("Main").Parse(spec.CommandLine)
	if err != nil {
		return "", err
	}
	tmpl.Option("missingkey=error")
	output := bytes.NewBufferString("")
	err = tmpl.Execute(output, arguments)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func SpecRunSync(spec Spec) error {
	return nil
}
