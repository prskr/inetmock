package state

import (
	"reflect"

	"github.com/dgraph-io/badger/v3"
)

type BadgerTransaction interface {
	NewIterator(opt badger.IteratorOptions) *badger.Iterator
	Get(key []byte) (item *badger.Item, rerr error)
	SetEntry(e *badger.Entry) error
}

func newBadgerTx(prefix string, txn *badger.Txn, encoding EncoderDecoder) *badgerTxn {
	return &badgerTxn{
		prefix:   prefix,
		txn:      txn,
		encoding: encoding,
	}
}

type badgerTxn struct {
	prefix   string
	txn      BadgerTransaction
	encoding EncoderDecoder
}

func (t *badgerTxn) Get(key string, v interface{}) (err error) {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return ErrReceiverNotAPointer
	}
	itemKey := itemKey(t.prefix, key)
	var item *badger.Item
	if item, err = t.txn.Get(itemKey); err != nil {
		return err
	}

	return item.Value(func(val []byte) error {
		return t.encoding.Decode(val, v)
	})
}

func (t *badgerTxn) GetAll(prefix string, into interface{}) error {
	sliceType := reflect.TypeOf(into)
	if sliceType.Kind() != reflect.Ptr {
		return ErrReceiverNotAPointer
	}

	if sliceType.Elem().Kind() != reflect.Slice {
		return ErrReceiverNotASlice
	}

	slicePtr := reflect.ValueOf(into)
	sliceVal := slicePtr.Elem()
	elemType := sliceType.Elem().Elem()

	keyPrefix := itemKey(t.prefix, prefix)
	it := t.txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	for it.Seek(keyPrefix); it.ValidForPrefix(keyPrefix); it.Next() {
		item := it.Item()
		if err := item.Value(func(val []byte) error {
			entity := reflect.New(elemType)
			if err := t.encoding.Decode(val, entity.Interface()); err != nil {
				return err
			}

			sliceVal.Set(reflect.Append(sliceVal, entity.Elem()))

			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func (t *badgerTxn) Set(key string, v interface{}, opts ...SetOption) error {
	itemKey := itemKey(t.prefix, key)
	data, err := t.encoding.Encode(v)
	if err != nil {
		return err
	}
	e := badger.NewEntry(itemKey, data)
	for idx := range opts {
		opts[idx].Apply(e)
	}
	return t.txn.SetEntry(e)
}
