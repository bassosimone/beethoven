package main

import (
	"errors"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"log"
	"os/exec"
	"os"
	"time"
)

func open_file(workdir string, nettest_name string, file_type string) (
		*os.File, error) {
	uuid4, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	prefix := "neubot-" + nettest_name + "-" + file_type + "-" + uuid4.String()
	file, err := ioutil.TempFile(workdir, prefix)
	if err != nil {
		return nil, err
	}
	return file, nil
}

type Runner struct {
	Status string
	Error error
	Process *exec.Cmd
}

func RunnerStart(nettest_name string, cmdline []string, workdir string) (
		Runner, error) {
	var runner Runner
	if len(cmdline) == 0 {
		return runner, errors.New("invalid command line")
	}
	runner.Process = exec.Command(cmdline[0])
	runner.Process.Args = cmdline
	stdout, err := open_file(workdir, nettest_name, "stdout")
	if err != nil {
		return runner, err
	}
	runner.Process.Stdout = stdout
	stderr, err := open_file(workdir, nettest_name, "stderr")
	if err != nil {
		return runner, err
	}
	runner.Process.Stderr = stderr
	err = runner.Process.Start()
	if err != nil {
		return runner, err
	}
	runner.Status = "running"
	return runner, nil
}

func RunnerWait(runner Runner, timeout time.Duration) (chan error) {
	done := make(chan error, 1)
	// See: <http://stackoverflow.com/a/11886829>
	internal := make(chan error, 1)
	go func() {
		internal <- runner.Process.Wait()
	}()
	select {
	case <-time.After(timeout):
		log.Printf("kill process that is running for too long\n")
		err := runner.Process.Process.Kill()
		if (err != nil) {
			// My understanding is that it should not happen, and that ot
			// would also lead to a leak of the process control struct
			runner.Status = "kill-failed"
			err = errors.New("failed to kill process")
		} else {
			err = <-internal // Make sure we don't leak control struct
			runner.Status = "killed"
		}
		done <-err
	case err := <-internal:
		log.Printf("Process terminated; exit code: %s\n", err)
		runner.Status = "terminated"
		done <-err
	}
	return done
}
