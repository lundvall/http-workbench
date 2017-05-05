package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Our data source
type stopwatch interface {
	elapsed() time.Duration
}

type uptime struct {
	Booted time.Time
}

// Implement stopwatch
func (up *uptime) elapsed() time.Duration {
	return time.Since(up.Booted)
}

// Compose stopwatch and handler
type upHandler struct {
	watch  stopwatch
	handle func(watch stopwatch, w http.ResponseWriter, r *http.Request)
}

func (handler upHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.handle(handler.watch, w, r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<a href="uptime">Uptime</a> (<a href="uptime/json">json</a>)`)
}

func htmlHandler(watch stopwatch, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<p>Up %s</p>", watch.elapsed())
}

func jsonHandler(watch stopwatch, w http.ResponseWriter, r *http.Request) {
	elapsed := struct {
		Uptime time.Duration `json:"uptime"`
	}{watch.elapsed() / time.Second}

	result, _ := json.Marshal(elapsed)
	w.Write(result)
}

func main() {
	up := &uptime{time.Now()}

	http.HandleFunc("/", indexHandler)
	http.Handle("/uptime", upHandler{up, htmlHandler})
	http.Handle("/uptime/json", upHandler{up, jsonHandler})
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
