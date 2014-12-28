package web

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gchaincl/delay"
)

type Server struct {
	mux *mux.Router
}

func NewServer() *Server {
	return &Server{mux.NewRouter()}
}

func (s *Server) Handle(route string, delayer *delay.Delayer) {
	pattern := "/" + route + "/{key}"
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		key := mux.Vars(r)["key"]
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		delayer.Register(key, string(payload))
	}).Methods("POST")

	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]
		if ok := delayer.Cancel(key); !ok {
			http.NotFound(w, r)
			return
		}
	}).Methods("DELETE")
}

func (s Server) Listen(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}