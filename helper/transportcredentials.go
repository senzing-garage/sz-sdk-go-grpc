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
based on the value of the SENZING_TOOLS_SERVER_CA_CERTIFICATE_PATH,
SENZING_TOOLS_CLIENT_CERTIFICATE_PATH, and SENZING_TOOLS_CLIENT_KEY_PATH OS environment
variables.  If only SENZING_TOOLS_SERVER_CA_CERTIFICATE_PATH is set, then the transport
credential is configured for "Server-side TLS".  If all three environment variables are set,
then the transport credential is configured for "Mutual TLS".  Otherwise an insecure transport
credential is returned.

Output
  - Transport Credential calculated by OS environment variables.
*/
func GetGrpcTransportCredentials() (credentials.TransportCredentials, error) {
	var certificates []tls.Certificate
	result := insecure.NewCredentials()
	serverCaCertificatePath, isSet := os.LookupEnv("SENZING_TOOLS_SERVER_CA_CERTIFICATE_PATH")
	if isSet {

		// Server-side TLS.

		pemServerCA, err := os.ReadFile(serverCaCertificatePath)
		if err != nil {
			return result, err
		}
		rootCAs := x509.NewCertPool()
		if !rootCAs.AppendCertsFromPEM(pemServerCA) {
			return result, fmt.Errorf("failed to add server CA's certificate")
		}

		// Mutual TLS.

		clientCertificatePath, isClientCertificatePathSet := os.LookupEnv("SENZING_TOOLS_CLIENT_CERTIFICATE_PATH")
		clientKeyPath, isClientKeyPathSet := os.LookupEnv("SENZING_TOOLS_CLIENT_KEY_PATH")
		if isClientCertificatePathSet && isClientKeyPathSet {
			clientCertificate, err := tls.LoadX509KeyPair(clientCertificatePath, clientKeyPath)
			if err != nil {
				return result, err
			}
			certificates = []tls.Certificate{clientCertificate}
		}

		// Create TLS configuration.

		config := &tls.Config{
			Certificates: certificates,
			MaxVersion:   tls.VersionTLS13,
			MinVersion:   tls.VersionTLS12, // See https://pkg.go.dev/crypto/tls#pkg-constants
			RootCAs:      rootCAs,
		}
		result = credentials.NewTLS(config)
	}
	return result, nil
}
