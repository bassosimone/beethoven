package agent

import (
	"log"
	"net/http"
	"neubot/common"
	"neubot/director"
)

func ApiRunnerGet(r *http.Request) (int, string) {
	test_name, err := GetTest(r)
	if err != nil {
		return 500, "{}"
	}
	log.Printf("test name: %s", test_name)

	settings := make(map[string]string)
	// TODO: check whether we have settings in the body

	// TODO: make sure no more than one test can run at once

	runner, err := director.DirectorStart(common.DefaultNeubotHome(),
		test_name, settings);
	if err != nil {
		return 500, "{}"
	}


	go func() {
		channel := director.DirectorWaitAsync(runner, func() {
			// XXX: not simple to do streaming here, perhaps it would otherwise
			// make sense to do streaming using another API
		})
		_ = <-channel
	}()
	return 200, "{}"
}
