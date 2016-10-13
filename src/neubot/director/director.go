package director

import (
	"log"
	"neubot/common"
	"os"
)

func DirectorStart(neubot_home string, nettest_name string,
	arguments map[string]string) (*Runner, error) {
	log.Printf("neubot_home: %s\n", neubot_home)
	log.Printf("nettest_name: %s\n", nettest_name)
	log.Printf("arguments: %s\n", arguments)
	spec, err := SpecLoad(neubot_home, nettest_name)
	if err != nil {
		return nil, err
	}
	cmdline, err := SpecCmdline(spec, arguments)
	if err != nil {
		return nil, err
	}
	log.Printf("cmdline: %s\n", cmdline)
	runner, err := RunnerStart(nettest_name, cmdline, common.DefaultWorkdir())
	if err != nil {
		return nil, err
	}
	log.Printf("command running")
	return runner, nil
}

func DirectorWaitAsync(runner *Runner, callback func()) chan error {
	channel := make(chan error, 1)
	go func() {
		err := <-RunnerWaitAsync(runner, common.DefaultProcTimeout(), callback)
		if err != nil {
			log.Printf("Command failed: %s\n", err)
		}
		err2 := MeasurementsAppend(&runner.M)
		if err == nil && err2 != nil {
			err = err2
		}
		channel <- err
	}()
	return channel
}

func DirectorRun(neubot_home string, nettest_name string,
	arguments map[string]string) error {
	runner, err := DirectorStart(neubot_home, nettest_name, arguments)
	if err != nil {
		return err
	}
	stderr, err := StreamingOpenStderr(runner)
	if err != nil {
		return err
	}
	channel := DirectorWaitAsync(runner, func() {
		StreamingForward(stderr, os.Stderr)
	})
	err = <-channel
	stderr.Close()
	return err
}
