//go:build linux

package g2engine

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2engine_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleG2engine_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2config/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	result := g2engine.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleG2engine_AddRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	loadID := "G2Engine_test"
	err := g2engine.AddRecord(ctx, dataSourceCode, recordID, jsonData, loadID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_AddRecord_secondRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1002"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`
	loadID := "G2Engine_test"
	err := g2engine.AddRecord(ctx, dataSourceCode, recordID, jsonData, loadID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_AddRecordWithInfo() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "PRIMARY_NAME_MIDDLE": "J", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "4/9/16", "STATUS": "Inactive", "AMOUNT": "300"}`
	loadID := "G2Engine_test"
	flags := int64(0)
	result, err := g2engine.AddRecordWithInfo(ctx, dataSourceCode, recordID, jsonData, loadID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_CloseExport() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	g2engine.CloseExport(ctx, responseHandle)
	// Output:
}

func ExampleG2engine_CountRedoRecords() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.CountRedoRecords(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: 1
}

func ExampleG2engine_ExportCSVEntityReport() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	csvColumnList := ""
	flags := int64(0)
	responseHandle, err := g2engine.ExportCSVEntityReport(ctx, csvColumnList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseHandle > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_ExportConfig() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.ExportConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 42))
	// Output: {"G2_CONFIG":{"CFG_ETYPE":[{"ETYPE_ID":...
}

func ExampleG2engine_ExportConfigAndConfigID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	_, configId, err := g2engine.ExportConfigAndConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configId > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_ExportJSONEntityReport() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseHandle > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_FetchNext() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	anEntity, _ := g2engine.FetchNext(ctx, responseHandle)
	fmt.Println(len(anEntity) >= 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_FindInterestingEntitiesByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.FindInterestingEntitiesByEntityID(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_FindInterestingEntitiesByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := int64(0)
	result, err := g2engine.FindInterestingEntitiesByRecordID(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_FindNetworkByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	maxDegree := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	result, err := g2engine.FindNetworkByEntityID(ctx, entityList, maxDegree, buildOutDegree, maxEntities)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 175))
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","RECORD_COUNT":3,"FIRST_SEEN_DT":...
}

func ExampleG2engine_FindNetworkByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	maxDegree := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	flags := int64(0)
	result, err := g2engine.FindNetworkByEntityID_V2(ctx, entityList, maxDegree, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindNetworkByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`
	maxDegree := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	result, err := g2engine.FindNetworkByRecordID(ctx, recordList, maxDegree, buildOutDegree, maxEntities)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 175))
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","RECORD_COUNT":3,"FIRST_SEEN_DT":...
}

func ExampleG2engine_FindNetworkByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`
	maxDegree := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	flags := int64(0)
	result, err := g2engine.FindNetworkByRecordID_V2(ctx, recordList, maxDegree, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	result, err := g2engine.FindPathByEntityID(ctx, entityID1, entityID2, maxDegree)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2engine_FindPathByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	flags := int64(0)
	result, err := g2engine.FindPathByEntityID_V2(ctx, entityID1, entityID2, maxDegree, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := int64(1)
	result, err := g2engine.FindPathByRecordID(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 87))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":...
}

func ExampleG2engine_FindPathByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := int64(1)
	flags := int64(0)
	result, err := g2engine.FindPathByRecordID_V2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathExcludingByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	result, err := g2engine.FindPathExcludingByEntityID(ctx, entityID1, entityID2, maxDegree, excludedEntities)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2engine_FindPathExcludingByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	flags := int64(0)
	result, err := g2engine.FindPathExcludingByEntityID_V2(ctx, entityID1, entityID2, maxDegree, excludedEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathExcludingByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := int64(1)
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003"}]}`
	result, err := g2engine.FindPathExcludingByRecordID(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedRecords)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2engine_FindPathExcludingByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := int64(1)
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003"}]}`
	flags := int64(0)
	result, err := g2engine.FindPathExcludingByRecordID_V2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedRecords, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathIncludingSourceByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	result, err := g2engine.FindPathIncludingSourceByEntityID(ctx, entityID1, entityID2, maxDegree, excludedEntities, requiredDsrcs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 106))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2engine_FindPathIncludingSourceByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	flags := int64(0)
	result, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, entityID1, entityID2, maxDegree, excludedEntities, requiredDsrcs, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathIncludingSourceByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	result, err := g2engine.FindPathIncludingSourceByRecordID(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedEntities, requiredDsrcs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 119))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2engine_FindPathIncludingSourceByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	flags := int64(0)
	result, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedEntities, requiredDsrcs, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_GetActiveConfigID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetActiveConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_GetEntityByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	result, err := g2engine.GetEntityByEntityID(ctx, entityID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2engine_GetEntityByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.GetEntityByEntityID_V2(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2engine_GetEntityByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	result, err := g2engine.GetEntityByRecordID(ctx, dataSourceCode, recordID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 35))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2engine_GetEntityByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := int64(0)
	result, err := g2engine.GetEntityByRecordID_V2(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2engine_GetRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	result, err := g2engine.GetRecord(ctx, dataSourceCode, recordID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","JSON_DATA":{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","RECORD_TYPE":"PERSON","PRIMARY_NAME_LAST":"Smith","PRIMARY_NAME_FIRST":"Robert","DATE_OF_BIRTH":"12/11/1978","ADDR_TYPE":"MAILING","ADDR_LINE1":"123 Main Street, Las Vegas NV 89132","PHONE_TYPE":"HOME","PHONE_NUMBER":"702-919-1300","EMAIL_ADDRESS":"bsmith@work.com","DATE":"1/2/18","STATUS":"Active","AMOUNT":"100"}}
}

func ExampleG2engine_GetRecord_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := int64(0)
	result, err := g2engine.GetRecord_V2(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}
}

func ExampleG2engine_GetRedoRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetRedoRecord(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"REASON":"deferred delete","DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002","ENTITY_TYPE":"GENERIC","DSRC_ACTION":"X"}
}

func ExampleG2engine_GetRepositoryLastModifiedTime() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetRepositoryLastModifiedTime(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_GetVirtualEntityByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	result, err := g2engine.GetVirtualEntityByRecordID(ctx, recordList)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":...
}

func ExampleG2engine_GetVirtualEntityByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	flags := int64(0)
	result, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, recordList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":2}}
}

func ExampleG2engine_HowEntityByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	result, err := g2engine.HowEntityByEntityID(ctx, entityID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[{"STEP":1,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V8","MEMBER_RECORDS":[{"INTERNAL_ID":8,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2","RESULT_VIRTUAL_ENTITY_ID":"V2-S1","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+PHONE","ERRULE_CODE":"CNAME_CFF_CEXCL","FEATURE_SCORES":{"ADDRESS":[{"INBOUND_FEAT_ID":20,"INBOUND_FEAT":"1515 Adela Lane Las Vegas NV 89111","INBOUND_FEAT_USAGE_TYPE":"HOME","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT":"123 Main Street, Las Vegas NV 89132","CANDIDATE_FEAT_USAGE_TYPE":"MAILING","FULL_SCORE":42,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"FF"}],"DOB":[{"INBOUND_FEAT_ID":19,"INBOUND_FEAT":"11/12/1978","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":95,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"FMES"}],"NAME":[{"INBOUND_FEAT_ID":18,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":1,"CANDIDATE_FEAT":"Robert Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":97,"GNR_SN":100,"GNR_GN":95,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"}],"PHONE":[{"INBOUND_FEAT_ID":4,"INBOUND_FEAT":"702-919-1300","INBOUND_FEAT_USAGE_TYPE":"MOBILE","CANDIDATE_FEAT_ID":4,"CANDIDATE_FEAT":"702-919-1300","CANDIDATE_FEAT_USAGE_TYPE":"HOME","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FF"}],"RECORD_TYPE":[{"INBOUND_FEAT_ID":16,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}]}}},{"STEP":2,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V2-S1","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":8,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V100001","MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2-S1","RESULT_VIRTUAL_ENTITY_ID":"V2-S2","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+EMAIL","ERRULE_CODE":"SF1_PNAME_CSTAB","FEATURE_SCORES":{"DOB":[{"INBOUND_FEAT_ID":2,"INBOUND_FEAT":"12/11/1978","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FMES"}],"EMAIL":[{"INBOUND_FEAT_ID":5,"INBOUND_FEAT":"bsmith@work.com","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT":"bsmith@work.com","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"F1"}],"NAME":[{"INBOUND_FEAT_ID":18,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":93,"GNR_SN":100,"GNR_GN":93,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"},{"INBOUND_FEAT_ID":1,"INBOUND_FEAT":"Robert Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":90,"GNR_SN":100,"GNR_GN":88,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"}],"RECORD_TYPE":[{"INBOUND_FEAT_ID":16,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}]}}}],"FINAL_STATE":{"NEED_REEVALUATION":1,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V2-S2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":8,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]}]}}}
}

func ExampleG2engine_HowEntityByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.HowEntityByEntityID_V2(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[{"STEP":1,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V8","MEMBER_RECORDS":[{"INTERNAL_ID":8,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2","RESULT_VIRTUAL_ENTITY_ID":"V2-S1","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+PHONE","ERRULE_CODE":"CNAME_CFF_CEXCL"}},{"STEP":2,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V2-S1","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":8,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V100001","MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2-S1","RESULT_VIRTUAL_ENTITY_ID":"V2-S2","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+EMAIL","ERRULE_CODE":"SF1_PNAME_CSTAB"}}],"FINAL_STATE":{"NEED_REEVALUATION":1,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V2-S2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":8,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]}]}}}
}

func ExampleG2engine_PrimeEngine() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.PrimeEngine(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_SearchByAttributes() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	jsonData := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	result, err := g2engine.SearchByAttributes(ctx, jsonData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 1962))
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PNAME+EMAIL","ERRULE_CODE":"SF1","FEATURE_SCORES":{"EMAIL":[{"INBOUND_FEAT":"bsmith@work.com","CANDIDATE_FEAT":"bsmith@work.com","FULL_SCORE":100}],"NAME":[{"INBOUND_FEAT":"Smith","CANDIDATE_FEAT":"Bob J Smith","GNR_FN":83,"GNR_SN":100,"GNR_GN":40,"GENERATION_MATCH":-1,"GNR_ON":-1},{"INBOUND_FEAT":"Smith","CANDIDATE_FEAT":"Robert Smith","GNR_FN":88,"GNR_SN":100,"GNR_GN":40,"GENERATION_MATCH":-1,"GNR_ON":-1}]}},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","FEATURES":{"ADDRESS":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":20,"USAGE_TYPE":"HOME","FEAT_DESC_VALUES":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":20}]},{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":3,"USAGE_TYPE":"MAILING","FEAT_DESC_VALUES":[{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":3}]}],"DOB":[{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":2},{"FEAT_DESC":"11/12/1978","LIB_FEAT_ID":19}]}],"EMAIL":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":5}]}],"NAME":[{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":1,"USAGE_TYPE":"PRIMARY","FEAT_DESC_VALUES":[{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":1},{"FEAT_DESC":"Bob J Smith","LIB_FEAT_ID":32},{"FEAT_DESC":"Bob Smith","LIB_FEAT_ID":18}]}],"PHONE":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4,"USAGE_TYPE":"HOME","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}]},{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4,"USAGE_TYPE":"MOBILE","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}]}],"RECORD_TYPE":[{"FEAT_DESC":"PERSON","LIB_FEAT_ID":16,"FEAT_DESC_VALUES":[{"FEAT_DESC":"PERSON","LIB_FEAT_ID":16}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","RECORD_COUNT":3,"FIRST_SEEN_DT":...
}

func ExampleG2engine_SearchByAttributes_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	jsonData := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	flags := int64(0)
	result, err := g2engine.SearchByAttributes_V2(ctx, jsonData, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PNAME+EMAIL","ERRULE_CODE":"SF1"},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1}}}]}
}

func ExampleG2engine_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go/blob/main/g2config/g2config_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_Stats() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.Stats(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 16))
	// Output: { "workload":...
}

// FIXME: Remove after GDEV-3576 is fixed
func ExampleG2engine_WhyEntities() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	result, err := g2engine.WhyEntities(ctx, entityID1, entityID2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 74))
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":1,"MATCH_INFO":{"WHY_KEY":...
}

// FIXME: Remove after GDEV-3576 is fixed
func ExampleG2engine_WhyEntities_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityId(truthset.CustomerRecords["1001"])
	entityID2 := getEntityId(truthset.CustomerRecords["1002"])
	flags := int64(0)
	result, err := g2engine.WhyEntities_V2(ctx, entityID1, entityID2, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":1,"MATCH_INFO":{"WHY_KEY":"+NAME+DOB+ADDRESS+PHONE+EMAIL","WHY_ERRULE_CODE":"SF1_SNAME_CFF_CSTAB","MATCH_LEVEL_CODE":"RESOLVED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_WhyRecords() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	result, err := g2engine.WhyRecords(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 115))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":8,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],...
}

func ExampleG2engine_WhyRecords_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	flags := int64(0)
	result, err := g2engine.WhyRecords_V2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":8,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],"INTERNAL_ID_2":2,"ENTITY_ID_2":1,"FOCUS_RECORDS_2":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}],"MATCH_INFO":{"WHY_KEY":"+NAME+DOB+PHONE","WHY_ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_LEVEL_CODE":"RESOLVED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_ReevaluateEntity() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	err := g2engine.ReevaluateEntity(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
func ExampleG2engine_ReevaluateEntityWithInfo() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.ReevaluateEntityWithInfo(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":2}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_ReevaluateRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := int64(0)
	err := g2engine.ReevaluateRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_ReevaluateRecordWithInfo() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := int64(0)
	result, err := g2engine.ReevaluateRecordWithInfo(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":2}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_ReplaceRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	loadID := "G2Engine_test"
	err := g2engine.ReplaceRecord(ctx, dataSourceCode, recordID, jsonData, loadID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_ReplaceRecordWithInfo() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	loadID := "G2Engine_test"
	flags := int64(0)
	result, err := g2engine.ReplaceRecordWithInfo(ctx, dataSourceCode, recordID, jsonData, loadID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_DeleteRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	loadID := "G2Engine_test"
	err := g2engine.DeleteRecord(ctx, dataSourceCode, recordID, loadID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_DeleteRecordWithInfo() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	loadID := "G2Engine_test"
	flags := int64(0)
	result, err := g2engine.DeleteRecordWithInfo(ctx, dataSourceCode, recordID, loadID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_Init() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	moduleName := "Test module name"
	iniParams := "{}"
	verboseLogging := int64(0)
	err := g2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		// This should produce a "senzing-60144002" error.
	}
	// Output:
}

func ExampleG2engine_InitWithConfigID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	moduleName := "Test module name"
	iniParams := "{}"
	initConfigID := int64(1)
	verboseLogging := int64(0)
	err := g2engine.InitWithConfigID(ctx, moduleName, iniParams, initConfigID, verboseLogging)
	if err != nil {
		// This should produce a "senzing-60144003" error.
	}
	// Output:
}

func ExampleG2engine_Reinit() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	initConfigID, _ := g2engine.GetActiveConfigID(ctx) // Example initConfigID.
	err := g2engine.Reinit(ctx, initConfigID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_Destroy() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-grpc/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.Destroy(ctx)
	if err != nil {
		// This should produce a "senzing-60164001" error.
	}
	// Output:
}
