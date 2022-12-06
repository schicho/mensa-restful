package internal

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (s *server) handleDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		ts, err := time.Parse("2006-01-02", vars["date"])
		if err != nil {
			http.Error(w, "invalid date format. expect YYYY-MM-DD", http.StatusBadRequest)
			return
		} else {
			s.respond(w, vars["university"], ts)
		}
	}
}
