//go:build linux

package szabstractfactory_test

import (
	"context"
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzabstractfactory_Close() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)

	err := szAbstractFactory.Close(ctx)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzabstractfactory_CreateConfigManager() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}

	_ = szConfigManager // szConfigManager can now be used.
	// Output:
}

func ExampleSzabstractfactory_CreateDiagnostic() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}

	_ = szDiagnostic // szDiagnostic can now be used.
	// Output:
}

func ExampleSzabstractfactory_CreateEngine() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}

	_ = szEngine // szEngine can now be used.
	// Output:
}

func ExampleSzabstractfactory_CreateProduct() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	if err != nil {
		handleError(err)
	}

	_ = szProduct // szProduct can now be used.
	// Output:
}

func ExampleSzabstractfactory_Reinitialize() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}

	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
	}

	err = szAbstractFactory.Reinitialize(ctx, configID)
	if err != nil {
		handleError(err)
	}
	// Output:
}
