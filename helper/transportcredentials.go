package helper

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

/*
The GetGrpcTransportCredentials function returns a gRPC credentials.TransportCredentials
based on the value of the value of the SENZING_TOOLS_CA_CERTIFICATE_PATH OS environment variable.
If the environment variable does not exist, an insecure transport credential is returned.

Output
  - Transport Credential calculated by OS environment variables.
*/
func GetGrpcTransportCredentials() (credentials.TransportCredentials, error) {
	var result credentials.TransportCredentials
	certFile, isSet := os.LookupEnv("SENZING_TOOLS_CA_CERTIFICATE_PATH")
	if isSet {
		pemServerCA, err := os.ReadFile(certFile)
		if err != nil {
			return result, err
		}
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(pemServerCA) {
			return result, fmt.Errorf("failed to add server CA's certificate")
		}
		config := &tls.Config{
			RootCAs:    certPool,
			MinVersion: tls.VersionTLS12, // See https://pkg.go.dev/crypto/tls#pkg-constants
			MaxVersion: tls.VersionTLS13,
		}
		result = credentials.NewTLS(config)
	} else {
		result = insecure.NewCredentials()
	}
	return result, nil
}
