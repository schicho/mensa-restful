package datastore

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/text/encoding/charmap"
)

var ErrDownloadFromSourceFail error = errors.New("could not download data")
var ErrInvalidUniversityRequest error = errors.New("invalid university provided")

var removeMultiNewline = regexp.MustCompile(`\n+;`)

// example url = `https://www.stwno.de/infomax/daten-extern/csv/UNI-P/11.csv`
const url = "https://www.stwno.de/infomax/daten-extern/csv/%v/%v.csv"

// downloadCSV downloads the source CSV file from STWNO.
// The returned data is in UTF8 format.
// The function may return NO error, whilst the data is no valid CSV.
// This is due to STWNO's inconsistent CSV format.
func (d *Datastore) downloadCSV(university string, weeknumber int) ([]byte, error) {
	if _, ok := universities[university]; !ok {
		return nil, ErrInvalidUniversityRequest
	}

	url := fmt.Sprintf(url, university, weeknumber)
	respData, err := d.makeRequest(url)
	if err != nil {
		return nil, err
	}

	data, err := convertWindows1252ToUTF8(respData)
	if err != nil {
		return nil, err
	}

	// finally remove the arbitrary newlines, which are littered in the CSV.
	// Fortunately those seem to appear in a pattern, being right before a semicolon.
	// This is strange, but I don't serve the CSV, I just need to be able to read it.
	// Better approach would be to have a CSV Reader, which just ignores newlines until a complete entry is full.
	data = removeMultiNewline.ReplaceAll(data, []byte{';'})

	return data, nil
}

// makeRequest fetches the data from the given url.
// The user agent is set to mensa_json_api_crawler to identify the automatic request.
func (d *Datastore) makeRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	const userAgent = "mensa_json_api_crawler"
	req.Header.Set("User-Agent", userAgent)

	resp, err := d.client.Do(req)
	if err != nil {
		log.Println("download of CSV failed", err)
		return nil, ErrDownloadFromSourceFail
	}
	defer resp.Body.Close()

	respData, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Println("reading response failed", readErr)
		return nil, ErrDownloadFromSourceFail
	}

	return respData, nil
}

func convertWindows1252ToUTF8(data []byte) ([]byte, error) {
	decoder := charmap.Windows1252.NewDecoder()
	return decoder.Bytes(data)
}
