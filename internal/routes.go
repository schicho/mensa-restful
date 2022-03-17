package internal

func (s *server) routes() {
	// this prefix is always required.
	u := s.mux.PathPrefix("/api/{university}").Subrouter()

	// unix time
	u.HandleFunc("/timestamp/{ts:[0-9]+}", s.handleTimestamp()).Methods("GET")

	// date
	u.HandleFunc("/date/{date}", s.handleDate()).Methods("GET")

	// syntactic sugar
	u.HandleFunc("/today", s.handleToday()).Methods("GET")
	u.HandleFunc("/tomorrow",s.handleTomorrow()).Methods("GET")
	u.HandleFunc("/week", s.handleWeek()).Methods("GET")
}