package internal

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *server) handleTimestamp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		unixts, _ := strconv.Atoi(chi.URLParam(r, "ts"))
		ts := time.Unix(int64(unixts), 0)

		s.respond(w, chi.URLParam(r, "university"), ts)
	}
}
