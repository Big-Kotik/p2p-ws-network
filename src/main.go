package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net-backend/src/hub"
	"net-backend/src/workers"
	"net/http"
)

var port = flag.String("port", ":8080", "http service port")

// ServeHandler type for serve hub functions
type ServeHandler func(hub.Hub, http.ResponseWriter, *http.Request)

func applyServeFunc(h hub.Hub, serve ServeHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) { serve(h, writer, request) }
}

func newServer() *http.Server {
	r := mux.NewRouter()
	debug := flag.Bool("debug", false, "run debug mode")
	h := hub.GetHub()
	flag.Parse()
	if *debug {
		r.HandleFunc("/", serveHome)
	}
	r.HandleFunc("/create_room", applyServeFunc(h, serveRoom))
	r.Path("/check_id").Methods("GET").Queries("id", "{id}").HandlerFunc(serveCheckIDExist)
	r.HandleFunc("/ws/client", applyServeFunc(h, serveClientWs))
	r.HandleFunc("/ws/node", applyServeFunc(h, serveNodeWs))

	return &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0" + *port,
		WriteTimeout: workers.WriteWait,
		ReadTimeout:  workers.WriteWait,
	}
}

func main() {
	srv := newServer()
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("can't create server")
		return
	}
}
