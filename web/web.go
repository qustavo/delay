package web

import (
	"io/ioutil"
	"net/http"

	"github.com/gchaincl/delay"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	return &Server{http.NewServeMux()}
}

func (s *Server) Handle(route string, delayer *delay.Delayer) {
	pattern := "/" + route + "/"
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len(pattern):]

		switch r.Method {
		case "POST":
			defer r.Body.Close()
			payload, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			delayer.Register(key, string(payload))
		case "DELETE":
			if ok := delayer.Cancel(key); !ok {
				http.NotFound(w, r)
				return
			}
		default:
			http.Error(w, "Method Not Allowed", 405)
		}
	})
}

func (s Server) Listen(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}