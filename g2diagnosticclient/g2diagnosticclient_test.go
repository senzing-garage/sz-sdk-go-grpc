package g2diagnosticclient

import (
	"context"
	"fmt"
	"os"
	"testing"

	g2diagnosticsdk "github.com/senzing/g2-sdk-go/g2diagnostic"
	"github.com/senzing/go-helpers/g2engineconfigurationjson"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcAddress    = "localhost:50051"
	grpcConnection *grpc.ClientConn
	g2diagnostic   g2diagnosticsdk.G2diagnostic
)

// ----------------------------------------------------------------------------
// Internal functions - names begin with lowercase letter
// ----------------------------------------------------------------------------

func getGrpcConnection() *grpc.ClientConn {
	var err error
	if grpcConnection == nil {
		fmt.Println(">>>>>>>>>>>>>>> Getting connection")
		grpcConnection, err = grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Did not connect: %v\n", err)
		}
		//		defer grpcConnection.Close()
	}
	return grpcConnection
}

func getTestObject(ctx context.Context) g2diagnosticsdk.G2diagnostic {
	// FIXME: Make a reusable G2diagnosticClient.
	// if g2diagnostic == nil {

	// 	fmt.Println(">>>>>>>>>>>>>>> Getting G2diagnosticClient")

	// 	grpcConnection := getGrpcConnection()
	// 	g2diagnostic = &G2diagnosticClient{
	// 		G2DiagnosticGrpcClient: pb.NewG2DiagnosticClient(grpcConnection),
	// 	}

	// 	moduleName := "Test module name"
	// 	verboseLogging := 0 // 0 for no Senzing logging; 1 for logging
	// 	iniParams, jsonErr := g2configuration.BuildSimpleSystemConfigurationJson("")
	// 	if jsonErr != nil {
	// 		logger.Fatalf("Cannot construct system configuration: %v", jsonErr)
	// 	}

	// 	initErr := g2diagnostic.Init(ctx, moduleName, iniParams, verboseLogging)
	// 	if initErr != nil {
	// 		logger.Fatalf("Cannot Init: %v", initErr)
	// 	}
	// }
	return g2diagnostic
}

func testError(test *testing.T, ctx context.Context, g2diagnostic g2diagnosticsdk.G2diagnostic, err error) {
	if err != nil {
		assert.FailNow(test, err.Error())
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

func TestGetObject(test *testing.T) {
	ctx := context.TODO()
	getTestObject(ctx)
}

// ----------------------------------------------------------------------------
// Test interface functions - names begin with "Test"
// ----------------------------------------------------------------------------

func TestCheckDBPerf(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	secondsToRun := 1
	actual, err := g2diagnostic.CheckDBPerf(ctx, secondsToRun)
	testError(test, ctx, g2diagnostic, err)
	test.Log("Actual:", actual)
}

// FAIL:
//func TestEntityListBySize(test *testing.T) {
//	ctx := context.TODO()
//	g2diagnostic := getTestObject(ctx)
//	aSize := 10
//
//	aHandle, err := g2diagnostic.GetEntityListBySize(ctx, aSize)
//	testError(test, ctx, g2diagnostic, err)
//
//	anEntity, err := g2diagnostic.FetchNextEntityBySize(ctx, aHandle)
//	testError(test, ctx, g2diagnostic, err)
//	test.Log("Entity:", anEntity)
//
//	err = g2diagnostic.CloseEntityListBySize(ctx, aHandle)
//	testError(test, ctx, g2diagnostic, err)
//}

func TestFindEntitiesByFeatureIDs(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	features := "{\"ENTITY_ID\":1,\"LIB_FEAT_IDS\":[1,3,4]}"
	actual, err := g2diagnostic.FindEntitiesByFeatureIDs(ctx, features)
	testError(test, ctx, g2diagnostic, err)
	test.Log("len(Actual):", len(actual))
}

func TestGetAvailableMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	actual, err := g2diagnostic.GetAvailableMemory(ctx)
	testError(test, ctx, g2diagnostic, err)
	assert.Greater(test, actual, int64(0))
	test.Log("Actual:", actual)
}

func TestGetDataSourceCounts(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	actual, err := g2diagnostic.GetDataSourceCounts(ctx)
	testError(test, ctx, g2diagnostic, err)
	test.Log("Data Source counts:", actual)
}

func TestGetDBInfo(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	actual, err := g2diagnostic.GetDBInfo(ctx)
	testError(test, ctx, g2diagnostic, err)
	test.Log("Actual:", actual)
}

func TestGetEntityDetails(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	entityID := int64(1)
	includeInternalFeatures := 1
	actual, err := g2diagnostic.GetEntityDetails(ctx, entityID, includeInternalFeatures)
	testError(test, ctx, g2diagnostic, err)
	test.Log("Actual:", actual)
}

func TestGetEntityResume(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	entityID := int64(1)
	actual, err := g2diagnostic.GetEntityResume(ctx, entityID)
	testError(test, ctx, g2diagnostic, err)
	test.Log("Actual:", actual)
}

func TestGetEntitySizeBreakdown(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	minimumEntitySize := 1
	includeInternalFeatures := 1
	actual, err := g2diagnostic.GetEntitySizeBreakdown(ctx, minimumEntitySize, includeInternalFeatures)
	testError(test, ctx, g2diagnostic, err)
	test.Log("Actual:", actual)
}

func TestGetFeature(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	libFeatID := int64(1)
	actual, err := g2diagnostic.GetFeature(ctx, libFeatID)
	testError(test, ctx, g2diagnostic, err)
	test.Log("Actual:", actual)
}

func TestGetGenericFeatures(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	featureType := "PHONE"
	maximumEstimatedCount := 10
	actual, err := g2diagnostic.GetGenericFeatures(ctx, featureType, maximumEstimatedCount)
	testError(test, ctx, g2diagnostic, err)
	test.Log("Actual:", actual)
}

func TestGetLogicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	actual, err := g2diagnostic.GetLogicalCores(ctx)
	testError(test, ctx, g2diagnostic, err)
	assert.Greater(test, actual, 0)
	test.Log("Actual:", actual)
}

func TestGetMappingStatistics(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	includeInternalFeatures := 1
	actual, err := g2diagnostic.GetMappingStatistics(ctx, includeInternalFeatures)
	testError(test, ctx, g2diagnostic, err)
	test.Log("Actual:", actual)
}

func TestGetPhysicalCores(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	actual, err := g2diagnostic.GetPhysicalCores(ctx)
	testError(test, ctx, g2diagnostic, err)
	assert.Greater(test, actual, 0)
	test.Log("Actual:", actual)
}

func TestGetRelationshipDetails(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	relationshipID := int64(1)
	includeInternalFeatures := 1
	actual, err := g2diagnostic.GetRelationshipDetails(ctx, relationshipID, includeInternalFeatures)
	testError(test, ctx, g2diagnostic, err)
	test.Log("Actual:", actual)
}

func TestGetResolutionStatistics(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	actual, err := g2diagnostic.GetResolutionStatistics(ctx)
	testError(test, ctx, g2diagnostic, err)
	test.Log("Actual:", actual)
}

func TestGetTotalSystemMemory(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	actual, err := g2diagnostic.GetTotalSystemMemory(ctx)
	testError(test, ctx, g2diagnostic, err)
	assert.Greater(test, actual, int64(0))
	test.Log("Actual:", actual)
}

func TestInit(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	testError(test, ctx, g2diagnostic, jsonErr)
	err := g2diagnostic.Init(ctx, moduleName, iniParams, verboseLogging)
	testError(test, ctx, g2diagnostic, err)
}

func TestInitWithConfigID(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	moduleName := "Test module name"
	initConfigID := int64(1)
	verboseLogging := 0
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	testError(test, ctx, g2diagnostic, jsonErr)
	err := g2diagnostic.InitWithConfigID(ctx, moduleName, iniParams, initConfigID, verboseLogging)
	testError(test, ctx, g2diagnostic, err)
}

func TestReinit(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	initConfigID := int64(4019066234)
	err := g2diagnostic.Reinit(ctx, initConfigID)
	testError(test, ctx, g2diagnostic, err)
}

// ----------------------------------------------------------------------------
// Finis
// ----------------------------------------------------------------------------

func TestDestroy(test *testing.T) {
	ctx := context.TODO()
	g2diagnostic := getTestObject(ctx)
	err := g2diagnostic.Destroy(ctx)
	testError(test, ctx, g2diagnostic, err)
}
