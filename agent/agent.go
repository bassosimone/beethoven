package agent

import (
	"fmt"
	"net/http"
)

func Run(address string, port int, verbose bool) error {
	http.HandleFunc("/api/runner", ApiRunner)
	http.HandleFunc("/api/state", DefaultStateHandler().ServeHTTP)
	endpoint := fmt.Sprintf("%s:%d", address, port)
	return http.ListenAndServe(endpoint, nil)
}
