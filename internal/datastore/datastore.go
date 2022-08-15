package datastore

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/allegro/bigcache/v3"
)

var universities = map[string]struct{}{
	"UNI-R": {},
	"UNI-R-Gs": {},
	"Cafeteria-PT": {},
	"Cafeteria-Chemie": {},
	"Cafeteria-Sport": {},
	"HS-R-tag": {},
	"HS-R-abend": {},
	"Cafeteria-pruefening": {},
	"UNI-P": {},
	"Cafeteria-Nikolakloster": {},
	"HS-DEG": {},
	"HS-LA": {},
	"HS-SR": {},
	"HS-PAN": {},
}

type dish struct {
	Date          string `json:"date"`           // csv column 0
	Type          string `json:"type"`           // csv column 2
	Name          string `json:"name"`           // csv column 3
	Tags          string `json:"tags"`           // csv column 4
	PriceStudent  string `json:"price_student"`  // csv column 6
	PriceEmployee string `json:"price_employee"` // csv column 7
	PriceGuest    string `json:"price_guest"`    // csv column 8
}

type Datastore struct {
	cache *bigcache.BigCache
	// client is thread save and can be used by multiple goroutines.
	client *http.Client
}

func NewDatastore() *Datastore {
	// use simple initialization
	bc, _ := bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))

	d := &Datastore{
		cache: bc,
		client: http.DefaultClient,
	}
	return d
}

// GetJsonDay returns the JSON data of dishes for the given university and day of the timestamp.
// The slice may be empty, with no error.
func (d *Datastore) GetJsonDay(university string, ts time.Time) ([]byte, error) {
	return d.getJson(university, ts, true)
}

// GetJsonWeek returns the JSON data of dishes for the given university and week of the timestamp.
// The slice may be empty, with no error.
func (d *Datastore) GetJsonWeek(university string, ts time.Time) ([]byte, error) {
	return d.getJson(university, ts, false)
}

// getJson returns the JSON data for the given university and timestamp.
// It may return cached data, if the data has been recently requested 
func (d *Datastore) getJson(university string, ts time.Time, filterDay bool) ([]byte, error) {
	if _, ok := universities[university]; !ok {
		return nil, ErrInvalidUniversityRequest
	}

	requestKey := fmt.Sprint(university, ts.YearDay(), filterDay)

	// quick respond with cached value
	if cachedJson, err := d.cache.Get(requestKey); err == nil {
		// log.Println("cache hit", requestKey)
		return cachedJson, nil
	}
	// log.Println("cache miss", requestKey)

	dishes, err := d.getDishes(university, ts)
	if err != nil {
		return nil, err
	}

	if filterDay {
		dishes = filterDishesDay(dishes, ts)
	}

	json, err := json.Marshal(dishes)
	if err != nil {
		// report and assume invalid CSV when marshaling error happens.
		log.Println("marshalling of dishes failed", err)
		return nil, ErrInvalidCSVData
	}

	// populate cache with new data for future requests.
	d.cache.Set(requestKey, json)
	return json, nil
}

// filterDishesDay returns only the dishes matching the timestamps date.
// May return an empty slice.
func filterDishesDay(dishes []dish, ts time.Time) []dish {
	var dishesForDate []dish = make([]dish, 0, 10)
	date := ts.Format("02.01.2006")
	
	for _, dish := range dishes {
		if dish.Date == date {
			dishesForDate = append(dishesForDate, dish)
		}
	}
	return dishesForDate
}
