package main

import (
	"fmt"
	"math/rand"
	"github.com/bassosimone/beethoven/cmdline"
	"github.com/bassosimone/beethoven/common"
	"os"
	"time"
)

const usage = `usage:
  neubot command [options...]
  neubot --help
  neubot --version
commands:
  agent - runs neubot's api
  run - runs a specific neubot test from command line`

func main() {
	// Make sure we seed the random number generator properly
	//   see <http://stackoverflow.com/a/12321192>
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) < 2 {
		fmt.Printf("%s\n", usage)
		os.Exit(0)
	}

	if len(os.Args) == 2 {
		if os.Args[1] == "--help" {
			fmt.Printf("%s\n", usage)
			os.Exit(0)
		}
		if os.Args[1] == "--version" {
			fmt.Printf("%s\n", common.Version)
			os.Exit(0)
		}
		// FALLTHROUGH
	}

	if os.Args[1] == "agent" {
		cmdline.CmdAgentMain()
		os.Exit(0)
	}

	if os.Args[1] == "run" {
		cmdline.CmdRunMain()
		os.Exit(0)
	}

	fmt.Printf("%s\n", usage)
	os.Exit(1)
}
