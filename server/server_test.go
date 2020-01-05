package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/fullstorydev/grpcurl"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func init() {
	logger.Logger.SetOutput(ioutil.Discard)
}

func TestServerClose(t *testing.T) {
	const value = "hello"
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, value)
	})
	s := New("localhost:50001", "localhost:8080", nil, nil, f)
	if err := s.StartGrpcServer(func(*grpc.Server) error { return nil }); err != nil {
		t.Fatal(err)
	}
	httpF := func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error { return nil }
	if err := s.StartHttpServer(httpF); err != nil {
		t.Fatal(err)
	}
	s.Close()

	assert := assert.New(t)

	success, err := grpcTestConnection("localhost:50001")
	assert.Error(err)
	assert.False(success)

	r, err := http.Get("http://localhost:8080/")
	assert.Error(err)
	assert.Nil(r)
}

func TestServerGrpc(t *testing.T) {
	s := New("localhost:50001", "localhost:8080", nil, nil, nil)
	if err := s.StartGrpcServer(func(*grpc.Server) error { return nil }); err != nil {
		t.Fatal(err)
	}
	httpF := func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error { return nil }
	if err := s.StartHttpServer(httpF); err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	success, err := grpcTestConnection("localhost:50001")

	assert := assert.New(t)
	assert.NoError(err)
	assert.True(success)
}

func TestServerHttp(t *testing.T) {
	const value = "hello"
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, value)
	})
	s := New("localhost:50001", "localhost:8080", nil, nil, f)
	if err := s.StartGrpcServer(func(*grpc.Server) error { return nil }); err != nil {
		t.Fatal(err)
	}
	httpF := func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error { return nil }
	if err := s.StartHttpServer(httpF); err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	r, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(http.StatusOK, r.StatusCode)
	assert.Equal([]byte(value), data)
}

func grpcTestConnection(addr string) (bool, error) {
	ctx := context.Background()

	var creds credentials.TransportCredentials
	conn, err := grpcurl.BlockingDial(ctx, "tcp", addr, creds)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	refClient := grpcreflect.NewClient(ctx, reflectpb.NewServerReflectionClient(conn))
	descSource := grpcurl.DescriptorSourceFromServer(ctx, refClient)
	allServices, err := descSource.ListServices()
	if err != nil {
		return false, err
	}
	for _, svc := range allServices {
		if svc == "grpc.reflection.v1alpha.ServerReflection" {
			return true, nil
		}
	}
	return false, nil
}
