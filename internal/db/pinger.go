package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Pinger - interface to check db health.
type Pinger interface {
	// Ping - pong.
	Ping(ctx context.Context) error
}

type mongoPinger struct {
	client *mongo.Client
}

// Ping - ping your db.
func (p *mongoPinger) Ping(ctx context.Context) error {
	const pingTimeout = 2 * time.Second

	ctx, cancel := context.WithTimeout(ctx, pingTimeout)

	err := p.client.Ping(ctx, readpref.PrimaryPreferred())
	cancel()

	return err
}

// NewMongoPinger - mongo pinger constructor.
func NewMongoPinger(opts ...MongoPingerOption) Pinger {
	m := &mongoPinger{}

	for i := range opts {
		opts[i](m)
	}

	return m
}
