package agent

import (
	"fmt"
	"github.com/go-martini/martini"
)

func Run(address string, port int, verbose bool) {
	m := martini.Classic()

	m.Group("/api", func(r martini.Router) {
		r.Get("/runner", ApiRunnerGet);
	})

	endpoint := fmt.Sprintf("%s:%d", address, port)
	m.RunOnAddr(endpoint)
}
