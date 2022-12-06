package internal

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (s *server) handleTimestamp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		// No need to check error here. gorilla mux already does regex check this.
		unixts, _ := strconv.Atoi(vars["ts"])
		ts := time.Unix(int64(unixts), 0)

		s.respond(w, vars["university"], ts)
	}
}
