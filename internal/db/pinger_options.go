package db

import "go.mongodb.org/mongo-driver/mongo"

type MongoPingerOption func(*mongoPinger)

func WithMongoClient(client *mongo.Client) MongoPingerOption {
	return func(m *mongoPinger) {
		m.client = client
	}
}
