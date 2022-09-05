package mongo

import (
	"context"
	"fmt"

	"auth-poc/svc/auth/config"
	"auth-poc/svc/auth/constants"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	Client *mongo.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}

func getURI(config config.Config) string {
	return fmt.Sprintf("mongodb://%s:%s/", config.MongoDBHost, config.MongoDBPort)
}

func New(config config.Config) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DEFAULT_TIMEOUT)
	clientOptions := options.Client().ApplyURI(getURI(config))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	return &DB{Client: client, Ctx: ctx, Cancel: cancel}, nil
}

func (d *DB) Close() {
	defer d.Cancel()
	defer func() {
		if err := d.Client.Disconnect(d.Ctx); err != nil {
			panic(err)
		}
	}()
}

func (d *DB) GetColl() *mongo.Collection {
	return d.Client.Database(constants.DB_NAME).Collection(constants.DB_COLL)
}
