//go:build linux

package szengine_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go-grpc/szengine"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szenginepb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GetEntityByRecordIDResponse struct {
	ResolvedEntity struct {
		EntityID int64 `json:"ENTITY_ID"`
	} `json:"RESOLVED_ENTITY"`
}

const (
	baseTen = 10
)

var (
	grpcAddress = "localhost:8261"
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzengine() {}

func ExampleSzengine_AddRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	flags := senzing.SzWithoutInfo
	result, err := szEngine.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleSzengine_AddRecord_secondRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "CUSTOMERS"
	recordID := "1002"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`
	flags := senzing.SzWithoutInfo
	result, err := szEngine.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleSzengine_AddRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "PRIMARY_NAME_MIDDLE": "J", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "4/9/16", "STATUS": "Inactive", "AMOUNT": "300"}`
	flags := senzing.SzWithInfo
	result, err := szEngine.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[{"ENTITY_ID":100001}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzengine_CloseExport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	flags := senzing.SzNoFlags
	exportHandle, err := szEngine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		handleError(err)
	}
	err = szEngine.CloseExport(ctx, exportHandle)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzengine_CountRedoRecords() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	result, err := szEngine.CountRedoRecords(ctx)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: 4
}

func ExampleSzengine_DeleteRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	flags := senzing.SzWithoutInfo
	result, err := szEngine.DeleteRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleSzengine_DeleteRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	flags := senzing.SzWithInfo
	result, err := szEngine.DeleteRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzengine_ExportCsvEntityReport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	csvColumnList := ""
	flags := senzing.SzNoFlags
	exportHandle, err := szEngine.ExportCsvEntityReport(ctx, csvColumnList, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(exportHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzengine_ExportCsvEntityReportIterator() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	csvColumnList := ""
	flags := senzing.SzNoFlags
	for result := range szEngine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
		if result.Error != nil {
			handleError(err)
			break
		}
		fmt.Println(result.Value)
	}
	// Output: RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID
}

func ExampleSzengine_ExportJSONEntityReport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	flags := senzing.SzNoFlags
	exportHandle, err := szEngine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(exportHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzengine_ExportJSONEntityReportIterator() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	flags := senzing.SzNoFlags
	for result := range szEngine.ExportJSONEntityReportIterator(ctx, flags) {
		if result.Error != nil {
			handleError(err)
			break
		}
		fmt.Println(result.Value)
	}
	// Output:
}

func ExampleSzengine_FetchNext() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	flags := senzing.SzNoFlags
	exportHandle, err := szEngine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		handleError(err)
	}
	defer func() {
		err = szEngine.CloseExport(ctx, exportHandle)
	}()
	jsonEntityReport := ""
	for {
		jsonEntityReportFragment, err := szEngine.FetchNext(ctx, exportHandle)
		if err != nil {
			handleError(err)
		}
		if len(jsonEntityReportFragment) == 0 {
			break
		}
		jsonEntityReport += jsonEntityReportFragment
	}
}

func ExampleSzengine_FindInterestingEntitiesByEntityID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	entityID, err := getEntityIDForRecord("CUSTOMERS", "1001")
	if err != nil {
		handleError(err)
	}
	flags := senzing.SzNoFlags
	result, err := szEngine.FindInterestingEntitiesByEntityID(ctx, entityID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzengine_FindInterestingEntitiesByRecordID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := senzing.SzNoFlags
	result, err := szEngine.FindInterestingEntitiesByRecordID(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzengine_FindNetworkByEntityID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	entityID1, err := getEntityIDStringForRecord("CUSTOMERS", "1001")
	if err != nil {
		handleError(err)
	}
	entityID2, err := getEntityIDStringForRecord("CUSTOMERS", "1002")
	if err != nil {
		handleError(err)
	}
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	maxDegrees := int64(2)
	buildOutDegrees := int64(1)
	maxEntities := int64(10)
	flags := senzing.SzNoFlags
	result, err := szEngine.FindNetworkByEntityID(ctx, entityList, maxDegrees, buildOutDegrees, maxEntities, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleSzengine_FindNetworkByRecordID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`
	maxDegrees := int64(1)
	buildOutDegrees := int64(2)
	maxEntities := int64(10)
	flags := senzing.SzNoFlags
	result, err := szEngine.FindNetworkByRecordID(ctx, recordList, maxDegrees, buildOutDegrees, maxEntities, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleSzengine_FindPathByEntityID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	startEntityID, err := getEntityIDForRecord("CUSTOMERS", "1001")
	if err != nil {
		handleError(err)
	}
	endEntityID, err := getEntityIDForRecord("CUSTOMERS", "1002")
	if err != nil {
		handleError(err)
	}
	maxDegrees := int64(1)
	avoidEntityIDs := ""
	requiredDataSources := ""
	flags := senzing.SzNoFlags
	result, err := szEngine.FindPathByEntityID(ctx, startEntityID, endEntityID, maxDegrees, avoidEntityIDs, requiredDataSources, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[100001]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleSzengine_FindPathByEntityID_avoiding() {
	// TODO: Implement ExampleSzEngine_FindPathByEntityID_avoiding
	// // For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	// ctx := context.TODO()
	// szEngine := getSzEngineExample(ctx)
	// startEntityID := getEntityIDForRecord("CUSTOMERS", "1001")
	// endEntityID := getEntityIDForRecord("CUSTOMERS", "1002")
	// maxDegrees := int64(1)
	// avoidEntityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDStringForRecord("CUSTOMERS", "1003") + `}]}`
	// requiredDataSources := ""
	// flags := senzing.SzNoFlags
	// result, err := szEngine.FindPathByEntityID(ctx, startEntityID, endEntityID, maxDegrees, avoidEntityIDs, requiredDataSources, flags)
	// if err != nil {
	// 	handleError(err)
	// }
	// fmt.Println(truncate(result, 107))
	// // Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngine_FindPathByEntityID_avoidingAndIncluding() {
	// TODO: Implement ExampleSzEngine_FindPathByEntityID_avoidingAndIncluding
}

func ExampleSzengine_FindPathByEntityID_including() {
	// TODO: Implement ExampleSzEngine_FindPathByEntityID_including
	// // For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	// ctx := context.TODO()
	// szEngine := getSzEngineExample(ctx)
	// startEntityID := getEntityIDForRecord("CUSTOMERS", "1001")
	// endEntityID := getEntityIDForRecord("CUSTOMERS", "1002")
	// maxDegree := int64(1)
	// avoidEntityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDStringForRecord("CUSTOMERS", "1003") + `}]}`
	// requiredDataSources := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	// flags := senzing.SzNoFlags
	// result, err := szEngine.FindPathByEntityID(ctx, startEntityID, endEntityID, maxDegree, avoidEntityIDs, requiredDataSources, flags)
	// if err != nil {
	// 	handleError(err)
	// }
	// fmt.Println(truncate(result, 106))
	// // Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzengine_FindPathByRecordID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	startDataSourceCode := "CUSTOMERS"
	startRecordID := "1001"
	endDataSourceCode := "CUSTOMERS"
	endRecordID := "1002"
	maxDegrees := int64(1)
	avoidRecordKeys := ""
	requiredDataSources := ""
	flags := senzing.SzNoFlags
	result, err := szEngine.FindPathByRecordID(ctx, startDataSourceCode, startRecordID, endDataSourceCode, endRecordID, maxDegrees, avoidRecordKeys, requiredDataSources, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[100001]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleSzengine_FindPathByRecordID_avoiding() {
	// TODO: Implement ExampleSzEngine_FindPathByRecordID_avoiding
	// // For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	// ctx := context.TODO()
	// szEngine := getSzEngineExample(ctx)
	// startDataSourceCode := "CUSTOMERS"
	// startRecordID := "1001"
	// endDataSourceCode := "CUSTOMERS"
	// endRecordID := "1002"
	// maxDegree := int64(1)
	// avoidRecordKeys := `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003"}]}`
	// requiredDataSources := ""
	// flags := senzing.SzNoFlags
	// result, err := szEngine.FindPathByRecordID(ctx, startDataSourceCode, startRecordID, endDataSourceCode, endRecordID, maxDegree, avoidRecordKeys, requiredDataSources, flags)
	// if err != nil {
	// 	handleError(err)
	// }
	// fmt.Println(truncate(result, 107))
	// // Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngine_FindPathByRecordID_avoidingAndIncluding() {
	// TODO: Implement ExampleSzEngine_FindPathByRecordID_avoidingAndIncluding
}

func ExampleSzengine_FindPathByRecordID_including() {
	// TODO: Implement ExampleSzEngine_FindPathByRecordID_including
	// // For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	// ctx := context.TODO()
	// szEngine := getSzEngineExample(ctx)
	// startDataSourceCode := "CUSTOMERS"
	// startRecordID := "1001"
	// endDataSourceCode := "CUSTOMERS"
	// endRecordID := "1002"
	// maxDegrees := int64(1)
	// avoidRecordKeys := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIDStringForRecord("CUSTOMERS", "1003") + `}]}`
	// requiredDataSources := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	// flags := senzing.SzNoFlags
	// result, err := szEngine.FindPathByRecordID(ctx, startDataSourceCode, startRecordID, endDataSourceCode, endRecordID, maxDegrees, avoidRecordKeys, requiredDataSources, flags)
	// if err != nil {
	// 	handleError(err)
	// }
	// fmt.Println(truncate(result, 119))
	// // Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleSzengine_GetActiveConfigID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	result, err := szEngine.GetActiveConfigID(ctx)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleSzengine_GetEntityByEntityID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	entityID, err := getEntityIDForRecord("CUSTOMERS", "1001")
	if err != nil {
		handleError(err)
	}
	flags := senzing.SzNoFlags
	result, err := szEngine.GetEntityByEntityID(ctx, entityID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":100001}}
}

func ExampleSzengine_GetEntityByRecordID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := senzing.SzNoFlags
	result, err := szEngine.GetEntityByRecordID(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":100001}}
}

func ExampleSzengine_GetRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := senzing.SzNoFlags
	result, err := szEngine.GetRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Flatten(jsonutil.Normalize(result)))
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}
}

func ExampleSzengine_GetRedoRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	result, err := szEngine.GetRedoRecord(ctx)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"REASON":"deferred delete","DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","REEVAL_ITERATION":1,"DSRC_ACTION":"X"}
}

func ExampleSzengine_GetStats() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	result, err := szEngine.GetStats(ctx)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Truncate(result, 5))
	// Output: {"workload":{"abortedUnresolve":0,"actualAmbiguousTest":0,"addedRecords":3,...
}

func ExampleSzengine_GetVirtualEntityByRecordID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	flags := senzing.SzNoFlags
	result, err := szEngine.GetVirtualEntityByRecordID(ctx, recordList, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":100001}}
}

func ExampleSzengine_HowEntityByEntityID() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	entityID, err := getEntityIDForRecord("CUSTOMERS", "1001")
	if err != nil {
		handleError(err)
	}
	flags := senzing.SzNoFlags
	result, err := szEngine.HowEntityByEntityID(ctx, entityID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Flatten(jsonutil.NormalizeAndSort(jsonutil.Flatten(jsonutil.Redact(result, "RECORD_ID", "INBOUND_FEAT_USAGE_TYPE")))))
	// Output: {"HOW_RESULTS":{"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100001-S1"}]},"RESOLUTION_STEPS":[{"INBOUND_VIRTUAL_ENTITY_ID":"V100002","MATCH_INFO":{"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE"},"RESULT_VIRTUAL_ENTITY_ID":"V100001-S1","STEP":1,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100001"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100002"}}]}}
}

func ExampleSzengine_PreprocessRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	flags := senzing.SzWithoutInfo
	result, err := szEngine.PreprocessRecord(ctx, recordDefinition, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleSzengine_PrimeEngine() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	err = szEngine.PrimeEngine(ctx)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzEngine_ProcessRedoRecord() {
	// TODO: Uncomment after it has been implemented.
	// // For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	// ctx := context.TODO()
	// szEngine := getSzEngineExample(ctx)
	// redoRecord, err := szEngine.GetRedoRecord(ctx)
	// if err != nil {
	// 	handleError(err)
	// }
	// flags := senzing.SzWithoutInfo
	// result, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
	// if err != nil {
	// 	handleError(err)
	// }
	// fmt.Println(result)
	// // Output:
}

func ExampleSzEngine_ProcessRedoRecord_withInfo() {
	// TODO: Uncomment after it has been implemented.
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	// ctx := context.TODO()
	// szEngine := getSzEngineExample(ctx)
	// redoRecord, err := szEngine.GetRedoRecord(ctx)
	// if err != nil {
	// 	handleError(err)
	// }
	// flags := senzing.SzWithInfo
	// result, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
	// if err != nil {
	// 	handleError(err)
	// }
	// fmt.Println(result)
	// // Output: {}
}

func ExampleSzengine_ReevaluateEntity() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	entityID, err := getEntityIDForRecord("CUSTOMERS", "1001")
	if err != nil {
		handleError(err)
	}
	flags := senzing.SzWithoutInfo
	result, err := szEngine.ReevaluateEntity(ctx, entityID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleSzengine_ReevaluateEntity_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	entityID, err := getEntityIDForRecord("CUSTOMERS", "1001")
	if err != nil {
		handleError(err)
	}
	flags := senzing.SzWithInfo
	result, err := szEngine.ReevaluateEntity(ctx, entityID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":100001}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzengine_ReevaluateRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := senzing.SzWithoutInfo
	result, err := szEngine.ReevaluateRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleSzengine_ReevaluateRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := senzing.SzWithInfo
	result, err := szEngine.ReevaluateRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":100001}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzengine_SearchByAttributes() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	searchProfile := ""
	flags := senzing.SzNoFlags
	result, err := szEngine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Flatten(jsonutil.Redact(jsonutil.Flatten(jsonutil.NormalizeAndSort(result)), "FIRST_SEEN_DT", "LAST_SEEN_DT")))
	// Output: {"RESOLVED_ENTITIES":[{"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":100001}},"MATCH_INFO":{"ERRULE_CODE":"SF1","MATCH_KEY":"+PNAME+EMAIL","MATCH_LEVEL_CODE":"POSSIBLY_RELATED"}}]}
}

func ExampleSzEngine_SearchByAttributes_searchProfile() {
	// TODO: Implement ExampleSzEngine_SearchByAttributes_searchProfile
}

func ExampleSzengine_WhyEntities() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	entityID1, err := getEntityID(truthset.CustomerRecords["1001"])
	if err != nil {
		handleError(err)
	}
	entityID2, err := getEntityID(truthset.CustomerRecords["1002"])
	if err != nil {
		handleError(err)
	}
	flags := senzing.SzNoFlags
	result, err := szEngine.WhyEntities(ctx, entityID1, entityID2, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":100001,"ENTITY_ID_2":100001,"MATCH_INFO":{"WHY_KEY":"+NAME+DOB+ADDRESS+PHONE+EMAIL","WHY_ERRULE_CODE":"SF1_SNAME_CFF_CSTAB","MATCH_LEVEL_CODE":"RESOLVED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleSzengine_WhyRecordInEntity() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := senzing.SzNoFlags
	result, err := szEngine.WhyRecordInEntity(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":100001,"ENTITY_ID":100001,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],"MATCH_INFO":{"WHY_KEY":"+NAME+DOB+PHONE","WHY_ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_LEVEL_CODE":"RESOLVED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleSzengine_WhyRecords() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := getSzAbstractFactory(ctx)
	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
	}
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	flags := senzing.SzNoFlags
	result, err := szEngine.WhyRecords(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, flags)
	if err != nil {
		handleError(err)
	}
	fmt.Println(jsonutil.Truncate(result, 121))
	// Output: {"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}],"WHY_RESULTS":[{"ENTITY_ID":100001,"ENTITY_ID_2":100001,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],"FOCUS_RECORDS_2":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}],"INTERNAL_ID":100001,"INTERNAL_ID_2":100002,"MATCH_INFO":{"MATCH_LEVEL_CODE":"RESOLVED","WHY_ERRULE_CODE":"CNAME_CFF_CEXCL","WHY_KEY":"+NAME+DOB+PHONE"}}]}
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzengine_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szenzing/szengine_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	err := szEngine.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		handleError(err)
	}
	// Output:
}

func ExampleSzengine_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzengine_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szenzing/szengine_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
	result := szEngine.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func getEntityID(record record.Record) (int64, error) {
	return getEntityIDForRecord(record.DataSource, record.ID)
}

func getEntityIDForRecord(datasource string, id string) (int64, error) {
	var result int64
	var err error
	ctx := context.TODO()
	// var result int64
	szEngine := getSzEngine(ctx)
	response, err := szEngine.GetEntityByRecordID(ctx, datasource, id, senzing.SzWithoutInfo)
	if err != nil {
		return result, err
	}
	getEntityByRecordIDResponse := &GetEntityByRecordIDResponse{}
	err = json.Unmarshal([]byte(response), &getEntityByRecordIDResponse)
	if err != nil {
		return result, err
	}

	result = getEntityByRecordIDResponse.ResolvedEntity.EntityID
	return result, err
}

func getEntityIDStringForRecord(datasource string, id string) (string, error) {
	var result string
	var err error
	entityID, err := getEntityIDForRecord(datasource, id)
	if err != nil {
		return result, err
	}
	result = strconv.FormatInt(entityID, baseTen)
	return result, err
}

func getSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	var err error
	var result senzing.SzAbstractFactory
	_ = ctx
	grpcConnection, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	result = &szabstractfactory.Szabstractfactory{
		GrpcConnection: grpcConnection,
	}
	return result
}

func getSzEngine(ctx context.Context) *szengine.Szengine {
	_ = ctx
	grpcConnection, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	result := &szengine.Szengine{
		GrpcClient: szenginepb.NewSzEngineClient(grpcConnection),
	}
	return result
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}
