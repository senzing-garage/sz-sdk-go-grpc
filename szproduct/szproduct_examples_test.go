//go:build linux

package szproduct_test

import (
	"context"
	"fmt"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-grpc/szproduct"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szproductpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
	"google.golang.org/grpc"
)

var (
	grpcAddress    = "0.0.0.0:8261"
	grpcConnection *grpc.ClientConn
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzproduct_GetLicense() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)

	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	if err != nil {
		handleError(err)
	}

	result, err := szProduct.GetLicense(ctx)
	if err != nil {
		handleError(err)
	}

	fmt.Println(jsonutil.Truncate(result, 4))
	// Output: {"billing":"YEARLY","contract":"Senzing Public Test License","customer":"Senzing Public Test License",...
}

func ExampleSzproduct_GetVersion() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)

	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	if err != nil {
		handleError(err)
	}

	result, err := szProduct.GetVersion(ctx)
	if err != nil {
		handleError(err)
	}

	fmt.Println(truncate(result, 43))
	// Output: {"PRODUCT_NAME":"Senzing SDK","VERSION":...
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzproduct_SetLogLevel() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)

	err := szProduct.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzproduct_SetObserverOrigin() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzproduct_GetObserverOrigin() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szproduct/szproduct_examples_test.go
	ctx := context.TODO()
	szProduct := getSzProduct(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	result := szProduct.GetObserverOrigin(ctx)

	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
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
