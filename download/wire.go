package download

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// WireXKCD contains all the data and metadata about an XKCD strip you can get
// by calling its json endpoint.
type WireXKCD struct {
	// The number of the XKCD strip
	ID int `json:"num"`

	// The title of the strip
	Title string `json:"safe_title"`

	// The URL of the image for the strip
	Img string `json:"img"`

	// The alternative text
	Alt string `json:"alt"`

	// The transcript of the strip
	Transcript string

	// An external link, if any.
	Link string `json:"link"`

	// News from the author
	News string `json:"news"`

	// Date information.
	Year  int `json:"year,string"`
	Day   int `json:"day,string"`
	Month int `json:"month,string"`

	// Store the date information in a dedicated field
	DateTime time.Time `json:"date"`
}

// NewFromWire creates a new struct based on the json data coming from the server.
func NewFromWire(r io.Reader) (*WireXKCD, error) {
	var w WireXKCD
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&w)
	if err != nil {
		logger.Debugw("Error decoding the server response", "error", err)
		return nil, err
	}
	w.DateTime, err = w.GetTime()
	if err != nil {
		w.DateTime = time.Time{}
	}
	return &w, nil
}

// GetTime returns a time object for the date of the strip.
func (w WireXKCD) GetTime() (time.Time, error) {
	return time.Parse("2006-01-02", w.Date())
}

// Date returns a string representation of the date of the strip
func (w WireXKCD) Date() string {
	return fmt.Sprintf("%d-%02d-%02d", w.Year, w.Month, w.Day)
}
