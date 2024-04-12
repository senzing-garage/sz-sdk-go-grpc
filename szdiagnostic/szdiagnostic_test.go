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
	"github.com/senzing-garage/g2-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/g2-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/g2-sdk-go-grpc/szengine"
	"github.com/senzing-garage/g2-sdk-go/g2api"
	g2diagnosticapi "github.com/senzing-garage/g2-sdk-go/g2diagnostic"
	"github.com/senzing-garage/g2-sdk-go/g2error"
	g2configpb "github.com/senzing-garage/g2-sdk-proto/go/g2config"
	g2configmgrpb "github.com/senzing-garage/g2-sdk-proto/go/g2configmgr"
	g2pb "github.com/senzing-garage/g2-sdk-proto/go/g2diagnostic"
	g2enginepb "github.com/senzing-garage/g2-sdk-proto/go/g2engine"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	g2configSingleton     g2api.G2config
	g2configmgrSingleton  g2api.G2configmgr
	g2diagnosticSingleton g2api.G2diagnostic
	g2engineSingleton     g2api.G2engine
	grpcAddress           = "localhost:8261"
	grpcConnection        *grpc.ClientConn
	localLogger           logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return g2error.Cast(localLogger.NewError(errorId, err), err)
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

func getTestObject(ctx context.Context, test *testing.T) g2api.G2diagnostic {
	if g2diagnosticSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2diagnosticSingleton = &G2diagnostic{
			GrpcClient: g2pb.NewG2DiagnosticClient(grpcConnection),
		}
	}
	return g2diagnosticSingleton
}

func getG2Config(ctx context.Context) g2api.G2config {
	if g2configSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2configSingleton = &szconfig.SzConfig{
			GrpcClient: g2configpb.NewG2ConfigClient(grpcConnection),
		}
	}
	return g2configSingleton
}

func getG2Configmgr(ctx context.Context) g2api.G2configmgr {
	if g2configmgrSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2configmgrSingleton = &szconfigmanager.G2configmgr{
			GrpcClient: g2configmgrpb.NewG2ConfigMgrClient(grpcConnection),
		}
	}
	return g2configmgrSingleton
}

func getG2Diagnostic(ctx context.Context) g2api.G2diagnostic {
	if g2diagnosticSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2diagnosticSingleton = &G2diagnostic{
			GrpcClient: g2pb.NewG2DiagnosticClient(grpcConnection),
		}
	}
	return g2diagnosticSingleton
}

func getG2Engine(ctx context.Context) g2api.G2engine {
	if g2engineSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2engineSingleton = &szengine.G2engine{
			GrpcClient: g2enginepb.NewG2EngineClient(grpcConnection),
		}
	}
	return g2engineSingleton
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

func testError(test *testing.T, ctx context.Context, g2diagnostic g2api.G2diagnostic, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, g2diagnostic g2api.G2diagnostic, err error, messageId string) {
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

func testErrorNoFail(test *testing.T, ctx context.Context, g2diagnostic g2api.G2diagnostic, err error) {
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
		if g2error.Is(err, g2error.G2Unrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if g2error.Is(err, g2error.G2Retryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if g2error.Is(err, g2error.G2BadInput) {
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

	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		return createError(5907, err)
	}

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, datasourceName := range datasourceNames {
		datasource := truthset.TruthsetDataSources[datasourceName]
		_, err := g2config.AddDataSource(ctx, configHandle, datasource.Json)
		if err != nil {
			return createError(5908, err)
		}
	}

	configStr, err := g2config.Save(ctx, configHandle)
	if err != nil {
		return createError(5909, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	g2configmgr := getG2Configmgr(ctx)
	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := g2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return createError(5913, err)
	}

	err = g2configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return createError(5914, err)
	}

	g2diagnostic := getG2Diagnostic(ctx)
	err = g2diagnostic.Reinit(ctx, configID)

	return err
}

func setupPurgeRepository(ctx context.Context) error {
	g2diagnostic := getG2Diagnostic(ctx)
	err := g2diagnostic.PurgeRepository(ctx)
	return err
}

func setupAddRecords(ctx context.Context) error {
	var err error = nil
	g2engine := getG2Engine(ctx)
	testRecordIds := []string{"1001", "1002", "1003", "1004", "1005", "1039", "1040"}
	for _, testRecordId := range testRecordIds {
		testRecord := truthset.CustomerRecords[testRecordId]
		err := g2engine.AddRecord(ctx, testRecord.DataSource, testRecord.Id, testRecord.Json, "G2Diagnostic_test")
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
	localLogger, err = logging.NewSenzingSdkLogger(ComponentId, g2diagnosticapi.IdMessages, options...)
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

func TestG2diagnostic_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2diagnostic.SetObserverOrigin(ctx, origin)
}

func TestG2diagnostic_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2diagnostic.SetObserverOrigin(ctx, origin)
	actual := g2diagnostic.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestG2diagnostic_CheckDBPerf(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	secondsToRun := 1
	actual, err := g2diagnostic.CheckDBPerf(ctx, secondsToRun)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnostic_Init(test *testing.T) {
	ctx := context.TODO()
	grpcConnection := getGrpcConnection()
	g2diagnostic := &G2diagnostic{
		GrpcClient: g2pb.NewG2DiagnosticClient(grpcConnection),
	}
	moduleName := "Test module name"
	iniParams := "{}"
	verboseLogging := int64(0)
	err := g2diagnostic.Init(ctx, moduleName, iniParams, verboseLogging)
	expectError(test, ctx, g2diagnostic, err, "senzing-60134002")
}

func TestG2diagnostic_InitWithConfigID(test *testing.T) {
	ctx := context.TODO()
	grpcConnection := getGrpcConnection()
	g2diagnostic := &G2diagnostic{
		GrpcClient: g2pb.NewG2DiagnosticClient(grpcConnection),
	}
	moduleName := "Test module name"
	initConfigID := int64(1)
	iniParams := "{}"
	verboseLogging := int64(0)
	err := g2diagnostic.InitWithConfigID(ctx, moduleName, iniParams, initConfigID, verboseLogging)
	expectError(test, ctx, g2diagnostic, err, "senzing-60134003")
}

func TestG2diagnostic_Reinit(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	g2Configmgr := getG2Configmgr(ctx)
	initConfigID, err := g2Configmgr.GetDefaultConfigID(ctx)
	testError(test, ctx, g2diagnostic, err)
	err = g2diagnostic.Reinit(ctx, initConfigID)
	testErrorNoFail(test, ctx, g2diagnostic, err)
}

func TestG2diagnostic_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	err := g2diagnostic.Destroy(ctx)
	expectError(test, ctx, g2diagnostic, err, "senzing-60134001")
	g2diagnosticSingleton = nil
}
