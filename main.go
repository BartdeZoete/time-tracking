package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time-tracking/trie"

	"github.com/julienschmidt/httprouter"
)

var t = trie.MakeTrie()

const (
	prefix = "/api"
)

func list(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(t)
}

func create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := strings.Split(ps[0].Value, "+")
	t.Add(path)
	// TODO: res
}

func start(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: take res
	path := strings.Split(ps[0].Value, "+")
	t.Record(path)
}

func stop(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: take res
	t.Stop()
}

func main() {
	router := httprouter.New()

	router.GET(prefix, list)
	router.POST(prefix+"/projects/:path", create)
	router.POST(prefix+"/projects/:path/start", start)
	router.POST(prefix+"/stop", stop)

	log.Fatal(http.ListenAndServe(":8080", router))
}
