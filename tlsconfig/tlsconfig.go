package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/hellupline/samples-golang-grpc-server/static"
)

func LoadKeyPair(fs http.FileSystem) (*tls.Config, error) {
	rootCert, err := static.ReadAll(fs, "/rootca.cert")
	if err != nil {
		return nil, err
	}
	certPEMBlock, err := static.ReadAll(fs, "/service.pem")
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := static.ReadAll(fs, "/service.key")
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, fmt.Errorf("faield to create tls key pair: %w", err)
	}
	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(rootCert)
	return &tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: rootCAs}, nil
}
