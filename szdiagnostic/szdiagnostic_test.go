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
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-grpc/szdiagnostic"
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
	jsonIndentation   = "    "
	observerOrigin    = "SzDiagnostic observer"
	originMessage     = "Machine: nn; Task: UnitTest"
	printErrors       = false
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
	szDiagnosticSingleton    *szdiagnostic.Szdiagnostic
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
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzdiagnostic_CheckDatastorePerformance_badSecondsToRun(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, badSecondsToRun)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzdiagnostic_CheckDatastorePerformance_nilSecondsToRun(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.CheckDatastorePerformance(ctx, nilSecondsToRun)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzdiagnostic_GetDatastoreInfo(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	actual, err := szDiagnostic.GetDatastoreInfo(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
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
	printDebug(test, err, actual)
	require.NoError(test, err)
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
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szdiagnostic.(*Szdiagnostic).GetFeature","error":{"function":"szdiagnosticserver.(*SzDiagnosticServer).GetFeature","error":{"function":"szdiagnostic.(*Szdiagnostic).GetFeature","error":{"id":"SZSDK60034004","reason":"SENZ0057|Unknown feature ID value '-1'"}}}}`
	require.JSONEq(test, expectedErr, err.Error())
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
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSz)

	expectedErr := `{"function":"szdiagnostic.(*Szdiagnostic).GetFeature","error":{"function":"szdiagnosticserver.(*SzDiagnosticServer).GetFeature","error":{"function":"szdiagnostic.(*Szdiagnostic).GetFeature","error":{"id":"SZSDK60034004","reason":"SENZ0057|Unknown feature ID value '0'"}}}}`
	require.JSONEq(test, expectedErr, err.Error())
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
	szDiagnostic.SetObserverOrigin(ctx, originMessage)
}

func TestSzdiagnostic_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	szDiagnostic.SetObserverOrigin(ctx, originMessage)
	actual := szDiagnostic.GetObserverOrigin(ctx)
	assert.Equal(test, originMessage, actual)
}

func TestSzdiagnostic_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szDiagnostic := getTestObject(test)
	err := szDiagnostic.UnregisterObserver(ctx, observerSingleton)
	printDebug(test, err)
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
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzdiagnostic_Initialize(test *testing.T) {
	ctx := test.Context()
	grpcConnection := getGrpcConnection()
	szDiagnostic := &szdiagnostic.Szdiagnostic{
		GrpcClient: szpb.NewSzDiagnosticClient(grpcConnection),
	}
	settings := getSettings()
	configID := senzing.SzInitializeWithDefaultConfiguration
	err := szDiagnostic.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzdiagnostic_Initialize_error(test *testing.T) {
	// IMPROVE: Implement TestSzdiagnostic_Initialize_error
	_ = test
}

func TestSzdiagnostic_Initialize_withConfigId(test *testing.T) {
	ctx := test.Context()
	grpcConnection := getGrpcConnection()
	szDiagnostic := &szdiagnostic.Szdiagnostic{
		GrpcClient: szpb.NewSzDiagnosticClient(grpcConnection),
	}
	settings := getSettings()
	configID := getDefaultConfigID()
	err := szDiagnostic.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzdiagnostic_Initialize_withConfigId_badConfigID(test *testing.T) {
	// IMPROVE: Implement TestSzdiagnostic_Initialize_withConfigId_badConfigID
	_ = test
}

// func TestSzdiagnostic_Reinitialize(test *testing.T) {
// 	ctx := test.Context()
// 	szDiagnostic := getTestObject(ctx, test)
// 	configID := getDefaultConfigID()
// 	err := szDiagnostic.Reinitialize(ctx, configID)
// 	require.NoError(test, err)
// }

func TestSzdiagnostic_Reinitialize_error(test *testing.T) {
	// IMPROVE: Implement TestSzdiagnostic_Reinitialize_error
	_ = test
}

func TestSzdiagnostic_Destroy(test *testing.T) {
	ctx := test.Context()
	szDiagnosticSingleton = nil
	szDiagnostic := getTestObject(test)
	err := szDiagnostic.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzdiagnostic_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szDiagnosticSingleton = nil
	szDiagnostic := getTestObject(test)
	err := szDiagnostic.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzdiagnostic_Destroy_error(test *testing.T) {
	// IMPROVE: Implement TestSzdiagnostic_Destroy_error
	_ = test
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func addRecords(ctx context.Context, records []record.Record) {
	szEngine := getSzEngine(ctx)
	flags := senzing.SzWithoutInfo

	for _, record := range records {
		_, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		panicOnError(err)
	}
}

func deleteRecords(ctx context.Context, records []record.Record) {
	szEngine := getSzEngine(ctx)
	flags := senzing.SzWithoutInfo

	for _, record := range records {
		_, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		panicOnError(err)
	}
}

func getDefaultConfigID() int64 {
	return defaultConfigID
}

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

func getSzConfigManager(ctx context.Context) senzing.SzConfigManager {
	var err error

	if szConfigManagerSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigManagerSingleton = &szconfigmanager.Szconfigmanager{
			GrpcClient:         szconfigmanagerpb.NewSzConfigManagerClient(grpcConnection),
			GrpcClientSzConfig: szconfigpb.NewSzConfigClient(grpcConnection),
		}
		err = szConfigManagerSingleton.SetLogLevel(ctx, logLevel)

		panicOnError(err)

		if logLevel == "TRACE" {
			szConfigManagerSingleton.SetObserverOrigin(ctx, observerOrigin)

			err = szConfigManagerSingleton.RegisterObserver(ctx, observerSingleton)
			panicOnError(err)

			err = szConfigManagerSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			panicOnError(err)
		}
	}

	return szConfigManagerSingleton
}

func getSzDiagnostic(ctx context.Context) *szdiagnostic.Szdiagnostic {
	if szDiagnosticSingleton == nil {
		grpcConnection := getGrpcConnection()
		szDiagnosticSingleton = &szdiagnostic.Szdiagnostic{
			GrpcClient: szpb.NewSzDiagnosticClient(grpcConnection),
		}
		err := szDiagnosticSingleton.SetLogLevel(ctx, logLevel)
		panicOnError(err)

		if logLevel == "TRACE" {
			szDiagnosticSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szDiagnosticSingleton.RegisterObserver(ctx, observerSingleton)
			panicOnError(err)
			err = szDiagnosticSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			panicOnError(err)
		}
	}

	return szDiagnosticSingleton
}

func getSzDiagnosticAsInterface(ctx context.Context) senzing.SzDiagnostic {
	return getSzDiagnostic(ctx)
}

func getSzEngine(ctx context.Context) senzing.SzEngine {
	if szEngineSingleton == nil {
		grpcConnection := getGrpcConnection()
		szEngineSingleton = &szengine.Szengine{
			GrpcClient: szenginepb.NewSzEngineClient(grpcConnection),
		}
		err := szEngineSingleton.SetLogLevel(ctx, logLevel)
		panicOnError(err)

		if logLevel == "TRACE" {
			szEngineSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szEngineSingleton.RegisterObserver(ctx, observerSingleton)
			panicOnError(err)
			err = szEngineSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			panicOnError(err)
		}
	}

	return szEngineSingleton
}

func getTestObject(t *testing.T) *szdiagnostic.Szdiagnostic {
	t.Helper()

	return getSzDiagnostic(t.Context())
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

func printDebug(t *testing.T, err error, items ...any) {
	t.Helper()

	if printErrors {
		if err != nil {
			t.Logf("Error: %s\n", err.Error())
		}
	}

	if printResults {
		for _, item := range items {
			outLine := truncator.Truncate(fmt.Sprintf("%v", item), defaultTruncation, "...", truncator.PositionEnd)
			t.Logf("Result: %s\n", outLine)
		}
	}
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	os.Exit(code)
}

func setup() {
	setupSenzingConfiguration()
	setupPurgeRepository()
}

func setupSenzingConfiguration() {
	ctx := context.TODO()
	now := time.Now()

	// Create sz objects.

	szConfigManager := getSzConfigManager(ctx)
	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	panicOnError(err)

	// Add data sources to Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, dataSourceCode)
		panicOnError(err)
	}

	// Persist the Senzing configuration to the Senzing repository as default.

	configComment := fmt.Sprintf("Created by szdiagnostic_test at %s", now.UTC())
	configDefinition, err := szConfig.Export(ctx)
	panicOnError(err)

	configID, err := szConfigManager.RegisterConfig(ctx, configDefinition, configComment)
	panicOnError(err)

	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	panicOnError(err)
}

func setupPurgeRepository() {
	ctx := context.TODO()
	szDiagnostic := getSzDiagnostic(ctx)
	err := szDiagnostic.PurgeRepository(ctx)
	panicOnError(err)
}
