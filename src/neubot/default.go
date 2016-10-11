package main

import (
	"path/filepath"
)

func DefaultNeubotHome() string {
	return filepath.Join("var", "lib", "neubot")
}
