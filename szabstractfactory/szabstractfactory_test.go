package szabstractfactory_test

import (
	"context"
	"fmt"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	baseCallerSkip    = 4
	defaultTruncation = 76
	instanceName      = "SzAbstractFactory Test"
	printResults      = false
	verboseLogging    = senzing.SzNoLogging
)

var (
	grpcAddress    = "0.0.0.0:8261"
	grpcConnection *grpc.ClientConn
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzAbstractFactory_CreateConfigManager(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	require.NoError(test, err)
	configList, err := szConfigManager.GetConfigs(ctx)
	require.NoError(test, err)
	printActual(test, configList)
}

func TestSzAbstractFactory_CreateDiagnostic(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
	require.NoError(test, err)
	result, err := szDiagnostic.CheckDatastorePerformance(ctx, 1)
	require.NoError(test, err)
	printActual(test, result)
}

func TestSzAbstractFactory_CreateEngine(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	require.NoError(test, err)
	stats, err := szEngine.GetStats(ctx)
	require.NoError(test, err)
	printActual(test, stats)
}

func TestSzAbstractFactory_CreateProduct(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	require.NoError(test, err)
	version, err := szProduct.GetVersion(ctx)
	require.NoError(test, err)
	printActual(test, version)
}

func TestSzAbstractFactory_Destroy(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
}

func TestSzAbstractFactory_Reinitialize(test *testing.T) {
	ctx := context.TODO()
	szAbstractFactory := getTestObject(ctx, test)
	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
	_, err := szAbstractFactory.CreateDiagnostic(ctx)
	require.NoError(test, err)
	_, err = szAbstractFactory.CreateEngine(ctx)
	require.NoError(test, err)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	require.NoError(test, err)
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	require.NoError(test, err)
	err = szAbstractFactory.Reinitialize(ctx, configID)
	require.NoError(test, err)
}

// func TestSzAbstractFactory_Reinitialize_extended(test *testing.T) {
// 	ctx := context.TODO()
// 	newDataSourceName := "BOB"
// 	newRecordID := "9999"
// 	newRecord := `{"DATA_SOURCE": "BOB", "RECORD_ID": "9999", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`

// 	szAbstractFactory := getTestObject(ctx, test)
// 	defer func() { handleError(szAbstractFactory.Destroy(ctx)) }()
// 	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
// 	require.NoError(test, err)
// 	szConfig, err := szAbstractFactory.CreateConfig(ctx)
// 	require.NoError(test, err)
// 	szDiagnostic, err := szAbstractFactory.CreateDiagnostic(ctx)
// 	require.NoError(test, err)
// 	szEngine, err := szAbstractFactory.CreateEngine(ctx)
// 	require.NoError(test, err)

// 	oldConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
// 	require.NoError(test, err)

// 	oldJSONConfig, err := szConfigManager.GetConfig(ctx, oldConfigID)
// 	require.NoError(test, err)

// 	configHandle, err := szConfig.ImportConfig(ctx, oldJSONConfig)
// 	require.NoError(test, err)

// 	_, err = szConfig.AddDataSource(ctx, configHandle, newDataSourceName)
// 	require.NoError(test, err)

// 	newJSONConfig, err := szConfig.ExportConfig(ctx, configHandle)
// 	require.NoError(test, err)

// 	newConfigID, err := szConfigManager.AddConfig(ctx, newJSONConfig, "Add TruthSet datasources")
// 	require.NoError(test, err)

// 	err = szConfigManager.ReplaceDefaultConfigID(ctx, oldConfigID, newConfigID)
// 	require.NoError(test, err)

// 	err = szAbstractFactory.Reinitialize(ctx, newConfigID)
// 	require.NoError(test, err)

// 	_, err = szEngine.AddRecord(ctx, newDataSourceName, newRecordID, string(newRecord), senzing.SzWithInfo)
// 	require.NoError(test, err)

// 	_, err = szEngine.DeleteRecord(ctx, newDataSourceName, newRecordID, senzing.SzWithInfo)
// 	require.NoError(test, err)

// 	err = szDiagnostic.PurgeRepository(ctx)
// 	require.NoError(test, err)

// 	err = szConfigManager.ReplaceDefaultConfigID(ctx, newConfigID, oldConfigID)
// 	require.NoError(test, err)

// 	err = szAbstractFactory.Reinitialize(ctx, oldConfigID)
// 	require.NoError(test, err)

// 	_, err = szDiagnostic.CheckDatastorePerformance(ctx, 1)
// 	require.NoError(test, err)
// }

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getGrpcConnection() *grpc.ClientConn {
	if grpcConnection == nil {
		transportCredentials, err := helper.GetGrpcTransportCredentials()
		if err != nil {
			panic(err)
		}
		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(transportCredentials),
		}
		grpcConnection, err = grpc.NewClient(grpcAddress, dialOptions...)
		if err != nil {
			panic(err)
		}
	}
	return grpcConnection
}

func getSzAbstractFactory(ctx context.Context) (senzing.SzAbstractFactory, error) {
	_ = ctx
	result := &Szabstractfactory{
		GrpcConnection: getGrpcConnection(),
	}
	return result, nil
}

func getTestObject(ctx context.Context, test *testing.T) senzing.SzAbstractFactory {
	result, err := getSzAbstractFactory(ctx)
	require.NoError(test, err)
	return result
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
