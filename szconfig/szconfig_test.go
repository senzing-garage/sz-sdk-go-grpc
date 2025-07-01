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
	jsonIndentation   = "    "
	observerOrigin    = "SzConfig observer"
	originMessage     = "Machine: nn; Task: UnitTest"
	printErrors       = false
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

func TestSzconfig_RegisterDataSource(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.RegisterDataSource(ctx, dataSourceCode)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzconfig_RegisterDataSource_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.RegisterDataSource(ctx, badDataSourceCode)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfig.(*Szconfig).RegisterDataSource","error":{"function":"szconfigserver.(*SzConfigServer).RegisterDataSource","text":"RegisterDataSource: \n\tGO_TEST","error":{"function":"szconfig.(*Szconfig).RegisterDataSource","error":{"function":"szconfig.(*Szconfig).registerDataSource","text":"registerDataSource: \n\tGO_TEST","error":{"id":"SZSDK60014001","reason":"SENZ3121|JSON Parsing Failure [code=12,offset=15]"}}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfig_RegisterDataSource_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.RegisterDataSource(ctx, nilDataSourceCode)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfig.(*Szconfig).RegisterDataSource","error":{"function":"szconfigserver.(*SzConfigServer).RegisterDataSource","text":"RegisterDataSource: ","error":{"function":"szconfig.(*Szconfig).RegisterDataSource","error":{"function":"szconfig.(*Szconfig).registerDataSource","text":"registerDataSource: ","error":{"id":"SZSDK60014001","reason":"SENZ7313|A non-empty value for [DSRC_CODE] must be specified."}}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfig_UnregisterDataSource(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.GetDataSourceRegistry(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)

	_, _ = szConfig.RegisterDataSource(ctx, dataSourceCode)
	actual, err = szConfig.GetDataSourceRegistry(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)

	_, err = szConfig.UnregisterDataSource(ctx, dataSourceCode)
	printDebug(test, err, actual)
	require.NoError(test, err)
	actual, err = szConfig.GetDataSourceRegistry(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzconfig_UnregisterDataSource_badDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.UnregisterDataSource(ctx, badDataSourceCode)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfig.(*Szconfig).UnregisterDataSource","error":{"function":"szconfigserver.(*SzConfigServer).UnregisterDataSource","text":"UnregisterDataSource: \n\tGO_TEST","error":{"function":"szconfig.(*Szconfig).UnregisterDataSource","error":{"function":"szconfig.(*Szconfig).unregisterDataSourceChoreography","text":"unregisterDataSource(\n\tGO_TEST)","error":{"id":"SZSDK60014004","reason":"SENZ3121|JSON Parsing Failure [code=12,offset=15]"}}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfig_UnregisterDataSource_nilDataSourceCode(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.UnregisterDataSource(ctx, nilDataSourceCode)
	printDebug(test, err, actual)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfig.(*Szconfig).UnregisterDataSource","error":{"function":"szconfigserver.(*SzConfigServer).UnregisterDataSource","text":"UnregisterDataSource: ","error":{"function":"szconfig.(*Szconfig).UnregisterDataSource","error":{"function":"szconfig.(*Szconfig).unregisterDataSourceChoreography","text":"unregisterDataSource()","error":{"id":"SZSDK60014004","reason":"SENZ7313|A non-empty value for [DSRC_CODE] must be specified."}}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfig_Export(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.Export(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzconfig_GetDataSourceRegistry(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	actual, err := szConfig.GetDataSourceRegistry(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

func TestSzconfig_Import(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	configDefinition, err := szConfig.Export(ctx)
	printDebug(test, err)
	require.NoError(test, err)
	err = szConfig.Import(ctx, configDefinition)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_Import_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Import(ctx, badConfigDefinition)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_Import_nilConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Import(ctx, nilConfigDefinition)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_VerifyConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	configDefinition, err := szConfig.Export(ctx)
	printDebug(test, err)
	require.NoError(test, err)
	err = szConfig.VerifyConfigDefinition(ctx, configDefinition)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_VerifyConfigDefinition_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.VerifyConfigDefinition(ctx, badConfigDefinition)
	printDebug(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfig.(*Szconfig).VerifyConfigDefinition","error":{"function":"szconfigserver.(*SzConfigServer).createSzConfig","error":{"function":"szconfigmanager.(*Szconfigmanager).CreateConfigFromStringChoreography","text":"VerifyConfigDefinition","error":{"function":"szconfig.(*Szconfig).VerifyConfigDefinition","error":{"function":"szconfig.(*Szconfig).verifyConfigDefinitionChoreography","text":"load","error":{"id":"SZSDK60014009","reason":"SENZ3121|JSON Parsing Failure [code=3,offset=0]"}}}}}}`
	require.JSONEq(test, expectedErr, err.Error())
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
	origin := originMessage
	szConfig.SetObserverOrigin(ctx, origin)
}

func TestSzconfig_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	origin := originMessage
	szConfig.SetObserverOrigin(ctx, origin)
	actual := szConfig.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzconfig_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.UnregisterObserver(ctx, observerSingleton)
	printDebug(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzconfig_AsInterface(test *testing.T) {
	ctx := test.Context()
	szConfig := getSzConfigAsInterface(ctx)
	actual, err := szConfig.GetDataSourceRegistry(ctx)
	printDebug(test, err, actual)
	require.NoError(test, err)
}

func TestSzconfig_Initialize(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	settings := getSettings()
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_Initialize_badSettings(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Initialize(ctx, instanceName, badSettings, verboseLogging)
	assert.NoError(test, err)
}

func TestSzconfig_Initialize_error(test *testing.T) {
	// IMPROVE: Implement TestSzconfig_Initialize_error
	_ = test
}

func TestSzconfig_Initialize_again(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	settings := getSettings()
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_Destroy(test *testing.T) {
	ctx := test.Context()
	szConfig := getTestObject(test)
	err := szConfig.Destroy(ctx)
	printDebug(test, err)
	require.NoError(test, err)
}

func TestSzconfig_Destroy_error(test *testing.T) {
	// IMPROVE: Implement TestSzconfig_Destroy_error
	_ = test
}

func TestSzconfig_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szConfigSingleton = nil
	szConfig := getTestObject(test)
	err := szConfig.Destroy(ctx)
	printDebug(test, err)
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
			GrpcClient:         szconfigmanagerpb.NewSzConfigManagerClient(grpcConnection),
			GrpcClientSzConfig: szpb.NewSzConfigClient(grpcConnection),
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

	return getSzConfig(t.Context())
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
