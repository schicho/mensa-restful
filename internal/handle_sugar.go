package internal

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (s *server) handleToday() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, mux.Vars(r)["university"], time.Now())
	}
}

func (s *server) handleTomorrow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, mux.Vars(r)["university"], time.Now().Add(24*time.Hour))
	}
}

func (s *server) handleWeek() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respondWeek(w, mux.Vars(r)["university"], time.Now())
	}
}
