package download

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/blevesearch/bleve"
)

// WireXKCD contains all the data and metadata about an XKCD strip you can get
// by calling its json endpoint.
type WireXKCD struct {
	// The number of the XKCD strip
	ID int `json:"num"`

	// The title of the strip
	Title string

	// The URL of the image for the strip
	Img string `json:"img"`

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

// Summary offers a formatted output.
func (w WireXKCD) Summary() string {
	return fmt.Sprintf("XKCD %d (%s): %s\n\tstrip: %s\n", w.ID, w.Date(), w.Title, w.Img)
}

// Index performs the indexing of this resource in a bleve index.
func (w *WireXKCD) Index(idx bleve.Index) error {
	err := idx.Index(w.Title, w)
	if err != nil {
		logger.Errorw("Error indexing", "id", w.ID, "error", err.Error())
	}
	return err
}
