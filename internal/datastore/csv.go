package datastore

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"time"
)

var ErrInvalidCSVData error = errors.New("invalid CSV from STWNO")

// getDishes gets all dishes for a certain week, specified by the timestamp.
// The slice of dishes may be empty, if there are no dishes for the given week.
func (d *Datastore) getDishes(university string, ts time.Time) ([]dish, error) {
	if _, ok := universities[university]; !ok {
		return nil, ErrInvalidUniversityRequest
	}
	_, weeknumber := ts.ISOWeek()

	data, err := d.downloadCSV(university, weeknumber)
	if err != nil {
		return nil, err
	}

	dishes, err := parseCSVToDishes(data)
	if err != nil {
		log.Printf("%v, for %v week %v", err, university, weeknumber)
		return nil, ErrInvalidCSVData
	}

	return dishes, nil
}

// parseCSVToDishes parses the given CSV data and returns a slice of dishes.
// The first line is the header line and is ignored.
// The parsing may fail if the CSV data is invalid. This happens every once in a while.
func parseCSVToDishes(data []byte) ([]dish, error) {
	reader := newStwnoReader(bytes.NewReader(data))
	csvData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return convertParsedCSVToDishSlice(csvData), nil
}

// stwnoReader is a CSV reader that can be used to read the CSV data from STWNO.
func newStwnoReader(r io.Reader) *csv.Reader {
	csv := csv.NewReader(r)
	csv.Comma = ';'
	csv.LazyQuotes = true
	return csv
}

// convertParsedCSVToDishSlice converts the parsed CSV data to a slice of dishes.
// The conversion is done manually and depends on the correct order of the CSV data.
func convertParsedCSVToDishSlice(csvData [][]string) []dish {
	// guesstimate of around 10 dishes per weekday.
	var dishes []dish = make([]dish, 0, 50)
	for i, line := range csvData {
		if i > 0 { // omit header line
			var rec dish
			for j, val := range line {
				// for this the layout of the csv needs to be known.
				// see dish struct in datastore.go.
				switch j {
				case 0:
					rec.Date = val
				case 2:
					rec.Type = val
				case 3:
					rec.Name = val
				case 4:
					rec.Tags = val
				case 6:
					rec.PriceStudent = val
				case 7:
					rec.PriceEmployee = val
				case 8:
					rec.PriceGuest = val
				// do not care about the other fields.
				default:
					continue
				}
			}
			dishes = append(dishes, rec)
		}
	}
	return dishes
}
