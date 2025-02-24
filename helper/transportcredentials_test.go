package helper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestHelpers_GetGrpcTransportCredentials_ServerSideTLS(test *testing.T) {
	envVar := "SENZING_TOOLS_CA_CERTIFICATE_PATH"
	_, isSet := os.LookupEnv(envVar)
	if !isSet {
		os.Setenv(envVar, "../testdata/certificates/certificate-authority/certificate.pem")
		defer os.Unsetenv(envVar)
	}
	actual, err := GetGrpcTransportCredentials()
	require.NoError(test, err)
	assert.NotEmpty(test, actual)
}

func TestHelpers_GetGrpcTransportCredentials_Insecure(test *testing.T) {
	envVar := "SENZING_TOOLS_CA_CERTIFICATE_PATH"
	value, isSet := os.LookupEnv(envVar)
	if isSet {
		os.Unsetenv(envVar)
		defer os.Setenv(envVar, value)
	}
	actual, err := GetGrpcTransportCredentials()
	require.NoError(test, err)
	assert.Empty(test, actual)
}
