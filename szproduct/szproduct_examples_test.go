//go:build linux

package szproduct

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go/senzing"
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzproduct_GetLicense() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProductExample(ctx)
	result, err := szProduct.GetLicense(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"customer":"","contract":"","issueDate":"2024-09-26","licenseType":"EVAL (Solely for non-productive use)","licenseLevel":"","billing":"","expireDate":"2025-09-27","recordLimit":500}
}

func ExampleSzproduct_GetVersion() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProductExample(ctx)
	result, err := szProduct.GetVersion(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 43))
	// Output: {"PRODUCT_NAME":"Senzing API","VERSION":...
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzproduct_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct, err := getSzProduct(ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = szProduct.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzproduct_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct, err := getSzProduct(ctx)
	if err != nil {
		fmt.Println(err)
	}
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzproduct_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct, err := getSzProduct(ctx)
	if err != nil {
		fmt.Println(err)
	}
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	result := szProduct.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func ExampleSzproduct_Initialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct, err := getSzProduct(ctx)
	if err != nil {
		fmt.Println(err)
	}
	instanceName := "Test name"
	settings, err := getSettings()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := senzing.SzNoLogging
	err = szProduct.Initialize(ctx, instanceName, settings, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func getSzProductExample(ctx context.Context) senzing.SzProduct {
	result, err := getSzProduct(ctx)
	if err != nil {
		panic(err)
	}
	return result
}
