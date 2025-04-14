package szconfig_test

import (
	"context"
	"fmt"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/env"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	baseTen           = 10
	dataSourceCode    = "GO_TEST"
	defaultTruncation = 76
	instanceName      = "SzConfig Test"
	observerOrigin    = "SzConfig observer"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

// Bad parameters

const (
	badConfigDefinition = "}{"
	badConfigHandle     = uintptr(0)
	badDataSourceCode   = "\n\tGO_TEST"
	badLogLevelName     = "BadLogLevelName"
	badSettings         = "{]"
)

// Nil/empty parameters

var (
	nilConfigDefinition string
	nilDataSourceCode   string
)

var (
	grpcAddress       = "0.0.0.0:8261"
	grpcConnection    *grpc.ClientConn
	logLevel          = env.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szConfigSingleton        *szconfig.Szconfig
	szConfigManagerSingleton *szconfigmanager.Szconfigmanager
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func TestSzconfig_AddDataSource(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.AddDataSource(ctx, dataSourceCode)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_AddDataSource_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.AddDataSource(ctx, badDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzconfig_AddDataSource_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.AddDataSource(ctx, nilDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzconfig_DeleteDataSource(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.GetDataSources(ctx)
	require.NoError(test, err)
	printResult(test, "Original", actual)

	_, _ = szConfig.AddDataSource(ctx, dataSourceCode)
	actual, err = szConfig.GetDataSources(ctx)
	require.NoError(test, err)
	printResult(test, "     Add", actual)

	_, err = szConfig.DeleteDataSource(ctx, dataSourceCode)
	require.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx)
	require.NoError(test, err)
	printResult(test, "  Delete", actual)
}

func TestSzconfig_DeleteDataSource_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	_, err := szConfig.DeleteDataSource(ctx, badDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
}

func TestSzconfig_DeleteDataSource_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	_, err := szConfig.DeleteDataSource(ctx, nilDataSourceCode)
	require.NoError(test, err)
}

func TestSzconfig_Export(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.Export(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_GetDataSources(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.GetDataSources(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

func TestSzconfig_Import(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	configDefinition, err := szConfig.Export(ctx)
	require.NoError(test, err)
	err = szConfig.Import(ctx, configDefinition)
	require.NoError(test, err)
}

func TestSzconfig_Import_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Import(ctx, badConfigDefinition)
	require.NoError(test, err)
}

func TestSzconfig_Import_nilConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Import(ctx, nilConfigDefinition)
	require.NoError(test, err)
}

func TestSzconfig_VerifyConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	configDefinition, err := szConfig.Export(ctx)
	require.NoError(test, err)
	err = szConfig.VerifyConfigDefinition(ctx, configDefinition)
	require.NoError(test, err)
}

func TestSzconfig_VerifyConfigDefinition_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.VerifyConfigDefinition(ctx, badConfigDefinition)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzconfig_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzconfig_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
}

func TestSzconfig_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	actual := szConfig.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzconfig_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzconfig_AsInterface(test *testing.T) {
	ctx := test.Context()
	szConfig := getSzConfigAsInterface(ctx)
	actual, err := szConfig.GetDataSources(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_Initialize(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	settings := getSettings()
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

func TestSzconfig_Initialize_badSettings(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Initialize(ctx, instanceName, badSettings, verboseLogging)
	assert.NoError(test, err)
}

// TODO: Implement TestSzconfig_Initialize_error
// func TestSzconfig_Initialize_error(test *testing.T) {}

func TestSzconfig_Initialize_again(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	settings := getSettings()
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

func TestSzconfig_Destroy(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Destroy(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzconfig_Destroy_error
// func TestSzconfig_Destroy_error(test *testing.T) {}

func TestSzconfig_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szConfigSingleton = nil
	szConfig := getTestObject(test)
	err := szConfig.Destroy(ctx)
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

func getSzConfigAsInterface(ctx context.Context) senzing.SzConfig {
	return getSzConfig(ctx)
}

func getSzConfigManager(ctx context.Context) *szconfigmanager.Szconfigmanager {
	var err error
	if szConfigManagerSingleton == nil {

		grpcConnection := getGrpcConnection()
		szConfigManagerSingleton = &szconfigmanager.Szconfigmanager{
			GrpcClient: szconfigmanagerpb.NewSzConfigManagerClient(grpcConnection),
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

func getSzConfig(ctx context.Context) *szconfig.Szconfig {
	var szConfig *szconfig.Szconfig

	szConfigManager := getSzConfigManager(ctx)
	szConfigForExport, err := szConfigManager.CreateConfigFromTemplate(ctx)
	panicOnError(err)

	configDefinition, err := szConfigForExport.Export(ctx)
	panicOnError(err)

	grpcConnection := getGrpcConnection()
	szConfig = &szconfig.Szconfig{
		GrpcClient: szpb.NewSzConfigClient(grpcConnection),
	}
	err = szConfig.SetLogLevel(ctx, logLevel)
	panicOnError(err)

	err = szConfig.Import(ctx, configDefinition)
	panicOnError(err)

	if logLevel == "TRACE" {
		szConfig.SetObserverOrigin(ctx, observerOrigin)

		err = szConfig.RegisterObserver(ctx, observerSingleton)
		panicOnError(err)

		err = szConfig.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
		panicOnError(err)
	}

	return szConfig
}

func getTestObject(t *testing.T) *szconfig.Szconfig {
	t.Helper()
	ctx := t.Context()

	return getSzConfig(ctx)
}

func handleError(err error) {
	if err != nil {
		safePrintln("Error:", err)
	}
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
