package internal

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *server) routes() {
	s.mux.Use(middleware.RequestID)
	s.mux.Use(middleware.RealIP)
	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)

	// we use the root directory for health checks.
	s.mux.Use(middleware.Heartbeat("/"))

	// this prefix is always required.
	s.mux.Route("/api/{university}", func(r chi.Router) {
		// trailing slashes are ignored in this subrouter.
		r.Use(middleware.StripSlashes)

		r.Get("/timestamp/{ts:[0-9]+}", s.handleTimestamp())
		r.Get("/date/{date}", s.handleDate())
		r.Get("/today", s.handleToday())
		r.Get("/tomorrow", s.handleTomorrow())
		r.Get("/week", s.handleWeek())
	})
}
