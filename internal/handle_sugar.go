package internal

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *server) handleToday() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, chi.URLParam(r, "university"), time.Now())
	}
}

func (s *server) handleTomorrow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, chi.URLParam(r, "university"), time.Now().Add(24*time.Hour))
	}
}

func (s *server) handleWeek() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respondWeek(w, chi.URLParam(r, "university"), time.Now())
	}
}
