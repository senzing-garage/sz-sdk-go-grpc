package helper

import (
	"errors"
	"fmt"

	"github.com/senzing-garage/sz-sdk-go/szerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleConvertGrpcError() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/helper/helper_test.go
	senzingErrorMessage := `{"reason": "SENZ0033|Test message"}`        // Example message from Senzing Sz engine.
	grpcStatusError := status.Error(codes.Unknown, senzingErrorMessage) // Create a gRPC *status.Error

	err := ConvertGrpcError(grpcStatusError)
	if err != nil {
		if errors.Is(err, szerror.ErrSzNotFound) {
			fmt.Println("Is an ErrSzNotFound")
		}

		if errors.Is(err, szerror.ErrSzBadInput) {
			fmt.Println("Is an ErrSzBadInput")
		}

		if errors.Is(err, szerror.ErrSzRetryable) {
			fmt.Println("Is an ErrSzRetryable.")
		}
	}
	// Output:
	// Is an ErrSzNotFound
	// Is an ErrSzBadInput
}
