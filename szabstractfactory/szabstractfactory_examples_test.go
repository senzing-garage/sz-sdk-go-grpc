//go:build linux

package szabstractfactory

import (
	"context"
	"fmt"
)

// ----------------------------------------------------------------------------
// Interface functions - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzAbstractFactory_CreateSzConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfig, err := szAbstractFactory.CreateSzConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = szConfig.Destroy(ctx) }()
	// Output:
}

func ExampleSzAbstractFactory_CreateSzConfigManager() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfigManager, err := szAbstractFactory.CreateSzConfigManager(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = szConfigManager.Destroy(ctx) }()
	// Output:
}

func ExampleSzAbstractFactory_CreateSzDiagnostic() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szDiagnostic, err := szAbstractFactory.CreateSzDiagnostic(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = szDiagnostic.Destroy(ctx) }()
	// Output:
}

func ExampleSzAbstractFactory_CreateSzEngine() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateSzEngine(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = szEngine.Destroy(ctx) }()
	// Output:
}

func ExampleSzAbstractFactory_CreateSzProduct() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szProduct, err := szAbstractFactory.CreateSzProduct(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = szProduct.Destroy(ctx) }()
	// Output:
}
