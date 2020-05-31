package database

import (
	"github.com/blevesearch/bleve"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

type SearchOpts struct {
	Fields     []string
	SortBy     []string
	MaxRecords int
}

func (s SearchOpts) Apply(request *bleve.SearchRequest) {
	if s.Fields != nil {
		request.Fields = s.Fields
	}
	if s.SortBy != nil {
		request.SortBy(s.SortBy)
	}
	if s.MaxRecords != 0 {
		request.Size = s.MaxRecords
	}
}

// SetLogger sets the logger for the package
func SetLogger(l *zap.SugaredLogger) {
	logger = l
}

// Open returns a new database if none exists, or the existing one.
func Open(path string) (bleve.Index, error) {
	idx, err := bleve.Open(path)
	if err != nil {
		logger.Infow("Creating the database", "path", path)
		mapping := bleve.NewIndexMapping()
		mapping.AddDocumentMapping("xkcd", DocMapping())
		idx, err = bleve.New(path, mapping)
	}
	return idx, err
}

// GetAll fetches all records, according to the search options
func GetAll(idx bleve.Index, opts *SearchOpts) (*bleve.SearchResult, error) {
	query := bleve.NewMatchAllQuery()
	search := bleve.NewSearchRequest(query)
	if opts != nil {
		opts.Apply(search)
	}
	return idx.Search(search)
}

// GetLatestID returns the highest ID recorded in the database
func GetLatestID(idx bleve.Index) int {
	// Note: we're limiting to 1 result because that's all we need.
	allRecords, err := GetAll(idx, &SearchOpts{MaxRecords: 1, SortBy: []string{"-id"}})
	if err != nil {
		logger.Errorw("Could not retreive all data from the datastore", "error", err)
		return 0
	}
	if allRecords.Total == 0 {
		return 0
	}
	latest := NewStripFromDb(allRecords.Hits[0])
	// If errors occur getting the document from the database, we just return 0
	// as we assume the database is corrupted.
	if latest == nil {
		return 0
	}
	return latest.ID
}

// GetAllIDs will return a list of IDs of all strips that have been downloaded, up to maxID.
func GetAllIDs(idx bleve.Index, maxID int) map[int]bool {
	allRecords, err := GetAll(idx, &SearchOpts{MaxRecords: maxID, SortBy: []string{"id"}})
	if err != nil {
		logger.Errorw("Could not retreive all data from the datastore", "error", err)
		return make(map[int]bool)
	}
	results := make(map[int]bool, allRecords.Total)
	for _, record := range allRecords.Hits {
		strip := NewStripFromDb(record)
		if strip != nil {
			results[strip.ID] = true
		}
	}
	return results
}

// SearchStr will perform a string query on the datastore
func SearchStr(idx bleve.Index, queryStr string, opts *SearchOpts) (*bleve.SearchResult, error) {
	query := bleve.NewQueryStringQuery(queryStr)
	search := bleve.NewSearchRequest(query)
	if opts != nil {
		opts.Apply(search)
	} else {
		DefaultSearchOpts.Apply(search)
	}
	searchResults, err := idx.Search(search)
	if err != nil {
		logger.Errorw("Error querying the search index", "error", err)
		return nil, err
	}
	return searchResults, nil
}

var DefaultSearchOpts = &SearchOpts{
	Fields: allFields,
}
