package szdiagnostic_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/env"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-core/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-grpc/szengine"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
	szenginepb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	defaultTruncation = 76
	instanceName      = "SzDiagnostic Test"
	observerOrigin    = "SzDiagnostic observer"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

// Bad parameters

const (
	badFeatureID    = int64(-1)
	badLogLevelName = "BadLogLevelName"
	badSecondsToRun = -1
)

// Nil/empty parameters

var (
	nilSecondsToRun int
	nilFeatureID    int64
)

var (
	defaultConfigID   int64
	grpcAddress       = "0.0.0.0:8261"
	grpcConnection    *grpc.ClientConn
	logLevel          = env.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szConfigManagerSingleton *szconfigmanager.Szconfigmanager
	szConfigSingleton        *szconfig.Szconfig
	szDiagnosticSingleton    *Szdiagnostic
	szEngineSingleton        *szengine.Szengine
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzdiagnostic_CheckDatastorePerformance(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_CheckDatastorePerformance_badSecondsToRun(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, badSecondsToRun)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_CheckDatastorePerformance_nilSecondsToRun(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, nilSecondsToRun)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_GetDatastoreInfo(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.GetDatastoreInfo(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_GetFeature(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szDiagnostic := getTestObject(test)
	featureID := int64(1)
	actual, err := szDiagnostic.GetFeature(ctx, featureID)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_GetFeature_badFeatureID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.GetFeature(ctx, badFeatureID)
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

func TestSzdiagnostic_GetFeature_nilFeatureID(test *testing.T) {
	ctx := test.Context()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}

	defer func() { deleteRecords(ctx, records) }()

	addRecords(ctx, records)

	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.GetFeature(ctx, nilFeatureID)
	require.ErrorIs(test, err, szerror.ErrSz)
	printActual(test, actual)
}

// PurgeRepository is tested in szdiagnostic_examples_test.go
// func TestSzdiagnostic_PurgeRepository(test *testing.T) {}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzdiagnostic_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzdiagnostic_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
}

func TestSzdiagnostic_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	actual := szDiagnostic.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzdiagnostic_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	err := szDiagnostic.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzdiagnostic_AsInterface(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getSzDiagnosticAsInterface(ctx)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, secondsToRun)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzdiagnostic_Initialize(test *testing.T) {
	ctx := test.Context()
	grpcConnection := getGrpcConnection()
	szDiagnostic := &Szdiagnostic{
		GrpcClient: szpb.NewSzDiagnosticClient(grpcConnection),
	}
	settings, err := getSettings()
	require.NoError(test, err)

	configID := senzing.SzInitializeWithDefaultConfiguration
	err = szDiagnostic.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	require.NoError(test, err)
}

// TODO: Implement TestSzdiagnostic_Initialize_error
// func TestSzdiagnostic_Initialize_error(test *testing.T) {}

func TestSzdiagnostic_Initialize_withConfigId(test *testing.T) {
	ctx := test.Context()
	grpcConnection := getGrpcConnection()
	szDiagnostic := &Szdiagnostic{
		GrpcClient: szpb.NewSzDiagnosticClient(grpcConnection),
	}
	settings, err := getSettings()
	require.NoError(test, err)

	configID := getDefaultConfigID()
	err = szDiagnostic.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	require.NoError(test, err)
}

// TODO: Implement TestSzdiagnostic_Initialize_withConfigId_badConfigID
// func TestSzdiagnostic_Initialize_withConfigId_badConfigID(test *testing.T) {}

// func TestSzdiagnostic_Reinitialize(test *testing.T) {
// 	ctx := test.Context()
// 	szDiagnostic := getTestObject(ctx, test)
// 	configID := getDefaultConfigID()
// 	err := szDiagnostic.Reinitialize(ctx, configID)
// 	require.NoError(test, err)
// }

// TODO: Implement TestSzdiagnostic_Reinitialize_error
// func TestSzdiagnostic_Reinitialize_error(test *testing.T) {}

func TestSzdiagnostic_Destroy(test *testing.T) {
	ctx := test.Context()
	szDiagnosticSingleton = nil
	szDiagnostic := getTestObject(test)
	err := szDiagnostic.Destroy(ctx)
	require.NoError(test, err)
}

func TestSzdiagnostic_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szDiagnosticSingleton = nil
	szDiagnostic := getTestObject(test)
	err := szDiagnostic.Destroy(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzdiagnostic_Destroy_error
// func TestSzdiagnostic_Destroy_error(test *testing.T) {}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func addRecords(ctx context.Context, records []record.Record) {
	szEngine := getSzEngine(ctx)
	flags := senzing.SzWithoutInfo

	for _, record := range records {
		_, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		handleErrorWithPanic(err)
	}
}

func deleteRecords(ctx context.Context, records []record.Record) {
	szEngine := getSzEngine(ctx)
	flags := senzing.SzWithoutInfo

	for _, record := range records {
		_, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		handleErrorWithPanic(err)
	}
}

func getDefaultConfigID() int64 {
	return defaultConfigID
}

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

func getSettings() (string, error) {
	return "{}", nil
}

func getSzConfig(ctx context.Context) (senzing.SzConfig, error) {
	var err error
	if szConfigSingleton == nil {
		settings, err := getSettings()
		handleErrorWithPanic(err)

		grpcConnection := getGrpcConnection()
		szConfigSingleton = &szconfig.Szconfig{
			GrpcClient: szconfigpb.NewSzConfigClient(grpcConnection),
		}
		err = szConfigSingleton.SetLogLevel(ctx, logLevel)
		handleErrorWithPanic(err)

		if logLevel == "TRACE" {
			szConfigSingleton.SetObserverOrigin(ctx, observerOrigin)

			err = szConfigSingleton.RegisterObserver(ctx, observerSingleton)
			handleErrorWithPanic(err)

			err = szConfigSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			handleErrorWithPanic(err)
		}

		err = szConfigSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		handleErrorWithPanic(err)
	}

	return szConfigSingleton, err
}

func getSzConfigManager(ctx context.Context) senzing.SzConfigManager {
	var err error
	if szConfigManagerSingleton == nil {
		settings, err := getSettings()
		handleErrorWithPanic(err)

		grpcConnection := getGrpcConnection()
		szConfigManagerSingleton = &szconfigmanager.Szconfigmanager{
			GrpcClient: szconfigmanagerpb.NewSzConfigManagerClient(grpcConnection),
		}
		err = szConfigManagerSingleton.SetLogLevel(ctx, logLevel)
		handleErrorWithPanic(err)

		if logLevel == "TRACE" {
			szConfigManagerSingleton.SetObserverOrigin(ctx, observerOrigin)

			err = szConfigManagerSingleton.RegisterObserver(ctx, observerSingleton)
			handleErrorWithPanic(err)

			err = szConfigManagerSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			handleErrorWithPanic(err)
		}

		err = szConfigManagerSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		handleErrorWithPanic(err)
	}

	return szConfigManagerSingleton
}

func getSzDiagnostic(ctx context.Context) *Szdiagnostic {
	var err error

	if szDiagnosticSingleton == nil {
		grpcConnection := getGrpcConnection()
		szDiagnosticSingleton = &Szdiagnostic{
			GrpcClient: szpb.NewSzDiagnosticClient(grpcConnection),
		}
		err = szDiagnosticSingleton.SetLogLevel(ctx, logLevel)

		handleErrorWithPanic(err)

		if logLevel == "TRACE" {
			szDiagnosticSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szDiagnosticSingleton.RegisterObserver(ctx, observerSingleton)
			handleErrorWithPanic(err)
			err = szDiagnosticSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			handleErrorWithPanic(err)
		}
	}

	return szDiagnosticSingleton
}

func getSzDiagnosticAsInterface(ctx context.Context) senzing.SzDiagnostic {
	return getSzDiagnostic(ctx)
}

func getSzEngine(ctx context.Context) senzing.SzEngine {
	var err error

	if szEngineSingleton == nil {
		grpcConnection := getGrpcConnection()
		szEngineSingleton = &szengine.Szengine{
			GrpcClient: szenginepb.NewSzEngineClient(grpcConnection),
		}
		err = szEngineSingleton.SetLogLevel(ctx, logLevel)

		handleErrorWithPanic(err)

		if logLevel == "TRACE" {
			szEngineSingleton.SetObserverOrigin(ctx, observerOrigin)

			err = szEngineSingleton.RegisterObserver(ctx, observerSingleton)
			handleErrorWithPanic(err)

			err = szEngineSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			handleErrorWithPanic(err)

		}
	}

	return szEngineSingleton
}

func getTestObject(t *testing.T) *szdiagnostic.Szdiagnostic {
	t.Helper()
	ctx := t.Context()

	return getSzDiagnostic(ctx)
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

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	setup()
	code := m.Run()

	err := teardown()
	if err != nil {
		fmt.Print(err)
	}

	os.Exit(code)
}

func setup() {
	setupSenzingConfiguration()
	setupPurgeRepository()
}

func setupSenzingConfiguration() {
	ctx := test.Context()
	now := time.Now()

	// Create sz objects.

	szConfig, err := getSzConfig(ctx)
	handleErrorWithPanic(err)

	szConfigManager, err := getSzConfigManager(ctx)
	handleErrorWithPanic(err)

	// Create an in memory Senzing configuration.

	configHandle, err := szConfig.CreateConfig(ctx)
	handleErrorWithPanic(err)

	// Add data sources to in-memory Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		handleErrorWithPanic(err)

	}

	// Create a string representation of the in-memory configuration.

	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	handleErrorWithPanic(err)

	// Close szConfig in-memory object.

	err = szConfig.CloseConfig(ctx, configHandle)
	handleErrorWithPanic(err)

	// Persist the Senzing configuration to the Senzing repository as default.

	configComment := fmt.Sprintf("Created by szdiagnostic_test at %s", now.UTC())

	configID, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	handleErrorWithPanic(err)

	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	handleErrorWithPanic(err)

	szDiagnostic, err := getSzDiagnostic(ctx)
	handleErrorWithPanic(err)

	err = szDiagnostic.Reinitialize(ctx, configID)
	handleErrorWithPanic(err)

}

func setupPurgeRepository() {
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	err := szDiagnostic.PurgeRepository(ctx)
	handleErrorWithPanic(err)

}

func teardown() error {
	var err error
	return err
}
