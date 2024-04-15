package szconfigmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/sz-sdk-go-grpc/szengine"
	"github.com/senzing-garage/sz-sdk-go/sz"
	szconfigmanagerapi "github.com/senzing-garage/sz-sdk-go/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
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
	localLogger              logging.LoggingInterface
	szConfigManagerSingleton *SzConfigManager
	szConfigSingleton        sz.SzConfig
	szEngineSingleton        sz.SzEngine
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(localLogger.NewError(errorId, err), err)
}

func getGrpcConnection() *grpc.ClientConn {
	var err error
	if grpcConnection == nil {
		grpcConnection, err = grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Did not connect: %v\n", err)
		}
	}
	return grpcConnection
}

func getTestObject(ctx context.Context, test *testing.T) *SzConfigManager {
	return getSzConfigManager(ctx)
}

func getSzConfig(ctx context.Context) sz.SzConfig {
	if szConfigSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigSingleton = &szconfig.SzConfig{
			GrpcClient: szconfigpb.NewSzConfigClient(grpcConnection),
		}
	}
	return szConfigSingleton
}

func getSzConfigManager(ctx context.Context) *SzConfigManager {
	if szConfigManagerSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigManagerSingleton = &SzConfigManager{
			GrpcClient: szpb.NewSzConfigManagerClient(grpcConnection),
		}
	}
	return szConfigManagerSingleton
}

func getSzEngine(ctx context.Context) sz.SzEngine {
	if szEngineSingleton == nil {
		grpcConnection := getGrpcConnection()
		szEngineSingleton = &szengine.SzEngine{
			GrpcClient: szenginepb.NewSzEngineClient(grpcConnection),
		}
	}
	return szEngineSingleton
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

func printResult(test *testing.T, title string, result interface{}) {
	if printResults {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, err error, messageId string) {
	if err != nil {
		errorMessage := err.Error()[strings.Index(err.Error(), "{"):]
		var dictionary map[string]interface{}
		unmarshalErr := json.Unmarshal([]byte(errorMessage), &dictionary)
		if unmarshalErr != nil {
			test.Log("Unmarshal Error:", unmarshalErr.Error())
		}
		assert.Equal(test, messageId, dictionary["id"].(string))
	} else {
		assert.FailNow(test, "Should have failed with", messageId)
	}
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
	configId, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		return createError(5913, err)
	}

	err = szConfigManager.SetDefaultConfigId(ctx, configId)
	if err != nil {
		return createError(5914, err)
	}

	return err
}

func setup() error {
	ctx := context.TODO()
	var err error = nil

	options := []interface{}{
		&logging.OptionCallerSkip{Value: 4},
	}
	localLogger, err = logging.NewSenzingSdkLogger(ComponentId, szconfigmanagerapi.IdMessages, options...)
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
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSzConfigManager_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
}

func TestSzConfigManager_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	actual := szConfigManager.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzConfigManager_AddConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	now := time.Now()
	szConfig := getSzConfig(ctx)
	configHandle, err1 := szConfig.CreateConfig(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfig.CreateConfig()")
	}
	dataSourceCode := "GO_TEST_" + strconv.FormatInt(now.Unix(), 10)
	_, err2 := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "szConfig.AddDataSource()")
	}
	configDefinition, err3 := szConfig.ExportConfig(ctx, configHandle)
	if err3 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, configDefinition)
	}
	configComment := fmt.Sprintf("szconfigmanager_test at %s", now.UTC())
	actual, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	testError(test, ctx, err)
	printActual(test, actual)
}

func TestSzConfigManager_GetConfig(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	configID, err1 := szConfigManager.GetDefaultConfigId(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigId()")
	}
	actual, err := szConfigManager.GetConfig(ctx, configID)
	testError(test, ctx, err)
	printActual(test, actual)
}

func TestSzConfigManager_GetConfigList(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	actual, err := szConfigManager.GetConfigList(ctx)
	testError(test, ctx, err)
	printActual(test, actual)
}

func TestSzConfigManager_GetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	actual, err := szConfigManager.GetDefaultConfigId(ctx)
	testError(test, ctx, err)
	printActual(test, actual)
}

func TestSzConfigManager_ReplaceDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	currentDefaultConfigId, err1 := szConfigManager.GetDefaultConfigId(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()")
	}

	// TODO: This is kind of a cheeter.

	newDefaultConfigId, err2 := szConfigManager.GetDefaultConfigId(ctx)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigID()-2")
	}

	err := szConfigManager.ReplaceDefaultConfigId(ctx, currentDefaultConfigId, newDefaultConfigId)
	testError(test, ctx, err)
}

func TestSzConfigManager_SetDefaultConfigId(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	configId, err1 := szConfigManager.GetDefaultConfigId(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "szConfigManager.GetDefaultConfigId()")
	}
	err := szConfigManager.SetDefaultConfigId(ctx, configId)
	testError(test, ctx, err)
}

func TestSzConfigManager_Init(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	instanceName := "Test name"
	verboseLogging := sz.SZ_NO_LOGGING
	settings := "{}"
	err := szConfigManager.Initialize(ctx, instanceName, settings, verboseLogging)
	expectError(test, ctx, err, "senzing-60124002")
}

func TestSzConfigManager_Destroy(test *testing.T) {
	ctx := context.TODO()
	szConfigManager := getTestObject(ctx, test)
	err := szConfigManager.Destroy(ctx)
	expectError(test, ctx, err, "senzing-60124001")
}
