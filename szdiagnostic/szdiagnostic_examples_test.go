//go:build linux

package szdiagnostic_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-grpc/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szdiagnosticpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcAddress = "localhost:8261"
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzdiagnostic_CheckDatastorePerformance() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}
	secondsToRun := 1
	result, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Truncate(result, 2))
	// Output: {"insertTime":1000,...
}

func ExampleSzdiagnostic_GetDatastoreInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}
	result, err := szDiagnostic.GetDatastoreInfo(ctx)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"dataStores":[{"id":"CORE","type":"sqlite3","location":"nowhere"}]}
}

func ExampleSzdiagnostic_GetFeature() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}
	featureID := int64(1)
	result, err := szDiagnostic.GetFeature(ctx, featureID)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Truncate(result, 5))
	// Output: {"ELEMENTS":[{"FELEM_CODE":"CATEGORY","FELEM_VALUE":"PERSON"...
}

func ExampleSzdiagnostic_PurgeRepository() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	if err != nil {
		handleError(err)
	}
	err = szDiagnostic.PurgeRepository(ctx)
	if err != nil {
		handleError(err)
	}
	// Output:
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzdiagnostic_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	err := szDiagnostic.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzdiagnostic_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzdiagnostic_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szdiagnostic/szdiagnostic_examples_test.go
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	result := szDiagnostic.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func getSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	var err error
	var result senzing.SzAbstractFactory
	_ = ctx
	grpcConnection, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	result = &szabstractfactory.Szabstractfactory{
		GrpcConnection: grpcConnection,
	}
	return result
}

func getSzDiagnostic(ctx context.Context) *szdiagnostic.Szdiagnostic {
	_ = ctx
	grpcConnection, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	result := &szdiagnostic.Szdiagnostic{
		GrpcClient: szdiagnosticpb.NewSzDiagnosticClient(grpcConnection),
	}
	return result
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}
