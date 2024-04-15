//go:build linux

package szengine

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go/sz"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzEngine_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzEngine_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szenzing/szengine_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
	result := szEngine.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleSzEngine_AddRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	flags := sz.SZ_WITHOUT_INFO
	_, err := szEngine.AddRecord(ctx, dataSourceCode, recordId, recordDefinition, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_AddRecord_secondRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1002"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`
	flags := sz.SZ_WITHOUT_INFO
	_, err := szEngine.AddRecord(ctx, dataSourceCode, recordId, recordDefinition, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_AddRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1003"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "PRIMARY_NAME_MIDDLE": "J", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "4/9/16", "STATUS": "Inactive", "AMOUNT": "300"}`
	flags := sz.SZ_WITH_INFO
	result, err := szEngine.AddRecord(ctx, dataSourceCode, recordId, recordDefinition, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzEngine_CloseExport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	flags := sz.SZ_NO_FLAGS
	responseHandle, err := szEngine.ExportJsonEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	szEngine.CloseExport(ctx, responseHandle)
	// Output:
}

func ExampleSzEngine_CountRedoRecords() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	result, err := szEngine.CountRedoRecords(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: 1
}

func ExampleSzEngine_ExportCsvEntityReport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	csvColumnList := ""
	flags := sz.SZ_NO_FLAGS
	responseHandle, err := szEngine.ExportCsvEntityReport(ctx, csvColumnList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngine_ExportJsonEntityReport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	flags := sz.SZ_NO_FLAGS
	responseHandle, err := szEngine.ExportJsonEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngine_FetchNext() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	flags := sz.SZ_NO_FLAGS
	responseHandle, err := szEngine.ExportJsonEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	anEntity, _ := szEngine.FetchNext(ctx, responseHandle)
	fmt.Println(len(anEntity) >= 0) // Dummy output.
	// Output: true
}

func ExampleSzEngine_FindNetworkByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	maxDegrees := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.FindNetworkByEntityId(ctx, entityList, maxDegrees, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 175))
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","RECORD_COUNT":3,"FIRST_SEEN_DT":...
}

func ExampleSzEngine_FindNetworkByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`
	maxDegrees := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.FindNetworkByRecordId(ctx, recordList, maxDegrees, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 175))
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","RECORD_COUNT":3,"FIRST_SEEN_DT":...
}

func ExampleSzEngine_FindPathByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	startEntityId := getEntityIdForRecord("CUSTOMERS", "1001")
	endEntityId := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegrees := int64(1)
	exclusions := ""
	requiredDataSources := ""
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngine_FindPathByEntityId_excluding() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	startEntityId := getEntityIdForRecord("CUSTOMERS", "1001")
	endEntityId := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegrees := int64(1)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDataSources := ""
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngine_FindPathByEntityId_including() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	startEntityId := getEntityIdForRecord("CUSTOMERS", "1001")
	endEntityId := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegree, exclusions, requiredDataSources, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 106))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngine_FindPathByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	startDataSourceCode := "CUSTOMERS"
	startRecordId := "1001"
	endDataSourceCode := "CUSTOMERS"
	endRecordId := "1002"
	maxDegrees := int64(1)
	exclusions := ""
	requiredDataSources := ""
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.FindPathByRecordId(ctx, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, requiredDataSources, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 87))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":...
}

func ExampleSzEngine_FindPathByRecordId_excluding() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	startDataSourceCode := "CUSTOMERS"
	startRecordId := "1001"
	endDataSourceCode := "CUSTOMERS"
	endRecordId := "1002"
	maxDegree := int64(1)
	exclusions := `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003"}]}`
	requiredDataSources := ""
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.FindPathByRecordId(ctx, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegree, exclusions, requiredDataSources, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngine_FindPathByRecordId_including() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	startDataSourceCode := "CUSTOMERS"
	startRecordId := "1001"
	endDataSourceCode := "CUSTOMERS"
	endRecordId := "1002"
	maxDegrees := int64(1)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.FindPathByRecordId(ctx, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, requiredDataSources, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 119))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleSzEngine_GetActiveConfigId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	result, err := szEngine.GetActiveConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngine_GetEntityByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityId := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.GetEntityByEntityId(ctx, entityId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleSzEngine_GetEntityByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.GetEntityByRecordId(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 35))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleSzEngine_GetRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.GetRecord(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","JSON_DATA":{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","RECORD_TYPE":"PERSON","PRIMARY_NAME_LAST":"Smith","PRIMARY_NAME_FIRST":"Robert","DATE_OF_BIRTH":"12/11/1978","ADDR_TYPE":"MAILING","ADDR_LINE1":"123 Main Street, Las Vegas NV 89132","PHONE_TYPE":"HOME","PHONE_NUMBER":"702-919-1300","EMAIL_ADDRESS":"bsmith@work.com","DATE":"1/2/18","STATUS":"Active","AMOUNT":"100"}}
}

func ExampleSzEngine_GetRedoRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	result, err := szEngine.GetRedoRecord(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"REASON":"deferred delete","DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002","ENTITY_TYPE":"GENERIC","DSRC_ACTION":"X"}
}

func ExampleSzEngine_GetRepositoryLastModifiedTime() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	result, err := szEngine.GetRepositoryLastModifiedTime(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngine_GetStats() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	result, err := szEngine.GetStats(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 16))
	// Output: { "workload":...
}

func ExampleSzEngine_GetVirtualEntityByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.GetVirtualEntityByRecordId(ctx, recordList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":...
}

func ExampleSzEngine_HowEntityByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityId := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.HowEntityByEntityId(ctx, entityId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"HOW_RESULTS":{"RESOLUTION_STEPS":[{"STEP":1,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V8","MEMBER_RECORDS":[{"INTERNAL_ID":8,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2","RESULT_VIRTUAL_ENTITY_ID":"V2-S1","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+PHONE","ERRULE_CODE":"CNAME_CFF_CEXCL","FEATURE_SCORES":{"ADDRESS":[{"INBOUND_FEAT_ID":20,"INBOUND_FEAT":"1515 Adela Lane Las Vegas NV 89111","INBOUND_FEAT_USAGE_TYPE":"HOME","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT":"123 Main Street, Las Vegas NV 89132","CANDIDATE_FEAT_USAGE_TYPE":"MAILING","FULL_SCORE":42,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"FF"}],"DOB":[{"INBOUND_FEAT_ID":19,"INBOUND_FEAT":"11/12/1978","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":95,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"FMES"}],"NAME":[{"INBOUND_FEAT_ID":18,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":1,"CANDIDATE_FEAT":"Robert Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":97,"GNR_SN":100,"GNR_GN":95,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"}],"PHONE":[{"INBOUND_FEAT_ID":4,"INBOUND_FEAT":"702-919-1300","INBOUND_FEAT_USAGE_TYPE":"MOBILE","CANDIDATE_FEAT_ID":4,"CANDIDATE_FEAT":"702-919-1300","CANDIDATE_FEAT_USAGE_TYPE":"HOME","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FF"}],"RECORD_TYPE":[{"INBOUND_FEAT_ID":16,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}]}}},{"STEP":2,"VIRTUAL_ENTITY_1":{"VIRTUAL_ENTITY_ID":"V2-S1","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":8,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]}]},"VIRTUAL_ENTITY_2":{"VIRTUAL_ENTITY_ID":"V100001","MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]},"INBOUND_VIRTUAL_ENTITY_ID":"V2-S1","RESULT_VIRTUAL_ENTITY_ID":"V2-S2","MATCH_INFO":{"MATCH_KEY":"+NAME+DOB+EMAIL","ERRULE_CODE":"SF1_PNAME_CSTAB","FEATURE_SCORES":{"DOB":[{"INBOUND_FEAT_ID":2,"INBOUND_FEAT":"12/11/1978","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FMES"}],"EMAIL":[{"INBOUND_FEAT_ID":5,"INBOUND_FEAT":"bsmith@work.com","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT":"bsmith@work.com","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"F1"}],"NAME":[{"INBOUND_FEAT_ID":18,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":93,"GNR_SN":100,"GNR_GN":93,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"},{"INBOUND_FEAT_ID":1,"INBOUND_FEAT":"Robert Smith","INBOUND_FEAT_USAGE_TYPE":"PRIMARY","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GNR_FN":90,"GNR_SN":100,"GNR_GN":88,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"CLOSE","SCORE_BEHAVIOR":"NAME"}],"RECORD_TYPE":[{"INBOUND_FEAT_ID":16,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}]}}}],"FINAL_STATE":{"NEED_REEVALUATION":1,"VIRTUAL_ENTITIES":[{"VIRTUAL_ENTITY_ID":"V2-S2","MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}]},{"INTERNAL_ID":8,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}]},{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003"}]}]}]}}}
}

func ExampleSzEngine_PrimeEngine() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	err := szEngine.PrimeEngine(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_SearchByAttributes() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	searchProfile := ""
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 1962))
	// Output: {"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PNAME+EMAIL","ERRULE_CODE":"SF1","FEATURE_SCORES":{"EMAIL":[{"INBOUND_FEAT":"bsmith@work.com","CANDIDATE_FEAT":"bsmith@work.com","FULL_SCORE":100}],"NAME":[{"INBOUND_FEAT":"Smith","CANDIDATE_FEAT":"Bob J Smith","GNR_FN":83,"GNR_SN":100,"GNR_GN":40,"GENERATION_MATCH":-1,"GNR_ON":-1},{"INBOUND_FEAT":"Smith","CANDIDATE_FEAT":"Robert Smith","GNR_FN":88,"GNR_SN":100,"GNR_GN":40,"GENERATION_MATCH":-1,"GNR_ON":-1}]}},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","FEATURES":{"ADDRESS":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":20,"USAGE_TYPE":"HOME","FEAT_DESC_VALUES":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":20}]},{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":3,"USAGE_TYPE":"MAILING","FEAT_DESC_VALUES":[{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":3}]}],"DOB":[{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":2},{"FEAT_DESC":"11/12/1978","LIB_FEAT_ID":19}]}],"EMAIL":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":5}]}],"NAME":[{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":1,"USAGE_TYPE":"PRIMARY","FEAT_DESC_VALUES":[{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":1},{"FEAT_DESC":"Bob J Smith","LIB_FEAT_ID":32},{"FEAT_DESC":"Bob Smith","LIB_FEAT_ID":18}]}],"PHONE":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4,"USAGE_TYPE":"HOME","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}]},{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4,"USAGE_TYPE":"MOBILE","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}]}],"RECORD_TYPE":[{"FEAT_DESC":"PERSON","LIB_FEAT_ID":16,"FEAT_DESC_VALUES":[{"FEAT_DESC":"PERSON","LIB_FEAT_ID":16}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","RECORD_COUNT":3,"FIRST_SEEN_DT":...
}

func ExampleSzEngine_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szenzing/szengine_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	err := szEngine.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_WhyEntities() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityId1 := getEntityId(truthset.CustomerRecords["1001"])
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.WhyEntities(ctx, entityId1, entityId2, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 74))
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":1,"MATCH_INFO":{"WHY_KEY":...
}

func ExampleSzEngine_WhyRecordInEntity() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.WhyRecordInEntity(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleSzEngine_WhyRecords() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordId1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordId2 := "1002"
	flags := sz.SZ_NO_FLAGS
	result, err := szEngine.WhyRecords(ctx, dataSourceCode1, recordId1, dataSourceCode2, recordId2, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 115))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":8,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],...
}

func ExampleSzEngine_ReevaluateEntity() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityId := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := sz.SZ_WITHOUT_INFO
	_, err := szEngine.ReevaluateEntity(ctx, entityId, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
func ExampleSzEngine_ReevaluateEntity_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := sz.SZ_WITH_INFO
	result, err := szEngine.ReevaluateEntity(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":2}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzEngine_ReevaluateRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := sz.SZ_WITHOUT_INFO
	_, err := szEngine.ReevaluateRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_ReevaluateRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := sz.SZ_WITH_INFO
	result, err := szEngine.ReevaluateRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":2}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzEngine_ReplaceRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	flags := sz.SZ_WITHOUT_INFO
	_, err := szEngine.ReplaceRecord(ctx, dataSourceCode, recordId, recordDefinition, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_ReplaceRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	flags := sz.SZ_WITH_INFO
	result, err := szEngine.ReplaceRecord(ctx, dataSourceCode, recordId, recordDefinition, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzEngine_DeleteRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	flags := sz.SZ_WITHOUT_INFO
	_, err := szEngine.DeleteRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_DeleteRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	flags := sz.SZ_WITH_INFO
	result, err := szEngine.DeleteRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzEngine_Initialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	instanceName := "Test module name"
	settings := "{}"
	verboseLogging := sz.SZ_NO_LOGGING
	configId := sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION
	err := szEngine.Initialize(ctx, instanceName, settings, verboseLogging, configId)
	if err != nil {
		// This should produce a "senzing-60144002" error.
	}
	// Output:
}

func ExampleSzEngine_Reinitialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	configId, _ := szEngine.GetActiveConfigId(ctx) // Example initConfigID.
	err := szEngine.Reinitialize(ctx, configId)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	err := szEngine.Destroy(ctx)
	if err != nil {
		// This should produce a "senzing-60164001" error.
	}
	// Output:
}
