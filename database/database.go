package database

import (
	"github.com/blevesearch/bleve"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

type SearchOpts struct {
	Fields []string
	SortBy []string
}

func (s SearchOpts) Apply(request *bleve.SearchRequest) {
	if s.Fields != nil {
		request.Fields = s.Fields
	}
	if s.SortBy != nil {
		request.SortBy(s.SortBy)
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
	allRecords, err := GetAll(idx, &SearchOpts{Fields: []string{"id"}, SortBy: []string{"-id"}})
	if err != nil {
		logger.Errorw("Could not retreive all data from the datastore", "error", err)
		return 0
	}
	if allRecords.Total == 0 {
		return 0
	}
	maxidFloat := allRecords.Hits[0].Fields["id"].(float64)
	return int(maxidFloat)
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
