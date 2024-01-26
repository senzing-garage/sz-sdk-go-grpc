//go:build linux

package g2diagnostic

import (
	"context"
	"fmt"

	g2pb "github.com/senzing-garage/g2-sdk-proto/go/g2diagnostic"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2diagnostic_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2diagnostic.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleG2diagnostic_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2config/g2diagnostic_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2diagnostic.SetObserverOrigin(ctx, origin)
	result := g2diagnostic.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleG2diagnostic_CheckDBPerf() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
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

// func ExampleG2diagnostic_FetchNextEntityBySize() {
// 	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
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

func ExampleG2diagnostic_GetAvailableMemory() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	result, err := g2diagnostic.GetAvailableMemory(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleG2diagnostic_GetDBInfo() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	result, err := g2diagnostic.GetDBInfo(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 52))
	// Output: {"Hybrid Mode":false,"Database Details":[{"Name":...
}

func ExampleG2diagnostic_GetLogicalCores() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	result, err := g2diagnostic.GetLogicalCores(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleG2diagnostic_GetPhysicalCores() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	result, err := g2diagnostic.GetPhysicalCores(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleG2diagnostic_GetTotalSystemMemory() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	result, err := g2diagnostic.GetTotalSystemMemory(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

// func ExampleG2diagnostic_SetLogLevel() {
// 	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
// 	g2diagnostic := &G2diagnosticClient{}
// 	ctx := context.TODO()
// 	err := g2diagnostic.SetLogLevel(ctx, logger.LevelInfo)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// Output:
// }

func ExampleG2diagnostic_Init() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	grpcConnection := getGrpcConnection()
	g2diagnostic := &G2diagnostic{
		GrpcClient: g2pb.NewG2DiagnosticClient(grpcConnection),
	}
	moduleName := "Test module name"
	iniParams := "{}"
	verboseLogging := int64(0)
	err := g2diagnostic.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		// This should produce a "senzing-60134002" error.
	}
	// Output:
}

func ExampleG2diagnostic_InitWithConfigID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	grpcConnection := getGrpcConnection()
	g2diagnostic := &G2diagnostic{
		GrpcClient: g2pb.NewG2DiagnosticClient(grpcConnection),
	}
	moduleName := "Test module name"
	iniParams := "{}"
	initConfigID := int64(1)
	verboseLogging := int64(0)
	err := g2diagnostic.InitWithConfigID(ctx, moduleName, iniParams, initConfigID, verboseLogging)
	if err != nil {
		// This should produce a "senzing-60134003" error.
	}
	// Output:
}

func ExampleG2diagnostic_Reinit() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
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

func ExampleG2diagnostic_Destroy() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2diagnostic/g2diagnostic_examples_test.go
	ctx := context.TODO()
	g2diagnostic := getG2Diagnostic(ctx)
	err := g2diagnostic.Destroy(ctx)
	if err != nil {
		// This should produce a "senzing-60134001" error.
	}
	// Output:
}
