package szconfig

import (
	"context"
	"fmt"
	"testing"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	badConfigDefinition = "}{"
	badConfigHandle     = uintptr(0)
	badDataSourceCode   = "\n\tGO_TEST"
	badLogLevelName     = "BadLogLevelName"
	badSettings         = "{]"
	defaultTruncation   = 76
	instanceName        = "SzConfig Test"
	observerOrigin      = "SzConfig observer"
	printResults        = false
	verboseLogging      = senzing.SzNoLogging
)

var (
	grpcAddress       = "localhost:8261"
	grpcConnection    *grpc.ClientConn
	logLevel          = helper.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szConfigSingleton *Szconfig
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzconfig_AddDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	dataSourceCode := "GO_TEST"
	actual, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	require.NoError(test, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
}

func TestSzconfig_AddDataSource_withLoad(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	require.NoError(test, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
	configHandle2, err := szConfig.ImportConfig(ctx, configDefinition)
	require.NoError(test, err)
	dataSourceCode := "GO_TEST"
	actual, err := szConfig.AddDataSource(ctx, configHandle2, dataSourceCode)
	require.NoError(test, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle2)
	require.NoError(test, err)
}

func TestSzconfig_AddDataSource_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	actual, err := szConfig.AddDataSource(ctx, configHandle, badDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzconfig_CloseConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
}

func TestSzconfig_CloseConfig_badConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.CloseConfig(ctx, badConfigHandle)
	require.NoError(test, err) // TODO: Is this the correct behavior?
	// require.ErrorIs(test, err, szerror.ErrSzBase) // TODO: Update to correct error once sz-sdk-go/szerror/ is updated.
}

// TODO: Implement TestSzconfig_CloseConfig_error
// func TestSzconfig_CloseConfig_error(test *testing.T) {}

func TestSzconfig_CreateConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// TODO: Implement TestSzconfig_CreateConfig_error
// func TestSzconfig_CreateConfig_error(test *testing.T) {}

func TestSzconfig_DeleteDataSource(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	require.NoError(test, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printResult(test, "Original", actual)
	dataSourceCode := "GO_TEST"
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	require.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printResult(test, "     Add", actual)
	err = szConfig.DeleteDataSource(ctx, configHandle, dataSourceCode)
	require.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printResult(test, "  Delete", actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
}

func TestSzconfig_DeleteDataSource_withLoad(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printResult(test, "Original", actual)
	dataSourceCode := "GO_TEST"
	_, err = szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
	require.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printResult(test, "     Add", actual)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	require.NoError(test, err)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
	configHandle2, err := szConfig.ImportConfig(ctx, configDefinition)
	require.NoError(test, err)
	err = szConfig.DeleteDataSource(ctx, configHandle2, dataSourceCode)
	require.NoError(test, err)
	actual, err = szConfig.GetDataSources(ctx, configHandle2)
	require.NoError(test, err)
	printResult(test, "  Delete", actual)
	err = szConfig.CloseConfig(ctx, configHandle2)
	require.NoError(test, err)
}

func TestSzconfig_DeleteDataSource_badConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	dataSourceCode := "GO_TEST"
	err := szConfig.DeleteDataSource(ctx, badConfigHandle, dataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBase)
}

func TestSzconfig_DeleteDataSource_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	err = szConfig.DeleteDataSource(ctx, configHandle, badDataSourceCode)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
}

func TestSzconfig_ExportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	actual, err := szConfig.ExportConfig(ctx, configHandle)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_ExportConfig_badConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.ExportConfig(ctx, badConfigHandle)
	assert.Equal(test, "", actual)
	require.ErrorIs(test, err, szerror.ErrSzBase)
}

func TestSzconfig_GetDataSources(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
}

func TestSzconfig_GetDataSources_badConfigHandle(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	actual, err := szConfig.GetDataSources(ctx, badConfigHandle)
	assert.Equal(test, "", actual)
	require.ErrorIs(test, err, szerror.ErrSzBase)
}

func TestSzconfig_ImportConfig(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	require.NoError(test, err)
	actual, err := szConfig.ImportConfig(ctx, configDefinition)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzconfig_ImportConfig_badConfigDefinition(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_, err := szConfig.ImportConfig(ctx, badConfigDefinition)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzconfig_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzconfig_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
}

func TestSzconfig_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szConfig.SetObserverOrigin(ctx, origin)
	actual := szConfig.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
}

func TestSzconfig_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzconfig_AsInterface(test *testing.T) {
	ctx := context.TODO()
	szConfig := getSzConfigAsInterface(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	require.NoError(test, err)
	actual, err := szConfig.GetDataSources(ctx, configHandle)
	require.NoError(test, err)
	printActual(test, actual)
	err = szConfig.CloseConfig(ctx, configHandle)
	require.NoError(test, err)
}

func TestSzconfig_Initialize(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	settings, err := getSettings()
	require.NoError(test, err)
	err = szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

func TestSzconfig_Initialize_badSettings(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.Initialize(ctx, instanceName, badSettings, verboseLogging)
	assert.NoError(test, err)
}

// TODO: Implement TestSzconfig_Initialize_error
// func TestSzconfig_Initialize_error(test *testing.T) {}

func TestSzconfig_Initialize_again(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	settings, err := getSettings()
	require.NoError(test, err)
	err = szConfig.Initialize(ctx, instanceName, settings, verboseLogging)
	require.NoError(test, err)
}

func TestSzconfig_Destroy(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	err := szConfig.Destroy(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzconfig_Destroy_error
// func TestSzconfig_Destroy_error(test *testing.T) {}

func TestSzconfig_Destroy_withObserver(test *testing.T) {
	ctx := context.TODO()
	szConfigSingleton = nil
	szConfig := getTestObject(ctx, test)
	err := szConfig.Destroy(ctx)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getGrpcConnection() *grpc.ClientConn {
	var err error
	if grpcConnection == nil {
		grpcConnection, err = grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Did not connect: %v\n", err)
		}
		//		defer grpcConnection.Close()
	}
	return grpcConnection
}

func getSettings() (string, error) {
	return "{}", nil
}

func getSzConfig(ctx context.Context) (*Szconfig, error) {
	var err error
	if szConfigSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			return szConfigSingleton, fmt.Errorf("getSettings() Error: %w", err)
		}
		grpcConnection := getGrpcConnection()
		szConfigSingleton = &Szconfig{
			GrpcClient: szpb.NewSzConfigClient(grpcConnection),
		}
		err = szConfigSingleton.SetLogLevel(ctx, logLevel)
		if err != nil {
			return szConfigSingleton, fmt.Errorf("SetLogLevel() Error: %w", err)
		}
		if logLevel == "TRACE" {
			szConfigSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szConfigSingleton.RegisterObserver(ctx, observerSingleton)
			if err != nil {
				return szConfigSingleton, fmt.Errorf("RegisterObserver() Error: %w", err)
			}
			err = szConfigSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			if err != nil {
				return szConfigSingleton, fmt.Errorf("SetLogLevel() - 2 Error: %w", err)
			}
		}
		err = szConfigSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			return szConfigSingleton, fmt.Errorf("Initialize() Error: %w", err)
		}
	}
	return szConfigSingleton, err
}

func getSzConfigAsInterface(ctx context.Context) senzing.SzConfig {
	result, err := getSzConfig(ctx)
	if err != nil {
		panic(err)
	}
	return result
}

func getTestObject(ctx context.Context, test *testing.T) *Szconfig {
	result, err := getSzConfig(ctx)
	require.NoError(test, err)
	return result
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
