package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/hellupline/samples-golang-grpc-server/server"
	"github.com/hellupline/samples-golang-grpc-server/tlsconfig"
)

const (
	dbName = "booru.db"
)

var (
	grpcAddr = flag.String("http-addr", "localhost:50051", "endpoint of the gRPC service")
	httpAddr = flag.String("grpc-addr", "localhost:8080", "endpoint of the http service")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		logrus.Error(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		signal.Notify(c, syscall.SIGINT)
		<-c
		cancel()
	}()

	tlsConfig, err := tlsconfig.LoadKeyPair()
	if err != nil {
		return err
	}
	s := server.New(*grpcAddr, *httpAddr, tlsConfig, http.DefaultServeMux)
	if err := s.StartGrpcServer(func(s *grpc.Server) error { return nil }); err != nil {
		return err
	}
	if err := s.StartHttpServer(); err != nil {
		return err
	}
	defer s.Close()

	<-ctx.Done()
	return ctx.Err()
}
