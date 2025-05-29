package helper

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/senzing-garage/go-helpers/jsonutil"
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
		return nil
	}

	senzingErrorCode, err := strconv.Atoi(reason[4:8])
	if err != nil {
		return nil
	}

	return szerror.New(senzingErrorCode, errorMessage)

}

// ----------------------------------------------------------------------------
// Old stuff
// ----------------------------------------------------------------------------

// func convertGrpcError2(originalError error) error {
// 	var result error

// 	// Determine if error is an RPC error.

// 	if reflect.TypeOf(originalError).String() == "*status.Error" {
// 		errorMessage := originalError.Error()
// 		if strings.HasPrefix(errorMessage, "rpc error:") {
// 			result = extractErrorFromRPCError(originalError, errorMessage)
// 		}
// 	}

// 	return result
// }

// func extractErrorFromRPCError(originalError error, errorMessage string) error {
// 	var result error

// 	// IMPROVE: Improve the fragile method of pulling out the Senzing JSON error.

// 	indexOfDesc := strings.Index(errorMessage, " desc = ")
// 	senzingErrorMessage := errorMessage[indexOfDesc+8:] // Implicitly safe from "0+8" because of "rpc error:" prefix.
// 	indexOfBrace := strings.Index(senzingErrorMessage, "{")

// 	if indexOfBrace >= 0 {
// 		senzingErrorMessage = senzingErrorMessage[indexOfBrace:]
// 		if jsonutil.IsJSON(senzingErrorMessage) {
// 			result = recurseThroughErrors(originalError, senzingErrorMessage)
// 		}
// 	}

// 	return result
// }

// func recurseThroughErrors(originalError error, errorMessage string) error {
// 	var target map[string]any

// 	fmt.Printf(">>>>>> in recurseThroughErrors:  %s\n", errorMessage)

// 	err := json.Unmarshal([]byte(errorMessage), &target)
// 	if err != nil {
// 		fmt.Printf(">>>>>> Cannot Unmarshal error: %v\n", err)

// 		// log.Fatalf("Unable to marshal JSON due to %s", err)
// 	}

// 	errorValue, isOK := target["error"]
// 	if isOK {
// 		fmt.Printf(">>>>>> found errorValue. Type: %T Value: %+v\n", errorValue, errorValue)

// 		errorValueString, isAlsoOK := errorValue.(string)
// 		if isAlsoOK {
// 			return recurseThroughErrors(originalError, errorValueString)
// 		} else {
// 			fmt.Printf(">>>>>> returning originalError: %s\n", originalError)

// 			return originalError
// 		}
// 	} else {
// 		fmt.Printf(">>>>>> not OK: %s\n", errorMessage)
// 	}

// 	return extractErrorFromJSON(originalError, errorMessage)

// }

// func recurseDown(originalError error, errorMap map[string]any) error {
// 	errorValue, isOK := errorMap["reason"]
// 	if isOK {
// 		errorValueString, isOK := errorValue.(string)
// 	}

// 	errorValue, isOK = errorMap["error"]
// 	if isOK {
// 		errorValueMap, isAlsoOK := errorValue.(map[string]any)
// 		if isAlsoOK {
// 			return recurseDown(originalError, errorValueMap)
// 		} else {
// 			return originalError
// 		}
// 	}

// 	return originalError

// }

// func extractErrorFromJSON(originalError error, errorMessage string) error {
// 	var result error

// 	// IMPROVE: Add information about any gRPC error.
// 	// Status: https://pkg.go.dev/google.golang.org/grpc/status
// 	// Codes: https://pkg.go.dev/google.golang.org/grpc/codes
// 	// Create a new Senzing nested error.

// 	fmt.Printf(">>>>>>>> errorMessage: %s\n", errorMessage)

// 	parsedMessage, err := parser.Parse(errorMessage)
// 	if err != nil {
// 		return wraperror.Errorf(err, "parse(%s) Original Error: %s", errorMessage, originalError.Error())
// 	}

// 	reason := parsedMessage.Reason
// 	if len(reason) < maxReasons {
// 		return wraperror.Errorf(errForPackage, "len(%s) Original Error: %s", reason, originalError.Error())
// 	}

// 	senzingErrorCode, err := strconv.Atoi(reason[4:8])
// 	if err != nil {
// 		return wraperror.Errorf(err, "strconv.Atoi(%s) Original Error: %s", reason, originalError.Error())
// 	}

// 	result = szerror.New(senzingErrorCode, errorMessage)

// 	return wraperror.Errorf(result, wraperror.NoMessage)
// }
