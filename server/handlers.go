package server

import (
	"compress/gzip"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// gatewayHandler create a http.Handler for server middleware
func gatewayHandler(gwmux *runtime.ServeMux, next http.Handler) http.Handler {
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	r := mux.NewRouter()
	r.PathPrefix("/").Handler(next)

	var handler http.Handler = r
	handler = handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(handler)
	handler = handlers.CombinedLoggingHandler(os.Stdout, handler)
	handler = handlers.CompressHandlerLevel(handler, gzip.BestCompression)
	handler = handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(handler)
	return handler
}
