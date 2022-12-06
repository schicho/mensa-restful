package internal

import "github.com/go-chi/chi/v5"

func (s *server) routes() {
	// this prefix is always required.
	s.mux.Route("/api/{university}", func(r chi.Router) {
		r.Get("/timestamp/{ts:[0-9]+}", s.handleTimestamp())
		r.Get("/date/{date}", s.handleDate())
		r.Get("/today", s.handleToday())
		r.Get("/tomorrow", s.handleTomorrow())
		r.Get("/week", s.handleWeek())
	})
}
