package g2configmgr

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
	"github.com/senzing/g2-sdk-go-grpc/g2config"
	"github.com/senzing/g2-sdk-go-grpc/g2engine"
	"github.com/senzing/g2-sdk-go/g2api"
	g2configmgrapi "github.com/senzing/g2-sdk-go/g2configmgr"
	"github.com/senzing/g2-sdk-go/g2error"
	g2configpb "github.com/senzing/g2-sdk-proto/go/g2config"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	g2enginepb "github.com/senzing/g2-sdk-proto/go/g2engine"
	"github.com/senzing/go-common/truthset"
	"github.com/senzing/go-logging/logging"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	g2configSingleton    g2api.G2config
	g2configmgrSingleton g2api.G2configmgr
	g2engineSingleton    g2api.G2engine
	grpcAddress          = "localhost:8258"
	grpcConnection       *grpc.ClientConn
	localLogger          logging.LoggingInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func createError(errorId int, err error) error {
	return g2error.Cast(localLogger.Error(errorId, err), err)
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

func getTestObject(ctx context.Context, test *testing.T) g2api.G2configmgr {
	if g2configmgrSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2configmgrSingleton = &G2configmgr{
			GrpcClient: g2pb.NewG2ConfigMgrClient(grpcConnection),
		}
	}
	return g2configmgrSingleton
}

func getG2Config(ctx context.Context) g2api.G2config {
	if g2configSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2configSingleton = &g2config.G2config{
			GrpcClient: g2configpb.NewG2ConfigClient(grpcConnection),
		}
	}
	return g2configSingleton
}

func getG2Configmgr(ctx context.Context) g2api.G2configmgr {
	if g2configmgrSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2configmgrSingleton = &G2configmgr{
			GrpcClient: g2pb.NewG2ConfigMgrClient(grpcConnection),
		}
	}
	return g2configmgrSingleton
}

func getG2Engine(ctx context.Context) g2api.G2engine {
	if g2engineSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2engineSingleton = &g2engine.G2engine{
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

func testError(test *testing.T, ctx context.Context, g2configmgr g2api.G2configmgr, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, g2configmgr g2api.G2configmgr, err error, messageId string) {
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
		if g2error.Is(err, g2error.G2Unrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if g2error.Is(err, g2error.G2Retryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if g2error.Is(err, g2error.G2BadUserInput) {
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

	return err
}

func setupPurgeRepository(ctx context.Context) error {
	g2engine := getG2Engine(ctx)
	err := g2engine.PurgeRepository(ctx)
	return err
}

func setup() error {
	ctx := context.TODO()
	var err error = nil

	options := []interface{}{
		&logging.OptionCallerSkip{Value: 4},
	}
	localLogger, err = logging.NewSenzingSdkLogger(ComponentId, g2configmgrapi.IdMessages, options...)
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

	return err
}

func teardown() error {
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2configmgr_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2configmgr.SetObserverOrigin(ctx, origin)
}

func TestG2configmgr_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	g2configmgr.SetObserverOrigin(ctx, origin)
	actual := g2configmgr.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestG2configmgr_AddConfig(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	now := time.Now()
	g2config := getG2Config(ctx)
	configHandle, err1 := g2config.Create(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "g2config.Create()")
	}
	inputJson := `{"DSRC_CODE": "GO_TEST_` + strconv.FormatInt(now.Unix(), 10) + `"}`
	_, err2 := g2config.AddDataSource(ctx, configHandle, inputJson)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "g2config.AddDataSource()")
	}
	configStr, err3 := g2config.Save(ctx, configHandle)
	if err3 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, configStr)
	}
	configComments := fmt.Sprintf("g2configmgr_test at %s", now.UTC())
	actual, err := g2configmgr.AddConfig(ctx, configStr, configComments)
	testError(test, ctx, g2configmgr, err)
	printActual(test, actual)
}

func TestG2configmgr_GetConfig(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	configID, err1 := g2configmgr.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "g2configmgr.GetDefaultConfigID()")
	}
	actual, err := g2configmgr.GetConfig(ctx, configID)
	testError(test, ctx, g2configmgr, err)
	printActual(test, actual)
}

func TestG2configmgr_GetConfigList(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	actual, err := g2configmgr.GetConfigList(ctx)
	testError(test, ctx, g2configmgr, err)
	printActual(test, actual)
}

func TestG2configmgr_GetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	actual, err := g2configmgr.GetDefaultConfigID(ctx)
	testError(test, ctx, g2configmgr, err)
	printActual(test, actual)
}

func TestG2configmgr_ReplaceDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	oldConfigID, err1 := g2configmgr.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "g2configmgr.GetDefaultConfigID()")
	}

	// FIXME: This is kind of a cheeter.

	newConfigID, err2 := g2configmgr.GetDefaultConfigID(ctx)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "g2configmgr.GetDefaultConfigID()-2")
	}

	err := g2configmgr.ReplaceDefaultConfigID(ctx, oldConfigID, newConfigID)
	testError(test, ctx, g2configmgr, err)
}

func TestG2configmgr_SetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	configID, err1 := g2configmgr.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "g2configmgr.GetDefaultConfigID()")
	}
	err := g2configmgr.SetDefaultConfigID(ctx, configID)
	testError(test, ctx, g2configmgr, err)
}

func TestG2configmgr_Init(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams := "{}"
	err := g2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	expectError(test, ctx, g2configmgr, err, "senzing-60124002")
}

func TestG2configmgr_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	err := g2configmgr.Destroy(ctx)
	expectError(test, ctx, g2configmgr, err, "senzing-60124001")
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2configmgr_SetObserverOrigin() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2configmgr.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleG2configmgr_GetObserverOrigin() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2config/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2configmgr.SetObserverOrigin(ctx, origin)
	result := g2configmgr.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleG2configmgr_AddConfig() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	g2configmgr := getG2Configmgr(ctx)
	configStr, err := g2config.Save(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	configComments := "Example configuration"
	configID, err := g2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleG2configmgr_GetConfig() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	configID, err := g2configmgr.GetDefaultConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configStr, err := g2configmgr.GetConfig(ctx, configID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(configStr, defaultTruncation))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR...
}

func ExampleG2configmgr_GetConfigList() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	jsonConfigList, err := g2configmgr.GetConfigList(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(jsonConfigList, 28))
	// Output: {"CONFIGS":[{"CONFIG_ID":...
}

func ExampleG2configmgr_GetDefaultConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	configID, err := g2configmgr.GetDefaultConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleG2configmgr_ReplaceDefaultConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	oldConfigID, err := g2configmgr.GetDefaultConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configStr, err := g2config.Save(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	configComments := "Example configuration"
	newConfigID, err := g2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		fmt.Println(err)
	}
	err = g2configmgr.ReplaceDefaultConfigID(ctx, oldConfigID, newConfigID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configmgr_SetDefaultConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	configID, err := g2configmgr.GetDefaultConfigID(ctx) // For example purposes only. Normally would use output from GetConfigList()
	if err != nil {
		fmt.Println(err)
	}
	err = g2configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configmgr_SetLogLevel() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	err := g2configmgr.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configmgr_Init() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	grpcConnection := getGrpcConnection()
	g2configmgr := &G2configmgr{
		GrpcClient: g2pb.NewG2ConfigMgrClient(grpcConnection),
	}
	moduleName := "Test module name"
	iniParams := "{}"
	verboseLogging := 0
	err := g2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		// This should produce a "senzing-60124002" error.
	}
	// Output:
}

func ExampleG2configmgr_Destroy() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	err := g2configmgr.Destroy(ctx)
	if err != nil {
		// This should produce a "senzing-60124001" error.
	}
	// Output:
}
