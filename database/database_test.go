package database

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/blevesearch/bleve"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var fixturedb bleve.Index

func setup() {
	l, _ := zap.NewDevelopment()
	logger = l.Sugar()
	var err error
	fixturedb, err = Open("../fixtures/sample.db")
	if err != nil {
		panic(err)
	}
	// populate the default database
	for i := 1; i < 11; i++ {
		strip := XKCDStrip{
			Title:   fmt.Sprintf("test strip %d", i),
			ID:      2300 + i,
			Img:     fmt.Sprintf("https://example.com/%d/image.png", i),
			Date:    fmt.Sprintf("2020-01-%02d", i),
			Comment: fmt.Sprintf("Comment #%d", i),
		}
		fixturedb.Index(strip.Title, strip)
	}
}

func teardown() {
	logger.Sync()
	fixturedb.Close()
	os.RemoveAll("../fixtures/sample.db")
}
func TestOpen(t *testing.T) {
	setup()
	defer teardown()
	tempdir, err := ioutil.TempDir("", "xkcli-test")
	if err != nil {
		t.Errorf("Unable to create the temporary directory")
	}
	defer os.RemoveAll(tempdir)
	// Try opening an inexistent database.
	db, err := Open(path.Join(tempdir, "xkcli.bleve"))
	if err != nil {
		t.Errorf("Error opening a new database: %s", err)
	}
	db.Close()
	//Now open an existing database.
	db, err = Open(path.Join(tempdir, "xkcli.bleve"))
	if err != nil {
		t.Errorf("Error opening an existing database: %s", err)
	}
	db.Close()
}

// Test GetAll fetches all the records from the database
func TestGetAll(t *testing.T) {
	setup()
	defer teardown()
	results, err := GetAll(fixturedb, nil)
	assert.Equal(t, nil, err, "Error getting all records: %s", err)
	assert.Equal(t, 10, int(results.Total), "Wrong number of results from GetAll")
}

// Test GetLatestID actually returns the expected ID
func TestGetLatestID(t *testing.T) {
	setup()
	defer teardown()
	assert.Equal(t, 2310, GetLatestID(fixturedb))
}
