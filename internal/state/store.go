package state

import (
	"errors"
	"io"
	"path"
	"time"

	"github.com/dgraph-io/badger/v4"
	"go.uber.org/zap"
)

const DefaultPrefix = "/"

var (
	ErrReceiverNotAPointer = errors.New("receiver is not a pointer")
	ErrReceiverNotASlice   = errors.New("receiver is not a slice")

	_ KVStore         = (*Store)(nil)
	_ TxnReaderWriter = (*badgerTxn)(nil)
)

type (
	TxnReader interface {
		Get(key string, v any) error
		GetAll(prefix string, into any) error
	}
	TxnWriter interface {
		Set(key string, v any, opts ...SetOption) error
	}
	TxnReaderWriter interface {
		TxnReader
		TxnWriter
	}
	KVStore interface {
		TxnReaderWriter
		io.Closer
		ReadOnlyTransaction(ops func(reader TxnReader) error) error
		ReadWriteTransaction(ops func(rw TxnReaderWriter) error) error
		WithSuffixes(suffixes ...string) KVStore
	}
)

type (
	StoreOptions struct {
		FilePath string
		InMemory bool
		Encoding EncoderDecoder
		Logger   badger.Logger
	}
	StoreOption interface {
		Apply(opt *StoreOptions)
	}
	StoreOptionFunc func(opt *StoreOptions)
)

func (f StoreOptionFunc) Apply(opt *StoreOptions) {
	f(opt)
}

func WithPath(filePath string) StoreOption {
	return StoreOptionFunc(func(opt *StoreOptions) {
		opt.FilePath = filePath
	})
}

func WithInMemory() StoreOption {
	return StoreOptionFunc(func(opt *StoreOptions) {
		opt.FilePath = ""
		opt.InMemory = true
	})
}

func WithEncoding(encoding EncoderDecoder) StoreOption {
	return StoreOptionFunc(func(opt *StoreOptions) {
		opt.Encoding = encoding
	})
}

func WithLogger(logger badger.Logger) StoreOption {
	return StoreOptionFunc(func(opt *StoreOptions) {
		opt.Logger = logger
	})
}

func NewDefault(opts ...StoreOption) (*Store, error) {
	const callerSkip = 2
	options := &StoreOptions{
		Encoding: MsgPackEncoding{},
		Logger:   Logger{SugaredLogger: zap.L().WithOptions(zap.AddCallerSkip(callerSkip)).Sugar()},
	}

	for idx := range opts {
		opts[idx].Apply(options)
	}

	if db, err := badger.Open(
		badger.DefaultOptions(options.FilePath).
			WithInMemory(options.InMemory).
			WithLogger(options.Logger),
	); err != nil {
		return nil, err
	} else {
		return &Store{
			db:       db,
			prefix:   DefaultPrefix,
			encoding: options.Encoding,
		}, nil
	}
}

type Store struct {
	db       *badger.DB
	encoding EncoderDecoder
	prefix   string
}

func (s *Store) ReadOnlyTransaction(ops func(reader TxnReader) error) error {
	return s.db.View(func(txn *badger.Txn) error {
		return ops(newBadgerTx(s.prefix, txn, s.encoding))
	})
}

func (s *Store) ReadWriteTransaction(ops func(rw TxnReaderWriter) error) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return ops(newBadgerTx(s.prefix, txn, s.encoding))
	})
}

func (s *Store) WithSuffixes(suffixes ...string) KVStore {
	joined := path.Join(suffixes...)
	return &Store{
		db:       s.db,
		encoding: s.encoding,
		prefix:   path.Join(s.prefix, joined),
	}
}

func (s *Store) Get(key string, v any) error {
	return s.db.View(func(txn *badger.Txn) (err error) {
		return newBadgerTx(s.prefix, txn, s.encoding).Get(key, v)
	})
}

func (s *Store) GetAll(prefix string, into any) error {
	return s.db.View(func(txn *badger.Txn) error {
		return newBadgerTx(s.prefix, txn, s.encoding).GetAll(prefix, into)
	})
}

type (
	SetOption interface {
		Apply(e *badger.Entry)
	}
	SetOptionFunc func(e *badger.Entry)
	WithTTL       time.Duration
)

func (f SetOptionFunc) Apply(e *badger.Entry) {
	f(e)
}

func (t WithTTL) Apply(e *badger.Entry) {
	e.WithTTL(time.Duration(t))
}

func (s *Store) Set(key string, v any, opts ...SetOption) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return newBadgerTx(s.prefix, txn, s.encoding).Set(key, v, opts...)
	})
}

func (s *Store) Close() error {
	return s.db.Close()
}

func itemKey(prefix, key string) []byte {
	return []byte(path.Join(prefix, key))
}
