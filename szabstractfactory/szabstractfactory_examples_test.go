//go:build linux

package szabstractfactory_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"google.golang.org/grpc"
)

var (
	grpcAddress    = "0.0.0.0:8261"
	grpcConnection *grpc.ClientConn
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzabstractfactory_CreateConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	_ = szConfig // szConfig can now be used.
	// Output:
}

func ExampleSzabstractfactory_CreateConfigManager() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	_ = szConfigManager // szConfigManager can now be used.
	// Output:
}

func ExampleSzabstractfactory_CreateDiagnostic() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}
	_ = szDiagnostic // szDiagnostic can now be used.
	// Output:
}

func ExampleSzabstractfactory_CreateEngine() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	_ = szEngine // szEngine can now be used.
	// Output:
}

func ExampleSzabstractfactory_CreateProduct() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	if err != nil {
		handleError(err)
	}
	_ = szProduct // szProduct can now be used.
	// Output:
}

func ExampleSzabstractfactory_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	err := szAbstractFactory.Destroy(ctx)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzabstractfactory_Reinitialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szabstractfactory/szabstractfactory_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
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

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func getGrpcConnection() *grpc.ClientConn {
	if grpcConnection == nil {
		transportCredentials, err := helper.GetGrpcTransportCredentials()
		if err != nil {
			panic(err)
		}
		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(transportCredentials),
		}
		grpcConnection, err = grpc.NewClient(grpcAddress, dialOptions...)
		if err != nil {
			panic(err)
		}
	}
	return grpcConnection
}

func getSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	_ = ctx
	return &szabstractfactory.Szabstractfactory{
		GrpcConnection: getGrpcConnection(),
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}
