package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search"
	"github.com/lavagetto/xkcli/download"
)

// XKCDStrip is the data structure we save to the index
type XKCDStrip struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Date       string `json:"date"`
	Img        string `json:"img"`
	Comment    string `json:"comment"`
}

// NewStrip transforms what we got from the wire into a document
// we can index in bleve.
func NewStrip(w *download.WireXKCD) *XKCDStrip {
	strip := XKCDStrip{
		ID:         w.ID,
		Title:      w.Title,
		Transcript: w.Transcript,
		Date:       w.Date(),
		Img:        w.Img,
		Comment:    w.Alt,
	}
	return &strip
}

// NewStripFromDb returns an xkcd strip from data recovered from the database.
func NewStripFromDb(result *search.DocumentMatch) *XKCDStrip {
	id, err := strconv.Atoi(result.ID)
	if err != nil {
		logger.Errorw("Non-numeric ID found", "error", err)
		return nil
	}
	strip := XKCDStrip{ID: id}
	if datetime, ok := result.Fields["date"]; ok {
		// This is funnily stored as a full RFC3339 time string
		t, err := time.Parse(time.RFC3339, datetime.(string))
		if err == nil {
			strip.Date = t.Format("2006-01-02")
		}
	}
	// This is verbose for code readability. Resist the temptation
	// of making this code shorter by 2 lines.
	if title, ok := result.Fields["title"]; ok {
		strip.Title = title.(string)
	}
	if comment, ok := result.Fields["comment"]; ok {
		strip.Comment = comment.(string)
	}
	if img, ok := result.Fields["img"]; ok {
		strip.Img = img.(string)
	}
	if transcript, ok := result.Fields["transcript"]; ok {
		strip.Transcript = transcript.(string)
	}
	return &strip
}

// BleveType implements the BleveClassifier interface
func (x *XKCDStrip) BleveType() string {
	return "xkcd"
}

// Index performs the indexing of this resource in a bleve index.
func (x *XKCDStrip) Index(idx bleve.Index) error {
	err := idx.Index(strconv.Itoa(x.ID), x)
	if err != nil {
		logger.Errorw("Error indexing", "id", x.ID, "error", err.Error())
	}
	return err
}

// Summary offers a formatted output.
func (x XKCDStrip) Summary() string {
	return fmt.Sprintf("XKCD %d (%s): %s\n\tstrip: %s\n", x.ID, x.Date, x.Title, x.Img)
}

// URL returns the full url of a strip
func (x XKCDStrip) URL() string {
	return fmt.Sprintf("https://xkcd.com/%d", x.ID)
}

var allFields = []string{"title", "id", "img", "comment", "transcript", "date"}

// DocMapping returns a bleve document mapping suitable to store this object
// and attaches it to a main index mapping.
func DocMapping() *mapping.DocumentMapping {
	docmap := bleve.NewDocumentMapping()
	title := bleve.NewTextFieldMapping()
	docmap.AddFieldMappingsAt("title", title)
	title.Store = true
	id := bleve.NewNumericFieldMapping()
	docmap.AddFieldMappingsAt("id", id)
	for _, label := range []string{"img", "comment", "transcript"} {
		fm := bleve.NewTextFieldMapping()
		fm.Store = true
		docmap.AddFieldMappingsAt(label, fm)
	}
	datemap := bleve.NewDateTimeFieldMapping()
	datemap.Store = true
	docmap.AddFieldMappingsAt("date", datemap)
	return docmap
}
