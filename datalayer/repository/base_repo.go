package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CreateIndexer interface {
	CreateIndex(context.Context, ...options.Lister[options.CreateIndexesOptions]) ([]string, error)
}

type Collection[T any] interface {
	CreateIndexer
	Name() string
	Database() *mongo.Database
	Collection() *mongo.Collection
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...options.Lister[options.BulkWriteOptions]) (*mongo.BulkWriteResult, error)
	InsertOne(ctx context.Context, doc *T, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, docs []*T, opts ...options.Lister[options.InsertManyOptions]) (*mongo.InsertManyResult, error)
	DeleteOne(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter any, opts ...options.Lister[options.DeleteManyOptions]) (*mongo.DeleteResult, error)
	UpdateByID(ctx context.Context, id, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
	UpdateOne(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error)
	ReplaceOne(ctx context.Context, filter, replacement any, opts ...options.Lister[options.ReplaceOptions]) (*mongo.UpdateResult, error)
	Aggregate(ctx context.Context, pipe any, opts ...options.Lister[options.AggregateOptions]) (*mongo.Cursor, error)
	CountDocuments(ctx context.Context, filter any, opts ...options.Lister[options.CountOptions]) (int64, error)
	EstimatedDocumentCount(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error)
	Distinct(ctx context.Context, fieldName string, filter any, opts ...options.Lister[options.DistinctOptions]) *mongo.DistinctResult
	Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) ([]*T, error)
	FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) (*T, error)
	FindOneAndDelete(ctx context.Context, filter any, opts ...options.Lister[options.FindOneAndDeleteOptions]) (*T, error)
	FindOneAndReplace(ctx context.Context, filter, replacement any, opts ...options.Lister[options.FindOneAndReplaceOptions]) (*T, error)
	FindOneAndUpdate(ctx context.Context, filter, update any, opts ...options.Lister[options.FindOneAndUpdateOptions]) (*T, error)
	FindTo(ctx context.Context, filter, result any, opts ...options.Lister[options.FindOptions]) error
	Watch(ctx context.Context, pipeline any, opts ...options.Lister[options.ChangeStreamOptions]) (*mongo.ChangeStream, error)
	Indexes() mongo.IndexView
	SearchIndexes() mongo.SearchIndexView
	Drop(ctx context.Context, opts ...options.Lister[options.DropCollectionOptions]) error
	FindByID(ctx context.Context, id any, opts ...options.Lister[options.FindOneOptions]) (*T, error)
	DeleteByID(ctx context.Context, id any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error)
	DistinctString(ctx context.Context, fieldName string, filter any, opts ...options.Lister[options.DistinctOptions]) ([]string, error)
	DistinctObjectID(ctx context.Context, fieldName string, filter any, opts ...options.Lister[options.DistinctOptions]) ([]bson.ObjectID, error)
	AggregateTo(ctx context.Context, pipe, result any, opts ...options.Lister[options.AggregateOptions]) error
	Page(ctx context.Context, filter any, page, size int64, opts ...options.Lister[options.FindOptions]) (*Pages[T], error)
}

type Pages[T any] struct {
	Page    int64 `json:"page"`
	Size    int64 `json:"size"`
	Total   int64 `json:"total"`
	Records []*T  `json:"records,omitzero"`
}

func NewEmptyPages[T any]() *Pages[T] {
	return &Pages[T]{
		Page: 1,
		Size: 10,
	}
}
