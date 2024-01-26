//go:build linux

package g2product

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2product_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2product.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleG2product_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2product.SetObserverOrigin(ctx, origin)
	result := g2product.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleG2product_Init() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	moduleName := "Test module name"
	iniParams := "{}"
	verboseLogging := int64(0)
	err := g2product.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		// This should produce a "senzing-60164002" error.
	}
	// Output:
}

func ExampleG2product_License() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	result, err := g2product.License(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"customer":"Senzing Public Test License","contract":"Senzing Public Test - 50K records test","issueDate":"2023-11-02","licenseType":"EVAL (Solely for non-productive use)","licenseLevel":"STANDARD","billing":"YEARLY","expireDate":"2024-11-02","recordLimit":50000}
}

func ExampleG2product_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	err := g2product.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2product_Version() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	result, err := g2product.Version(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 43))
	// Output: {"PRODUCT_NAME":"Senzing API","VERSION":...
}

func ExampleG2product_Destroy() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2product/g2product_examples_test.go
	ctx := context.TODO()
	g2product := getG2Product(ctx)
	err := g2product.Destroy(ctx)
	if err != nil {
		// This should produce a "senzing-60164001" error.
	}
	// Output:
}
