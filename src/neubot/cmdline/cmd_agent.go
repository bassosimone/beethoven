package cmdline

import (
	"fmt"
	"github.com/pborman/getopt"
	"neubot/agent"
	"os"
)

const agent_usage = `usage:
  neubot agent [-v] [-A address] [-p port]
  neubot agent -h|--help`

func CmdAgentMain() {
	address := getopt.String('A', "127.0.0.1", "Set listen address")
	display_help := getopt.BoolLong("help", 'h', "Display help")
	port := getopt.Int('p', 9774, "Set default port")
	verbose := getopt.Bool('v', "Be verbose")

	os.Args = os.Args[1:]
	if err := getopt.Getopt(nil); err != nil {
		fmt.Printf("%s\n", agent_usage)
		os.Exit(1)
	}
	optarg := getopt.Args()
	if len(optarg) != 0 {
		fmt.Printf("%s\n", agent_usage)
		os.Exit(1)
	}

	if *display_help {
		fmt.Printf("%s\n", agent_usage)
		os.Exit(0)
	}

	agent.Run(*address, *port, *verbose)
}
