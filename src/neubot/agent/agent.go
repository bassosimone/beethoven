package agent

import (
	"fmt"
	"github.com/go-martini/martini"
)

func Run(address string, port int, verbose bool) {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	endpoint := fmt.Sprintf("%s:%d", address, port)
	m.RunOnAddr(endpoint)
}
