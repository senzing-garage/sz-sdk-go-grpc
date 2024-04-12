package szconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	truncator "github.com/aquilax/truncate"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	szConfigSingleton *SzConfig
	grpcAddress       = "localhost:8261"
	grpcConnection    *grpc.ClientConn
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

func getTestObject(ctx context.Context, test *testing.T) *SzConfig {
	if szConfigSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigSingleton = &SzConfig{
			GrpcClient: szpb.NewSzConfigClient(grpcConnection),
		}
	}
	return szConfigSingleton
}

func getSzConfig(ctx context.Context) *SzConfig {
	if szConfigSingleton == nil {
		grpcConnection := getGrpcConnection()
		szConfigSingleton = &SzConfig{
			GrpcClient: szpb.NewSzConfigClient(grpcConnection),
		}
	}
	return szConfigSingleton
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

func setup() error {
	var err error = nil
	return err
}

func teardown() error {
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestSzConfig_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
}

func TestSzConfig_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	actual := szConfig.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzConfig_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, err)
	dataSourceCode := "GO_TEST"
	actual, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	testError(test, ctx, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	testError(test, ctx, err)
}

func TestSzConfig_AddDataSource_withLoad(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, err)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	testError(test, ctx, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	testError(test, ctx, err)
	configHandle2, err := szConfig.ImportConfig(ctx, configDefinition)
	testError(test, ctx, err)
	inputJson := "GO_TEST"
	actual, err := szConfig.AddDataSource(ctx, configHandle2, inputJson)
	testError(test, ctx, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle2)
	testError(test, ctx, err)
}

func TestSzConfig_Close(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	testError(test, ctx, err)
}

func TestSzConfig_Create(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, err)
	printActual(test, actual)
}

func TestSzConfig_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, err)
	printResult(test, "Original", actual)
	dataSourceCode := "GO_TEST"
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	testError(test, ctx, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, err)
	printResult(test, "     Add", actual)
	err = szConfig.DeleteDataSource(ctx, configHandle, dataSourceCode)
	testError(test, ctx, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, err)
	printResult(test, "  Delete", actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	testError(test, ctx, err)
}

func TestSzConfig_DeleteDataSource_withLoad(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, err)
	printResult(test, "Original", actual)
	dataSourceCode := "GO_TEST"
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	testError(test, ctx, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, err)
	printResult(test, "     Add", actual)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	testError(test, ctx, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	testError(test, ctx, err)
	configHandle2, err := szConfig.ImportConfig(ctx, configDefinition)
	testError(test, ctx, err)
	err = szConfig.DeleteDataSource(ctx, configHandle2, dataSourceCode)
	testError(test, ctx, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle2)
	testError(test, ctx, err)
	printResult(test, "  Delete", actual)
	err = szConfig.CloseConfig(ctx, configHandle2)
	testError(test, ctx, err)
}

func TestSzConfig_GetDataSources(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	testError(test, ctx, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	testError(test, ctx, err)
}

func TestSzConfig_ImportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, err)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	testError(test, ctx, err)
	actual, err := szConfig.ImportConfig(ctx, configDefinition)
	testError(test, ctx, err)
	printActual(test, actual)
}

func TestSzConfig_ExportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	testError(test, ctx, err)
	actual, err := szConfig.ExportConfig(ctx, configHandle)
	testError(test, ctx, err)
	printActual(test, actual)
}

func TestSzConfig_Initialize(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	instanceName := "Test module name"
	verboseLogging := int64(0)
	settings := "{}"
	err := szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	expectError(test, ctx, err, "senzing-60114002")
}

func TestSzConfig_Destroy(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.Destroy(ctx)
	expectError(test, ctx, err, "senzing-60114001")
}
