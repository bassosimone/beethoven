package agent

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"neubot/common"
	"neubot/director"
)

const MaximumBodyLength = 1024 * 1024 * 1024

func ApiRunnerGet(w http.ResponseWriter, r *http.Request) {
	test_name, err := GetTest(r)
	if err != nil {
		WriteResponseJson(w, 500, EmptyJson)
		return
	}
	log.Printf("test name: %s", test_name)

	streaming, err := GetOptionalInt64(r, "streaming", 0)
	if err != nil {
		WriteResponseJson(w, 500, EmptyJson)
		return
	}
	log.Printf("streaming: %s", streaming)

	settings := make(map[string]string)
	reader := http.MaxBytesReader(w, r.Body, MaximumBodyLength)
	request_body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("cannot read request body")
		WriteResponseJson(w, 500, EmptyJson)
		return
	}
	if len(request_body) > 0 {
		err := json.Unmarshal(request_body, &settings)
		if err != nil {
			log.Printf("cannot unmarshal request body")
			WriteResponseJson(w, 500, EmptyJson)
			return
		}
	}

	dir := director.Get(common.DefaultNeubotHome())

	runner, err := dir.Start(test_name, settings);
	if err != nil {
		log.Printf("cannot start the selected test")
		WriteResponseJson(w, 500, EmptyJson)
		return
	}

	if streaming == 0 {
		go func() { <-dir.WaitAsync(runner, func() {}) }()
		WriteResponseJson(w, 200, EmptyJson)
		return
	}

	// Implementation of test's standard error streaming

	stderr, err := dir.OpenStderr(runner)
	if err != nil {
		log.Printf("cannot open test standard error")
		go func() { <-dir.WaitAsync(runner, func() {}) }()
		WriteResponseJson(w, 500, EmptyJson)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain; encoding=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")

	log.Printf("start streaming test stderr")
	channel := dir.WaitAsync(runner, func() {
		dir.Forward(stderr, w)
		// Note: be cautious here because it's not granted that all
		// available response handler would be flushers. Note that
		// this implies that, if the handler does not implement the
		// flusher interface then output would not be seen by the
		// client until a sufficient amount of bytes has been sent.
		if flusher, okay := w.(http.Flusher); okay {
			flusher.Flush()
		}
	})
	<-channel
	log.Printf("end streaming test stderr")
	stderr.Close()
}
