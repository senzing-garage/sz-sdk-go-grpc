package szdiagnostic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
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
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	badFeatureID      = int64(-1)
	badLogLevelName   = "BadLogLevelName"
	badSecondsToRun   = -1
	defaultTruncation = 76
	instanceName      = "SzDiagnostic Test"
	observerOrigin    = "SzDiagnostic observer"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

var (
	grpcAddress       = "localhost:8261"
	grpcConnection    *grpc.ClientConn
	logger            logging.Logging
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szConfigManagerSingleton senzing.SzConfigManager
	szConfigSingleton        senzing.SzConfig
	szDiagnosticSingleton    *Szdiagnostic
	szEngineSingleton        senzing.SzEngine
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzdiagnostic_CheckDatastorePerformance(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_CheckDatastorePerformance_badSecondsToRun(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, badSecondsToRun)
	require.NoError(test, err) // TODO: TestSzdiagnostic_CheckDatastorePerformance_badSecondsToRun should fail.
	printActual(test, actual)
}

// TODO: Implement TestSzdiagnostic_CheckDatastorePerformance_error
// func TestSzdiagnostic_CheckDatastorePerformance_error(test *testing.T) {}

func TestSzdiagnostic_GetDatastoreInfo(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	actual, err := szDiagnostic.GetDatastoreInfo(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// TODO: Implement TestSzdiagnostic_GetDatastoreInfo_error
// func TestSzdiagnostic_GetDatastoreInfo_error(test *testing.T) {}

func TestSzdiagnostic_GetFeature(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	featureID := int64(1)
	actual, err := szDiagnostic.GetFeature(ctx, featureID)
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

func TestSzdiagnostic_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

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

func TestSzdiagnostic_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	err := szDiagnostic.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzdiagnostic_AsInterface(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getSzDiagnosticAsInterface(ctx)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	require.NoError(test, err)
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
	configID := senzing.SzInitializeWithDefaultConfiguration
	err := szDiagnostic.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	testError(test, err)
}

func TestSzdiagnostic_Reinitialize(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	szConfigManager := getSzConfigManager(ctx)
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	testError(test, err)
	err = szDiagnostic.Reinitialize(ctx, configID)
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

func createError(errorID int, err error) error {
	return logger.NewError(errorID, err)
}

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
		if errors.Is(err, szerror.ErrSzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if errors.Is(err, szerror.ErrSzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if errors.Is(err, szerror.ErrSzBadInput) {
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
	configID, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		return createError(5913, err)
	}

	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return createError(5914, err)
	}

	szDiagnostic := getSzDiagnostic(ctx)
	err = szDiagnostic.Reinitialize(ctx, configID)

	return err
}

func setupPurgeRepository(ctx context.Context) error {
	szDiagnostic := getSzDiagnostic(ctx)
	err := szDiagnostic.PurgeRepository(ctx)
	return err
}

func setupAddRecords(ctx context.Context) error {
	var err error
	szEngine := getSzEngine(ctx)
	testRecordIDs := []string{"1001", "1002", "1003", "1004", "1005", "1039", "1040"}
	flags := senzing.SzWithoutInfo
	for _, testRecordID := range testRecordIDs {
		testRecord := truthset.CustomerRecords[testRecordID]
		_, err := szEngine.AddRecord(ctx, testRecord.DataSource, testRecord.ID, testRecord.JSON, flags)
		if err != nil {
			return createError(5917, err)
		}
	}
	return err
}

func setup() error {
	ctx := context.TODO()
	var err error

	options := []interface{}{
		&logging.OptionCallerSkip{Value: 4},
	}
	logger, err = logging.NewSenzingLogger(ComponentID, szdiagnosticapi.IDMessages, options...)
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
	var err error
	return err
}
