package g2configclient

import (
	"context"
	"fmt"
	"testing"

	truncator "github.com/aquilax/truncate"
	pb "github.com/senzing/g2-sdk-proto/go/g2config"
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
	grpcAddress             = "localhost:8258"
	grpcConnection          *grpc.ClientConn
	g2configClientSingleton *G2configClient
	localLogger             messagelogger.MessageLoggerInterface
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
func getTestObject(ctx context.Context, test *testing.T) G2configClient {
	if g2configClientSingleton == nil {

		grpcConnection := getGrpcConnection()
		g2configClientSingleton = &G2configClient{
			G2ConfigGrpcClient: pb.NewG2ConfigClient(grpcConnection),
		}

		moduleName := "Test module name"
		verboseLogging := 0
		iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if jsonErr != nil {
			logger.Fatalf("Cannot construct system configuration: %v", jsonErr)
		}

		initErr := g2configClientSingleton.Init(ctx, moduleName, iniParams, verboseLogging)
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

func testError(test *testing.T, ctx context.Context, g2config G2configClient, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func testErrorNoFail(test *testing.T, ctx context.Context, g2config G2configClient, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
	}
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestBuildSimpleSystemConfigurationJson(test *testing.T) {
	actual, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, actual)
	}
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2configClient_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	inputJson := `{"DSRC_CODE": "GO_TEST"}`
	actual, err := g2config.AddDataSource(ctx, configHandle, inputJson)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
}

func TestG2configClient_Close(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
}

func TestG2configClient_Create(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	actual, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
}

func TestG2configClient_DeleteDataSource(test *testing.T) {

	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)

	actual, err := g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printResult(test, "Original", actual)

	inputJson := `{"DSRC_CODE": "GO_TEST"}`
	_, err = g2config.AddDataSource(ctx, configHandle, inputJson)
	testError(test, ctx, g2config, err)

	actual, err = g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printResult(test, "     Add", actual)

	err = g2config.DeleteDataSource(ctx, configHandle, inputJson)
	testError(test, ctx, g2config, err)

	actual, err = g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printResult(test, "  Delete", actual)

	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
}

func TestG2configClient_Init(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	testError(test, ctx, g2config, jsonErr)
	err := g2config.Init(ctx, moduleName, iniParams, verboseLogging)
	testError(test, ctx, g2config, err)
}

func TestG2configClient_ListDataSources(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	actual, err := g2config.ListDataSources(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
	err = g2config.Close(ctx, configHandle)
	testError(test, ctx, g2config, err)
}

func TestG2configClient_Load(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	jsonConfig, err := g2config.Save(ctx, configHandle)
	testError(test, ctx, g2config, err)
	err = g2config.Load(ctx, configHandle, jsonConfig)
	testError(test, ctx, g2config, err)
}

func TestG2configClient_Save(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	configHandle, err := g2config.Create(ctx)
	testError(test, ctx, g2config, err)
	actual, err := g2config.Save(ctx, configHandle)
	testError(test, ctx, g2config, err)
	printActual(test, actual)
}

func TestG2configClient_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2config := getTestObject(ctx, test)
	err := g2config.Destroy(ctx)
	testError(test, ctx, g2config, err)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2configClient_AddDataSource() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configclient/g2configclient_test.go
	grpcConnection := getGrpcConnection()
	g2config := &G2configClient{
		G2ConfigGrpcClient: pb.NewG2ConfigClient(grpcConnection),
	}
	ctx := context.TODO()
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	inputJson := `{"DSRC_CODE": "GO_TEST"}`
	result, err := g2config.AddDataSource(ctx, configHandle, inputJson)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DSRC_ID":1001}
}

func ExampleG2configClient_Close() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configclient/g2configclient_test.go
	grpcConnection := getGrpcConnection()
	g2config := &G2configClient{
		G2ConfigGrpcClient: pb.NewG2ConfigClient(grpcConnection),
	}
	ctx := context.TODO()
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = g2config.Close(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configClient_Create() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configclient/g2configclient_test.go
	grpcConnection := getGrpcConnection()
	g2config := &G2configClient{
		G2ConfigGrpcClient: pb.NewG2ConfigClient(grpcConnection),
	}
	ctx := context.TODO()
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configHandle > 0) // Dummy output.
	// Output: true
}

func ExampleG2configClient_DeleteDataSource() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configclient/g2configclient_test.go
	grpcConnection := getGrpcConnection()
	g2config := &G2configClient{
		G2ConfigGrpcClient: pb.NewG2ConfigClient(grpcConnection),
	}
	ctx := context.TODO()
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	inputJson := `{"DSRC_CODE": "TEST"}`
	err = g2config.DeleteDataSource(ctx, configHandle, inputJson)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configClient_Destroy() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configclient/g2configclient_test.go
	grpcConnection := getGrpcConnection()
	g2config := &G2configClient{
		G2ConfigGrpcClient: pb.NewG2ConfigClient(grpcConnection),
	}
	ctx := context.TODO()
	err := g2config.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configClient_Init() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configclient/g2configclient_test.go
	grpcConnection := getGrpcConnection()
	g2config := &G2configClient{
		G2ConfigGrpcClient: pb.NewG2ConfigClient(grpcConnection),
	}
	ctx := context.TODO()
	moduleName := "Test module name"
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := 0
	err = g2config.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configClient_ListDataSources() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configclient/g2configclient_test.go
	grpcConnection := getGrpcConnection()
	g2config := &G2configClient{
		G2ConfigGrpcClient: pb.NewG2ConfigClient(grpcConnection),
	}
	ctx := context.TODO()
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	result, err := g2config.ListDataSources(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCES":[{"DSRC_ID":1,"DSRC_CODE":"TEST"},{"DSRC_ID":2,"DSRC_CODE":"SEARCH"}]}
}

func ExampleG2configClient_Load() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configclient/g2configclient_test.go
	grpcConnection := getGrpcConnection()
	g2config := &G2configClient{
		G2ConfigGrpcClient: pb.NewG2ConfigClient(grpcConnection),
	}
	ctx := context.TODO()
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	jsonConfig, err := g2config.Save(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	err = g2config.Load(ctx, configHandle, jsonConfig)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2configClient_Save() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configclient/g2configclient_test.go
	grpcConnection := getGrpcConnection()
	g2config := &G2configClient{
		G2ConfigGrpcClient: pb.NewG2ConfigClient(grpcConnection),
	}
	ctx := context.TODO()
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	jsonConfig, err := g2config.Save(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(jsonConfig, 207))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR_CLASS":"OBSERVATION","FTYPE_CODE":null,"FELEM_CODE":null,"FELEM_REQ":"Yes","DEFAULT_VALUE":null,"ADVANCED":"Yes","INTERNAL":"No"},...
}

func ExampleG2configClient_SetLogLevel() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2configclient/g2configclient_test.go
	grpcConnection := getGrpcConnection()
	g2config := &G2configClient{
		G2ConfigGrpcClient: pb.NewG2ConfigClient(grpcConnection),
	}
	ctx := context.TODO()
	err := g2config.SetLogLevel(ctx, logger.LevelInfo)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
