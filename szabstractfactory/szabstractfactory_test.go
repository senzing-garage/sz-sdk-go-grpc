package szabstractfactory

import (
	"context"
	"fmt"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szconfig"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	baseCallerSkip    = 4
	defaultTruncation = 76
	instanceName      = "SzAbstractFactory Test"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

var (
	grpcAddress = "localhost:8261"
	logger      logging.Logging
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzAbstractFactory_CreateSzConfig(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	szConfig, err := szAbstractFactory.CreateSzConfig(ctx)
	require.NoError(test, err)
	defer func() { handleError(szConfig.Destroy(ctx)) }()
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	dataSources, err := szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printActual(test, dataSources)
}

func TestSzAbstractFactory_CreateSzConfigManager(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	szConfigManager, err := szAbstractFactory.CreateSzConfigManager(ctx)
	require.NoError(test, err)
	defer func() { handleError(szConfigManager.Destroy(ctx)) }()
	configList, err := szConfigManager.GetConfigs(ctx)
	require.NoError(test, err)
	printActual(test, configList)
}

func TestSzAbstractFactory_CreateSzDiagnostic(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	szDiagnostic, err := szAbstractFactory.CreateSzDiagnostic(ctx)
	require.NoError(test, err)
	defer func() { handleError(szDiagnostic.Destroy(ctx)) }()
	result, err := szDiagnostic.CheckDatastorePerformance(ctx, 1)
	require.NoError(test, err)
	printActual(test, result)
}

func TestSzAbstractFactory_CreateSzEngine(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	szEngine, err := szAbstractFactory.CreateSzEngine(ctx)
	require.NoError(test, err)
	defer func() { handleError(szEngine.Destroy(ctx)) }()
	stats, err := szEngine.GetStats(ctx)
	require.NoError(test, err)
	printActual(test, stats)
}

func TestSzAbstractFactory_CreateSzProduct(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	szProduct, err := szAbstractFactory.CreateSzProduct(ctx)
	require.NoError(test, err)
	defer func() { handleError(szProduct.Destroy(ctx)) }()
	version, err := szProduct.GetVersion(ctx)
	require.NoError(test, err)
	printActual(test, version)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getSzAbstractFactory(ctx context.Context) (senzing.SzAbstractFactory, error) {
	var err error
	var result senzing.SzAbstractFactory
	_ = ctx
	grpcConnection, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return result, err
	}
	result = &Szabstractfactory{
		GrpcConnection: grpcConnection,
	}
	return result, err
}

func getTestObject(ctx context.Context, test *testing.T) senzing.SzAbstractFactory {
	result, err := getSzAbstractFactory(ctx)
	require.NoError(test, err)
	return result
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
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

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setup() error {
	var err error
	logger = helper.GetLogger(ComponentID, szconfig.IDMessages, baseCallerSkip)
	return err
}

func teardown() error {
	return nil
}
