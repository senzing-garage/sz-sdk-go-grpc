package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/sz-sdk-go/szerror"
)

const maxReasons = 10

var errPackage = errors.New("helper")

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

func convertGrpcError(grpcError error) error {
	var result error

	// Verify that it is a gRPC error.

	if reflect.TypeOf(grpcError).String() != "*status.Error" {
		return result
	}

	// Make sure there is a "desc" field.

	grpcErrorMessage := grpcError.Error()
	desc := " desc = "

	indexOfDesc := strings.Index(grpcErrorMessage, desc)
	if indexOfDesc < 0 {
		return result
	}

	// Get the JSON string.

	senzingErrorMessage := grpcErrorMessage[indexOfDesc+len(desc):]

	indexOfBrace := strings.Index(senzingErrorMessage, "{")
	if indexOfBrace < 0 {
		return result
	}

	senzingErrorJSON := senzingErrorMessage[indexOfBrace:]
	if !jsonutil.IsJSON(senzingErrorJSON) {
		return result
	}

	// Inspect the JSON for the "reason" field in a "message".

	reason := extractReasonFromJSON(senzingErrorJSON)
	if len(reason) == 0 {
		return result
	}

	return createErrorFromReason(senzingErrorJSON, reason)
}

func extractReasonFromJSON(message string) string {
	var (
		result   string
		errorMap map[string]any
	)

	err := json.Unmarshal([]byte(message), &errorMap)
	if err != nil {
		return result
	}

	return extractReasonFromAny(errorMap)
}

func extractReasonFromAny(errorMap map[string]any) string {
	reasonValue, isOK := errorMap["reason"]
	if isOK {
		reasonValueString, isOK := reasonValue.(string)
		if isOK {
			return reasonValueString
		}
	}

	errorValue, isOK := errorMap["error"]
	if isOK {
		newErrorMap, isOK := errorValue.(map[string]any)
		if isOK {
			return extractReasonFromAny(newErrorMap)
		}
	}

	result, err := json.Marshal(errorMap)
	if err != nil {
		panic(err)
	}

	return string(result)
}

func createErrorFromReason(errorMessage string, reason string) error {
	if len(reason) < maxReasons {
		return wraperror.Errorf(
			errPackage,
			wraperror.Quote(fmt.Sprintf("errorMessage: %s; reason: %s", errorMessage, reason)),
		)
	}

	senzingErrorCode, err := strconv.Atoi(reason[4:8])
	if err != nil {
		return wraperror.Errorf(err, wraperror.Quote(fmt.Sprintf("errorMessage: %s; reason: %s", errorMessage, reason)))
	}

	return szerror.New(senzingErrorCode, errorMessage) //nolint
}
