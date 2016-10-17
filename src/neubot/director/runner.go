package director

import (
	"errors"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type Runner struct {
	Process    *exec.Cmd
	M          Measurement
}

func open_file(runner *Runner, file_type string) (
	*os.File, error) {
	prefix := runner.M.TestName + "-" + file_type + "-" + runner.M.TestId
	file, err := ioutil.TempFile(runner.M.Workdir, prefix)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func RunnerStart(nettest_name string, cmdline []string, workdir string) (
		*Runner, error) {
	var runner Runner
	runner.M.TestName = nettest_name
	if len(cmdline) == 0 {
		return nil, errors.New("invalid command line")
	}
	runner.M.CmdLine = cmdline
	runner.Process = exec.Command(cmdline[0])
	runner.Process.Args = cmdline
	runner.M.Workdir = workdir
	test_id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	runner.M.TestId = test_id.String()
	stdout, err := open_file(&runner, "stdout")
	if err != nil {
		return nil, err
	}
	runner.M.StdoutPath = stdout.Name()
	runner.Process.Stdout = stdout
	stderr, err := open_file(&runner, "stderr")
	if err != nil {
		return nil, err
	}
	runner.M.StderrPath = stderr.Name()
	runner.Process.Stderr = stderr
	err = runner.Process.Start()
	if err != nil {
		return nil, err
	}
	runner.M.Timestamp = time.Now()
	runner.M.Status = "running"
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
		delta := time.Since(runner.M.Timestamp)
		if delta < timeout {
			goto again
		}
		err := runner.Process.Process.Kill()
		if err != nil {
			// My understanding is that it should not happen, and that ot
			// would also lead to a leak of the process control struct
			runner.M.Status = "kill-failed"
			err = errors.New("failed to kill process")
		} else {
			err = <-internal // Make sure we don't leak control struct
			runner.M.Status = "killed"
		}
		done <- err
	case err := <-internal:
		if err == nil {
			runner.M.Status = "succeded"
		} else {
			runner.M.Status = "failed"
		}
		done <- err
	}
	// Call periodic one last time to read what may have been written
	// before exiting and therefore could still be in the buffers
	periodic()
	return done
}
