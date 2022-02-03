package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// DB - interface for k/v db.
// TODO: fuck mongo, use leveldb instead.
type DB interface {
	// Set - k, v.
	Set(context.Context, string, string) error
	// Get - key -> value.
	Get(context.Context, string) (string, error)
}

type mongoDB struct {
	collection *mongo.Collection
}

func NewMongoDB(opts ...MongoDBOption) DB {
	m := &mongoDB{}

	for i := range opts {
		opts[i](m)
	}

	return m
}
