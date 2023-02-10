package g2diagnosticclient

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing/g2-sdk-go-grpc/g2configclient"
	"github.com/senzing/g2-sdk-go-grpc/g2configmgrclient"
	"github.com/senzing/g2-sdk-go-grpc/g2engineclient"
	pbg2config "github.com/senzing/g2-sdk-proto/go/g2config"
	pbg2configmgr "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	pb "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
	pbg2engine "github.com/senzing/g2-sdk-proto/go/g2engine"
	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-common/truthset"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTruncation = 76
	printResults      = false
)

var (
	g2configClientSingleton     *g2configclient.G2configClient
	g2configmgrClientSingleton  *g2configmgrclient.G2configmgrClient
	g2diagnosticClientSingleton *G2diagnosticClient
	g2engineClientSingleton     *g2engineclient.G2engineClient
	grpcAddress                 = "localhost:8258"
	grpcConnection              *grpc.ClientConn
	localLogger                 messagelogger.MessageLoggerInterface
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

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

func getTestObject(ctx context.Context, test *testing.T) G2diagnosticClient {
	if g2diagnosticClientSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2diagnosticClientSingleton = &G2diagnosticClient{
			GrpcClient: pb.NewG2DiagnosticClient(grpcConnection),
		}
	}
	return *g2diagnosticClientSingleton
}

func getG2Config(ctx context.Context) g2configclient.G2configClient {
	if g2configClientSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2configClientSingleton = &g2configclient.G2configClient{
			GrpcClient: pbg2config.NewG2ConfigClient(grpcConnection),
		}
	}
	return *g2configClientSingleton
}

func getG2Configmgr(ctx context.Context) g2configmgrclient.G2configmgrClient {
	if g2configmgrClientSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2configmgrClientSingleton = &g2configmgrclient.G2configmgrClient{
			GrpcClient: pbg2configmgr.NewG2ConfigMgrClient(grpcConnection),
		}
	}
	return *g2configmgrClientSingleton
}

func getG2Diagnostic(ctx context.Context) G2diagnosticClient {
	if g2diagnosticClientSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2diagnosticClientSingleton = &G2diagnosticClient{
			GrpcClient: pb.NewG2DiagnosticClient(grpcConnection),
		}
	}
	return *g2diagnosticClientSingleton
}

func getG2Engine(ctx context.Context) g2engineclient.G2engineClient {
	if g2engineClientSingleton == nil {
		grpcConnection := getGrpcConnection()
		g2engineClientSingleton = &g2engineclient.G2engineClient{
			GrpcClient: pbg2engine.NewG2EngineClient(grpcConnection),
		}
	}
	return *g2engineClientSingleton
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

func testError(test *testing.T, ctx context.Context, g2diagnostic G2diagnosticClient, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

func expectError(test *testing.T, ctx context.Context, g2diagnostic G2diagnosticClient, err error, messageId string) {
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

func testErrorNoFail(test *testing.T, ctx context.Context, g2diagnostic G2diagnosticClient, err error) {
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

func setupSenzingConfig(ctx context.Context) error {
	now := time.Now()

	// Create a fresh Senzing configuration.

	g2config := getG2Config(ctx)
	configHandle, err := g2config.Create(ctx)
	if err != nil {
		return localLogger.Error(5907, err)
	}

	datasourceNames := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, datasourceName := range datasourceNames {
		datasource := truthset.TruthsetDataSources[datasourceName]
		_, err := g2config.AddDataSource(ctx, configHandle, datasource.Json)
		if err != nil {
			return localLogger.Error(5908, err)
		}
	}

	configStr, err := g2config.Save(ctx, configHandle)
	if err != nil {
		return localLogger.Error(5909, err)
	}

	// Persist the Senzing configuration to the Senzing repository.

	g2configmgr := getG2Configmgr(ctx)
	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := g2configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return localLogger.Error(5913, err)
	}

	err = g2configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return localLogger.Error(5914, err)
	}

	return err
}

func setupPurgeRepository(ctx context.Context) error {
	g2engine := getG2Engine(ctx)
	err := g2engine.PurgeRepository(ctx)
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
			return localLogger.Error(5917, err)
		}
	}
	return err
}

func setup() error {
	ctx := context.TODO()
	var err error = nil

	localLogger, err = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	if err != nil {
		return localLogger.Error(5901, err)
	}

	// Add Data Sources to Senzing configuration.

	err = setupSenzingConfig(ctx)
	if err != nil {
		return localLogger.Error(5920, err)
	}

	// Purge repository.

	err = setupPurgeRepository(ctx)
	if err != nil {
		return localLogger.Error(5921, err)
	}

	// Add records.

	err = setupAddRecords(ctx)
	if err != nil {
		return localLogger.Error(5922, err)
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
// Test interface functions
// ----------------------------------------------------------------------------

func TestG2diagnosticClient_CheckDBPerf(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	secondsToRun := 1
	actual, err := g2diagnostic.CheckDBPerf(ctx, secondsToRun)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticClient_EntityListBySize(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	aSize := 1000
	aHandle, err := g2diagnostic.GetEntityListBySize(ctx, aSize)
	testError(test, ctx, g2diagnostic, err)
	anEntity, err := g2diagnostic.FetchNextEntityBySize(ctx, aHandle)
	testError(test, ctx, g2diagnostic, err)
	printResult(test, "Entity", anEntity)
	err = g2diagnostic.CloseEntityListBySize(ctx, aHandle)
	testError(test, ctx, g2diagnostic, err)
}

func TestG2diagnosticClient_FindEntitiesByFeatureIDs(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	features := "{\"ENTITY_ID\":1,\"LIB_FEAT_IDS\":[1,3,4]}"
	actual, err := g2diagnostic.FindEntitiesByFeatureIDs(ctx, features)
	testError(test, ctx, g2diagnostic, err)
	printResult(test, "len(Actual)", len(actual))
}

func TestG2diagnosticClient_GetAvailableMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	actual, err := g2diagnostic.GetAvailableMemory(ctx)
	testError(test, ctx, g2diagnostic, err)
	assert.Greater(test, actual, int64(0))
	printActual(test, actual)
}

func TestG2diagnosticClient_GetDataSourceCounts(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	actual, err := g2diagnostic.GetDataSourceCounts(ctx)
	testError(test, ctx, g2diagnostic, err)
	printResult(test, "Data Source counts", actual)
}

func TestG2diagnosticClient_GetDBInfo(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	actual, err := g2diagnostic.GetDBInfo(ctx)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticClient_GetEntityDetails(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	entityID := int64(1)
	includeInternalFeatures := 1
	actual, err := g2diagnostic.GetEntityDetails(ctx, entityID, includeInternalFeatures)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticClient_GetEntityResume(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	entityID := int64(1)
	actual, err := g2diagnostic.GetEntityResume(ctx, entityID)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticClient_GetEntitySizeBreakdown(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	minimumEntitySize := 1
	includeInternalFeatures := 1
	actual, err := g2diagnostic.GetEntitySizeBreakdown(ctx, minimumEntitySize, includeInternalFeatures)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticClient_GetFeature(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	libFeatID := int64(1)
	actual, err := g2diagnostic.GetFeature(ctx, libFeatID)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticClient_GetGenericFeatures(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	featureType := "PHONE"
	maximumEstimatedCount := 10
	actual, err := g2diagnostic.GetGenericFeatures(ctx, featureType, maximumEstimatedCount)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticClient_GetLogicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	actual, err := g2diagnostic.GetLogicalCores(ctx)
	testError(test, ctx, g2diagnostic, err)
	assert.Greater(test, actual, 0)
	printActual(test, actual)
}

func TestG2diagnosticClient_GetMappingStatistics(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	includeInternalFeatures := 1
	actual, err := g2diagnostic.GetMappingStatistics(ctx, includeInternalFeatures)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticClient_GetPhysicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	actual, err := g2diagnostic.GetPhysicalCores(ctx)
	testError(test, ctx, g2diagnostic, err)
	assert.Greater(test, actual, 0)
	printActual(test, actual)
}

func TestG2diagnosticClient_GetRelationshipDetails(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	relationshipID := int64(1)
	includeInternalFeatures := 1
	actual, err := g2diagnostic.GetRelationshipDetails(ctx, relationshipID, includeInternalFeatures)
	testErrorNoFail(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticClient_GetResolutionStatistics(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	actual, err := g2diagnostic.GetResolutionStatistics(ctx)
	testError(test, ctx, g2diagnostic, err)
	printActual(test, actual)
}

func TestG2diagnosticClient_GetTotalSystemMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	actual, err := g2diagnostic.GetTotalSystemMemory(ctx)
	testError(test, ctx, g2diagnostic, err)
	assert.Greater(test, actual, int64(0))
	printActual(test, actual)
}

func TestG2diagnosticClient_Init(test *testing.T) {
	ctx := context.TODO()
	grpcConnection := getGrpcConnection()
	g2diagnostic := &G2diagnosticClient{
		GrpcClient: pb.NewG2DiagnosticClient(grpcConnection),
	}
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	testError(test, ctx, *g2diagnostic, err)
	err = g2diagnostic.Init(ctx, moduleName, iniParams, verboseLogging)
	expectError(test, ctx, *g2diagnostic, err, "senzing-60134002")
}

func TestG2diagnosticClient_InitWithConfigID(test *testing.T) {
	ctx := context.TODO()
	grpcConnection := getGrpcConnection()
	g2diagnostic := &G2diagnosticClient{
		GrpcClient: pb.NewG2DiagnosticClient(grpcConnection),
	}
	moduleName := "Test module name"
	initConfigID := int64(1)
	verboseLogging := 0
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	testError(test, ctx, *g2diagnostic, err)
	err = g2diagnostic.InitWithConfigID(ctx, moduleName, iniParams, initConfigID, verboseLogging)
	expectError(test, ctx, *g2diagnostic, err, "senzing-60134003")
}

func TestG2diagnosticClient_Reinit(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	g2Configmgr := getG2Configmgr(ctx)
	initConfigID, err := g2Configmgr.GetDefaultConfigID(ctx)
	testError(test, ctx, g2diagnostic, err)
	err = g2diagnostic.Reinit(ctx, initConfigID)
	testErrorNoFail(test, ctx, g2diagnostic, err)
}

func TestG2diagnosticClient_Destroy(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx, test)
	err := g2diagnostic.Destroy(ctx)
	expectError(test, ctx, g2diagnostic, err, "senzing-60134001")
	g2diagnosticClientSingleton = nil
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2diagnosticClient_CheckDBPerf() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	secondsToRun := 1
	result, err := g2diagnostic.CheckDBPerf(ctx, secondsToRun)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 25))
	// Output: {"numRecordsInserted":...
}

func ExampleG2diagnosticClient_CloseEntityListBySize() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	aSize := 1000
	entityListBySizeHandle, err := g2diagnostic.GetEntityListBySize(ctx, aSize)
	if err != nil {
		fmt.Println(err)
	}
	err = g2diagnostic.CloseEntityListBySize(ctx, entityListBySizeHandle)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

// func ExampleG2diagnosticClient_FetchNextEntityBySize() {
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
// 	ctx := context.TODO()
// 	g2diagnostic := getG2Diagnostic(ctx)
// 	aSize := 1
// 	entityListBySizeHandle, err := g2diagnostic.GetEntityListBySize(ctx, aSize)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	anEntity, _ := g2diagnostic.FetchNextEntityBySize(ctx, entityListBySizeHandle)
// 	g2diagnostic.CloseEntityListBySize(ctx, entityListBySizeHandle)
// 	fmt.Println(anEntity)
// 	// Output: [{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","ENT_SRC_KEY":"EF75DB9728B437EEAD00889C077A7043B364269C","ENT_SRC_DESC":"John Smith","RECORD_ID":"1039","JSON_DATA":"{\"RECORD_TYPE\":\"PERSON\",\"PRIMARY_NAME_LAST\":\"Smith\",\"PRIMARY_NAME_FIRST\":\"John\",\"GENDER\":\"M\",\"DATE_OF_BIRTH\":\"10/10/70\",\"ADDR_TYPE\":\"HOME\",\"ADDR_LINE1\":\"3212 W. 32nd St Palm Harbor, FL 60527\",\"DATE\":\"1/28/18\",\"STATUS\":\"Active\",\"AMOUNT\":\"900\",\"DATA_SOURCE\":\"CUSTOMERS\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"1039\"}","OBS_ENT_ID":6,"ER_ID":0}]
// }

func ExampleG2diagnosticClient_FindEntitiesByFeatureIDs() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	features := `{"ENTITY_ID":1,"LIB_FEAT_IDS":[1,3,4]}`
	result, err := g2diagnostic.FindEntitiesByFeatureIDs(ctx, features)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: []
}

func ExampleG2diagnosticClient_GetAvailableMemory() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	result, err := g2diagnostic.GetAvailableMemory(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleG2diagnosticClient_GetDataSourceCounts() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	result, err := g2diagnostic.GetDataSourceCounts(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: [{"DSRC_ID":1001,"DSRC_CODE":"CUSTOMERS","ETYPE_ID":3,"ETYPE_CODE":"GENERIC","OBS_ENT_COUNT":7,"DSRC_RECORD_COUNT":7}]
}

func ExampleG2diagnosticClient_GetDBInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	result, err := g2diagnostic.GetDBInfo(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 52))
	// Output: {"Hybrid Mode":false,"Database Details":[{"Name":...
}

func ExampleG2diagnosticClient_GetEntityDetails() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	entityID := int64(1)
	includeInternalFeatures := 1
	result, err := g2diagnostic.GetEntityDetails(ctx, entityID, includeInternalFeatures)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: [{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"PRIMARY","FEAT_DESC":"Robert Smith"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"No","FTYPE_CODE":"DOB","USAGE_TYPE":"","FEAT_DESC":"12/11/1978"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"MAILING","FEAT_DESC":"123 Main Street, Las Vegas NV 89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"No","FTYPE_CODE":"PHONE","USAGE_TYPE":"HOME","FEAT_DESC":"702-919-1300"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"No","FTYPE_CODE":"EMAIL","USAGE_TYPE":"","FEAT_DESC":"bsmith@work.com"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB.MMYY_HASH=1278"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|ADDRESS.CITY_STD=LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|POST=89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|PHONE.PHONE_LAST_5=91300"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"123|MN||LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"123|MN||89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"PHONE_KEY","USAGE_TYPE":"","FEAT_DESC":"7029191300"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"No","FTYPE_CODE":"RECORD_TYPE","USAGE_TYPE":"","FEAT_DESC":"PERSON"},{"RES_ENT_ID":1,"OBS_ENT_ID":1,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","DERIVED":"Yes","FTYPE_CODE":"EMAIL_KEY","USAGE_TYPE":"","FEAT_DESC":"bsmith@WORK.COM"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"PRIMARY","FEAT_DESC":"Bob Smith"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"No","FTYPE_CODE":"DOB","USAGE_TYPE":"","FEAT_DESC":"11/12/1978"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"HOME","FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"No","FTYPE_CODE":"PHONE","USAGE_TYPE":"MOBILE","FEAT_DESC":"702-919-1300"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|ADDRESS.CITY_STD=LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|PHONE.PHONE_LAST_5=91300"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB.MMYY_HASH=1178"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|POST=89111"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|ADDRESS.CITY_STD=LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|POST=89111"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|DOB.MMYY_HASH=1178"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|PHONE.PHONE_LAST_5=91300"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"1515|ATL||89111"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"1515|ATL||LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"Yes","FTYPE_CODE":"PHONE_KEY","USAGE_TYPE":"","FEAT_DESC":"7029191300"},{"RES_ENT_ID":1,"OBS_ENT_ID":2,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","DERIVED":"No","FTYPE_CODE":"RECORD_TYPE","USAGE_TYPE":"","FEAT_DESC":"PERSON"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"PRIMARY","FEAT_DESC":"Bob J Smith"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"No","FTYPE_CODE":"DOB","USAGE_TYPE":"","FEAT_DESC":"12/11/1978"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"No","FTYPE_CODE":"EMAIL","USAGE_TYPE":"","FEAT_DESC":"bsmith@work.com"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB.MMYY_HASH=1278"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|J|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|J|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|J|SM0|DOB.MMYY_HASH=1278"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|SM0|DOB.MMYY_HASH=1278"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"BB|J|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"No","FTYPE_CODE":"RECORD_TYPE","USAGE_TYPE":"","FEAT_DESC":"PERSON"},{"RES_ENT_ID":1,"OBS_ENT_ID":3,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","DERIVED":"Yes","FTYPE_CODE":"EMAIL_KEY","USAGE_TYPE":"","FEAT_DESC":"bsmith@WORK.COM"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"PRIMARY","FEAT_DESC":"B Smith"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"No","FTYPE_CODE":"DOB","USAGE_TYPE":"","FEAT_DESC":"11/12/1979"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"HOME","FEAT_DESC":"1515 Adela Ln Las Vegas NV 89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"No","FTYPE_CODE":"EMAIL","USAGE_TYPE":"","FEAT_DESC":"bsmith@work.com"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"B|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"B|SM0|DOB=71211"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"B|SM0|DOB.MMYY_HASH=1179"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"B|SM0|ADDRESS.CITY_STD=LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"B|SM0|POST=89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"B|SM0|DOB.MMDD_HASH=1211"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"1515|ATL||LS FKS"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"1515|ATL||89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"No","FTYPE_CODE":"RECORD_TYPE","USAGE_TYPE":"","FEAT_DESC":"PERSON"},{"RES_ENT_ID":1,"OBS_ENT_ID":4,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","DERIVED":"Yes","FTYPE_CODE":"EMAIL_KEY","USAGE_TYPE":"","FEAT_DESC":"bsmith@WORK.COM"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"PRIMARY","FEAT_DESC":"Robbie Smith"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"MAILING","FEAT_DESC":"123 E Main St Henderson NV 89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"No","FTYPE_CODE":"DRLIC","USAGE_TYPE":"","FEAT_DESC":"112233 NV"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|POST=89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RBRT|SM0|ADDRESS.CITY_STD=HNTRSN"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RB|SM0|POST=89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RB|SM0"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","FEAT_DESC":"RB|SM0|ADDRESS.CITY_STD=HNTRSN"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"123|MN||89132"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","FEAT_DESC":"123|MN||HNTRSN"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"Yes","FTYPE_CODE":"ID_KEY","USAGE_TYPE":"","FEAT_DESC":"DRLIC=112233"},{"RES_ENT_ID":1,"OBS_ENT_ID":5,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","DERIVED":"No","FTYPE_CODE":"RECORD_TYPE","USAGE_TYPE":"","FEAT_DESC":"PERSON"}]
}

func ExampleG2diagnosticClient_GetEntityListBySize() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	entitySize := 1000
	entityListBySizeHandle, err := g2diagnostic.GetEntityListBySize(ctx, entitySize)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(entityListBySizeHandle > 0) // Dummy output.
	// Output: true
}

func ExampleG2diagnosticClient_GetEntityResume() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	entityID := int64(1)
	result, err := g2diagnostic.GetEntityResume(ctx, entityID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: [{"RES_ENT_ID":1,"REL_ENT_ID":0,"ERRULE_CODE":"","MATCH_KEY":"","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1001","ENT_SRC_DESC":"Robert Smith","JSON_DATA":"{\"RECORD_TYPE\":\"PERSON\",\"PRIMARY_NAME_LAST\":\"Smith\",\"PRIMARY_NAME_FIRST\":\"Robert\",\"DATE_OF_BIRTH\":\"12/11/1978\",\"ADDR_TYPE\":\"MAILING\",\"ADDR_LINE1\":\"123 Main Street, Las Vegas NV 89132\",\"PHONE_TYPE\":\"HOME\",\"PHONE_NUMBER\":\"702-919-1300\",\"EMAIL_ADDRESS\":\"bsmith@work.com\",\"DATE\":\"1/2/18\",\"STATUS\":\"Active\",\"AMOUNT\":\"100\",\"DATA_SOURCE\":\"CUSTOMERS\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"1001\"}"},{"RES_ENT_ID":1,"REL_ENT_ID":0,"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1002","ENT_SRC_DESC":"Bob Smith","JSON_DATA":"{\"RECORD_TYPE\":\"PERSON\",\"PRIMARY_NAME_LAST\":\"Smith\",\"PRIMARY_NAME_FIRST\":\"Bob\",\"DATE_OF_BIRTH\":\"11/12/1978\",\"ADDR_TYPE\":\"HOME\",\"ADDR_LINE1\":\"1515 Adela Lane\",\"ADDR_CITY\":\"Las Vegas\",\"ADDR_STATE\":\"NV\",\"ADDR_POSTAL_CODE\":\"89111\",\"PHONE_TYPE\":\"MOBILE\",\"PHONE_NUMBER\":\"702-919-1300\",\"DATE\":\"3/10/17\",\"STATUS\":\"Inactive\",\"AMOUNT\":\"200\",\"DATA_SOURCE\":\"CUSTOMERS\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"1002\"}"},{"RES_ENT_ID":1,"REL_ENT_ID":0,"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1003","ENT_SRC_DESC":"Bob J Smith","JSON_DATA":"{\"RECORD_TYPE\":\"PERSON\",\"PRIMARY_NAME_LAST\":\"Smith\",\"PRIMARY_NAME_FIRST\":\"Bob\",\"PRIMARY_NAME_MIDDLE\":\"J\",\"DATE_OF_BIRTH\":\"12/11/1978\",\"EMAIL_ADDRESS\":\"bsmith@work.com\",\"DATE\":\"4/9/16\",\"STATUS\":\"Inactive\",\"AMOUNT\":\"300\",\"DATA_SOURCE\":\"CUSTOMERS\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"1003\"}"},{"RES_ENT_ID":1,"REL_ENT_ID":0,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1004","ENT_SRC_DESC":"B Smith","JSON_DATA":"{\"RECORD_TYPE\":\"PERSON\",\"PRIMARY_NAME_LAST\":\"Smith\",\"PRIMARY_NAME_FIRST\":\"B\",\"DATE_OF_BIRTH\":\"11/12/1979\",\"ADDR_TYPE\":\"HOME\",\"ADDR_LINE1\":\"1515 Adela Ln\",\"ADDR_CITY\":\"Las Vegas\",\"ADDR_STATE\":\"NV\",\"ADDR_POSTAL_CODE\":\"89132\",\"EMAIL_ADDRESS\":\"bsmith@work.com\",\"DATE\":\"1/5/15\",\"STATUS\":\"Inactive\",\"AMOUNT\":\"400\",\"DATA_SOURCE\":\"CUSTOMERS\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"1004\"}"},{"RES_ENT_ID":1,"REL_ENT_ID":0,"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS","DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","RECORD_ID":"1005","ENT_SRC_DESC":"Robbie Smith","JSON_DATA":"{\"RECORD_TYPE\":\"PERSON\",\"PRIMARY_NAME_LAST\":\"Smith\",\"PRIMARY_NAME_FIRST\":\"Robbie\",\"DRIVERS_LICENSE_NUMBER\":\"112233\",\"DRIVERS_LICENSE_STATE\":\"NV\",\"ADDR_TYPE\":\"MAILING\",\"ADDR_LINE1\":\"123 E Main St\",\"ADDR_CITY\":\"Henderson\",\"ADDR_STATE\":\"NV\",\"ADDR_POSTAL_CODE\":\"89132\",\"DATE\":\"7/16/19\",\"STATUS\":\"Active\",\"AMOUNT\":\"500\",\"DATA_SOURCE\":\"CUSTOMERS\",\"ENTITY_TYPE\":\"GENERIC\",\"DSRC_ACTION\":\"A\",\"RECORD_ID\":\"1005\"}"}]
}

func ExampleG2diagnosticClient_GetEntitySizeBreakdown() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	minimumEntitySize := 1
	includeInternalFeatures := 1
	result, err := g2diagnostic.GetEntitySizeBreakdown(ctx, minimumEntitySize, includeInternalFeatures)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: [{"ENTITY_SIZE": 5,"ENTITY_COUNT": 1,"NAME": 5.00,"DOB": 3.00,"ADDRESS": 4.00,"PHONE": 2.00,"DRLIC": 1.00,"EMAIL": 1.00,"NAME_KEY": 31.00,"ADDR_KEY": 6.00,"ID_KEY": 1.00,"PHONE_KEY": 1.00,"RECORD_TYPE": 1.00,"EMAIL_KEY": 1.00,"MIN_RES_ENT_ID": 1,"MAX_RES_ENT_ID": 1},{"ENTITY_SIZE": 1,"ENTITY_COUNT": 2,"NAME": 1.00,"DOB": 1.00,"GENDER": 0.50,"ADDRESS": 1.00,"NAME_KEY": 6.00,"ADDR_KEY": 2.00,"RECORD_TYPE": 1.00,"MIN_RES_ENT_ID": 6,"MAX_RES_ENT_ID": 7}]
}

func ExampleG2diagnosticClient_GetFeature() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	libFeatID := int64(1)
	result, err := g2diagnostic.GetFeature(ctx, libFeatID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"LIB_FEAT_ID":1,"FTYPE_CODE":"NAME","ELEMENTS":[{"FELEM_CODE":"TOKENIZED_NM","FELEM_VALUE":"ROBERT|SMITH"},{"FELEM_CODE":"CATEGORY","FELEM_VALUE":"PERSON"},{"FELEM_CODE":"CULTURE","FELEM_VALUE":"ANGLO"},{"FELEM_CODE":"GIVEN_NAME","FELEM_VALUE":"Robert"},{"FELEM_CODE":"SUR_NAME","FELEM_VALUE":"Smith"},{"FELEM_CODE":"FULL_NAME","FELEM_VALUE":"Robert Smith"}]}
}

func ExampleG2diagnosticClient_GetGenericFeatures() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	featureType := "PHONE"
	maximumEstimatedCount := 10
	result, err := g2diagnostic.GetGenericFeatures(ctx, featureType, maximumEstimatedCount)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: []
}

func ExampleG2diagnosticClient_GetLogicalCores() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	result, err := g2diagnostic.GetLogicalCores(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleG2diagnosticClient_GetMappingStatistics() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	includeInternalFeatures := 1
	result, err := g2diagnostic.GetMappingStatistics(ctx, includeInternalFeatures)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: [{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"NAME","USAGE_TYPE":"PRIMARY","REC_COUNT":7,"REC_PCT":1.0,"UNIQ_COUNT":6,"UNIQ_PCT":0.8571428571428571,"MIN_FEAT_DESC":"B Smith","MAX_FEAT_DESC":"Robert Smith"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"DOB","USAGE_TYPE":"","REC_COUNT":6,"REC_PCT":0.8571428571428571,"UNIQ_COUNT":5,"UNIQ_PCT":0.8333333333333334,"MIN_FEAT_DESC":"10/10/70","MAX_FEAT_DESC":"3/15/90"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"GENDER","USAGE_TYPE":"","REC_COUNT":1,"REC_PCT":0.14285714285714286,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"M","MAX_FEAT_DESC":"M"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"HOME","REC_COUNT":4,"REC_PCT":0.5714285714285714,"UNIQ_COUNT":3,"UNIQ_PCT":0.75,"MIN_FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","MAX_FEAT_DESC":"3212 W. 32nd St Palm Harbor, FL 60527"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"ADDRESS","USAGE_TYPE":"MAILING","REC_COUNT":2,"REC_PCT":0.2857142857142857,"UNIQ_COUNT":2,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"123 E Main St Henderson NV 89132","MAX_FEAT_DESC":"123 Main Street, Las Vegas NV 89132"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"PHONE","USAGE_TYPE":"HOME","REC_COUNT":1,"REC_PCT":0.14285714285714286,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"702-919-1300","MAX_FEAT_DESC":"702-919-1300"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"PHONE","USAGE_TYPE":"MOBILE","REC_COUNT":1,"REC_PCT":0.14285714285714286,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"702-919-1300","MAX_FEAT_DESC":"702-919-1300"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"DRLIC","USAGE_TYPE":"","REC_COUNT":1,"REC_PCT":0.14285714285714286,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"112233 NV","MAX_FEAT_DESC":"112233 NV"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"EMAIL","USAGE_TYPE":"","REC_COUNT":3,"REC_PCT":0.42857142857142857,"UNIQ_COUNT":1,"UNIQ_PCT":0.3333333333333333,"MIN_FEAT_DESC":"bsmith@work.com","MAX_FEAT_DESC":"bsmith@work.com"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"NAME_KEY","USAGE_TYPE":"","REC_COUNT":57,"REC_PCT":8.142857142857143,"UNIQ_COUNT":40,"UNIQ_PCT":0.7017543859649122,"MIN_FEAT_DESC":"BB|J|SM0","MAX_FEAT_DESC":"RB|SM0|POST=89132"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"ADDR_KEY","USAGE_TYPE":"","REC_COUNT":12,"REC_PCT":1.7142857142857143,"UNIQ_COUNT":8,"UNIQ_PCT":0.6666666666666666,"MIN_FEAT_DESC":"123|MN||89132","MAX_FEAT_DESC":"3212|NT||PLM HRBR"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"ID_KEY","USAGE_TYPE":"","REC_COUNT":1,"REC_PCT":0.14285714285714286,"UNIQ_COUNT":1,"UNIQ_PCT":1.0,"MIN_FEAT_DESC":"DRLIC=112233","MAX_FEAT_DESC":"DRLIC=112233"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"PHONE_KEY","USAGE_TYPE":"","REC_COUNT":2,"REC_PCT":0.2857142857142857,"UNIQ_COUNT":1,"UNIQ_PCT":0.5,"MIN_FEAT_DESC":"7029191300","MAX_FEAT_DESC":"7029191300"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"No","FTYPE_CODE":"RECORD_TYPE","USAGE_TYPE":"","REC_COUNT":7,"REC_PCT":1.0,"UNIQ_COUNT":1,"UNIQ_PCT":0.14285714285714286,"MIN_FEAT_DESC":"PERSON","MAX_FEAT_DESC":"PERSON"},{"DSRC_CODE":"CUSTOMERS","ETYPE_CODE":"GENERIC","DERIVED":"Yes","FTYPE_CODE":"EMAIL_KEY","USAGE_TYPE":"","REC_COUNT":3,"REC_PCT":0.42857142857142857,"UNIQ_COUNT":1,"UNIQ_PCT":0.3333333333333333,"MIN_FEAT_DESC":"bsmith@WORK.COM","MAX_FEAT_DESC":"bsmith@WORK.COM"}]
}

func ExampleG2diagnosticClient_GetPhysicalCores() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	result, err := g2diagnostic.GetPhysicalCores(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleG2diagnosticClient_GetRelationshipDetails() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	relationshipID := int64(1)
	includeInternalFeatures := 1
	result, err := g2diagnostic.GetRelationshipDetails(ctx, relationshipID, includeInternalFeatures)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: [{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME","FEAT_DESC":"John Smith"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"DOB","FEAT_DESC":"10/10/70"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"GENDER","FEAT_DESC":"M"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"ADDRESS","FEAT_DESC":"3212 W. 32nd St Palm Harbor, FL 60527"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|ADDRESS.CITY_STD=PLM HRBR"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|DOB.MMDD_HASH=1010"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|DOB=71010"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|POST=60527"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|DOB.MMYY_HASH=1070"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"ADDR_KEY","FEAT_DESC":"3212|NT||60527"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"ADDR_KEY","FEAT_DESC":"3212|NT||PLM HRBR"},{"RES_ENT_ID":6,"ERRULE_CODE":"","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"RECORD_TYPE","FEAT_DESC":"PERSON"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME","FEAT_DESC":"John Smith"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"DOB","FEAT_DESC":"3/15/90"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"ADDRESS","FEAT_DESC":"3212 W. 32nd St Palm Harbor, FL 60527"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|ADDRESS.CITY_STD=PLM HRBR"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|POST=60527"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|DOB.MMYY_HASH=0390"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|DOB.MMDD_HASH=1503"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"NAME_KEY","FEAT_DESC":"JN|SM0|DOB=91503"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"ADDR_KEY","FEAT_DESC":"3212|NT||60527"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"ADDR_KEY","FEAT_DESC":"3212|NT||PLM HRBR"},{"RES_ENT_ID":7,"ERRULE_CODE":"CNAME_CFF_DEXCL","MATCH_KEY":"+NAME+ADDRESS-DOB","FTYPE_CODE":"RECORD_TYPE","FEAT_DESC":"PERSON"}]
}

func ExampleG2diagnosticClient_GetResolutionStatistics() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	result, err := g2diagnostic.GetResolutionStatistics(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: [{"MATCH_LEVEL":1,"MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME","RAW_MATCH_KEYS":[{"MATCH_KEY":"+DOB+ADDRESS+EMAIL+PNAME"}],"ERRULE_ID":110,"ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","IS_AMBIGUOUS":"No","RECORD_COUNT":1,"MIN_RES_ENT_ID":1,"MAX_RES_ENT_ID":1,"MIN_RES_REL_ID":0,"MAX_RES_REL_ID":0},{"MATCH_LEVEL":1,"MATCH_KEY":"+NAME+DOB+EMAIL","RAW_MATCH_KEYS":[{"MATCH_KEY":"+NAME+DOB+EMAIL"}],"ERRULE_ID":120,"ERRULE_CODE":"SF1_PNAME_CSTAB","IS_AMBIGUOUS":"No","RECORD_COUNT":1,"MIN_RES_ENT_ID":1,"MAX_RES_ENT_ID":1,"MIN_RES_REL_ID":0,"MAX_RES_REL_ID":0},{"MATCH_LEVEL":1,"MATCH_KEY":"+NAME+DOB+PHONE","RAW_MATCH_KEYS":[{"MATCH_KEY":"+NAME+DOB+PHONE"}],"ERRULE_ID":160,"ERRULE_CODE":"CNAME_CFF_CEXCL","IS_AMBIGUOUS":"No","RECORD_COUNT":1,"MIN_RES_ENT_ID":1,"MAX_RES_ENT_ID":1,"MIN_RES_REL_ID":0,"MAX_RES_REL_ID":0},{"MATCH_LEVEL":1,"MATCH_KEY":"+NAME+ADDRESS","RAW_MATCH_KEYS":[{"MATCH_KEY":"+NAME+ADDRESS"}],"ERRULE_ID":162,"ERRULE_CODE":"CNAME_CFF","IS_AMBIGUOUS":"No","RECORD_COUNT":1,"MIN_RES_ENT_ID":1,"MAX_RES_ENT_ID":1,"MIN_RES_REL_ID":0,"MAX_RES_REL_ID":0},{"MATCH_LEVEL":2,"MATCH_KEY":"+NAME+ADDRESS-DOB","RAW_MATCH_KEYS":[{"MATCH_KEY":"+NAME+ADDRESS-DOB"}],"ERRULE_ID":164,"ERRULE_CODE":"CNAME_CFF_DEXCL","IS_AMBIGUOUS":"No","RECORD_COUNT":1,"MIN_RES_ENT_ID":6,"MAX_RES_ENT_ID":7,"MIN_RES_REL_ID":1,"MAX_RES_REL_ID":1}]
}

func ExampleG2diagnosticClient_GetTotalSystemMemory() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	result, err := g2diagnostic.GetTotalSystemMemory(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

// func ExampleG2diagnosticClient_SetLogLevel() {
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
// 	g2diagnostic := &G2diagnosticClient{}
// 	ctx := context.TODO()
// 	err := g2diagnostic.SetLogLevel(ctx, logger.LevelInfo)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// Output:
// }

func ExampleG2diagnosticClient_Init() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	grpcConnection := getGrpcConnection()
	g2diagnostic := &G2diagnosticClient{
		GrpcClient: pb.NewG2DiagnosticClient(grpcConnection),
	}
	moduleName := "Test module name"
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("") // See https://pkg.go.dev/github.com/senzing/go-common
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := 0
	err = g2diagnostic.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		// This should produce a "senzing-60134002" error.
	}
	// Output:
}

func ExampleG2diagnosticClient_InitWithConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	grpcConnection := getGrpcConnection()
	g2diagnostic := &G2diagnosticClient{
		GrpcClient: pb.NewG2DiagnosticClient(grpcConnection),
	}
	moduleName := "Test module name"
	iniParams, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if err != nil {
		fmt.Println(err)
	}
	initConfigID := int64(1)
	verboseLogging := 0
	err = g2diagnostic.InitWithConfigID(ctx, moduleName, iniParams, initConfigID, verboseLogging)
	if err != nil {
		// This should produce a "senzing-60134003" error.
	}
	// Output:
}

func ExampleG2diagnosticClient_Reinit() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	g2Configmgr := getG2Configmgr(ctx)
	initConfigID, _ := g2Configmgr.GetDefaultConfigID(ctx)
	err := g2diagnostic.Reinit(ctx, initConfigID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2diagnosticClient_Destroy() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-grpc/blob/main/g2diagnosticclient/g2diagnosticclient_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	err := g2diagnostic.Destroy(ctx)
	if err != nil {
		// This should produce a "senzing-60134001" error.
	}
	// Output:
}
