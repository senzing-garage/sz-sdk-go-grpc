package szdiagnostic

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-grpc/szengine"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szdiagnosticapi "github.com/senzing-garage/sz-sdk-go/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
	szenginepb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	grpcAddress              = "localhost:8261"
	grpcConnection           *grpc.ClientConn
	logger                   logging.LoggingInterface
	szConfigManagerSingleton senzing.SzConfigManager
	szConfigSingleton        senzing.SzConfig
	szDiagnosticSingleton    *Szdiagnostic
	szEngineSingleton        senzing.SzEngine
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzdiagnostic_CheckDatastorePerformance(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	testError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_GetDatastoreInfo(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	actual, err := szDiagnostic.GetDatastoreInfo(ctx)
	testError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_GetFeature(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	featureId := int64(1)
	actual, err := szDiagnostic.GetFeature(ctx, featureId)
	testError(test, err)
	printActual(test, actual)
}

// func TestSzdiagnostic_PurgeRepository(test *testing.T) {
// 	ctx := context.TODO()
// 	szDiagnostic := getTestObject(ctx, test)
// 	err := szDiagnostic.PurgeRepository(ctx)
// 	testError(test, err)
// }

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzdiagnostic_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
}

func TestSzdiagnostic_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	actual := szDiagnostic.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzdiagnostic_AsInterface(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getSzDiagnosticAsInterface(ctx)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	testError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_Initialize(test *testing.T) {
	ctx := context.TODO()
	grpcConnection := getGrpcConnection()
	szDiagnostic := &Szdiagnostic{
		GrpcClient: szpb.NewSzDiagnosticClient(grpcConnection),
	}
	instanceName := "Test name"
	settings := "{}"
	verboseLogging := senzing.SzNoLogging
	configId := senzing.SzInitializeWithDefaultConfiguration
	err := szDiagnostic.Initialize(ctx, instanceName, settings, configId, verboseLogging)
	testError(test, err)
}

func TestSzdiagnostic_Reinitialize(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	szConfigManager := getSzConfigManager(ctx)
	configId, err := szConfigManager.GetDefaultConfigId(ctx)
	testError(test, err)
	err = szDiagnostic.Reinitialize(ctx, configId)
	testErrorNoFail(test, err)
}

func TestSzdiagnostic_Destroy(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	err := szDiagnostic.Destroy(ctx)
	testError(test, err)
	szDiagnosticSingleton = nil
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(logger.NewError(errorId, err), err)
}

func getGrpcConnection() *grpc.ClientConn {
	var err error = nil
	if grpcConnection == nil {
		grpcConnection, err = grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Did not connect: %v\n", err)
		}
		//		defer grpcConnection.Close()
	}
	return grpcConnection
}

func getTestObject(ctx context.Context, test *testing.T) *Szdiagnostic {
	_ = test
	return getSzDiagnostic(ctx)
}

func getSzConfig(ctx context.Context) senzing.SzConfig {
	_ = ctx
	if szConfigSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigSingleton = &szconfig.Szconfig{
			GrpcClient: szconfigpb.NewSzConfigClient(grpcConnection),
		}
	}
	return szConfigSingleton
}

func getSzConfigManager(ctx context.Context) senzing.SzConfigManager {
	_ = ctx
	if szConfigManagerSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigManagerSingleton = &szconfigmanager.Szconfigmanager{
			GrpcClient: szconfigmanagerpb.NewSzConfigManagerClient(grpcConnection),
		}
	}
	return szConfigManagerSingleton
}

func getSzDiagnostic(ctx context.Context) *Szdiagnostic {
	_ = ctx
	if szDiagnosticSingleton == nil {
		grpcConnection := getGrpcConnection()
		szDiagnosticSingleton = &Szdiagnostic{
			GrpcClient: szpb.NewSzDiagnosticClient(grpcConnection),
		}
	}
	return szDiagnosticSingleton
}

func getSzDiagnosticAsInterface(ctx context.Context) senzing.SzDiagnostic {
	return getSzDiagnostic(ctx)
}

func getSzEngine(ctx context.Context) senzing.SzEngine {
	_ = ctx
	if szEngineSingleton == nil {
		grpcConnection := getGrpcConnection()
		szEngineSingleton = &szengine.Szengine{
			GrpcClient: szenginepb.NewSzEngineClient(grpcConnection),
		}
	}
	return szEngineSingleton
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func testError(test *testing.T, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorNoFail(test *testing.T, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
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
		if szerror.Is(err, szerror.SzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if szerror.Is(err, szerror.SzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if szerror.Is(err, szerror.SzBadInput) {
			fmt.Printf("\nBad user input error detected. \n\n")
		}
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

func setupSenzingConfig(ctx context.Context) error {
	now := time.Now()

	// Create a fresh Senzing configuration.

	szCconfig := getSzConfig(ctx)
	configHandle, err := szCconfig.CreateConfig(ctx)
	if err != nil {
		return createError(5907, err)
	}

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range datasourceNames {
		_, err := szCconfig.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return createError(5908, err)
		}
	}

	configDefinition, err := szCconfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	szConfigManager := getSzConfigManager(ctx)
	configComment := fmt.Sprintf("Created by szdiagnostic_test at %s", now.UTC())
	configId, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		return createError(5913, err)
	}

	err = szConfigManager.SetDefaultConfigId(ctx, configId)
	if err != nil {
		return createError(5914, err)
	}

	szDiagnostic := getSzDiagnostic(ctx)
	err = szDiagnostic.Reinitialize(ctx, configId)

	return err
}

func setupPurgeRepository(ctx context.Context) error {
	szDiagnostic := getSzDiagnostic(ctx)
	err := szDiagnostic.PurgeRepository(ctx)
	return err
}

func setupAddRecords(ctx context.Context) error {
	var err error = nil
	szEngine := getSzEngine(ctx)
	testRecordIds := []string{"1001", "1002", "1003", "1004", "1005", "1039", "1040"}
	flags := senzing.SzWithoutInfo
	for _, testRecordId := range testRecordIds {
		testRecord := truthset.CustomerRecords[testRecordId]
		_, err := szEngine.AddRecord(ctx, testRecord.DataSource, testRecord.Id, testRecord.Json, flags)
		if err != nil {
			return createError(5917, err)
		}
	}
	return err
}

func setup() error {
	ctx := context.TODO()
	var err error = nil

	options := []interface{}{
		&logging.OptionCallerSkip{Value: 4},
	}
	logger, err = logging.NewSenzingSdkLogger(ComponentId, szdiagnosticapi.IdMessages, options...)
	if err != nil {
		return createError(5901, err)
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfig(ctx)
	if err != nil {
		return createError(5920, err)
	}

	// Purge repository.

	err = setupPurgeRepository(ctx)
	if err != nil {
		return createError(5921, err)
	}

	// Add records.

	err = setupAddRecords(ctx)
	if err != nil {
		return createError(5922, err)
	}

	return err
}

func teardown() error {
	var err error = nil
	return err
}
