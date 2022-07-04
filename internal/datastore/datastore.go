package datastore

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/text/encoding/charmap"
)

// example url = `https://www.stwno.de/infomax/daten-extern/csv/UNI-P/11.csv`
const url = "https://www.stwno.de/infomax/daten-extern/csv/%v/%v.csv"

var ErrDownloadFromSourceFail error = errors.New("could not download data")
var ErrInvalidUniversityRequest error = errors.New("invalid university provided")
var ErrInvalidCSVData error = errors.New("invalid CSV from STWNO")

type dish struct {
	Date          string `json:"date"`           // csv column 0
	Type          string `json:"type"`           // csv column 2
	Name          string `json:"name"`           // csv column 3
	Tags          string `json:"tags"`           // csv column 4
	PriceStudent  string `json:"price_student"`  // csv column 6
	PriceEmployee string `json:"price_employee"` // csv column 7
	PriceGuest    string `json:"price_guest"`    // csv column 8
}

type cacheddata struct {
	etag string
	dishes []dish
}

type Datastore struct {
	// thread save map is needed, as the cache might get updated by another goroutine
	// while data is being read.
	data map[string](*mutexMap[int, cacheddata])

	// client is thread save and can be used by multiple goroutines.
	client *http.Client
}

func NewDatastore() *Datastore {
	d := &Datastore{
		data: map[string](*mutexMap[int, cacheddata]){
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
		},
		client: http.DefaultClient,
	}
	for k := range d.data {
		d.data[k] = newMutexMap[int, cacheddata]()
	}
	return d
}

func (d *Datastore) getJson(university string, ts time.Time, filterDay bool) ([]byte, error) {
	if _, ok := d.data[university]; !ok {
		return nil, ErrInvalidUniversityRequest
	}

	dishes, err := d.getDishes(university, ts)
	if err != nil {
		// passthrough error
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
	return json, nil
}

func (d *Datastore) GetJsonDay(university string, ts time.Time) ([]byte, error) {
	return d.getJson(university, ts, true)
}

func (d *Datastore) GetJsonWeek(university string, ts time.Time) ([]byte, error) {
	return d.getJson(university, ts, false)
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

func convertStwnoCSVToDishSlice(csvData [][]string) []dish {
	// guesstimate of around 10 dishes per weekday.
	var dishes []dish = make([]dish, 0, 50)
	for i, line := range csvData {
        if i > 0 { // omit header line
            var rec dish
            for j, val := range line {
				// for this the layout of the csv needs to be known.
				// see struct on the beginning of this file.
                switch j {
				case 0: rec.Date = val
				case 2: rec.Type = val
				case 3: rec.Name = val
				case 4: rec.Tags = val
				case 6: rec.PriceStudent = val
				case 7: rec.PriceEmployee = val
				case 8: rec.PriceGuest = val
				// do not care about the other fields.
				default: continue
				}
            }
            dishes = append(dishes, rec)
        }
    }
    return dishes
}

func newStwnoReader(r io.Reader) *csv.Reader {
	csv := csv.NewReader(r)
	csv.Comma = ';'
	csv.LazyQuotes = true
	return csv
}

// getDishes gets all dishes for a certain week, specified by the timestamp.
// It reuses the cached data, if the data at STWNO has not changed.
// If it has changed it will return the updated list of dishes and update the cache.
func (d *Datastore) getDishes(university string, ts time.Time) ([]dish, error) {
	if _, ok := d.data[university]; !ok {
		return nil, ErrInvalidUniversityRequest
	}
	_, weeknumber := ts.ISOWeek()

	reuse, newData, etag, err := d.downloadCSV(university, weeknumber)
	if err != nil {
		return nil, err
	}

	if reuse {
		// no need to further check this as reuse is only set true,
		// if data was already present.
		data, _ :=  d.data[university].Load(weeknumber)
		return data.dishes, nil
	} else {
		// parse and convert data.
		csvReader := newStwnoReader(bytes.NewBuffer(newData))
		csvData, err := csvReader.ReadAll()
		if err != nil {
			log.Println("reading of CSV data failed weeknumber", weeknumber, ":", err)
			return nil, ErrInvalidCSVData
		}

		// cache data and return.
		dishes := convertStwnoCSVToDishSlice(csvData)
		d.data[university].Store(weeknumber, cacheddata{dishes: dishes, etag: etag})
		return dishes, nil
	}
}

// downloadCSV downloads the source CSV file. It uses a conditional HTTP GET request.
// reuse is true if the previous data can be reused. In this case newData is nil an may not be used.
// If reuse is false, newData will contain the new CSV data and etag the new etag string.
// If an error occurs the above does not hold. The return values may not be used.
//curl -I https://www.stwno.de/infomax/daten-extern/csv/UNI-P/11.csv --header 'If-None-Match: "af8-5da5643025f0b"'
func (d *Datastore) downloadCSV(university string, weeknumber int) (reuse bool, newData []byte, etag string, err error) {
	if _, ok := d.data[university]; !ok {
		err = ErrInvalidUniversityRequest
		return
	}
	url := fmt.Sprintf(url, university, weeknumber)

	req, _ := http.NewRequest("GET", url, nil)
	const userAgent = "mensa_json_api_crawler"
	req.Header.Set("User-Agent", userAgent)

	if cached, ok := d.data[university].Load(weeknumber); ok {
		req.Header.Add("If-None-Match", cached.etag)
	}
	
	resp, err := d.client.Do(req)
	if err != nil {
		log.Println("download of CSV failed", err)
		err = ErrDownloadFromSourceFail
		return
	}
	defer resp.Body.Close()

	// Not modified.
	if resp.StatusCode == 304 {
		reuse = true
		return
	} else {
		log.Println("cache outdated or non present for", university, weeknumber)
		respData, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Println("reading response failed", readErr)
			err = ErrDownloadFromSourceFail
			return
		}
		
		// convert this horrible format that is used by stwno to UTF8.
		decoder := charmap.Windows1252.NewDecoder()
		bufUTF8 := make([]byte, len(respData)*3)
		n, _, decoderErr := decoder.Transform(bufUTF8, respData, false)
		if decoderErr != nil {
			log.Println("conversion to UTF8 failed", err)
			err = ErrDownloadFromSourceFail
			return
		}
		
		newData = bufUTF8[:n]

		// finally remove the arbitrary newlines, which are littered in the CSV.
		// Fortunately those seem to appear in a pattern, being right before a semicolon.
		// This is strange, but I don't serve the CSV, I just need to be able to read it.
		// Better approach would be to have a CSV Reader which, just ignores newlines, until a complete entry is full.
		newData = bytes.ReplaceAll(newData, []byte{'\n',';'}, []byte{';'})

		etag = resp.Header.Get("Etag")
		return
	}
}