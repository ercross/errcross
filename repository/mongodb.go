package repository

import (
	"context"
	"time"

	errs "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/ercross/errcross"
)

type mongoRepo struct {
	client       *mongo.Client
	databaseName string
	timeout      time.Duration
}

//newMongoClient creates a new mongo client. @Param timeout is in nanosecond according to Duration() doc
func newMongoClient(mongoUrl string, timeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout))
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		return nil, err
	}

	//check that client can be pinged
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoRepo(mongoUrl, dbName string, timeout int) (errcross.ErrcrossRepository, error) {
	mongoRepo := &mongoRepo{
		timeout:      time.Duration(timeout) * time.Second,
		databaseName: dbName,
	}
	client, err := newMongoClient(mongoUrl, timeout)
	if err != nil {
		return nil, errs.Wrap(err, "repository.NewMongoRepo")
	}
	mongoRepo.client = client
	return mongoRepo, nil
}

const dbName = "errcrossDb"

func (m *mongoRepo) Find(key string) (*errcross.Errcross, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	e := &errcross.Errcross{}
	collection := m.client.Database(m.databaseName).Collection(dbName)
	filter := bson.M{"key": key}
	err := collection.FindOne(ctx, filter).Decode(&e)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.Wrap(errcross.ErrKeyNotFound, "repository.mongodb.Find")
		}
		return nil, errs.Wrap(err, "repository.mongodb.Find")
	}
	return e, nil
}

func (m *mongoRepo) Store(e *errcross.Errcross) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	collection := m.client.Database(m.databaseName).Collection(dbName)
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"key":       e.Key,
			"url":       e.URL,
			"timestamp": e.Timestamp,
		},
	)
	if err != nil {
		return errs.Wrap(err, "repository.mongodb.Store")
	}
	return nil
}
