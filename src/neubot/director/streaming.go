package director

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

func StreamingOpenStderr(runner *Runner) (*os.File, error) {
	filep, err := os.Open(runner.M.StderrPath)
	if err != nil {
		log.Printf("cannot open: %s", runner.M.StderrPath)
		return nil, err
	}
	return filep, nil
}

func StreamingForward(filep *os.File, writer io.Writer) error {
	buffer, err := ioutil.ReadAll(filep)
	if err != nil {
		return err
	}
	_, err = writer.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}
