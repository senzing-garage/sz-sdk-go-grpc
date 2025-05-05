package helper

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-messaging/parser"
	"github.com/senzing-garage/sz-sdk-go/szerror"
)

const maxReasons = 10

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

	for currentError := originalError; currentError != nil; currentError = errors.Unwrap(currentError) {
		err := convertGrpcError(currentError)
		if err != nil {
			return err
		}
	}

	return result
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func convertGrpcError(originalError error) error {
	var result error

	// Determine if error is an RPC error.

	if reflect.TypeOf(originalError).String() == "*status.Error" {
		errorMessage := originalError.Error()
		if strings.HasPrefix(errorMessage, "rpc error:") {
			result = extractErrorFromRPCError(originalError, errorMessage)
		}
	}

	return result
}

func extractErrorFromRPCError(originalError error, errorMessage string) error {
	var result error

	// IMPROVE: Improve the fragile method of pulling out the Senzing JSON error.

	indexOfDesc := strings.Index(errorMessage, " desc = ")
	senzingErrorMessage := errorMessage[indexOfDesc+8:] // Implicitly safe from "0+8" because of "rpc error:" prefix.
	indexOfBrace := strings.Index(senzingErrorMessage, "{")

	if indexOfBrace >= 0 {
		senzingErrorMessage = senzingErrorMessage[indexOfBrace:]
		if jsonutil.IsJSON(senzingErrorMessage) {
			result = extractErrorFromJSON(originalError, senzingErrorMessage)
		}
	}

	return result
}

func extractErrorFromJSON(originalError error, errorMessage string) error {
	var result error

	// IMPROVE: Add information about any gRPC error.
	// Status: https://pkg.go.dev/google.golang.org/grpc/status
	// Codes: https://pkg.go.dev/google.golang.org/grpc/codes
	// Create a new Senzing nested error.

	parsedMessage, err := parser.Parse(errorMessage)
	if err != nil {
		return wraperror.Errorf(
			err,
			"parse(%s) error: %w; Original Error: %w",
			errorMessage,
			err,
			originalError,
		)
	}

	reason := parsedMessage.Reason
	if len(reason) < maxReasons {
		return wraperror.Errorf(errForPackage, "len(%s) error: %w; Original Error: %w", reason, err, originalError)
	}

	senzingErrorCode, err := strconv.Atoi(reason[4:8])
	if err != nil {
		return wraperror.Errorf(
			err,
			"strconv.Atoi(%s) error %w; Original Error: %w",
			reason,
			err,
			originalError,
		)
	}

	result = szerror.New(senzingErrorCode, errorMessage)

	return wraperror.Errorf(result, "%w", result)
}
