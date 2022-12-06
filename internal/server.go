package internal

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/schicho/mensa-restful/internal/datastore"
)

type server struct {
	mux       *chi.Mux
	datastore *datastore.Datastore
}

func NewServer() (*server, error) {
	srv := &server{
		// StrictSlash(true) allows us to ignore trailing slashes.
		mux:       chi.NewRouter(),
		datastore: datastore.NewDatastore(),
	}
	srv.routes()
	return srv, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
