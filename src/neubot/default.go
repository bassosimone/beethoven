package main

import (
	"os/exec"
	"path/filepath"
)

func DefaultNeubotHome() string {
	return filepath.Join("var", "lib", "neubot")
}

func DefaultPython() (string, error) {
	return exec.LookPath("python")
}
