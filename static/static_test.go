package static

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const tmpDirPrefix = "samples-golang-grpc-server-static-test"

func TestReadAll(t *testing.T) {
	const value = "hello"
	dir, err := ioutil.TempDir("", tmpDirPrefix)
	if err != nil {
		t.Fatal(err)
	}
	fname := filepath.Join(dir, "success.txt")
	if err := ioutil.WriteFile(fname, []byte(value), 0644); err != nil {
		t.Fatal(err)
	}

	data, err := ReadAll(http.Dir(dir), "success.txt")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(data)
	assert.Equal([]byte(value), data)
}
