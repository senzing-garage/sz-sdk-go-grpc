//go:build linux

package szconfigmanager

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzConfigManager_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzConfigManager_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2config/g2configmgr_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szConfigManager.SetObserverOrigin(ctx, origin)
	result := szConfigManager.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleSzConfigManager_AddConfig() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	szConfigManager := getSzConfigManager(ctx)
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	configComment := "Example configuration"
	configId, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configId > 0) // Dummy output.
	// Output: true
}

func ExampleSzConfigManager_GetConfig() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	configId, err := szConfigManager.GetDefaultConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configDefinition, err := szConfigManager.GetConfig(ctx, configId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(configDefinition, defaultTruncation))
	// Output: {"G2_CONFIG":{"CFG_ATTR":[{"ATTR_ID":1001,"ATTR_CODE":"DATA_SOURCE","ATTR...
}

func ExampleSzConfigManager_GetConfigList() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	configList, err := szConfigManager.GetConfigList(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(configList, 28))
	// Output: {"CONFIGS":[{"CONFIG_ID":...
}

func ExampleSzConfigManager_GetDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	configID, err := szConfigManager.GetDefaultConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configID > 0) // Dummy output.
	// Output: true
}

func ExampleSzConfigManager_ReplaceDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	currentDefaultConfigId, err := szConfigManager.GetDefaultConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	szConfig := getSzConfig(ctx)
	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		fmt.Println(err)
	}
	configComment := "Example configuration"
	newDefaultConfigId, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		fmt.Println(err)
	}
	err = szConfigManager.ReplaceDefaultConfigId(ctx, currentDefaultConfigId, newDefaultConfigId)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzConfigManager_SetDefaultConfigId() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	configId, err := szConfigManager.GetDefaultConfigId(ctx) // For example purposes only. Normally would use output from GetConfigList()
	if err != nil {
		fmt.Println(err)
	}
	err = szConfigManager.SetDefaultConfigId(ctx, configId)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzConfigManager_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	err := szConfigManager.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

// func ExampleSzConfigManager_Init() {
// 	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_examples_test.go
// 	ctx := context.TODO()
// 	grpcConnection := getGrpcConnection()
// 	szConfigManager := &SzConfigManager{
// 		GrpcClient: szpb.NewG2ConfigMgrClient(grpcConnection),
// 	}
// 	moduleName := "Test module name"
// 	iniParams := "{}"
// 	verboseLogging := int64(0)
// 	err := szConfigManager.Init(ctx, moduleName, iniParams, verboseLogging)
// 	if err != nil {
// 		// This should produce a "senzing-60124002" error.
// 	}
// 	// Output:
// }

func ExampleSzConfigManager_Destroy() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2configmgr/g2configmgr_examples_test.go
	ctx := context.TODO()
	szConfigManager := getSzConfigManager(ctx)
	err := szConfigManager.Destroy(ctx)
	if err != nil {
		// This should produce a "senzing-60124001" error.
	}
	// Output:
}
