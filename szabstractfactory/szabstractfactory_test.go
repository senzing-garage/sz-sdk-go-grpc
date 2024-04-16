package szabstractfactory

import (
	"context"
	"fmt"
	"os"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go/sz"
	"github.com/senzing-garage/sz-sdk-go/szconfig"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTruncation = 76
	instanceName      = "SzAbstractFactory Test"
	printResults      = false
	verboseLogging    = sz.SZ_NO_LOGGING
)

var (
	grpcAddress = "localhost:8261"
	logger      logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzAbstractFactory_CreateConfig(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	config, err := szAbstractFactory.CreateConfig(ctx)
	testError(test, ctx, szAbstractFactory, err)
	defer config.Destroy(ctx)
	configHandle, err := config.CreateConfig(ctx)
	testError(test, ctx, szAbstractFactory, err)
	dataSources, err := config.GetDataSources(ctx, configHandle)
	testError(test, ctx, szAbstractFactory, err)
	printActual(test, dataSources)
}

func TestSzAbstractFactory_CreateConfigManager(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	configManager, err := szAbstractFactory.CreateConfigManager(ctx)
	testError(test, ctx, szAbstractFactory, err)
	defer configManager.Destroy(ctx)
	configList, err := configManager.GetConfigList(ctx)
	testError(test, ctx, szAbstractFactory, err)
	printActual(test, configList)
}

func TestSzAbstractFactory_CreateDiagnostic(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	diagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	testError(test, ctx, szAbstractFactory, err)
	defer diagnostic.Destroy(ctx)
	result, err := diagnostic.CheckDatabasePerformance(ctx, 1)
	testError(test, ctx, szAbstractFactory, err)
	printActual(test, result)
}

func TestSzAbstractFactory_CreateEngine(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	engine, err := szAbstractFactory.CreateEngine(ctx)
	testError(test, ctx, szAbstractFactory, err)
	defer engine.Destroy(ctx)
	stats, err := engine.GetStats(ctx)
	testError(test, ctx, szAbstractFactory, err)
	printActual(test, stats)
}

func TestSzAbstractFactory_CreateProduct(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	product, err := szAbstractFactory.CreateProduct(ctx)
	testError(test, ctx, szAbstractFactory, err)
	defer product.Destroy(ctx)
	version, err := product.GetVersion(ctx)
	testError(test, ctx, szAbstractFactory, err)
	printActual(test, version)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(logger.NewError(errorId, err), err)
}

func getSzAbstractFactory(ctx context.Context) sz.SzAbstractFactory {
	_ = ctx
	grpcConnection, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	result := &Szabstractfactory{
		GrpcConnection: grpcConnection,
	}
	return result
}

func getTestObject(ctx context.Context, test *testing.T) sz.SzAbstractFactory {
	_ = test
	return getSzAbstractFactory(ctx)
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func testError(test *testing.T, ctx context.Context, szAbstractFactory sz.SzAbstractFactory, err error) {
	_ = ctx
	_ = szAbstractFactory
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
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
	var err error = nil
	logger, err = logging.NewSenzingSdkLogger(ComponentId, szconfig.IdMessages)
	if err != nil {
		return createError(5901, err)
	}
	return err
}

func teardown() error {
	return nil
}
