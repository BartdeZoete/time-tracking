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
	port   = "8080"
)

func list(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(t)
}

func create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := strings.Split(ps[0].Value, "+")

	if err := t.Add(path); err != nil {
		w.WriteHeader(http.StatusInternalServerError) //REVIEW
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func start(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := strings.Split(ps[0].Value, "+")
	if t.IsRecording() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("already recording"))
		return
	}

	if ok := t.Record(path); !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("project doesn't exist"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func stop(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if ok := t.Stop(); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not recording"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	router := httprouter.New()

	router.GET(prefix, list)
	router.POST(prefix+"/projects/:path", create)
	router.POST(prefix+"/projects/:path/start", start)
	router.POST(prefix+"/stop", stop)

	log.Print("listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
