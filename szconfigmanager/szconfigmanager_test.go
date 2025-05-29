package szconfigmanager_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/env"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	defaultTruncation = 76
	instanceName      = "SzConfigManager Test"
	jsonIndentation   = "    "
	observerOrigin    = "SzConfigManager observer"
	originMessage     = "Machine: nn; Task: UnitTest"
	printErrors       = false
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

// Bad parameters

const (
	badConfigDefinition       = "\n\t"
	badConfigID               = int64(0)
	badCurrentDefaultConfigID = int64(0)
	badLogLevelName           = "BadLogLevelName"
	badNewDefaultConfigID     = int64(0)
	baseTen                   = 10
)

// Nil/empty parameters

var (
	nilConfigComment          string
	nilConfigDefinition       string
	nilConfigID               int64
	nilCurrentDefaultConfigID int64
	nilNewDefaultConfigID     int64
)

var (
	grpcAddress       = "0.0.0.0:8261"
	grpcConnection    *grpc.ClientConn
	logLevel          = env.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szConfigManagerSingleton *szconfigmanager.Szconfigmanager
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzconfigmanager_CreateConfigFromConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	configID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	printError(test, err1)
	require.NoError(test, err1)

	actual, err := szConfigManager.CreateConfigFromConfigID(ctx, configID)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_CreateConfigFromConfigID_badConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.CreateConfigFromConfigID(ctx, badConfigID)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).CreateConfigFromConfigID","error":{"function":"szconfigmanager.(*Szconfigmanager).createConfigFromConfigIDChoreography","text":"getConfig(0)","error":{"id":"SZSDK60024003","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}}`
	require.JSONEq(test, expectedErr, err.Error())
	assert.Nil(test, actual)
}

func TestSzconfigmanager_CreateConfigFromConfigID_nilConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.CreateConfigFromConfigID(ctx, nilConfigID)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).CreateConfigFromConfigID","error":{"function":"szconfigmanager.(*Szconfigmanager).createConfigFromConfigIDChoreography","text":"getConfig(0)","error":{"id":"SZSDK60024003","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}}`
	require.JSONEq(test, expectedErr, err.Error())
	assert.Nil(test, actual)
}

func TestSzconfigmanager_CreateConfigFromString(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	printError(test, err)
	require.NoError(test, err)
	configDefinition, err := szConfig.Export(ctx)
	printError(test, err)
	require.NoError(test, err)
	szConfig2, err := szConfigManager.CreateConfigFromString(ctx, configDefinition)
	printError(test, err)
	require.NoError(test, err)
	configDefinition2, err := szConfig2.Export(ctx)
	printError(test, err)
	require.NoError(test, err)
	assert.JSONEq(test, configDefinition, configDefinition2)
}

func TestSzconfigmanager_CreateConfigFromString_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	_, err := szConfigManager.CreateConfigFromString(ctx, badConfigDefinition)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).CreateConfigFromString","error":{"function":"szconfigmanager.(*Szconfigmanager).CreateConfigFromStringChoreography","text":"VerifyConfigDefinition","error":{"function":"szconfig.(*Szconfig).VerifyConfigDefinition","error":{"function":"szconfig.(*Szconfig).verifyConfigDefinitionChoreography","text":"load","error":{"id":"SZSDK60014009","reason":"SENZ3121|JSON Parsing Failure [code=1,offset=2]"}}}}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_CreateConfigFromTemplate(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.CreateConfigFromTemplate(ctx)
	printError(test, err)
	require.NoError(test, err)
	assert.NotEmpty(test, actual)
}

func TestSzconfigmanager_GetConfigs(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.GetConfigs(ctx)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_GetDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	actual, err := szConfigManager.GetDefaultConfigID(ctx)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_RegisterConfig(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	now := time.Now()
	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	printError(test, err)
	require.NoError(test, err)

	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), baseTen)
	_, err = szConfig.AddDataSource(ctx, dataSourceCode)
	printError(test, err)
	require.NoError(test, err)
	configDefinition, err := szConfig.Export(ctx)
	printError(test, err)
	require.NoError(test, err)

	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	actual, err := szConfigManager.RegisterConfig(ctx, configDefinition, configComment)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_RegisterConfig_badConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	now := time.Now()
	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	_, err := szConfigManager.RegisterConfig(ctx, badConfigDefinition, configComment)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).RegisterConfig","error":{"id":"SZSDK60024001","reason":"SENZ0028|Invalid JSON config document"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_RegisterConfig_nilConfigDefinition(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	now := time.Now()
	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	_, err := szConfigManager.RegisterConfig(ctx, nilConfigDefinition, configComment)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).RegisterConfig","error":{"id":"SZSDK60024001","reason":"SENZ0028|Invalid JSON config document"}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_RegisterConfig_nilConfigComment(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	printError(test, err)
	require.NoError(test, err)
	configDefinition, err := szConfig.Export(ctx)
	printError(test, err)
	require.NoError(test, err)
	actual, err := szConfigManager.RegisterConfig(ctx, configDefinition, nilConfigComment)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_ReplaceDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	currentDefaultConfigID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	printError(test, err1)
	require.NoError(test, err1)

	// IMPROVE: This is kind of a cheater.

	newDefaultConfigID, err2 := szConfigManager.GetDefaultConfigID(ctx)
	printError(test, err2)
	require.NoError(test, err2)

	err := szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, newDefaultConfigID)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzconfigmanager_ReplaceDefaultConfigID_badCurrentDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	newDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	printError(test, err)
	require.NoError(test, err)
	err = szConfigManager.ReplaceDefaultConfigID(ctx, badCurrentDefaultConfigID, newDefaultConfigID)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzReplaceConflict)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).ReplaceDefaultConfigID","error":{"id":"SZSDK60024007","reason":"SENZ7245|Current configuration ID does not match specified data ID [0]."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_ReplaceDefaultConfigID_badNewDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	currentDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	printError(test, err)
	require.NoError(test, err)
	err = szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, badNewDefaultConfigID)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).ReplaceDefaultConfigID","error":{"id":"SZSDK60024007","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_ReplaceDefaultConfigID_nilCurrentDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	newDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	printError(test, err)
	require.NoError(test, err)
	err = szConfigManager.ReplaceDefaultConfigID(ctx, nilCurrentDefaultConfigID, newDefaultConfigID)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzReplaceConflict)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).ReplaceDefaultConfigID","error":{"id":"SZSDK60024007","reason":"SENZ7245|Current configuration ID does not match specified data ID [0]."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_ReplaceDefaultConfigID_nilNewDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	currentDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	printError(test, err)
	require.NoError(test, err)
	err = szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, nilNewDefaultConfigID)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).ReplaceDefaultConfigID","error":{"id":"SZSDK60024007","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_SetDefaultConfig(test *testing.T) {
	ctx := test.Context()
	now := time.Now()
	szConfigManager := getTestObject(test)
	defaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	printError(test, err)
	require.NoError(test, err)
	szConfig, err := szConfigManager.CreateConfigFromConfigID(ctx, defaultConfigID)
	printError(test, err)
	require.NoError(test, err)

	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), baseTen)
	_, err = szConfig.AddDataSource(ctx, dataSourceCode)
	printError(test, err)
	require.NoError(test, err)
	configDefintion, err := szConfig.Export(ctx)
	printError(test, err)
	require.NoError(test, err)
	configID, err := szConfigManager.SetDefaultConfig(ctx, configDefintion, "Added "+dataSourceCode)
	printError(test, err)
	require.NoError(test, err)
	require.NotZero(test, configID)
}

func TestSzconfigmanager_SetDefaultConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	printError(test, err)
	require.NoError(test, err)
	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzconfigmanager_SetDefaultConfigID_badConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	err := szConfigManager.SetDefaultConfigID(ctx, badConfigID)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).SetDefaultConfigID","error":{"id":"SZSDK60024008","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

func TestSzconfigmanager_SetDefaultConfigID_nilConfigID(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	err := szConfigManager.SetDefaultConfigID(ctx, nilConfigID)
	printError(test, err)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)

	expectedErr := `{"function":"szconfigmanager.(*Szconfigmanager).SetDefaultConfigID","error":{"id":"SZSDK60024008","reason":"SENZ7221|No engine configuration registered with data ID [0]."}}`
	require.JSONEq(test, expectedErr, err.Error())
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzconfigmanager_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	_ = szConfigManager.SetLogLevel(ctx, badLogLevelName)
}

func TestSzconfigmanager_SetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	szConfigManager.SetObserverOrigin(ctx, originMessage)
}

func TestSzconfigmanager_GetObserverOrigin(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	szConfigManager.SetObserverOrigin(ctx, originMessage)
	actual := szConfigManager.GetObserverOrigin(ctx)
	assert.Equal(test, originMessage, actual)
}

func TestSzconfigmanager_UnregisterObserver(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	err := szConfigManager.UnregisterObserver(ctx, observerSingleton)
	printError(test, err)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzconfigmanager_AsInterface(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getSzConfigManagerAsInterface(ctx)
	actual, err := szConfigManager.GetConfigs(ctx)
	printError(test, err)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_Initialize(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	settings := getSettings()
	err := szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzconfigmanager_Initialize_error(test *testing.T) {
	// IMPROVE: Implement TestSzconfigmanager_Initialize_error
	_ = test
}

func TestSzconfigmanager_Destroy(test *testing.T) {
	ctx := test.Context()
	szConfigManager := getTestObject(test)
	err := szConfigManager.Destroy(ctx)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzconfigmanager_Destroy_withObserver(test *testing.T) {
	ctx := test.Context()
	szConfigManagerSingleton = nil
	szConfigManager := getTestObject(test)
	err := szConfigManager.Destroy(ctx)
	printError(test, err)
	require.NoError(test, err)
}

func TestSzconfigmanager_Destroy_error(test *testing.T) {
	// IMPROVE: Implement TestSzconfigmanager_Destroy_error
	_ = test
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

// func getSzConfig(ctx context.Context) *szconfig.Szconfig {
// 	var szConfig *szconfig.Szconfig

// 	szConfigManager := getSzConfigManager(ctx)
// 	szConfigForExport, err := szConfigManager.CreateConfigFromTemplate(ctx)
// 	panicOnError(err)

// 	configDefinition, err := szConfigForExport.Export(ctx)
// 	panicOnError(err)

// 	grpcConnection := getGrpcConnection()
// 	szConfig = &szconfig.Szconfig{
// 		GrpcClient: szconfigpb.NewSzConfigClient(grpcConnection),
// 	}
// 	err = szConfig.SetLogLevel(ctx, logLevel)
// 	panicOnError(err)

// 	err = szConfig.Import(ctx, configDefinition)
// 	panicOnError(err)

// 	if logLevel == "TRACE" {
// 		szConfig.SetObserverOrigin(ctx, observerOrigin)

// 		err = szConfig.RegisterObserver(ctx, observerSingleton)
// 		panicOnError(err)

// 		err = szConfig.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
// 		panicOnError(err)
// 	}

// 	return szConfig
// }

func getSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	_ = ctx

	return &szabstractfactory.Szabstractfactory{
		GrpcConnection: getGrpcConnection(),
	}
}

func getSzConfigManager(ctx context.Context) *szconfigmanager.Szconfigmanager {
	var err error

	if szConfigManagerSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigManagerSingleton = &szconfigmanager.Szconfigmanager{
			GrpcClient:         szpb.NewSzConfigManagerClient(grpcConnection),
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

func getSzConfigManagerAsInterface(ctx context.Context) senzing.SzConfigManager {
	return getSzConfigManager(ctx)
}

func getTestObject(t *testing.T) *szconfigmanager.Szconfigmanager {
	t.Helper()

	return getSzConfigManager(t.Context())
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
}

func setupSenzingConfiguration() {
	ctx := context.TODO()
	now := time.Now()

	// Create sz objects.

	// szConfig := getSzConfig(ctx)
	szConfigManager := getSzConfigManager(ctx)
	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	panicOnError(err)

	// Add data sources to in-memory Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, dataSourceCode)
		panicOnError(err)
	}

	// Persist the Senzing configuration to the Senzing repository as default.

	configComment := fmt.Sprintf("Created by szconfigmanager_test at %s", now.UTC())
	configDefinition, err := szConfig.Export(ctx)
	panicOnError(err)

	configID, err := szConfigManager.RegisterConfig(ctx, configDefinition, configComment)
	panicOnError(err)

	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	panicOnError(err)
}
