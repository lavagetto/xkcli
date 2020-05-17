package database

import (
	"testing"

	"github.com/lavagetto/xkcli/download"
	"github.com/stretchr/testify/assert"
)

func TestNewStrip(t *testing.T) {
	w := download.WireXKCD{
		ID:    1,
		Title: "title",
		Img:   "img",
		Alt:   "comment",
		Year:  2020,
		Month: 1,
		Day:   25,
	}
	strip := NewStrip(&w)
	assert.Equal(t, "title", strip.Title)
	assert.Equal(t, "comment", strip.Comment)
	assert.Equal(t, "img", strip.Img)
	assert.Equal(t, "2020-01-25", strip.Date)
}

// Test that we correctly build an object from a search result
func TestNewStripFromDb(t *testing.T) {
	setup()
	defer teardown()
	results, _ := GetAll(fixturedb, &SearchOpts{Fields: allFields, SortBy: []string{"-id"}})
	strip := NewStripFromDb(results.Hits[0])
	assert.Equal(t, 2310, strip.ID)
	assert.Equal(t, "https://example.com/10/image.png", strip.Img)
	assert.Equal(t, "2020-01-10", strip.Date)
}
func TestSummary(t *testing.T) {
	strip := XKCDStrip{
		Title: "test",
		Img:   "test.jpg",
		ID:    1,
		Date:  "2020-04-01",
	}
	expectedSummary := `XKCD 1 (2020-04-01): test
	strip: test.jpg
`
	assert.Equal(t, expectedSummary, strip.Summary())
}
