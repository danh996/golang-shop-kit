package repository

import "context"

// WithIndexes enforce indexes creation
type WithIndexes interface {
	CreateIndexes(context.Context) error
}
