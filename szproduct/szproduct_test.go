package szproduct_test

import (
	"context"
	"fmt"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/env"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
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
	jsonIndentation   = "    "
	observerOrigin    = "SzProduct observer"
	originMessage     = "Machine: nn; Task: UnitTest"
	printErrors       = false
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
	logLevel          = env.GetEnv("SENZING_LOG_LEVEL", "INFO")
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
	printError(test, err)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzproduct_GetVersion(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	actual, err := szProduct.GetVersion(ctx)
	printError(test, err)
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
	szProduct.SetObserverOrigin(ctx, originMessage)
}

func TestSzproduct_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	szProduct.SetObserverOrigin(ctx, originMessage)
	actual := szProduct.GetObserverOrigin(ctx)
	assert.Equal(test, originMessage, actual)
}

func TestSzproduct_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	err := szProduct.UnregisterObserver(ctx, observerSingleton)
	printError(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzproduct_AsInterface(test *testing.T) {
	ctx := test.Context()
	szProduct := getSzProductAsInterface(ctx)
	actual, err := szProduct.GetLicense(ctx)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzproduct_Initialize(test *testing.T) {
	ctx := test.Context()
	szProduct := &szproduct.Szproduct{}
	settings := getSettings()
	err := szProduct.Initialize(ctx, instanceName, settings, verboseLogging)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzproduct_Initialize_error(test *testing.T) {
	// IMPROVE: Implement TestSzengine_Initialize_error
	_ = test
}

func TestSzproduct_Destroy(test *testing.T) {
	ctx := test.Context()
	szProduct := getTestObject(test)
	err := szProduct.Destroy(ctx)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzproduct_Destroy_error(test *testing.T) {
	// IMPROVE: Implement TestSzengine_Destroy_error
	_ = test
}

func TestSzproduct_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szProductSingleton = nil
	szProduct := getTestObject(test)
	err := szProduct.Destroy(ctx)
	printError(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getGrpcConnection() *grpc.ClientConn {
	if grpcConnection == nil {
		transportCredentials, err := helper.GetGrpcTransportCredentials()
		panicOnError(err)

		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(transportCredentials),
		}

		grpcConnection, err = grpc.NewClient(grpcAddress, dialOptions...)
		panicOnError(err)
	}

	return grpcConnection
}

func getSettings() string {
	return "{}"
}

func getSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	_ = ctx

	return &szabstractfactory.Szabstractfactory{
		GrpcConnection: getGrpcConnection(),
	}
}

func getSzProduct(ctx context.Context) *szproduct.Szproduct {
	var err error

	if szProductSingleton == nil {
		settings := getSettings()

		grpcConnection := getGrpcConnection()
		szProductSingleton = &szproduct.Szproduct{
			GrpcClient: szpb.NewSzProductClient(grpcConnection),
		}
		err = szProductSingleton.SetLogLevel(ctx, logLevel)

		panicOnError(err)

		if logLevel == "TRACE" {
			szProductSingleton.SetObserverOrigin(ctx, observerOrigin)

			err = szProductSingleton.RegisterObserver(ctx, observerSingleton)
			panicOnError(err)

			err = szProductSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			panicOnError(err)
		}

		err = szProductSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		panicOnError(err)
	}

	return szProductSingleton
}

func getSzProductAsInterface(ctx context.Context) senzing.SzProduct {
	return getSzProduct(ctx)
}

func getTestObject(t *testing.T) *szproduct.Szproduct {
	t.Helper()

	return getSzProduct(t.Context())
}

func handleError(err error) {
	if err != nil {
		outputln("Error:", err)
	}
}

func outputln(message ...any) {
	fmt.Println(message...) //nolint
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func printActual(t *testing.T, actual interface{}) {
	t.Helper()
	printResult(t, "Actual", actual)
}

func printError(t *testing.T, err error) {
	t.Helper()

	if printErrors {
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}
}

func printResult(t *testing.T, title string, result interface{}) {
	t.Helper()

	if printResults {
		t.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}
