package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/google/logger"
	"github.com/nocquidant/go-hello/env"
)

func writeError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	io.WriteString(w, kvAsJson("error", msg))
}

func HandlerHealth(w http.ResponseWriter, r *http.Request) {
	// do not fill the logs here

	io.WriteString(w, kvAsJson("health", "UP"))
}

func HandlerHello(w http.ResponseWriter, r *http.Request) {
	logger.Infof("%s request to %s\n", r.Method, r.RequestURI)

	h, _ := os.Hostname()
	m := make(map[string]interface{})
	m["msg"] = fmt.Sprintf("Hello, my name is '%s', I'm served from '%s'", env.NAME, h)
	io.WriteString(w, mapAsJson(m))
}

func HandlerRemote(w http.ResponseWriter, r *http.Request) {
	logger.Infof("%s request to %s\n", r.Method, r.RequestURI)

	// Build the request
	req, err := http.NewRequest("GET", "http://"+env.REMOTE_URL, nil)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Error while building request")
		logger.Errorf("Error while building request: %s", err)
		return
	}

	// A Client is an HTTP client
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	// Send the request via a client
	resp, err := client.Do(req)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "Error while requesting backend")
		logger.Errorf("Error while requesting backend: %s", err)
		return
	}

	// Callers should close resp.Body
	defer resp.Body.Close()

	// Get body as string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			writeError(w, http.StatusInternalServerError, "Error while getting body from remote")
			logger.Errorf("Error while getting body: %s", err)
			return
		}
		var x map[string]interface{}
		err := json.Unmarshal(bodyBytes, &x)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Error while unmarshalling response from remote")
			logger.Errorf("Error while unmarshalling response from remote: %s", err)
		}
		io.WriteString(w, kmAsJson("fromRemote", x))
	} else {
		io.WriteString(w, fmt.Sprintf("Error while calling the back: %d", resp.StatusCode))
	}
}
