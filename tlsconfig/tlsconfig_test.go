package tlsconfig

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadKeyPair(t *testing.T) {
	const value = "hello"
	tlsconfig, err := LoadKeyPair(nil)
	if err != nil {
		t.Fatal(err)
	}

	s := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, value)
	}))
	s.TLS = tlsconfig
	s.StartTLS()
	defer s.Close()

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsconfig}}

	_, err = client.Get(s.URL)
	assert := assert.New(t)
	assert.NoError(err)
}
