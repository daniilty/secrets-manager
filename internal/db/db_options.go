package db

import "go.mongodb.org/mongo-driver/mongo"

// MongoDBOption - DI for mongo db.
type MongoDBOption func(*mongoDB)

// WithMongoCollection - set collection.
func WithMongoCollection(collection *mongo.Collection) MongoDBOption {
	return func(m *mongoDB) {
		m.collection = collection
	}
}
