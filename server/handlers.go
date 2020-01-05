package server

import (
	"compress/gzip"
	"net/http"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const apiBasePath = "/api"

// gatewayHandler create a http.Handler for server middleware
func gatewayHandler(gwmux *runtime.ServeMux, openapiData []byte, next http.Handler) http.Handler {
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	r := mux.NewRouter()
	r.PathPrefix(apiBasePath).Handler(openapiMiddleware(gwmux, openapiData))
	r.PathPrefix("/").Handler(next)

	var handler http.Handler = r
	handler = handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(handler)
	handler = handlers.CombinedLoggingHandler(os.Stdout, handler)
	handler = handlers.CompressHandlerLevel(handler, gzip.BestCompression)
	handler = handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(handler)
	return handler
}

// openapiMiddleware creates a http.Handler that serves an OpenAPI spec and documentation for that spec.
//
// The handler serves two endpoints:
// - `GET /swagger.json`: the OpenAPI spec as JSON
// - `GET /docs`: documentation for the OpenAPI spec
func openapiMiddleware(gwmux *runtime.ServeMux, data []byte) http.Handler {
	return middleware.Spec(apiBasePath, data, middleware.Redoc(middleware.RedocOpts{}, gwmux))
}
