package storage

import (
	"bytes"
	"fmt"

	"go.etcd.io/bbolt"
)

type Unmarshaler func(data []byte) error
type Marshaler func() ([]byte, error)

type Storage struct {
	db *bbolt.DB
}

func New(db *bbolt.DB) *Storage {
	return &Storage{db}
}

func (s *Storage) Scan(namespace, prefix string, unmarshal Unmarshaler) error {
	return s.objectScan(namespace, prefix, unmarshal)
}

func (s *Storage) Get(namespace, key string, unmarshal Unmarshaler) error {
	return s.objectGet(namespace, key, unmarshal)
}

func (s *Storage) Put(namespace, key string, marshal Marshaler) error {
	return s.objectPut(namespace, key, marshal)
}

func (s *Storage) Delete(namespace, key string) error {
	return s.objectDelete(namespace, key)
}

func (s *Storage) objectScan(namespace, prefix string, unmarshal Unmarshaler) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(namespace))
		if b == nil {
			return fmt.Errorf("bucket %s does not exists", namespace)
		}
		c := b.Cursor()
		k := []byte(prefix)
		for k, v := c.Seek(k); k != nil && bytes.HasPrefix(k, k); k, v = c.Next() {
			if err := unmarshal(v); err != nil {
				return fmt.Errorf("failed to unmarshal %w", err)
			}
		}
		return nil
	})
}

func (s *Storage) objectGet(namespace, key string, unmarshal Unmarshaler) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(namespace))
		if b == nil {
			return fmt.Errorf("bucket %s does not exists", namespace)
		}
		v := b.Get([]byte(key))
		if v == nil {
			return fmt.Errorf("object %s/%s not found", namespace, key)
		}
		if err := unmarshal(v); err != nil {
			return fmt.Errorf("failed to unmarshal %w", err)
		}
		return nil
	})
}

func (s *Storage) objectPut(namespace, key string, marshal Marshaler) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(namespace))
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", namespace, err)
		}
		v, err := marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
		if err := b.Put([]byte(key), v); err != nil {
			return fmt.Errorf("failed to save %s/%s: %w", namespace, key, err)
		}
		return nil
	})
}

func (s *Storage) objectDelete(namespace, key string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(namespace))
		if b == nil {
			return fmt.Errorf("bucket %s does not exists", namespace)
		}
		if err := b.Delete([]byte(key)); err != nil {
			return fmt.Errorf("failed to delete %s/%s: %w", namespace, key, err)
		}
		return nil
	})
}
