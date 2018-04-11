package main
import (
        "strconv"
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

        fmt.Fprintf(os.Stdout, "Running HTTP server, using endpoint: %s, port: %d", ep, port)

        http.HandleFunc("/" + ep, handler)
        http.ListenAndServe(":" + strconv.Itoa(port), nil)
}
