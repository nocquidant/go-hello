package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/google/logger"
	"github.com/nocquidant/go-hello/env"
)

func kvAsJson(key string, val string) string {
	m := make(map[string]string)
	m[key] = val
	data, err := json.Marshal(m)
	if err != nil {
		logger.Errorf("Error while serializing to json: %s", err)
		return ""
	}
	return string(data)
}

func kmAsJson(key string, v map[string]interface{}) string {
	m := make(map[string]interface{})
	m[key] = v
	data, err := json.Marshal(m)
	if err != nil {
		logger.Errorf("Error while serializing to json: %s", err)
		return ""
	}
	return string(data)
}

func mapAsJson(m map[string]interface{}) string {
	data, err := json.Marshal(m)
	if err != nil {
		logger.Errorf("Error while serializing to json: %s", err)
		return ""
	}
	return string(data)
}

func HandlerHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, kvAsJson("health", "UP"))
}

func HandlerHello(w http.ResponseWriter, r *http.Request) {
	h, _ := os.Hostname()
	m := make(map[string]interface{})
	m["msg"] = fmt.Sprintf("My name is '%s', I'm served from '%s'", env.NAME, h)

	fmt.Fprintf(w, mapAsJson(m))
}

func HandlerInfo(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["version"] = env.VERSION
	m["build"] = env.GITCOMMIT
	m["environment"] = map[string]interface{}{
		"name":      env.NAME,
		"port":      env.PORT,
		"remoteUrl": env.REMOTE_URL,
	}

	fmt.Fprintf(w, mapAsJson(m))
}

func HandlerRequest(w http.ResponseWriter, r *http.Request) {
	// Build the request
	req, err := http.NewRequest("GET", "http://"+env.REMOTE_URL, nil)
	if err != nil {
		http.Error(w, "Error while building request", http.StatusBadRequest)
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
		http.Error(w, "Error while requesting backend", http.StatusServiceUnavailable)
		logger.Errorf("Error while requesting backend: %s", err)
		return
	}

	// Callers should close resp.Body
	defer resp.Body.Close()

	// Get body as string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			http.Error(w, "Error while getting body", http.StatusInternalServerError)
			logger.Errorf("Error while getting body: %s", err)
			return
		}
		var x map[string]interface{}
		err := json.Unmarshal(bodyBytes, &x)
		if err != nil {
			http.Error(w, "Error while unmarshalling response from remote", http.StatusInternalServerError)
			logger.Errorf("Error while unmarshalling response from remote: %s", err)
		}
		fmt.Fprintf(w, kmAsJson("fromRemote", x))
	} else {
		fmt.Fprintf(w, "Error while calling the back: %d", resp.StatusCode)
	}
}
