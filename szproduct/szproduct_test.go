package szproduct_test

import (
	"context"
	"fmt"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go-grpc/szproduct"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
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
	grpcAddress       = "0.0.0.0:8261"
	grpcConnection    *grpc.ClientConn
	logLevel          = helper.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szProductSingleton *szproduct.Szproduct
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzproduct_GetLicense(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	actual, err := szProduct.GetLicense(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzproduct_GetVersion(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	actual, err := szProduct.GetVersion(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzproduct_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzproduct_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
}

func TestSzproduct_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	origin := "Machine: nn; Task: UnitTest"
	szProduct.SetObserverOrigin(ctx, origin)
	actual := szProduct.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzproduct_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	err := szProduct.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzproduct_AsInterface(test *testing.T) {
	ctx := test.Context()
	szProduct := getSzProductAsInterface(ctx)
	actual, err := szProduct.GetLicense(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzproduct_Initialize(test *testing.T) {
	ctx := test.Context()
	szProduct := &szproduct.Szproduct{}
	settings := getSettings()
	err := szProduct.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

// TODO: Implement TestSzengine_Initialize_error
// func TestSzproduct_Initialize_error(test *testing.T) {}

func TestSzproduct_Destroy(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	err := szProduct.Destroy(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzengine_Destroy_error
// func TestSzproduct_Destroy_error(test *testing.T) {}

func TestSzproduct_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szProductSingleton = nil
	szProduct := getTestObject(test)
	err := szProduct.Destroy(ctx)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getGrpcConnection() *grpc.ClientConn {
	if grpcConnection == nil {
		transportCredentials, err := helper.GetGrpcTransportCredentials()
		handleErrorWithPanic(err)

		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(transportCredentials),
		}

		grpcConnection, err = grpc.NewClient(grpcAddress, dialOptions...)
		handleErrorWithPanic(err)
	}

	return grpcConnection
}

func getSettings() string {
	return "{}"
}

func getSzProduct(ctx context.Context) *szproduct.Szproduct {
	if szProductSingleton == nil {
		settings, err := getSettings()
		handleErrorWithPanic(err)

		grpcConnection := getGrpcConnection()
		szProductSingleton = &Szproduct{
			GrpcClient: szpb.NewSzProductClient(grpcConnection),
		}
		err = szProductSingleton.SetLogLevel(ctx, logLevel)

		handleErrorWithPanic(err)

		if logLevel == "TRACE" {
			szProductSingleton.SetObserverOrigin(ctx, observerOrigin)

			err = szProductSingleton.RegisterObserver(ctx, observerSingleton)
			handleErrorWithPanic(err)

			err = szProductSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			handleErrorWithPanic(err)
		}

		err = szProductSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		handleErrorWithPanic(err)
	}

	return szProductSingleton
}

func getSzProductAsInterface(ctx context.Context) senzing.SzProduct {
	return getSzProduct(ctx)
}

func getTestObject(t *testing.T) *szproduct.Szproduct {
	t.Helper()
	ctx := t.Context()

	return getSzProduct(ctx)
}

func handleError(err error) {
	if err != nil {
		safePrintln("Error:", err)
	}
}

func handleErrorWithPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func printActual(t *testing.T, actual interface{}) {
	t.Helper()
	printResult(t, "Actual", actual)
}

func printResult(t *testing.T, title string, result interface{}) {
	t.Helper()

	if printResults {
		t.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func safePrintln(message ...any) {
	fmt.Println(message...) //nolint
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}
