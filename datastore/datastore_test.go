package datastore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
)

func TestPrefixOperations(t *testing.T) {
	const name = "hello"
	const value = "body"
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := bbolt.Open(path.Join(dir, "favorite-search.db"), 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	assert := assert.New(t)

	ds := New(db)

	for i := 0; i < 10; i++ {
		o := &TestData{Name: fmt.Sprintf("key-%03d", i), Value: value}
		err = ds.objectPut("testing", o.Name, func() ([]byte, error) {
			return json.Marshal(o)
		})
		assert.NoError(err)
	}

	objs := make([]*TestData, 0)
	ds.objectScan("testing", "key-", func(v []byte) error {
		obj := &TestData{}
		if err := json.Unmarshal(v, obj); err != nil {
			return err
		}
		objs = append(objs, obj)
		return nil
	})
	for i, row := range objs {
		assert.Equal(fmt.Sprintf("key-%03d", i), row.Name)
		assert.Equal(value, row.Value)
	}
	assert.Len(objs, 10)
}

func TestObjectOperations(t *testing.T) {
	const name = "hello"
	const value = "body"
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	db, err := bbolt.Open(path.Join(dir, "favorite-search.db"), 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	assert := assert.New(t)

	ds := New(db)

	{
		var err error
		err = ds.objectPut("testing", name, func() ([]byte, error) {
			return json.Marshal(TestData{Name: name, Value: value})
		})
		assert.NoError(err)
		obj := &TestData{}
		err = ds.objectGet("testing", name, func(v []byte) error {
			return json.Unmarshal(v, obj)
		})
		assert.NoError(err)
		assert.Equal(name, obj.Name)
		assert.Equal(value, obj.Value)
	}
	{
		var err error
		err = ds.objectDelete("testing", name)
		assert.NoError(err)
		obj := &TestData{}
		err = ds.objectGet("testing", name, func(v []byte) error {
			return json.Unmarshal(v, obj)
		})
		assert.Error(err)
		assert.NotEqual(name, obj.Name)
		assert.NotEqual(value, obj.Value)
	}
}

type TestData struct{ Name, Value string }
