package common

import (
	"os/exec"
	"path/filepath"
	"time"
)

func DefaultNeubotHome() string {
	return filepath.Join("var", "lib", "neubot")
}

func DefaultPython() (string, error) {
	return exec.LookPath("python")
}

func DefaultWorkdir() string {
	return filepath.Join(DefaultNeubotHome(), "data")
}

func DefaultMeasurementsDb() string {
	return filepath.Join(DefaultNeubotHome(), "measurements.sqlite")
}

func DefaultProcTimeout() time.Duration {
	return 60.0 * time.Second
}
