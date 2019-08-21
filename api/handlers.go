package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/nocquidant/go-hello/env"
	logger "github.com/sirupsen/logrus"
)

func writeError(w http.ResponseWriter, statusCode int, msg string) {
	data, err := json.Marshal(ErrorResponse{Code: statusCode, Error: msg})
	if err != nil {
		w.WriteHeader(statusCode)
		io.WriteString(w, fmt.Sprintf("Error while building response: %s", err))
		logger.Errorf("Error while building response: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	io.WriteString(w, string(data))
}

func writeJSON(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(json))
}

// HandlerHealth handles the 'health/' http endpoint
func HandlerHealth(w http.ResponseWriter, r *http.Request) {
	// This fuction is frequently used by K8S
	// -> do not fill the logs, do not record metrics neither

	data, err := json.Marshal(HealthResponse{"UP"})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error while marshalling HealthResponse")
	} else {
		writeJSON(w, data)
	}
}

// HandlerHello handles the 'hello/' http endpoint
func HandlerHello(w http.ResponseWriter, r *http.Request) {
	logger.Infof("%s request to %s\n", r.Method, r.RequestURI)

	h, _ := os.Hostname()
	hello := fmt.Sprintf("Hello, my name is '%s', I'm served from '%s'", env.NAME, h)
	data, err := json.Marshal(MsgResponse{hello})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error while marshalling MsgResponse")
	} else {
		writeJSON(w, data)
	}
}

// HandlerRemote handles the 'remote/' http endpoint
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

	if resp.StatusCode == http.StatusOK {
		// Get body as string
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			writeError(w, http.StatusInternalServerError, "Error while getting body from remote")
			logger.Errorf("Error while getting body: %s", err)
			return
		}

		// Get body as struct MsgResponse
		var respRemote MsgResponse
		err := json.Unmarshal(bodyBytes, &respRemote)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Error while unmarshalling response from remote")
			logger.Errorf("Error while unmarshalling response from remote: %s", err)
		}

		// Create struct for response
		h, _ := os.Hostname()
		respCurrent := MsgRemoteResponse{
			Msg:        fmt.Sprintf("Hello, my name is '%s', I'm served from '%s'", env.NAME, h),
			FromRemote: respRemote,
		}

		// Write response as Json
		data, err := json.Marshal(respCurrent)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Error while marshalling MsgRemoteResponse")
		} else {
			writeJSON(w, data)
		}
	} else {
		// Write error
		writeError(w, resp.StatusCode, "Error while calling the backend app")
	}
}
