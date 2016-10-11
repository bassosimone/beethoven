package main

import (
	"log"
)

func DirectorRun(neubot_home string, nettest_name string,
		arguments map[string]string) error {
	log.Printf("neubot_home: %s\n", neubot_home)
	log.Printf("nettest_name: %s\n", nettest_name)
	log.Printf("arguments: %s\n", arguments)
	spec, err := SpecLoad(neubot_home, nettest_name)
	if err != nil {
		return err
	}
	cmdline, err := SpecCmdline(spec, arguments)
	if err != nil {
		return err
	}
	log.Printf("cmdline: %s\n", cmdline)
	runner, err := RunnerStart(nettest_name, cmdline, DefaultWorkdir())
	if err != nil {
		return err
	}
	log.Printf("command running")
	channel := RunnerWait(runner, DefaultProcTimeout())
	return <-channel
}
