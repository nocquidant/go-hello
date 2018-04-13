package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func handlerHello(w http.ResponseWriter, r *http.Request) {
	h, _ := os.Hostname()
	fmt.Fprintf(w, "Hi there, I'm served from %s!\n", h)
}

func handlerInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL to back: %s!\n", backURL)
}

func backURL() string {
	name := os.Getenv("HOSTNAME")
	if name == "" {
		name = "hello-back-svc"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8485"
	}
	return "http://" + name + ":" + port + "/touch"
}

func handlerCallBack(w http.ResponseWriter, r *http.Request) {
	// Build the request
	req, err := http.NewRequest("GET", backURL(), nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	// Callers should close resp.Body
	defer resp.Body.Close()

	// Get body as string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			log.Fatal("Do: ", err)
			return
		}
		fmt.Fprintf(w, "Got response from the back => %s\n", string(bodyBytes))
	} else {
		fmt.Fprintf(w, "Error while calling the back: %d", resp.StatusCode)
	}
}

func main() {
	port := 8484

	fmt.Fprintf(os.Stdout, "HTTP server is running using port: %d\n", port)
	fmt.Fprintf(os.Stdout, "Available endpoints are: '/hello' and '/back'")

	http.HandleFunc("/hello", handlerHello)
	http.HandleFunc("/info", handlerInfo)
	http.HandleFunc("/back", handlerCallBack)

	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
