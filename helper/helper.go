package helper

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/senzing/g2-sdk-go/g2error"
)

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func isJson(unknownString string) bool {
	unknownStringUnescaped, err := strconv.Unquote(unknownString)
	if err != nil {
		unknownStringUnescaped = unknownString
	}
	var jsonString json.RawMessage
	return json.Unmarshal([]byte(unknownStringUnescaped), &jsonString) == nil
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// If appropriate, convert error to Senzing nested error.
/*
The ConvertGrpcError method transforms an error produced by google.golang.org/grpc/status
into a Senzing nested error.

Input
  - originalError: The error received from the gRPC call.

Output
  - A Senzing nested error.
*/
func ConvertGrpcError(originalError error) error {

	if originalError == nil {
		return originalError
	}

	result := originalError
	errorMessage := originalError.Error()

	// Determine if error is an RPC error.
	// TODO: This is an error prone method of determining if the error is a gRPC error.
	// Need to determine if error is of type *status.Error.

	if strings.HasPrefix(errorMessage, "rpc error:") {

		// TODO: Improve the fragile method of pulling out the Senzing JSON error.

		indexOfDesc := strings.Index(errorMessage, " desc = ")
		senzingErrorMessage := errorMessage[indexOfDesc+8:]

		if isJson(senzingErrorMessage) {

			// TODO: Add information about any gRPC error.
			// Status: https://pkg.go.dev/google.golang.org/grpc/status
			// Codes: https://pkg.go.dev/google.golang.org/grpc/codes

			// Create a new Senzing nested error.

			result = errors.New(senzingErrorMessage)
		}
	}
	return g2error.Convert(result)
}
