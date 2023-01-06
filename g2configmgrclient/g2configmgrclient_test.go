package g2configmgrclient

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/g2-sdk-go-grpc/g2configclient"
	"github.com/senzing/g2-sdk-go/g2config"
	"github.com/senzing/g2-sdk-go/g2configmgr"
	"github.com/senzing/g2-sdk-go/g2engine"
	"github.com/senzing/g2-sdk-go/testhelpers"
	pbg2config "github.com/senzing/g2-sdk-proto/go/g2config"
	pb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing/go-helpers/g2engineconfigurationjson"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTruncation = 76
)

var (
	grpcAddress                = "localhost:8258"
	grpcConnection             *grpc.ClientConn
	g2configClientSingleton    *g2configclient.G2configClient
	g2configmgrClientSingleton *G2configmgrClient
	localLogger                messagelogger.MessageLoggerInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getGrpcConnection() *grpc.ClientConn {
	var err error
	if grpcConnection == nil {
		grpcConnection, err = grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Did not connect: %v\n", err)
		}
		//		defer grpcConnection.Close()
	}
	return grpcConnection
}

func getTestObject(ctx context.Context, test *testing.T) G2configmgrClient {
	if g2configmgrClientSingleton == nil {

		grpcConnection := getGrpcConnection()
		g2configmgrClientSingleton = &G2configmgrClient{
			GrpcClient: pb.NewG2ConfigMgrClient(grpcConnection),
		}

		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if jsonErr != nil {
			logger.Fatalf("Cannot construct system configuration: %v", jsonErr)
		}

		initErr := g2configmgrClientSingleton.Init(ctx, moduleName, iniParams, verboseLogging)
		if initErr != nil {
			logger.Fatalf("Cannot Init: %v", initErr)
		}
	}
	return *g2configmgrClientSingleton
}

func getG2Configmgr(ctx context.Context) G2configmgrClient {
	grpcConnection := getGrpcConnection()
	g2configmgr := &G2configmgrClient{
		GrpcClient: pb.NewG2ConfigMgrClient(grpcConnection),
	}
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	g2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	return *g2configmgr
}

func getG2Config(ctx context.Context, test *testing.T) g2configclient.G2configClient {

	if g2configClientSingleton == nil {

		grpcConnection := getGrpcConnection()
		g2configClientSingleton = &g2configclient.G2configClient{
			GrpcClient: pbg2config.NewG2ConfigClient(grpcConnection),
		}

		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if jsonErr != nil {
			logger.Fatalf("Cannot construct system configuration: %v", jsonErr)
		}

		initErr := g2configmgrClientSingleton.Init(ctx, moduleName, iniParams, verboseLogging)
		if initErr != nil {
			logger.Fatalf("Cannot Init: %v", initErr)
		}
	}
	return *g2configClientSingleton
}

func truncate(aString string, length int) string {
	return truncator.Truncate(aString, length, "...", truncator.PositionEnd)
}

func printResult(test *testing.T, title string, result interface{}) {
	if 1 == 0 {
		test.Logf("%s: %v", title, truncate(fmt.Sprintf("%v", result), defaultTruncation))
	}
}

func printActual(test *testing.T, actual interface{}) {
	printResult(test, "Actual", actual)
}

func testError(test *testing.T, ctx context.Context, g2configmgr G2configmgrClient, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorNoFail(test *testing.T, ctx context.Context, g2configmgr G2configmgrClient, err error) {
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

func setupSenzingConfig(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	now := time.Now()

	aG2config := &g2config.G2configImpl{}
	err := aG2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5906, err)
	}

	configHandle, err := aG2config.Create(ctx)
	if err != nil {
		return localLogger.Error(5907, err)
	}

	for _, testDataSource := range testhelpers.TestDataSources {
		_, err := aG2config.AddDataSource(ctx, configHandle, testDataSource.Data)
		if err != nil {
			return localLogger.Error(5908, err)
		}
	}

	configStr, err := aG2config.Save(ctx, configHandle)
	if err != nil {
		return localLogger.Error(5909, err)
	}

	err = aG2config.Close(ctx, configHandle)
	if err != nil {
		return localLogger.Error(5910, err)
	}

	err = aG2config.Destroy(ctx)
	if err != nil {
		return localLogger.Error(5911, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	aG2configmgr := &g2configmgr.G2configmgrImpl{}
	err = aG2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5912, err)
	}

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := aG2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return localLogger.Error(5913, err)
	}

	err = aG2configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return localLogger.Error(5914, err)
	}

	err = aG2configmgr.Destroy(ctx)
	if err != nil {
		return localLogger.Error(5915, err)
	}
	return err
}

func setupPurgeRepository(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	aG2engine := &g2engine.G2engineImpl{}
	err := aG2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5903, err)
	}

	err = aG2engine.PurgeRepository(ctx)
	if err != nil {
		return localLogger.Error(5904, err)
	}

	err = aG2engine.Destroy(ctx)
	if err != nil {
		return localLogger.Error(5905, err)
	}
	return err
}

func setup() error {
	ctx := context.TODO()

	moduleName := "Test module name"
	verboseLogging := 0

	localLogger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	// if err != nil {
	// 	return logger.Error(5901, err)
	// }

	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		return localLogger.Error(5902, err)
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfig(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5920, err)
	}

	// Purge repository.

	err = setupPurgeRepository(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		return localLogger.Error(5921, err)
	}

	return err
}

func teardown() error {
	var err error = nil
	return err
}

func TestBuildSimpleSystemConfigurationJson(test *testing.T) {
	actual, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, actual)
	}
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Test interface functions - names begin with "Test"
// ----------------------------------------------------------------------------

func TestG2configmgrImpl_AddConfig(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	now := time.Now()

	// Create an in-memory configuration.

	g2config := getG2Config(ctx, test)
	configHandle, err1 := g2config.Create(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "g2config.Create()")
	}

	// Modify the in-memory configuration so it is different from the created configuration.
	// If not, on Save Senzing will detect that it is the same and no Save occurs.

	inputJson := `{"DSRC_CODE": "GO_TEST_` + strconv.FormatInt(now.Unix(), 10) + `"}`
	_, err2 := g2config.AddDataSource(ctx, configHandle, inputJson)
	if err2 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, "g2config.AddDataSource()")
	}

	// Create a JSON string from the in-memory version of the configuration.

	configStr, err3 := g2config.Save(ctx, configHandle)
	if err3 != nil {
		test.Log("Error:", err2.Error())
		assert.FailNow(test, configStr)
	}

	// Perform the test.

	configComments := fmt.Sprintf("g2configmgr_test at %s", now.UTC())
	actual, err := g2configmgr.AddConfig(ctx, configStr, configComments)
	testError(test, ctx, g2configmgr, err)
	printActual(test, actual)
}

func TestG2configmgrImpl_GetConfig(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)

	// Get a ConfigID.

	configID, err1 := g2configmgr.GetDefaultConfigID(ctx)
	if err1 != nil {
		test.Log("Error:", err1.Error())
		assert.FailNow(test, "g2configmgr.GetDefaultConfigID()")
	}

	actual, err := g2configmgr.GetConfig(ctx, configID)
	testError(test, ctx, g2configmgr, err)
	printActual(test, actual)
}

func TestG2configmgrImpl_GetConfigList(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	actual, err := g2configmgr.GetConfigList(ctx)
	testError(test, ctx, g2configmgr, err)
	printActual(test, actual)
}

func TestG2configmgrImpl_GetDefaultConfigID(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	actual, err := g2configmgr.GetDefaultConfigID(ctx)
	testError(test, ctx, g2configmgr, err)
	printActual(test, actual)
}

func TestG2configmgrImpl_Init(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if jsonErr != nil {
		test.Fatalf("Cannot construct system configuration: %v", jsonErr)
	}
	err := g2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	testError(test, ctx, g2configmgr, err)
}

func TestG2configmgrImpl_ReplaceDefaultConfigID(test *testing.T) {
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

func TestG2configmgrImpl_SetDefaultConfigID(test *testing.T) {
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

func TestG2configmgrImpl_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2configmgr := getTestObject(ctx, test)
	err := g2configmgr.Destroy(ctx)
	testError(test, ctx, g2configmgr, err)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2configmgrImpl_AddConfig() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2configmgr/g2configmgr_test.go

	// Create an in-memory configuration.
	ctx := context.TODO()
	g2config := &g2config.G2configImpl{}
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

func ExampleG2configmgrImpl_Destroy() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	err := g2configmgr.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configmgrImpl_GetConfig() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2configmgr/g2configmgr_test.go
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

func ExampleG2configmgrImpl_GetConfigList() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	jsonConfigList, err := g2configmgr.GetConfigList(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(jsonConfigList, 28))
	// Output: {"CONFIGS":[{"CONFIG_ID":...
}

func ExampleG2configmgrImpl_GetDefaultConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	configID, err := g2configmgr.GetDefaultConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleG2configmgrImpl_Init() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2configmgr/g2configmgr_test.go
	grpcConnection := getGrpcConnection()
	g2configmgr := &G2configmgrClient{
		GrpcClient: pb.NewG2ConfigMgrClient(grpcConnection),
	}
	ctx := context.TODO()
	moduleName := "Test module name"
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("") // See https://pkg.go.dev/github.com/senzing/go-helpers
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := 0
	err = g2configmgr.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configmgrImpl_ReplaceDefaultConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	oldConfigID, err := g2configmgr.GetDefaultConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Create an example configuration.
	g2config := &g2config.G2configImpl{}
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

func ExampleG2configmgrImpl_SetDefaultConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2configmgr/g2configmgr_test.go
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

func ExampleG2configmgrImpl_SetLogLevel() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go/blob/main/g2configmgr/g2configmgr_test.go
	ctx := context.TODO()
	g2configmgr := getG2Configmgr(ctx)
	err := g2configmgr.SetLogLevel(ctx, logger.LevelInfo)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
