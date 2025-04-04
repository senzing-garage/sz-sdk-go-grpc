//go:build linux

package szconfigmanager_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	"google.golang.org/grpc"
)

var (
	grpcAddress    = "0.0.0.0:8261"
	grpcConnection *grpc.ClientConn
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzconfigmanager_AddConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		handleError(err)
	}
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	configComment := "Example configuration"
	configID, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		handleError(err)
	}
	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_GetConfig() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
	}
	configDefinition, err := szConfigManager.GetConfig(ctx, configID)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Truncate(configDefinition, 10))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_CLASS":"ADDRESS","ATTR_CODE":"ADDR_CITY","ATTR_ID":1608,"DEFAULT_VALUE":null,"FELEM_CODE":"CITY","FELEM_REQ":"Any",...
}

func ExampleSzconfigmanager_GetConfigs() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	configList, err := szConfigManager.GetConfigs(ctx)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Truncate(configList, 3))
	// Output: {"CONFIGS":[{...
}

func ExampleSzconfigmanager_GetDefaultConfigID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
	}
	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzconfigmanager_ReplaceDefaultConfigID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfig, err := szAbstractFactory.CreateConfig(ctx)
	if err != nil {
		handleError(err)
	}
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	currentDefaultConfigID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		handleError(err)
	}
	configDefinition, err := szConfigManager.GetConfig(ctx, currentDefaultConfigID)
	if err != nil {
		handleError(err)
	}
	configHandle, err := szConfig.ImportConfig(ctx, configDefinition)
	if err != nil {
		handleError(err)
	}
	_, err = szConfig.AddDataSource(ctx, configHandle, "XXX")
	if err != nil {
		handleError(err)
	}
	newConfigDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		handleError(err)
	}
	err = szConfig.CloseConfig(ctx, configHandle)
	if err != nil {
		handleError(err)
	}
	newConfigID, err := szConfigManager.AddConfig(ctx, newConfigDefinition, "Command")
	if err != nil {
		handleError(err)
	}
	err = szConfigManager.ReplaceDefaultConfigID(ctx, currentDefaultConfigID, newConfigID)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzconfigmanager_SetDefaultConfigID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		handleError(err)
	}
	configID, err := szConfigManager.GetDefaultConfigID(ctx) // For example purposes only. Normally would use output from GetConfigList()
	if err != nil {
		handleError(err)
	}
	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	if err != nil {
		handleError(err)
	}
	// Output:
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzconfigmanager_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	err := szConfigManager.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzconfigmanager_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmanager_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzconfigmanager_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szconfigmanager/szconfigmananger_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	result := szConfigManager.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

// ----------------------------------------------------------------------------
// Helper functions
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

func getSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	_ = ctx
	return &szabstractfactory.Szabstractfactory{
		GrpcConnection: getGrpcConnection(),
	}
}

func getSzConfigManager(ctx context.Context) *szconfigmanager.Szconfigmanager {
	_ = ctx
	return &szconfigmanager.Szconfigmanager{
		GrpcClient: szconfigmanagerpb.NewSzConfigManagerClient(getGrpcConnection()),
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}
