package download

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var wire *WireXKCD

func wire_setup() {
	wire = &WireXKCD{
		ID:         1,
		Title:      "test",
		Img:        "test.jpg",
		Transcript: "yes",
		Year:       2020,
		Day:        1,
		Month:      4,
	}

}

func TestGetTime(t *testing.T) {
	wire_setup()
	date, err := wire.GetTime()
	assert.Equal(t, err, nil, "Error rendering a valid date")
	assert.Equal(t, date.Year(), 2020, "Date rendered incorrectly.")
}

func TestGetBadTime(t *testing.T) {
	wire_setup()
	wire.Day = 129
	_, err := wire.GetTime()
	assert.Error(t, err)
}

func TestDate(t *testing.T) {
	wire_setup()
	assert.Equal(t, wire.Date(), "2020-04-01")
}

func TestSummary(t *testing.T) {
	wire_setup()
	expectedSummary := `XKCD 1 (2020-04-01): test
	strip: test.jpg
`
	assert.Equal(t, wire.Summary(), expectedSummary)
}
