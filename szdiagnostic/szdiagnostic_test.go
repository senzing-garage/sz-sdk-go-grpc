package szdiagnostic

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-grpc/szengine"
	"github.com/senzing-garage/sz-sdk-go/sz"
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
	localLogger              logging.LoggingInterface
	szConfigManagerSingleton sz.SzConfigManager
	szConfigSingleton        sz.SzConfig
	szDiagnosticSingleton    *SzDiagnostic
	szEngineSingleton        sz.SzEngine
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return szerror.Cast(localLogger.NewError(errorId, err), err)
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

func getTestObject(ctx context.Context, test *testing.T) *SzDiagnostic {
	_ = test
	return getSzDiagnostic(ctx)
}

func getSzConfig(ctx context.Context) sz.SzConfig {
	_ = ctx
	if szConfigSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigSingleton = &szconfig.SzConfig{
			GrpcClient: szconfigpb.NewSzConfigClient(grpcConnection),
		}
	}
	return szConfigSingleton
}

func getSzConfigManager(ctx context.Context) sz.SzConfigManager {
	_ = ctx
	if szConfigManagerSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigManagerSingleton = &szconfigmanager.SzConfigManager{
			GrpcClient: szconfigmanagerpb.NewSzConfigManagerClient(grpcConnection),
		}
	}
	return szConfigManagerSingleton
}

func getSzDiagnostic(ctx context.Context) *SzDiagnostic {
	_ = ctx
	if szDiagnosticSingleton == nil {
		grpcConnection := getGrpcConnection()
		szDiagnosticSingleton = &SzDiagnostic{
			GrpcClient: szpb.NewSzDiagnosticClient(grpcConnection),
		}
	}
	return szDiagnosticSingleton
}

func getSzEngine(ctx context.Context) sz.SzEngine {
	_ = ctx
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

func testError(test *testing.T, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, err error, messageId string) {
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

func testErrorNoFail(test *testing.T, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
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
	flags := sz.SZ_WITHOUT_INFO
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
	localLogger, err = logging.NewSenzingSdkLogger(ComponentId, szdiagnosticapi.IdMessages, options...)
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

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSzDiagnostic_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
}

func TestSzDiagnostic_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szDiagnostic.SetObserverOrigin(ctx, origin)
	actual := szDiagnostic.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzDiagnostic_CheckDatabasePerformance(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	secondsToRun := 1
	actual, err := szDiagnostic.CheckDatabasePerformance(ctx, secondsToRun)
	testError(test, err)
	printActual(test, actual)
}

func TestSzDiagnostic_Initialize(test *testing.T) {
	ctx := context.TODO()
	grpcConnection := getGrpcConnection()
	szDiagnostic := &SzDiagnostic{
		GrpcClient: szpb.NewSzDiagnosticClient(grpcConnection),
	}
	instanceName := "Test name"
	settings := "{}"
	verboseLogging := sz.SZ_NO_LOGGING
	configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
	err := szDiagnostic.Initialize(ctx, instanceName, settings, verboseLogging, configId)
	expectError(test, err, "senzing-60134002")
}

func TestSzDiagnostic_Reinit(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	szConfigManager := getSzConfigManager(ctx)
	configId, err := szConfigManager.GetDefaultConfigId(ctx)
	testError(test, err)
	err = szDiagnostic.Reinitialize(ctx, configId)
	testErrorNoFail(test, err)
}

func TestSzDiagnostic_Destroy(test *testing.T) {
	ctx := context.TODO()
	szDiagnostic := getTestObject(ctx, test)
	err := szDiagnostic.Destroy(ctx)
	expectError(test, err, "senzing-60134001")
	szDiagnosticSingleton = nil
}
