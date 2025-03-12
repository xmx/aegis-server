package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func newBaseRepository[T any](db *mongo.Database, name string, opts ...options.Lister[options.CollectionOptions]) *baseRepository[T] {
	coll := db.Collection(name, opts...)
	return &baseRepository[T]{coll: coll}
}

type baseRepository[T any] struct {
	coll *mongo.Collection
}

func (br *baseRepository[T]) Clone(opts ...options.Lister[options.CollectionOptions]) *mongo.Collection {
	return br.coll.Clone(opts...)
}

func (br *baseRepository[T]) Name() string {
	return br.coll.Name()
}

func (br *baseRepository[T]) Database() *mongo.Database {
	return br.coll.Database()
}

func (br *baseRepository[T]) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...options.Lister[options.BulkWriteOptions]) (*mongo.BulkWriteResult, error) {
	return br.coll.BulkWrite(ctx, models, opts...)
}

func (br *baseRepository[T]) InsertOne(ctx context.Context, doc *T, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error) {
	return br.coll.InsertOne(ctx, doc, opts...)
}

func (br *baseRepository[T]) InsertMany(ctx context.Context, docs []*T, opts ...options.Lister[options.InsertManyOptions]) (*mongo.InsertManyResult, error) {
	return br.coll.InsertMany(ctx, docs, opts...)
}

func (br *baseRepository[T]) DeleteOne(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error) {
	return br.coll.DeleteOne(ctx, filter, opts...)
}

func (br *baseRepository[T]) DeleteMany(ctx context.Context, filter any, opts ...options.Lister[options.DeleteManyOptions]) (*mongo.DeleteResult, error) {
	return br.coll.DeleteMany(ctx, filter, opts...)
}

func (br *baseRepository[T]) UpdateByID(ctx context.Context, id, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	return br.coll.UpdateByID(ctx, id, update, opts...)
}

func (br *baseRepository[T]) UpdateOne(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	return br.coll.UpdateOne(ctx, filter, update, opts...)
}

func (br *baseRepository[T]) UpdateMany(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error) {
	return br.coll.UpdateMany(ctx, filter, update, opts...)
}

func (br *baseRepository[T]) ReplaceOne(ctx context.Context, filter, replacement any, opts ...options.Lister[options.ReplaceOptions]) (*mongo.UpdateResult, error) {
	return br.coll.ReplaceOne(ctx, filter, replacement, opts...)
}

func (br *baseRepository[T]) Aggregate(ctx context.Context, pipeline any, opts ...options.Lister[options.AggregateOptions]) (*mongo.Cursor, error) {
	return br.coll.Aggregate(ctx, pipeline, opts...)
}

func (br *baseRepository[T]) CountDocuments(ctx context.Context, filter any, opts ...options.Lister[options.CountOptions]) (int64, error) {
	return br.coll.CountDocuments(ctx, filter, opts...)
}

func (br *baseRepository[T]) EstimatedDocumentCount(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error) {
	return br.coll.EstimatedDocumentCount(ctx, opts...)
}

func (br *baseRepository[T]) Distinct(ctx context.Context, fieldName string, filter any, opts ...options.Lister[options.DistinctOptions]) *mongo.DistinctResult {
	return br.coll.Distinct(ctx, fieldName, filter, opts...)
}

func (br *baseRepository[T]) Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) ([]*T, error) {
	cur, err := br.coll.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	return br.decodeCursor(ctx, cur)
}

func (br *baseRepository[T]) FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) (*T, error) {
	return br.decodeSingleResult(br.coll.FindOne(ctx, filter, opts...))
}

func (br *baseRepository[T]) FindOneAndDelete(ctx context.Context, filter any, opts ...options.Lister[options.FindOneAndDeleteOptions]) (*T, error) {
	return br.decodeSingleResult(br.coll.FindOneAndDelete(ctx, filter, opts...))
}

func (br *baseRepository[T]) FindOneAndReplace(ctx context.Context, filter, replacement any, opts ...options.Lister[options.FindOneAndReplaceOptions]) (*T, error) {
	return br.decodeSingleResult(br.coll.FindOneAndReplace(ctx, filter, replacement, opts...))
}

func (br *baseRepository[T]) FindOneAndUpdate(ctx context.Context, filter, update any, opts ...options.Lister[options.FindOneAndUpdateOptions]) (*T, error) {
	return br.decodeSingleResult(br.coll.FindOneAndUpdate(ctx, filter, update, opts...))
}

func (br *baseRepository[T]) Watch(ctx context.Context, pipeline any, opts ...options.Lister[options.ChangeStreamOptions]) (*mongo.ChangeStream, error) {
	return br.coll.Watch(ctx, pipeline, opts...)
}

func (br *baseRepository[T]) Indexes() mongo.IndexView {
	return br.coll.Indexes()
}

func (br *baseRepository[T]) SearchIndexes() mongo.SearchIndexView {
	return br.coll.SearchIndexes()
}

func (br *baseRepository[T]) Drop(ctx context.Context, opts ...options.Lister[options.DropCollectionOptions]) error {
	return br.coll.Drop(ctx, opts...)
}

func (br *baseRepository[T]) CreateIndex(context.Context) error {
	return nil
}

func (br *baseRepository[T]) FindByID(ctx context.Context, id any, opts ...options.Lister[options.FindOneOptions]) (*T, error) {
	return br.FindOne(ctx, bson.D{{Key: "_id", Value: id}}, opts...)
}

func (br *baseRepository[T]) FindPage(ctx context.Context, filter any, page, size int64, opts ...options.Lister[options.FindOptions]) (*Pages[T], error) {
	cnt, err := br.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	pages := NewPages[T](page, size)
	if cnt == 0 {
		return pages, nil
	}

	skip := pages.Skip(cnt)
	opt := options.Find().SetSkip(skip).SetLimit(size)
	opts = append(opts, opt)
	ts, err := br.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	pages.Total, pages.Records = cnt, ts

	return pages, nil
}

func (*baseRepository[T]) decodeSingleResult(ret *mongo.SingleResult) (*T, error) {
	t := new(T)
	if err := ret.Decode(t); err != nil {
		return nil, err
	}

	return t, nil
}

func (*baseRepository[T]) decodeCursor(ctx context.Context, cur *mongo.Cursor) ([]*T, error) {
	//goland:noinspection GoUnhandledErrorResult
	defer cur.Close(context.Background())

	ts := make([]*T, 0, 10)
	if err := cur.All(ctx, &ts); err != nil {
		return nil, err
	}

	return ts, nil
}
