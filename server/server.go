package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/hellupline/samples-golang-grpc-server/server/health"
)

var logger = logrus.WithField("module", "server")

type grpcRegister func(*grpc.Server) error
type httpRegister func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

type Server struct {
	mux        *runtime.ServeMux
	grpcServer *grpc.Server
	httpServer *http.Server
	grpcAddr   string
	httpAddr   string
	tlsconfig  *tls.Config
}

func New(grpcAddr, httpAddr string, tlsconfig *tls.Config, next http.Handler) *Server {
	grpcopt := []grpc.ServerOption{}
	if tlsconfig != nil {
		grpcopt = append(grpcopt, grpc.Creds(credentials.NewTLS(tlsconfig)))
	}
	muxopt := []runtime.ServeMuxOption{
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
	}
	mux := runtime.NewServeMux(muxopt...)
	handler := gatewayHandler(mux, next)

	grpcServer := grpc.NewServer(grpcopt...)
	httpServer := &http.Server{TLSConfig: tlsconfig, Handler: handler}

	return &Server{mux, grpcServer, httpServer, grpcAddr, httpAddr, tlsconfig}
}

func (s *Server) Close() {
	logger.Info("Shutting down the http server")
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		logger.Errorf("Failed to shutdown http server: %v", err)
	}
	logger.Info("Shutting down the grpc server")
	s.grpcServer.GracefulStop()
}

func (s *Server) StartGrpcServer(register grpcRegister) error {
	if err := register(s.grpcServer); err != nil {
		return err
	}
	grpc_health_v1.RegisterHealthServer(s.grpcServer, health.New())
	reflection.Register(s.grpcServer)
	lis, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return fmt.Errorf("error listening on address %s: %w", s.grpcAddr, err)
	}
	go func() {
		logger.Infof("listening on %s", lis.Addr())
		if err := s.grpcServer.Serve(lis); err != nil {
			logger.WithError(err).Error("failed to serve grpc")
		}
	}()
	return nil
}

func (s *Server) StartHttpServer(register httpRegister) error {
	cctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	o := grpc.WithInsecure()
	if s.tlsconfig != nil {
		o = grpc.WithTransportCredentials(credentials.NewTLS(s.tlsconfig))
	}
	conn, err := grpc.DialContext(cctx, s.grpcAddr, o, grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to dial grpc server: %w", err)
	}
	if err := register(cctx, s.mux, conn); err != nil {
		return err
	}
	lis, err := net.Listen("tcp", s.httpAddr)
	if err != nil {
		return fmt.Errorf("error listening on address %s: %w", s.httpAddr, err)
	}
	go func() {
		logger.Infof("listening on %s", lis.Addr())
		if s.tlsconfig != nil {
			if err := s.httpServer.ServeTLS(lis, "", ""); err != http.ErrServerClosed {
				logger.WithError(err).Error("failed to serve https")
			}
		} else {
			if err := s.httpServer.Serve(lis); err != http.ErrServerClosed {
				logger.WithError(err).Error("failed to serve http")
			}
		}
	}()
	return nil
}
