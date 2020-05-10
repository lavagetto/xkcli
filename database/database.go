package database

import (
	"github.com/blevesearch/bleve"
	"go.uber.org/zap"
)

// New returns a new database
func Open(path string, logger *zap.SugaredLogger) (bleve.Index, error) {
	idx, err := bleve.Open(path)
	if err != nil {
		logger.Infow("Creating the database", "path", path)
		mapping := bleve.NewIndexMapping()
		idx, err = bleve.New(path, mapping)
	}
	return idx, err
}
