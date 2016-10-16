package agent

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var EmptyJson = []byte("{}")

func WriteResponseJson(w http.ResponseWriter, code int, body []byte) error {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; encoding=utf-8")
	_, err := w.Write(body)
	if err != nil {
		log.Printf("cannot write http response body")
		return err
	}
	return nil
}

func GetTest(r *http.Request) (string, error) {
	value, exists := r.URL.Query()["test"]
	if !exists {
		return "", errors.New("missing test")
	}
	if len(value) != 1 {
		return "", errors.New("too many arguments")
	}
	test := value[0]
	matched, err := regexp.MatchString("^[a-z_]+$", test)
	if err != nil {
		return "", errors.New("regexp does not compile")
	}
	if !matched {
		return "", errors.New("regexp does not match")
	}
	return test, nil
}

func GetOptionalInt(r *http.Request, name string, def_value int) (int, error) {
	value, exists := r.URL.Query()[name]
	if !exists {
		return def_value, nil
	}
	if len(value) != 1 {
		return 0, errors.New("too many arguments")
	}
	single_value := value[0]
	matched, err := regexp.MatchString("^[0-9]+$", single_value)
	if err != nil {
		return 0, errors.New("regexp does not compile")
	}
	if !matched {
		return 0, errors.New("regexp does not match")
	}
	number, err := strconv.Atoi(single_value)
	if err != nil {
		return 0, errors.New("strconv failed")
	}
	return number, nil
}
