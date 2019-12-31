package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGatewayHandlers(t *testing.T) {
	const value = "hello"
	s := httptest.NewServer(gatewayHandler(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, value)
	})))
	defer s.Close()

	r, err := http.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(data)
	assert.Equal([]byte("hello"), data)
}
