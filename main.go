package main
import (
        "fmt"
        "net/http"
        "os"
)
func handler(w http.ResponseWriter, r *http.Request) {
        h, _ := os.Hostname()
        fmt.Fprintf(w, "Hi there, I'm served from %s!", h)
}
func main() {
        port := 8484
        ep := "hello"

        fmt.Fprintf(w, "Running HTTP server, using endpoint: %s, port: %d", ep, port)

        http.HandleFunc("/" + ep, handler)
        http.ListenAndServe(":" + port, nil)
}
