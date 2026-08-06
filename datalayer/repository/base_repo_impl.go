package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BaseCollection[T any] struct {
	coll *mongo.Collection
}

func (bc *BaseCollection[T]) CreateIndex(ctx context.Context, opts ...options.Lister[options.CreateIndexesOptions]) ([]string, error) {
	return nil, nil
}

func (bc *BaseCollection[T]) Name() string {
	return bc.coll.Name()
}

func (bc *BaseCollection[T]) Database() *mongo.Database {
	return bc.coll.Database()
}

func (bc *BaseCollection[T]) Collection() *mongo.Collection {
	return bc.coll
}

func (bc *BaseCollection[T]) BulkWrite(ctx context.Context, mods []mongo.WriteModel, opts ...options.Lister[options.BulkWriteOptions]) (*mongo.BulkWriteResult, error) {
	return bc.coll.BulkWrite(ctx, mods, opts...)
}

func (bc *BaseCollection[T]) InsertOne(ctx context.Context, doc *T, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error) {
	return bc.coll.InsertOne(ctx, doc, opts...)
}

func (bc *BaseCollection[T]) InsertMany(ctx context.Context, docs []*T, opts ...options.Lister[options.InsertManyOptions]) (*mongo.InsertManyResult, error) {
	return bc.coll.InsertMany(ctx, docs, opts...)
}

func (bc *BaseCollection[T]) DeleteOne(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error) {
	return bc.coll.DeleteOne(ctx, filter, opts...)
}

func (bc *BaseCollection[T]) DeleteMany(ctx context.Context, filter any, opts ...options.Lister[options.DeleteManyOptions]) (*mongo.DeleteResult, error) {
	return bc.coll.DeleteMany(ctx, filter, opts...)
}

func (bc *BaseCollection[T]) UpdateByID(ctx context.Context, id, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	return bc.coll.UpdateByID(ctx, id, update, opts...)
}

func (bc *BaseCollection[T]) UpdateOne(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	return bc.coll.UpdateOne(ctx, filter, update, opts...)
}

func (bc *BaseCollection[T]) UpdateMany(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error) {
	return bc.coll.UpdateMany(ctx, filter, update, opts...)
}

func (bc *BaseCollection[T]) ReplaceOne(ctx context.Context, filter, replacement any, opts ...options.Lister[options.ReplaceOptions]) (*mongo.UpdateResult, error) {
	return bc.coll.ReplaceOne(ctx, filter, replacement, opts...)
}

func (bc *BaseCollection[T]) Aggregate(ctx context.Context, pipe any, opts ...options.Lister[options.AggregateOptions]) (*mongo.Cursor, error) {
	return bc.coll.Aggregate(ctx, pipe, opts...)
}

func (bc *BaseCollection[T]) CountDocuments(ctx context.Context, filter any, opts ...options.Lister[options.CountOptions]) (int64, error) {
	return bc.coll.CountDocuments(ctx, filter, opts...)
}

func (bc *BaseCollection[T]) EstimatedDocumentCount(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error) {
	return bc.coll.EstimatedDocumentCount(ctx, opts...)
}

func (bc *BaseCollection[T]) Distinct(ctx context.Context, fieldName string, filter any, opts ...options.Lister[options.DistinctOptions]) *mongo.DistinctResult {
	return bc.coll.Distinct(ctx, fieldName, filter, opts...)
}

func (bc *BaseCollection[T]) Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) ([]*T, error) {
	cur, err := bc.coll.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	return bc.cursorTS(ctx, cur)
}

func (bc *BaseCollection[T]) FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) (*T, error) {
	t := new(T)
	dec := bc.coll.FindOne(ctx, filter, opts...)
	if err := dec.Decode(t); err != nil {
		return nil, err
	}

	return t, nil
}

func (bc *BaseCollection[T]) FindOneAndDelete(ctx context.Context, filter any, opts ...options.Lister[options.FindOneAndDeleteOptions]) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (bc *BaseCollection[T]) FindOneAndReplace(ctx context.Context, filter, replacement any, opts ...options.Lister[options.FindOneAndReplaceOptions]) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (bc *BaseCollection[T]) FindOneAndUpdate(ctx context.Context, filter, update any, opts ...options.Lister[options.FindOneAndUpdateOptions]) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (bc *BaseCollection[T]) FindTo(ctx context.Context, filter, result any, opts ...options.Lister[options.FindOptions]) error {
	cur, err := bc.coll.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}

	return bc.cursorAll(ctx, cur, result)
}

func (bc *BaseCollection[T]) Watch(ctx context.Context, pipeline any, opts ...options.Lister[options.ChangeStreamOptions]) (*mongo.ChangeStream, error) {
	return bc.coll.Watch(ctx, pipeline, opts...)
}

func (bc *BaseCollection[T]) Indexes() mongo.IndexView {
	return bc.coll.Indexes()
}

func (bc *BaseCollection[T]) SearchIndexes() mongo.SearchIndexView {
	return bc.coll.SearchIndexes()
}

func (bc *BaseCollection[T]) Drop(ctx context.Context, opts ...options.Lister[options.DropCollectionOptions]) error {
	return bc.coll.Drop(ctx)
}

func (bc *BaseCollection[T]) FindByID(ctx context.Context, id any, opts ...options.Lister[options.FindOneOptions]) (*T, error) {
	t := new(T)
	dec := bc.coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}}, opts...)
	if err := dec.Decode(t); err != nil {
		return nil, err
	}

	return t, nil
}

func (bc *BaseCollection[T]) DeleteByID(ctx context.Context, id any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error) {
	return bc.coll.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}}, opts...)
}

func (bc *BaseCollection[T]) DistinctString(ctx context.Context, fieldName string, filter any, opts ...options.Lister[options.DistinctOptions]) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (bc *BaseCollection[T]) DistinctObjectID(ctx context.Context, fieldName string, filter any, opts ...options.Lister[options.DistinctOptions]) ([]bson.ObjectID, error) {
	//TODO implement me
	panic("implement me")
}

func (bc *BaseCollection[T]) AggregateTo(ctx context.Context, pipe, result any, opts ...options.Lister[options.AggregateOptions]) error {
	//TODO implement me
	panic("implement me")
}

func (bc *BaseCollection[T]) Page(ctx context.Context, filter any, page, size int64, opts ...options.Lister[options.FindOptions]) (*Pages[T], error) {
	//TODO implement me
	panic("implement me")
}

func (bc *BaseCollection[T]) cursorTS(ctx context.Context, cur *mongo.Cursor) ([]*T, error) {
	var ret []*T
	if err := bc.cursorAll(ctx, cur, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

//goland:noinspection GoUnhandledErrorResult
func (bc *BaseCollection[T]) cursorAll(ctx context.Context, cur *mongo.Cursor, result any) error {
	defer cur.Close(ctx)
	return cur.All(ctx, result)
}
