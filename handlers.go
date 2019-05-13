package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// handles "/"
func indexHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", htmlutf8)
	n, err := w.Write([]byte("getwtxt v" + getwtxt))
	if err != nil || n == 0 {
		log.Printf("Error writing to HTTP stream: %v\n", err)
	}

}

// handles "/api"
func apiBaseHandler(w http.ResponseWriter, r *http.Request) {
	timerfc3339, err := time.Now().MarshalText()
	if err != nil {
		log.Printf("Couldn't format time as RFC3339: %v\n", err)
	}
	etag := fmt.Sprintf("%x", sha256.Sum256(timerfc3339))
	w.Header().Set("ETag", etag)
	w.Header().Set("Content-Type", txtutf8)
	pathdata := []byte("\n\n" + r.URL.Path)
	timerfc3339 = append(timerfc3339, pathdata...)
	n, err := w.Write(timerfc3339)
	if err != nil || n == 0 {
		log.Printf("Error writing to HTTP stream: %v\n", err)
	}
}

// handles "/api/plain"
// maybe add json/xml support later
func apiFormatHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	format := vars["format"]

	w.Header().Set("Content-Type", txtutf8)
	n, err := w.Write([]byte(format + "\n"))
	if err != nil || n == 0 {
		log.Printf("Error writing to HTTP stream: %v\n", err)
	}
}

// handles "/api/plain/(users|mentions|tweets)"
func apiEndpointHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	format := vars["format"]
	endpoint := vars["endpoint"]

	w.Header().Set("Content-Type", htmlutf8)
	n, err := w.Write([]byte(format + "/" + endpoint))
	if err != nil || n == 0 {
		log.Printf("Error writing to HTTP stream: %v\n", err)
	}

}

// handles POST for "/api/plain/users"
func apiEndpointPOSTHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	format := vars["format"]
	endpoint := vars["endpoint"]

	w.Header().Set("Content-Type", htmlutf8)
	n, err := w.Write([]byte(format + "/" + endpoint))
	if err != nil || n == 0 {
		log.Printf("Error writing to HTTP stream: %v\n", err)
	}

}

// handles "/api/plain/tags"
func apiTagsBaseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	format := vars["format"]

	w.Header().Set("Content-Type", htmlutf8)
	n, err := w.Write([]byte("api/" + format + "/tags"))
	if err != nil || n == 0 {
		log.Printf("Error writing to HTTP stream: %v\n", err)
	}

}

// handles "/api/plain/tags/[a-zA-Z0-9]+"
func apiTagsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	format := vars["format"]
	tags := vars["tags"]

	w.Header().Set("Content-Type", htmlutf8)
	n, err := w.Write([]byte("api/" + format + "/tags/" + tags))
	if err != nil || n == 0 {
		log.Printf("Error writing to HTTP stream: %v\n", err)
	}

}
