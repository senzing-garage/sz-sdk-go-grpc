package szconfigmanager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szconfigmanagerapi "github.com/senzing-garage/sz-sdk-go/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	badConfigDefinition       = "\n\t"
	badConfigID               = int64(0)
	badCurrentDefaultConfigID = int64(0)
	badLogLevelName           = "BadLogLevelName"
	defaultTruncation         = 76
	instanceName              = "SzConfigManager Test"
	observerOrigin            = "SzConfigManager observer"
	printResults              = false
	verboseLogging            = senzing.SzNoLogging
)

var (
	grpcAddress       = "localhost:8261"
	grpcConnection    *grpc.ClientConn
	logger            logging.Logging
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szConfigManagerSingleton *Szconfigmanager
	szConfigSingleton        senzing.SzConfig
)

// ----------------------------------------------------------------------------
// Interface functions - test
// ----------------------------------------------------------------------------

func TestSzconfigmanager_AddConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	now := time.Now()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), 10)
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	require.NoError(test, err)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	require.NoError(test, err)
	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	actual, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_GetConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	configID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()")
	}
	actual, err := szConfigManager.GetConfig(ctx, configID)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_GetConfig_badConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	actual, err := szConfigManager.GetConfig(ctx, badConfigID)
	assert.Equal(test, "", actual)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_GetConfigs(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	actual, err := szConfigManager.GetConfigs(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// TODO: Implement TestSzconfigmanager_GetConfigs_error
// func TestSzconfigmanager_GetConfigs_error(test *testing.T) {}

func TestSzconfigmanager_GetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	actual, err := szConfigManager.GetDefaultConfigID(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// TODO: Implement TestSzconfigmanager_GetDefaultConfigID_error
// func TestSzconfigmanager_GetDefaultConfigID_error(test *testing.T) {}

func TestSzconfigmanager_ReplaceDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	currentDefaultConfigID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()")
	}

	// TODO: This is kind of a cheater.

	newDefaultConfigID, err2 := szConfigManager.GetDefaultConfigID(ctx)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()-2")
	}

	err := szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, newDefaultConfigID)
	require.NoError(test, err)
}

func TestSzconfigmanager_ReplaceDefaultConfigID_badCurrentDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	newDefaultConfigID, err2 := szConfigManager.GetDefaultConfigID(ctx)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()-2")
	}
	err := szConfigManager.ReplaceDefaultConfigID(ctx, badCurrentDefaultConfigID, newDefaultConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_ReplaceDefaultConfigID_badNewDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	currentDefaultConfigID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()")
	}
	newDefaultConfigID := int64(0)
	err := szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, newDefaultConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

func TestSzconfigmanager_SetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	configID, err1 := szConfigManager.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()")
	}
	err := szConfigManager.SetDefaultConfigID(ctx, configID)
	require.NoError(test, err)
}

func TestSzconfigmanager_SetDefaultConfigID_badConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	err := szConfigManager.SetDefaultConfigID(ctx, badConfigID)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzconfigmanager_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	_ = szConfigManager.SetLogLevel(ctx, badLogLevelName)
}

func TestSzconfigmanager_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
}

func TestSzconfigmanager_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	actual := szConfigManager.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzconfigmanager_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	err := szConfigManager.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzconfigmanager_AsInterface(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getSzConfigManagerAsInterface(ctx)
	actual, err := szConfigManager.GetConfigs(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfigmanager_Initialize(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	instanceName := "Test name"
	verboseLogging := senzing.SzNoLogging
	settings, err := getSettings()
	require.NoError(test, err)
	err = szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

// TODO: Implement TestSzconfigmanager_Initialize_error
// func TestSzconfigmanager_Initialize_error(test *testing.T) {}

func TestSzconfigmanager_Destroy(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	err := szConfigManager.Destroy(ctx)
	require.NoError(test, err)
}

func TestSzconfigmanager_Destroy_withObserver(test *testing.T) {
	ctx := context.TODO()
	szConfigManagerSingleton = nil
	szConfigManager := getTestObject(ctx, test)
	err := szConfigManager.Destroy(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzconfigmanager_Destroy_error
// func TestSzconfigmanager_Destroy_error(test *testing.T) {}

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
	}
	return grpcConnection
}

func getSettings() (string, error) {
	return "{}", nil
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

func getSzConfigManager(ctx context.Context) *Szconfigmanager {
	_ = ctx
	if szConfigManagerSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigManagerSingleton = &Szconfigmanager{
			GrpcClient: szpb.NewSzConfigManagerClient(grpcConnection),
		}
	}
	return szConfigManagerSingleton
}

func getSzConfigManagerAsInterface(ctx context.Context) senzing.SzConfigManager {
	return getSzConfigManager(ctx)
}

func getTestObject(ctx context.Context, test *testing.T) *Szconfigmanager {
	_ = test
	return getSzConfigManager(ctx)
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

	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		return createError(5907, err)
	}

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, datasourceName := range datasourceNames {
		_, err := szConfig.AddDataSource(ctx, configHandle, datasourceName)
		if err != nil {
			return createError(5908, err)
		}
	}

	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	szConfigManager := getSzConfigManager(ctx)
	configComment := fmt.Sprintf("Created by szconfigmanager_test at %s", now.UTC())
	configID, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		return createError(5913, err)
	}

	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return createError(5914, err)
	}

	return err
}

func setup() error {
	ctx := context.TODO()
	var err error

	options := []interface{}{
		&logging.OptionCallerSkip{Value: 4},
	}
	logger, err = logging.NewSenzingLogger(ComponentID, szconfigmanagerapi.IDMessages, options...)
	if err != nil {
		return createError(5901, err)
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfig(ctx)
	if err != nil {
		return createError(5920, err)
	}

	return err
}

func teardown() error {
	var err error
	return err
}
