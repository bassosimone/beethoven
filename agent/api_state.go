package agent

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type StateCtx struct {
	T       int64  `json:"t"`
	Pid     int    `json:"pid"`
	Since   int64  `json:"since"`
}

type StateHandler struct {
	Internal StateCtx
	Channel  chan bool
}

func NewStateHandler() *StateHandler {
	var handler StateHandler
	handler.Internal.T = 1 // So the client starts outdated
	handler.Internal.Pid = os.Getpid()
	handler.Internal.Since = time.Now().Unix()
	handler.Channel = make(chan bool)
	return &handler
}

var default_state_handler *StateHandler = nil

func DefaultStateHandler() *StateHandler {
	if default_state_handler == nil {
		default_state_handler = NewStateHandler()
	}
	return default_state_handler
}

func (self *StateHandler) WaitForChanges() {
	select {
		case <-self.Channel:
			log.Printf("received state change from channel")
		case <-time.After(15.0 * time.Second):
			log.Printf("timeout waiting for state change")
	}
}

func (self *StateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	latest, err := GetOptionalInt64(r, "t", 0)
	if err != nil {
		WriteResponseJson(w, 500, EmptyJson)
		return
	}
	log.Printf("/api/state: t = %d", latest)

	debugging, err := GetOptionalInt64(r, "debug", 0)
	if err != nil {
		WriteResponseJson(w, 500, EmptyJson)
		return
	}

	if latest >= self.Internal.T {
		self.WaitForChanges()
	}

	data, err := func () ([]byte, error) {
		if (debugging != 0) {
			return json.MarshalIndent(self.Internal, "", "    ")
		} else {
			return json.Marshal(self.Internal)
		}
	}()
	if err != nil {
		WriteResponseJson(w, 500, EmptyJson)
		return
	}

	WriteResponseJson(w, 200, data)
}
