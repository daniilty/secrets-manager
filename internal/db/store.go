package db

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoVal struct {
	AppName string `bson:"an"`
	Value   string `bson:"v"`
}

func (m *mongoVal) toBSOND() bson.D {
	return bson.D{
		{Key: "an", Value: m.AppName},
		{Key: "v", Value: m.Value},
	}
}

func (m *mongoDB) Set(ctx context.Context, appName string, value string) error {
	filter := bson.M{"an": appName}
	updateStruct := &mongoVal{
		AppName: appName,
		Value:   value,
	}
	update := bson.M{"$set": updateStruct.toBSOND()}
	opts := options.Update().SetUpsert(true)

	_, err := m.collection.UpdateOne(ctx, filter, update, opts)

	return err
}

func (m *mongoDB) Get(ctx context.Context, appName string) (string, error) {
	filter := bson.M{"an": appName}

	res := m.collection.FindOne(ctx, filter)
	v := &mongoVal{}

	err := res.Decode(v)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", nil
		}

		return "", err
	}

	return v.Value, nil
}
