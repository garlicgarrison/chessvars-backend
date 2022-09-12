package firestore

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// Exists is a Precondition that checks for the existence of a resource before
// writing to it. If the check fails, the write does not occur.
var Exists = firestore.Exists

// Desc sorts results from largest to smallest.
const Desc = firestore.Desc

type Query = firestore.Query

// Transaction represents a Firestore transaction.
type Transaction = firestore.Transaction

// A WriteBatch holds multiple database updates. Build a batch with the Create, Set,
// Update and Delete methods, then run it with the Commit method. Errors in Create,
// Set, Update or Delete are recorded instead of being returned immediately. The
// first such error is returned by Commit.
type WriteBatch = firestore.WriteBatch

// A WriteResult is returned by methods that write documents.
type WriteResult = firestore.WriteResult

type CollectionRef = firestore.CollectionRef
type DocumentRef = firestore.DocumentRef

// DocumentIterator is an iterator over documents returned by a query.
type DocumentIterator = firestore.DocumentIterator

// A DocumentSnapshot contains document data and metadata.
type DocumentSnapshot = firestore.DocumentSnapshot

// An Update describes an update to a value referred to by a path.
// An Update should have either a non-empty Path or a non-empty FieldPath,
// but not both.
//
// See DocumentRef.Create for acceptable values.
// To delete a field, specify firestore.Delete as the value.
type Update = firestore.Update

// A FieldPath is a non-empty sequence of non-empty fields that reference a value.
//
// A FieldPath value should only be necessary if one of the field names contains
// one of the runes ".˜*/[]". Most methods accept a simpler form of field path
// as a string in which the individual fields are separated by dots.
// For example,
//   []string{"a", "b"}
// is equivalent to the string form
//   "a.b"
// but
//   []string{"*"}
// has no equivalent string form.
type FieldPath = firestore.FieldPath

// Increment is an alias for FieldTransformIncrement.
//
// FieldTransformIncrement returns a special value that can be used with Set, Create, or
// Update that tells the server to transform the field's current value
// by the given value.
//
// The supported values are:
//
//    int, int8, int16, int32, int64
//    uint8, uint16, uint32
//    float32, float64
//
// If the field does not yet exist, the transformation will set the field to
// the given value.
var Increment = firestore.Increment

// MergeAll is a SetOption that causes all the field paths given in the data argument
// to Set to be overwritten. It is not supported for struct data.
var MergeAll = firestore.MergeAll

// Firestore is a fwaygo-kit wrapper over google firestore sdk
//
// It's main purpose is to make testing services that use it easier.
// For production firestore databases, Firestore should simply call the
// corresponding google firestore function.
type Firestore interface {
	// Collection creates a reference to a collection with the given path.
	// A path is a sequence of IDs separated by slashes.
	//
	// Collection returns nil if path contains an even number of IDs or any ID is empty.
	Collection(path string) *firestore.CollectionRef

	// RunTransaction runs f in a transaction. f should use the transaction it is given
	// for all Firestore operations. For any operation requiring a context, f should use
	// the context it is passed, not the first argument to RunTransaction.
	//
	// f must not call Commit or Rollback on the provided Transaction.
	//
	// If f returns nil, RunTransaction commits the transaction. If the commit fails due
	// to a conflicting transaction, RunTransaction retries f. It gives up and returns an
	// error after a number of attempts that can be configured with the MaxAttempts
	// option. If the commit succeeds, RunTransaction returns a nil error.
	//
	// If f returns non-nil, then the transaction will be rolled back and
	// this method will return the same error. The function f is not retried.
	//
	// Note that when f returns, the transaction is not committed. Calling code
	// must not assume that any of f's changes have been committed until
	// RunTransaction returns nil.
	//
	// Since f may be called more than once, f should usually be idempotent – that is, it
	// should have the same result when called multiple times.
	RunTransaction(context.Context, func(context.Context, *firestore.Transaction) error, ...firestore.TransactionOption) error

	// Batch returns a WriteBatch.
	Batch() *firestore.WriteBatch

	// Configure changes the firestore config
	//
	// This function unsafe for use in production since we don't currently keep track
	// of each request context.
	unsafe_configure(context.Context, Config) (*Document, error)
}

type Config struct {
	ProjectID  string            `envconfig:"FIRESTORE_PROJECT_ID"`
	Collection string            `envconfig:"FIRESTORE_COLLECTION"`
	Document   string            `envconfig:"FIRESTORE_DOCUMENT"`
	IsCreate   bool              `envconfig:"FIRESTORE_IS_CREATE"`
	Metadata   map[string]string `envconfig:"FIRESTORE_METADATA"`
}

type fs struct {
	client *firestore.Client

	isMock   bool
	mockLock sync.RWMutex

	collection string
	document   string
}

type Document struct {
	Collection string            `firestore:"collection"`
	Document   string            `firestore:"document"`
	Metadata   map[string]string `firestore:"metadata"`
	Timestamp  time.Time         `firestore:"timestamp"`
}

// NewClient returns a new Firestore instance.
//
// Firestore creates a new underlying firestore.Client. It also wraps all
// Collection calls such that is a level lower in the firestore database.
func NewClient(ctx context.Context, cfg *Config, opts ...option.ClientOption) (Firestore, error) {
	client, err := firestore.NewClient(ctx, cfg.ProjectID, opts...)
	if err != nil {
		return nil, err
	}

	var isMock bool
	if cfg.Collection != "" || cfg.Document != "" {
		if cfg.Collection != "" && cfg.Document != "" {
			isMock = true

			_, err = client.
				Collection(cfg.Collection).
				Doc(cfg.Document).
				Set(ctx, Document{
					Collection: cfg.Collection,
					Document:   cfg.Document,
					Metadata:   cfg.Metadata,
					Timestamp:  time.Now(),
				})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("invalid configs -- collection: %s, document: %s", cfg.Collection, cfg.Document)
		}
	}

	return &fs{
		client:     client,
		isMock:     isMock,
		collection: cfg.Collection,
		document:   cfg.Document,
	}, nil
}

func (fs *fs) Collection(name string) *firestore.CollectionRef {
	if fs.isMock {
		fs.mockLock.RLock()
		defer fs.mockLock.RUnlock()

		return fs.client.
			Collection(fs.collection).
			Doc(fs.document).
			Collection(name)
	}
	return fs.client.Collection(name)
}

func (fs *fs) RunTransaction(ctx context.Context, f func(context.Context, *firestore.Transaction) error, opts ...firestore.TransactionOption) error {
	fn := func(ctx context.Context, t *firestore.Transaction) error {
		span, ctx := tracer.StartSpanFromContext(ctx, "RunTransaction", tracer.ResourceName("inside"))
		defer span.Finish()

		return f(ctx, t)
	}

	span, ctx := tracer.StartSpanFromContext(ctx, "RunTransaction")
	defer span.Finish()

	return fs.client.RunTransaction(ctx, fn, opts...)
}

func (fs *fs) Batch() *firestore.WriteBatch {
	return fs.client.Batch()
}

func (fs *fs) unsafe_configure(ctx context.Context, config Config) (*Document, error) {
	if !fs.isMock {
		return nil, errors.New("not allowed -- only available for mock firestores")
	}

	fs.mockLock.Lock()
	defer fs.mockLock.Unlock()

	if config.Collection != "" {
		fs.collection = config.Collection
	}

	if config.Document != "" {
		fs.document = config.Document
	}

	doc := &Document{
		Collection: fs.collection,
		Document:   fs.document,
		Metadata:   config.Metadata,
		Timestamp:  time.Now(),
	}

	var err error
	if config.IsCreate {
		_, err = fs.client.
			Collection(fs.collection).
			Doc(fs.document).
			Create(ctx, Document{
				Collection: fs.collection,
				Document:   fs.document,
				Metadata:   config.Metadata,
				Timestamp:  time.Now(),
			})
	} else {
		_, err = fs.client.
			Collection(fs.collection).
			Doc(fs.document).
			Set(ctx, doc)
	}
	if err != nil {
		return nil, err
	}

	return doc, nil
}
