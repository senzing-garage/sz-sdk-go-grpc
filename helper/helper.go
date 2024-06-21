package helper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/senzing-garage/go-messaging/parser"
	"github.com/senzing-garage/sz-sdk-go/szerror"
)

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func isJSON(unknownString string) bool {
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
		return nil
	}

	result := originalError

	// Determine if error is an RPC error.

	if reflect.TypeOf(originalError).String() == "*status.Error" {
		errorMessage := originalError.Error()
		if strings.HasPrefix(errorMessage, "rpc error:") {

			// TODO: Improve the fragile method of pulling out the Senzing JSON error.

			indexOfDesc := strings.Index(errorMessage, " desc = ")
			senzingErrorMessage := errorMessage[indexOfDesc+8:] // Implicitly safe from "0+8" because of "rpc error:" prefix.

			if isJSON(senzingErrorMessage) {

				// TODO: Add information about any gRPC error.
				// Status: https://pkg.go.dev/google.golang.org/grpc/status
				// Codes: https://pkg.go.dev/google.golang.org/grpc/codes

				// Create a new Senzing nested error.

				parsedMessage, err := parser.Parse(senzingErrorMessage)
				if err != nil {
					return fmt.Errorf("parse(%s) error: %w; Original Error: %w", senzingErrorMessage, err, originalError)
				}
				reason := parsedMessage.Reason
				if len(reason) < 10 {
					return fmt.Errorf("len(%s) error: %w; Original Error: %w", reason, err, originalError)
				}
				senzingErrorCode, err := strconv.Atoi(reason[4:8])
				if err != nil {
					return fmt.Errorf("strconv.Atoi(%s) error %w; Original Error: %w", reason, err, originalError)
				}
				result = szerror.New(senzingErrorCode, senzingErrorMessage)
			}
		}
	}
	return result
}
