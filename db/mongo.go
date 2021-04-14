package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConf struct {
	Uri           string `json:"mongo.Uri"`
	TimeoutSecond int    `json:"mongo.TimeoutSecond"`
}

func NewVerifiedMongoCLI(ctx context.Context, conf *MongoConf) (*mongo.Client, error) {
	cli, err := mongo.NewClient(options.Client().ApplyURI(conf.Uri))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(conf.TimeoutSecond)*time.Second)
	defer cancel()
	if err = cli.Connect(ctx); err != nil {
		cli = nil
	}
	return cli, err
}
