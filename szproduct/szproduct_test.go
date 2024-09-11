package szproduct

import (
	"context"
	"fmt"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTruncation = 76
	instanceName      = "SzProduct Test"
	observerOrigin    = "SzProduct observer"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

// Bad parameters

const (
	badLogLevelName = "BadLogLevelName"
)

var (
	grpcAddress       = "localhost:8261"
	grpcConnection    *grpc.ClientConn
	logLevel          = helper.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szProductSingleton *Szproduct
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzproduct_GetLicense(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	actual, err := szProduct.GetLicense(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzproduct_GetVersion(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	actual, err := szProduct.GetVersion(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzproduct_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzproduct_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
}

func TestSzproduct_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	actual := szProduct.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzproduct_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	err := szProduct.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzproduct_AsInterface(test *testing.T) {
	ctx := context.TODO()
	szProduct := getSzProductAsInterface(ctx)
	actual, err := szProduct.GetLicense(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzproduct_Initialize(test *testing.T) {
	ctx := context.TODO()
	szProduct, err := getSzProduct(ctx)
	require.NoError(test, err)
	settings, err := getSettings()
	require.NoError(test, err)
	err = szProduct.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

// TODO: Implement TestSzengine_Initialize_error
// func TestSzproduct_Initialize_error(test *testing.T) {}

func TestSzproduct_Destroy(test *testing.T) {
	ctx := context.TODO()
	szProduct := getTestObject(ctx, test)
	err := szProduct.Destroy(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzengine_Destroy_error
// func TestSzproduct_Destroy_error(test *testing.T) {}

func TestSzproduct_Destroy_withObserver(test *testing.T) {
	ctx := context.TODO()
	szProductSingleton = nil
	szProduct := getTestObject(ctx, test)
	err := szProduct.Destroy(ctx)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getGrpcConnection() *grpc.ClientConn {
	var err error
	if grpcConnection == nil {
		grpcConnection, err = grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Did not connect: %v\n", err)
		}
		//		defer grpcConnection.Close()
	}
	return grpcConnection
}

func getSettings() (string, error) {
	return "{}", nil
}

func getSzProduct(ctx context.Context) (*Szproduct, error) {
	var err error
	if szProductSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			return szProductSingleton, fmt.Errorf("getSettings() Error: %w", err)
		}
		grpcConnection := getGrpcConnection()
		szProductSingleton = &Szproduct{
			GrpcClient: szpb.NewSzProductClient(grpcConnection),
		}
		err = szProductSingleton.SetLogLevel(ctx, logLevel)
		if err != nil {
			return szProductSingleton, fmt.Errorf("SetLogLevel() Error: %w", err)
		}
		if logLevel == "TRACE" {
			szProductSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szProductSingleton.RegisterObserver(ctx, observerSingleton)
			if err != nil {
				return szProductSingleton, fmt.Errorf("RegisterObserver() Error: %w", err)
			}
			err = szProductSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			if err != nil {
				return szProductSingleton, fmt.Errorf("SetLogLevel() - 2 Error: %w", err)
			}
		}
		err = szProductSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			return szProductSingleton, fmt.Errorf("Initialize() Error: %w", err)
		}
	}
	return szProductSingleton, err
}

func getSzProductAsInterface(ctx context.Context) senzing.SzProduct {
	result, err := getSzProduct(ctx)
	if err != nil {
		panic(err)
	}
	return result
}

func getTestObject(ctx context.Context, test *testing.T) *Szproduct {
	result, err := getSzProduct(ctx)
	require.NoError(test, err)
	return result
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}
