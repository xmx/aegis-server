package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository[T any] interface {
	// Clone creates a copy of the Collection configured with the given CollectionOptions.
	// The specified options are merged with the existing options on the collection, with the specified options taking
	// precedence.
	Clone(opts ...options.Lister[options.CollectionOptions]) *mongo.Collection

	// Name returns the name of the collection.
	Name() string

	// Database returns the Database that was used to create the Collection.
	Database() *mongo.Database

	// BulkWrite performs a bulk write operation (https://www.mongodb.com/docs/manual/core/bulk-write-operations/).
	//
	// The models parameter must be a slice of operations to be executed in this bulk write. It cannot be nil or empty.
	// All of the models must be non-nil. See the mongo.WriteModel documentation for a list of valid model types and
	// examples of how they should be used.
	//
	// The opts parameter can be used to specify options for the operation (see the options.BulkWriteOptions documentation.)
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...options.Lister[options.BulkWriteOptions]) (*mongo.BulkWriteResult, error)

	// InsertOne executes an insert command to insert a single document into the collection.
	//
	// The document parameter must be the document to be inserted. It cannot be nil. If the document does not have an _id
	// field when transformed into BSON, one will be added automatically to the marshalled document. The original document
	// will not be modified. The _id can be retrieved from the InsertedID field of the returned InsertOneResult.
	//
	// The opts parameter can be used to specify options for the operation (see the options.InsertOneOptions documentation.)
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/insert/.
	InsertOne(ctx context.Context, doc *T, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error)

	// InsertMany executes an insert command to insert multiple documents into the collection. If write errors occur
	// during the operation (e.g. duplicate key error), this method returns a BulkWriteException error.
	//
	// The documents parameter must be a slice of documents to insert. The slice cannot be nil or empty. The elements must
	// all be non-nil. For any document that does not have an _id field when transformed into BSON, one will be added
	// automatically to the marshalled document. The original document will not be modified. The _id values for the inserted
	// documents can be retrieved from the InsertedIDs field of the returned InsertManyResult.
	//
	// The opts parameter can be used to specify options for the operation (see the options.InsertManyOptions documentation.)
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/insert/.
	InsertMany(ctx context.Context, docs []*T, opts ...options.Lister[options.InsertManyOptions]) (*mongo.InsertManyResult, error)

	// DeleteOne executes a delete command to delete at most one document from the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// deleted. It cannot be nil. If the filter does not match any documents, the operation will succeed and a DeleteResult
	// with a DeletedCount of 0 will be returned. If the filter matches multiple documents, one will be selected from the
	// matched set.
	//
	// The opts parameter can be used to specify options for the operation (see the options.DeleteOptions documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/delete/.
	DeleteOne(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error)

	// DeleteMany executes a delete command to delete documents from the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the documents to
	// be deleted. It cannot be nil. An empty document (e.g. bson.D{}) should be used to delete all documents in the
	// collection. If the filter does not match any documents, the operation will succeed and a DeleteResult with a
	// DeletedCount of 0 will be returned.
	//
	// The opts parameter can be used to specify options for the operation (see the options.DeleteOptions documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/delete/.
	DeleteMany(ctx context.Context, filter any, opts ...options.Lister[options.DeleteManyOptions]) (*mongo.DeleteResult, error)

	// UpdateByID executes an update command to update the document whose _id value matches the provided ID in the collection.
	// This is equivalent to running UpdateOne(ctx, bson.D{{"_id", id}}, update, opts...).
	//
	// The id parameter is the _id of the document to be updated. It cannot be nil. If the ID does not match any documents,
	// the operation will succeed and an UpdateResult with a MatchedCount of 0 will be returned.
	//
	// The update parameter must be a document containing update operators
	// (https://www.mongodb.com/docs/manual/reference/operator/update/) and can be used to specify the modifications to be
	// made to the selected document. It cannot be nil or empty.
	//
	// The opts parameter can be used to specify options for the operation (see the options.UpdateOptions documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/update/.
	UpdateByID(ctx context.Context, id, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)

	// UpdateOne executes an update command to update at most one document in the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// updated. It cannot be nil. If the filter does not match any documents, the operation will succeed and an UpdateResult
	// with a MatchedCount of 0 will be returned. If the filter matches multiple documents, one will be selected from the
	// matched set and MatchedCount will equal 1.
	//
	// The update parameter must be a document containing update operators
	// (https://www.mongodb.com/docs/manual/reference/operator/update/) and can be used to specify the modifications to be
	// made to the selected document. It cannot be nil or empty.
	//
	// The opts parameter can be used to specify options for the operation (see the options.UpdateOptions documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/update/.
	UpdateOne(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)

	// UpdateMany executes an update command to update documents in the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the documents to be
	// updated. It cannot be nil. If the filter does not match any documents, the operation will succeed and an UpdateResult
	// with a MatchedCount of 0 will be returned.
	//
	// The update parameter must be a document containing update operators
	// (https://www.mongodb.com/docs/manual/reference/operator/update/) and can be used to specify the modifications to be made
	// to the selected documents. It cannot be nil or empty.
	//
	// The opts parameter can be used to specify options for the operation (see the options.UpdateOptions documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/update/.
	UpdateMany(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error)

	// ReplaceOne executes an update command to replace at most one document in the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// replaced. It cannot be nil. If the filter does not match any documents, the operation will succeed and an
	// UpdateResult with a MatchedCount of 0 will be returned. If the filter matches multiple documents, one will be
	// selected from the matched set and MatchedCount will equal 1.
	//
	// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
	// and cannot contain any update operators (https://www.mongodb.com/docs/manual/reference/operator/update/).
	//
	// The opts parameter can be used to specify options for the operation (see the options.ReplaceOptions documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/update/.
	ReplaceOne(ctx context.Context, filter, replacement any, opts ...options.Lister[options.ReplaceOptions]) (*mongo.UpdateResult, error)

	// Aggregate executes an aggregate command against the collection and returns a cursor over the resulting documents.
	//
	// The pipeline parameter must be an array of documents, each representing an aggregation stage. The pipeline cannot
	// be nil but can be empty. The stage documents must all be non-nil. For a pipeline of bson.D documents, the
	// mongo.Pipeline type can be used. See
	// https://www.mongodb.com/docs/manual/reference/operator/aggregation-pipeline/#db-collection-aggregate-stages for a list of
	// valid stages in aggregations.
	//
	// The opts parameter can be used to specify options for the operation (see the options.AggregateOptions documentation.)
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/aggregate/.
	Aggregate(ctx context.Context, pipeline any, opts ...options.Lister[options.AggregateOptions]) (*mongo.Cursor, error)

	// CountDocuments returns the number of documents in the collection. For a fast count of the documents in the
	// collection, see the EstimatedDocumentCount method.
	//
	// The filter parameter must be a document and can be used to select which documents contribute to the count. It
	// cannot be nil. An empty document (e.g. bson.D{}) should be used to count all documents in the collection. This will
	// result in a full collection scan.
	//
	// The opts parameter can be used to specify options for the operation (see the options.CountOptions documentation).
	CountDocuments(ctx context.Context, filter any, opts ...options.Lister[options.CountOptions]) (int64, error)

	// EstimatedDocumentCount executes a count command and returns an estimate of the number of documents in the collection
	// using collection metadata.
	//
	// The opts parameter can be used to specify options for the operation (see the options.EstimatedDocumentCountOptions
	// documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/count/.
	EstimatedDocumentCount(ctx context.Context, opts ...options.Lister[options.EstimatedDocumentCountOptions]) (int64, error)

	// Distinct executes a distinct command to find the unique values for a specified field in the collection.
	//
	// The fieldName parameter specifies the field name for which distinct values should be returned.
	//
	// The filter parameter must be a document containing query operators and can be used to select which documents are
	// considered. It cannot be nil. An empty document (e.g. bson.D{}) should be used to select all documents.
	//
	// The opts parameter can be used to specify options for the operation (see the options.DistinctOptions documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/distinct/.
	Distinct(ctx context.Context, fieldName string, filter any, opts ...options.Lister[options.DistinctOptions]) *mongo.DistinctResult

	// Find executes a find command and returns a Cursor over the matching documents in the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select which documents are
	// included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all documents.
	//
	// The opts parameter can be used to specify options for the operation (see the options.FindOptions documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/find/.
	Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) ([]*T, error)

	// FindOne executes a find command and returns a SingleResult for one document in the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// returned. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
	// ErrNoDocuments will be returned. If the filter matches multiple documents, one will be selected from the matched set.
	//
	// The opts parameter can be used to specify options for this operation (see the options.FindOneOptions documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/find/.
	FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) (*T, error)

	// FindOneAndDelete executes a findAndModify command to delete at most one document in the collection. and returns the
	// document as it appeared before deletion.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// deleted. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
	// ErrNoDocuments wil be returned. If the filter matches multiple documents, one will be selected from the matched set.
	//
	// The opts parameter can be used to specify options for the operation (see the options.FindOneAndDeleteOptions
	// documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/findAndModify/.
	FindOneAndDelete(ctx context.Context, filter any, opts ...options.Lister[options.FindOneAndDeleteOptions]) (*T, error)

	// FindOneAndReplace executes a findAndModify command to replace at most one document in the collection
	// and returns the document as it appeared before replacement.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// replaced. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
	// ErrNoDocuments wil be returned. If the filter matches multiple documents, one will be selected from the matched set.
	//
	// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
	// and cannot contain any update operators (https://www.mongodb.com/docs/manual/reference/operator/update/).
	//
	// The opts parameter can be used to specify options for the operation (see the options.FindOneAndReplaceOptions
	// documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/findAndModify/.
	FindOneAndReplace(ctx context.Context, filter, replacement any, opts ...options.Lister[options.FindOneAndReplaceOptions]) (*T, error)

	// FindOneAndUpdate executes a findAndModify command to update at most one document in the collection and returns the
	// document as it appeared before updating.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// updated. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
	// ErrNoDocuments wil be returned. If the filter matches multiple documents, one will be selected from the matched set.
	//
	// The update parameter must be a document containing update operators
	// (https://www.mongodb.com/docs/manual/reference/operator/update/) and can be used to specify the modifications to be made
	// to the selected document. It cannot be nil or empty.
	//
	// The opts parameter can be used to specify options for the operation (see the options.FindOneAndUpdateOptions
	// documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/findAndModify/.
	FindOneAndUpdate(ctx context.Context, filter, update any, opts ...options.Lister[options.FindOneAndUpdateOptions]) (*T, error)

	// Watch returns a change stream for all changes on the corresponding collection. See
	// https://www.mongodb.com/docs/manual/changeStreams/ for more information about change streams.
	//
	// The Collection must be configured with read concern majority or no read concern for a change stream to be created
	// successfully.
	//
	// The pipeline parameter must be an array of documents, each representing a pipeline stage. The pipeline cannot be
	// nil but can be empty. The stage documents must all be non-nil. See https://www.mongodb.com/docs/manual/changeStreams/ for
	// a list of pipeline stages that can be used with change streams. For a pipeline of bson.D documents, the
	// mongo.Pipeline{} type can be used.
	//
	// The opts parameter can be used to specify options for change stream creation (see the options.ChangeStreamOptions
	// documentation).
	Watch(ctx context.Context, pipeline any, opts ...options.Lister[options.ChangeStreamOptions]) (*mongo.ChangeStream, error)

	// Indexes returns an IndexView instance that can be used to perform operations on the indexes for the collection.
	Indexes() mongo.IndexView

	// SearchIndexes returns a SearchIndexView instance that can be used to perform operations on the search indexes for the collection.
	SearchIndexes() mongo.SearchIndexView

	// Drop drops the collection on the server. This method ignores "namespace not found" errors so it is safe to drop
	// a collection that does not exist on the server.
	Drop(ctx context.Context, opts ...options.Lister[options.DropCollectionOptions]) error

	FindByID(ctx context.Context, id any, opts ...options.Lister[options.FindOneOptions]) (*T, error)

	FindPage(ctx context.Context, filter any, page, size int64, opts ...options.Lister[options.FindOptions]) (*Pages[T], error)

	IndexCreator
}

type IndexCreator interface {
	CreateIndex(ctx context.Context) error
}

func CreateIndex(ctx context.Context, indexes []IndexCreator) error {
	for _, idx := range indexes {
		if err := idx.CreateIndex(ctx); err != nil {
			return err
		}
	}

	return nil
}
