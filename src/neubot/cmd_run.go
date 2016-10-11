package main

import (
	"fmt"
	"github.com/pborman/getopt"
	"log"
	"os"
	"strings"
)

const run_usage = `usage:
  neubot run [-v] [-D key=value] [-d neubot_dir] nettest_name
  neubot run -h|--help`

func CmdRunMain() {
	properties := getopt.List('D', "Set test specific properties")
	neubot_home := getopt.String('d', DefaultNeubotHome(), "Set Neubot home")
	verbose := getopt.Bool('v', "Be verbose")
	display_help := getopt.BoolLong("help", 'h', "Display help")

	os.Args = os.Args[1:]
	if err:= getopt.Getopt(nil); err != nil {
		fmt.Printf("%s\n", run_usage)
		os.Exit(1)
	}
	optarg := getopt.Args()
	if len(optarg) != 1 {
		fmt.Printf("%s\n", run_usage)
		os.Exit(1)
	}
	nettest_name := optarg[0]

	if *verbose {
		log.Println("Running in verbose mode")
	}

	properties_map := make(map[string]string)
	for _, elem := range *properties {
		result := strings.SplitN(elem, "=", 2)
		if len(result) != 2 {
			fmt.Printf("invalid property: %s\n", elem)
			fmt.Printf("%s\n", run_usage)
			os.Exit(1)
		}
		properties_map[result[0]] = result[1]
	}

	if *display_help {
		fmt.Printf("%s\n", run_usage)
		os.Exit(0)
	}

	DirectorRun(*neubot_home, nettest_name)
}
