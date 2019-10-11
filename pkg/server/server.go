package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Parameters is a list of stuff in order to customize the server
type Parameters struct {
	serviceName string
	serverPort  int
	remoteURL   string
	helloMsg    string
}

// NewParameters constructs a Parameters struct
func NewParameters(svc string, port int, url string) *Parameters {
	h, _ := os.Hostname()
	msg := fmt.Sprintf("Hello! I'm service '%s' hosted by '%s'", svc, h)
	return &Parameters{svc, port, url, msg}
}

// WithMessage sets a custom message to Parameters
func (p *Parameters) WithMessage(msg string) {
	p.helloMsg = msg
}

// Server is our server component as a struct (-> no global state)
type server struct {
	params *Parameters
	router *mux.Router
}

// ListenAndServe listens HTTP requests
func ListenAndServe(params *Parameters) {
	server := &server{params, mux.NewRouter()}
	server.routes()
	portWithColon := fmt.Sprintf(":%d", params.serverPort)
	log.Printf("Listen and serve at http://localhost%s\n", portWithColon)
	log.Fatal(http.ListenAndServe(portWithColon, server))
}

// A server is essentially an http.Handler
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// A handler to say we are alive
func (s *server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// This fuction is frequently used by K8S
		// -> do not fill the logs, do not record metrics neither
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"alive": true}`)
	}
}

// A handler to say hello
func (s *server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s request to %s\n", r.Method, r.RequestURI)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		io.WriteString(w, `{"message": "`+s.params.helloMsg+`"}`)
	}
}

// A handler to say hello
func (s *server) handleGetURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s request to %s\n", r.Method, r.RequestURI)

		w.Header().Set("Content-Type", "application/json")

		// Build the request
		req, err := http.NewRequest("GET", "http://"+s.params.remoteURL, nil)
		if err != nil {
			log.Errorf("Error while building response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, `{"message": "Error while building request"}`)
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
			log.Errorf("Error while requesting backend: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			io.WriteString(w, `{"message": "Error while requesting backend"}`)
			return
		}

		// Callers should close resp.Body
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			// Get body as string
			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				log.Errorf("Error while getting body: %v", err2)
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, `{"message": "Error while getting body"}`)
				return
			}

			// Write response as Json
			io.WriteString(w, `{"remote": "`+string(bodyBytes)+`"}`)
		} else {
			// Write error
			log.Errorf("Error while calling the backend app: %d", resp.StatusCode)
			w.WriteHeader(resp.StatusCode)
			io.WriteString(w, `{"message": "Error while calling the backend app"}`)
		}
	}
}
