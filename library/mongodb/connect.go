package mongodb

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

func Connect(opts ...*options.ClientOptions) (*mongo.Database, error) {
	opt := options.MergeClientOptions(opts...)
	uri := opt.GetURI()
	cs, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return nil, err
	}

	cli, err := mongo.Connect(opt)
	if err != nil {
		return nil, err
	}

	dbname := cs.Database
	db := cli.Database(dbname)

	return db, nil
}
