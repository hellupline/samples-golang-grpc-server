package static

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func ReadAll(fs http.FileSystem, fname string) ([]byte, error) {
	f, err := fs.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s: %w", fname, err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s: %w", fname, err)
	}

	return data, nil
}
