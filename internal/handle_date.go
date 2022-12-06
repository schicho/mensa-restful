package internal

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *server) handleDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := time.Parse("2006-01-02", chi.URLParam(r, "date"))
		if err != nil {
			http.Error(w, "invalid date format. expect YYYY-MM-DD", http.StatusBadRequest)
			return
		} else {
			s.respond(w, chi.URLParam(r, "university"), ts)
		}
	}
}
