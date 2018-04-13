package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func confServerPort() int {
	port := os.Getenv("PORT")
	if port == "" {
		return 8484
	}
	num, err := strconv.Atoi(port)
	if err != nil {
		log.Println("ERROR during conversion: ", err)
		return 8484
	}
	return num
}

func confBackURL() string {
	name := os.Getenv("BACK")
	if name == "" {
		name = "hello-back-svc"
	}
	port := os.Getenv("BACK_PORT")
	if port == "" {
		port = "8485"
	}
	return "http://" + name + ":" + port + "/touch"
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
	h, _ := os.Hostname()
	fmt.Fprintf(w, "Hi there, I'm served from %s!\n", h)
}

func handlerInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "- Used port is: %d!\n", confServerPort())
	fmt.Fprintf(w, "- Used back URL is: %s!\n", confBackURL)
}

func handlerCallBack(w http.ResponseWriter, r *http.Request) {
	// Build the request
	req, err := http.NewRequest("GET", confBackURL(), nil)
	if err != nil {
		log.Println("ERROR when building request: ", err)
		return
	}

	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR when requesting backend: ", err)
		return
	}

	// Callers should close resp.Body
	defer resp.Body.Close()

	// Get body as string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			log.Println("ERROR when getting body: ", err)
			return
		}
		fmt.Fprintf(w, "Got response from the back => %s\n", string(bodyBytes))
	} else {
		fmt.Fprintf(w, "Error while calling the back: %d", resp.StatusCode)
	}
}

func main() {
	port := confServerPort()

	log.Printf("HTTP server is running using port: %d\n", port)
	log.Println("Available endpoints are: '/hello' and '/back'")

	http.HandleFunc("/hello", handlerHello)
	http.HandleFunc("/info", handlerInfo)
	http.HandleFunc("/back", handlerCallBack)

	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
