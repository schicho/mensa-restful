package internal

import (
	"net/http"
	"time"

	"github.com/schicho/mensa-restful/internal/datastore"
)

func (s *server) respond(w http.ResponseWriter, university string, ts time.Time) {
	data, err := s.datastore.GetJsonDay(university, ts)
	if err != nil {
		respondError(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func (s *server) respondWeek(w http.ResponseWriter, university string, ts time.Time) {
	data, err := s.datastore.GetJsonWeek(university, ts)
	if err != nil {
		respondError(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func respondError(w http.ResponseWriter, err error) {
	switch err {
	case datastore.ErrInvalidUniversityRequest:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case datastore.ErrDownloadFromSourceFail:
		fallthrough
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
