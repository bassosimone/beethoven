package director

import (
	"io"
	"log"
	"neubot/common"
	"os"
)

type Director struct {
	NeubotHome string
}

func New(neubot_home string) *Director {
	var dir Director
	dir.NeubotHome = neubot_home
	return &dir
}

var mapping map[string]*Director = make(map[string]*Director)

func Get(neubot_home string) *Director {
	dir, okay := mapping[neubot_home]
	if !okay {
		dir = New(neubot_home)
		mapping[neubot_home] = dir
	}
	return dir
}

func (self Director) Start(nettest_name string,
		arguments map[string]string) (*Runner, error) {
	log.Printf("neubot_home: %s\n", self.NeubotHome)
	log.Printf("nettest_name: %s\n", nettest_name)
	log.Printf("arguments: %s\n", arguments)
	spec, err := SpecLoad(self.NeubotHome, nettest_name)
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

func (self Director) WaitAsync(runner *Runner, callback func()) chan error {
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

func (self Director) Run(nettest_name string,
		arguments map[string]string) error {
	runner, err := self.Start(nettest_name, arguments)
	if err != nil {
		return err
	}
	stderr, err := StreamingOpenStderr(runner)
	if err != nil {
		return err
	}
	channel := self.WaitAsync(runner, func() {
		StreamingForward(stderr, os.Stderr)
	})
	err = <-channel
	stderr.Close()
	return err
}

func (Director) OpenStderr(runner *Runner) (*os.File, error) {
	return StreamingOpenStderr(runner)
}

func (Director) Forward(filep *os.File, writer io.Writer) error {
	return StreamingForward(filep, writer)
}
