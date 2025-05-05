package helper

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"

	tlshelper "github.com/senzing-garage/go-helpers/tls"
	"github.com/senzing-garage/go-helpers/wraperror"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

/*
The GetGrpcTransportCredentials function returns a gRPC credentials.TransportCredentials
based on the value of the SENZING_TOOLS_SERVER_CA_CERTIFICATE_FILE,
SENZING_TOOLS_CLIENT_CERTIFICATE_FILE, and SENZING_TOOLS_CLIENT_KEY_FILE OS environment
variables.  If only SENZING_TOOLS_SERVER_CA_CERTIFICATE_FILE is set, then the transport
credential is configured for "Server-side TLS".  If all three environment variables are set,
then the transport credential is configured for "Mutual TLS".  Otherwise an insecure transport
credential is returned.

Output
  - Transport Credential calculated by OS environment variables.
*/
func GetGrpcTransportCredentials() (credentials.TransportCredentials, error) {
	var certificates []tls.Certificate

	result := insecure.NewCredentials()

	serverCaCertificatePath, isSet := os.LookupEnv("SENZING_TOOLS_SERVER_CA_CERTIFICATE_FILE")
	if isSet {
		// Server-side TLS.
		rootCAs, err := buildRootCAsFromFile(serverCaCertificatePath)
		if err != nil {
			return result, err
		}

		// Mutual TLS.

		clientCertificatePath, isClientCertificatePathSet := os.LookupEnv("SENZING_TOOLS_CLIENT_CERTIFICATE_FILE")
		clientKeyPath, isClientKeyPathSet := os.LookupEnv("SENZING_TOOLS_CLIENT_KEY_FILE")

		if isClientCertificatePathSet && isClientKeyPathSet {
			clientKeyPassPhrase, _ := os.LookupEnv("SENZING_TOOLS_CLIENT_KEY_PASSPHRASE")

			clientCertificate, err := tlshelper.LoadX509KeyPair(
				clientCertificatePath,
				clientKeyPath,
				clientKeyPassPhrase,
			)
			if err != nil {
				return result, wraperror.Errorf(
					err,
					"helper.GetGrpcTransportCredentials.LoadX509KeyPair error: %w",
					err,
				)
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

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func buildRootCAsFromFile(serverCaCertificatePath string) (*x509.CertPool, error) {
	var (
		err    error
		result *x509.CertPool
	)

	safeFilename := filepath.Clean(serverCaCertificatePath)
	pemServerCA, err := os.ReadFile(safeFilename)
	if err != nil {
		return result, wraperror.Errorf(err, "helper.buildRootCAsFromFile.os.Read error: %w", err)
	}

	result = x509.NewCertPool()
	if !result.AppendCertsFromPEM(pemServerCA) {
		return result, wraperror.Errorf(errForPackage, "failed to add server CA's certificate")
	}

	return result, wraperror.Errorf(err, "helper.buildRootCAsFromFile error: %w", err)
}
