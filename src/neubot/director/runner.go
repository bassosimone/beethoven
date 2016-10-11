package director

import (
	"encoding/json"
	"errors"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type Runner struct {
	Process    *exec.Cmd
	Status     string
	StdoutPath string
	StderrPath string
	Timestamp  time.Time
	TestName   string
	TestId	   string
	Workdir    string
	CmdLine    []byte
}

func save_cmdline(runner *Runner, cmdline []string) error {
	s, err := json.Marshal(cmdline)
	runner.CmdLine = s
	return err
}

func open_file(runner *Runner, file_type string) (
	*os.File, error) {
	prefix := runner.TestName + "-" + file_type + "-" + runner.TestId
	file, err := ioutil.TempFile(runner.Workdir, prefix)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func RunnerStart(nettest_name string, cmdline []string, workdir string) (
	*Runner, error) {
	var runner Runner
	runner.TestName = nettest_name
	if len(cmdline) == 0 {
		return nil, errors.New("invalid command line")
	}
	err := save_cmdline(&runner, cmdline)
	if err != nil {
		return nil, err
	}
	runner.Process = exec.Command(cmdline[0])
	runner.Process.Args = cmdline
	runner.Workdir = workdir
	test_id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	runner.TestId = test_id.String()
	stdout, err := open_file(&runner, "stdout")
	if err != nil {
		return nil, err
	}
	runner.StdoutPath = stdout.Name()
	runner.Process.Stdout = stdout
	stderr, err := open_file(&runner, "stderr")
	if err != nil {
		return nil, err
	}
	runner.StderrPath = stderr.Name()
	runner.Process.Stderr = stderr
	err = runner.Process.Start()
	if err != nil {
		return nil, err
	}
	runner.Timestamp = time.Now()
	runner.Status = "running"
	return &runner, nil
}

func RunnerWaitAsync(runner *Runner, timeout time.Duration,
	periodic func()) chan error {
	done := make(chan error, 1)
	// See: <http://stackoverflow.com/a/11886829>
	internal := make(chan error, 1)
	go func() {
		internal <- runner.Process.Wait()
	}()
again:
	periodic()
	select {
	case <-time.After(1.0 * time.Second):
		delta := time.Since(runner.Timestamp)
		if delta < timeout {
			goto again
		}
		err := runner.Process.Process.Kill()
		if err != nil {
			// My understanding is that it should not happen, and that ot
			// would also lead to a leak of the process control struct
			runner.Status = "kill-failed"
			err = errors.New("failed to kill process")
		} else {
			err = <-internal // Make sure we don't leak control struct
			runner.Status = "killed"
		}
		done <- err
	case err := <-internal:
		if err == nil {
			runner.Status = "succeded"
		} else {
			runner.Status = "failed"
		}
		done <- err
	}
	periodic()
	return done
}
