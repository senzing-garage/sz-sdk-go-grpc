package szengine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	truncator "github.com/aquilax/truncate"
	"github.com/senzing-garage/go-helpers/record"
	"github.com/senzing-garage/go-helpers/testfixtures"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-grpc/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	szdiagnosticpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	badAttributes          = "}{"
	badBuildOutDegree      = int64(-1)
	badBuildOutMaxEntities = int64(-1)
	badCsvColumnList       = "BAD, CSV, COLUMN, LIST"
	badDataSourceCode      = "BadDataSourceCode"
	badEntityID            = int64(0)
	badExclusions          = "}{"
	badExportHandle        = uintptr(0)
	badLogLevelName        = "BadLogLevelName"
	badMaxDegrees          = int64(-1)
	badRecordID            = "BadRecordID"
	badRedoRecord          = "{}"
	badRequiredDataSources = "}{"
	badSearchProfile       = "}{"
	defaultTruncation      = 76
	instanceName           = "SzEngine Test"
	observerOrigin         = "SzEngine observer"
	printResults           = false
	verboseLogging         = senzing.SzNoLogging
)

type GetEntityByRecordIDResponse struct {
	ResolvedEntity struct {
		EntityID int64 `json:"ENTITY_ID"`
	} `json:"RESOLVED_ENTITY"`
}

var (
	defaultConfigID   int64
	grpcAddress       = "localhost:8261"
	grpcConnection    *grpc.ClientConn
	logLevel          = helper.GetEnv("SENZING_LOG_LEVEL", "INFO")
	observerSingleton = &observer.NullObserver{
		ID:       "Observer 1",
		IsSilent: true,
	}
	szConfigManagerSingleton *szconfigmanager.Szconfigmanager
	szConfigSingleton        *szconfig.Szconfig
	szDiagnosticSingleton    *szdiagnostic.Szdiagnostic
	szEngineSingleton        *Szengine
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestSzengine_AddRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzWithoutInfo
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	for _, record := range records {
		actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		require.NoError(test, err)
		printActual(test, actual)
	}
	for _, record := range records {
		actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		require.NoError(test, err)
		printActual(test, actual)
	}
}

func TestG2engine_AddRecord_badDataSourceCodeInJSON(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzWithoutInfo
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record2Json := `{"DATA_SOURCE": "BOB", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`
	_, err := szEngine.AddRecord(ctx, record1.DataSource, record1.ID, record1.JSON, flags)
	require.NoError(test, err)
	_, err = szEngine.AddRecord(ctx, record2.DataSource, record2.ID, record2Json, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
}

func TestSzengine_AddRecord_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, badDataSourceCode, record.ID, record.JSON, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_AddRecord_badRecordID(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.AddRecord(ctx, record.DataSource, badRecordID, record.JSON, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_AddRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzWithInfo
	records := []record.Record{
		truthset.CustomerRecords["1003"],
		truthset.CustomerRecords["1004"],
	}
	for _, record := range records {
		actual, err := szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		require.NoError(test, err)
		printActual(test, actual)
	}
	for _, record := range records {
		actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		require.NoError(test, err)
		printActual(test, actual)
	}
}

func TestSzengine_AddRecord_withInfo_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzWithInfo
	records := []record.Record{
		truthset.CustomerRecords["1003"],
		truthset.CustomerRecords["1004"],
	}
	for _, record := range records {
		actual, err := szEngine.AddRecord(ctx, badDataSourceCode, record.ID, record.JSON, flags)
		require.ErrorIs(test, err, szerror.ErrSzBadInput)
		printActual(test, actual)
	}
	for _, record := range records {
		actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		require.NoError(test, err)
		printActual(test, actual)
	}
}
func TestSzengine_CloseExport(test *testing.T) {
	_ = test
	// Tested in:
	//  - TestSzengine_ExportCsvEntityReport
	//  - TestSzengine_ExportJSONEntityReport
}

func TestSzengine_CountRedoRecords(test *testing.T) {
	ctx := context.TODO()
	expected := int64(1)
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.CountRedoRecords(ctx)
	require.NoError(test, err)
	printActual(test, actual)
	assert.Equal(test, expected, actual)
}

func TestSzengine_DeleteRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	records := []record.Record{
		truthset.CustomerRecords["1005"],
	}
	err := addRecords(ctx, records)
	require.NoError(test, err)
	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	require.NoError(test, err)
	actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.DeleteRecord(ctx, badDataSourceCode, record.ID, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_badRecordID(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1005"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.DeleteRecord(ctx, record.DataSource, badRecordID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	records := []record.Record{
		truthset.CustomerRecords["1009"],
	}
	err := addRecords(ctx, records)
	require.NoError(test, err)
	record := truthset.CustomerRecords["1009"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_withInfo_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	records := []record.Record{
		truthset.CustomerRecords["1009"],
	}
	err := addRecords(ctx, records)
	require.NoError(test, err)
	record := truthset.CustomerRecords["1009"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.DeleteRecord(ctx, badDataSourceCode, record.ID, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_DeleteRecord_withInfo_badDataSourceCode_fix(test *testing.T) {
	_ = test
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1009"],
	}
	handleError(deleteRecords(ctx, records))
}

func TestSzengine_ExportCsvEntityReport(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	expected := []string{
		`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
		`8,0,"","","CUSTOMERS","1001"`,
		`8,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
		`8,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
	}
	szEngine := getTestObject(ctx, test)
	csvColumnList := ""
	flags := senzing.SzExportIncludeAllEntities
	exportHandle, err := szEngine.ExportCsvEntityReport(ctx, csvColumnList, flags)
	defer func() {
		err := szEngine.CloseExport(ctx, exportHandle)
		require.NoError(test, err)
	}()
	require.NoError(test, err)
	actualCount := 0
	for actual := range szEngine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))
		actualCount++
	}
	assert.Equal(test, len(expected), actualCount)
}

func TestSzengine_ExportCsvEntityReport_badCsvColumnList(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzExportIncludeAllEntities
	exportHandle, err := szEngine.ExportCsvEntityReport(ctx, badCsvColumnList, flags)
	defer func() {
		err := szEngine.CloseExport(ctx, exportHandle)
		require.ErrorIs(test, err, szerror.ErrSzBase)
	}()
	require.ErrorIs(test, err, szerror.ErrSzBadInput)

}

func TestSzengine_ExportCsvEntityReportIterator(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	expected := []string{
		`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
		`12,0,"","","CUSTOMERS","1001"`,
		`12,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
		`12,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
	}
	szEngine := getTestObject(ctx, test)
	csvColumnList := ""
	flags := senzing.SzExportIncludeAllEntities
	actualCount := 0
	for actual := range szEngine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
		require.NoError(test, actual.Error)
		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))
		actualCount++
	}
	assert.Equal(test, len(expected), actualCount)
}

func TestSzengine_ExportCsvEntityReportIterator_badCsvColumnList(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	expected := []string{
		``,
	}
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzExportIncludeAllEntities
	actualCount := 0
	for actual := range szEngine.ExportCsvEntityReportIterator(ctx, badCsvColumnList, flags) {
		require.ErrorIs(test, actual.Error, szerror.ErrSzBadInput)
		assert.Equal(test, expected[actualCount], strings.TrimSpace(actual.Value))
		actualCount++
	}
	assert.Equal(test, len(expected), actualCount)
}

func TestSzengine_ExportJSONEntityReport(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	aRecord := testfixtures.FixtureRecords["65536-periods"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.AddRecord(ctx, aRecord.DataSource, aRecord.ID, aRecord.JSON, flags)
	require.NoError(test, err)
	printActual(test, actual)
	defer func() {
		handleErrorWithString(szEngine.DeleteRecord(ctx, aRecord.DataSource, aRecord.ID, senzing.SzWithoutInfo))
	}()
	// TODO: Figure out correct flags.
	// flags := senzing.Flags(senzing.SZ_EXPORT_DEFAULT_FLAGS, senzing.SZ_EXPORT_INCLUDE_ALL_HAVING_RELATIONSHIPS, senzing.SZ_EXPORT_INCLUDE_ALL_HAVING_RELATIONSHIPS)
	flags = int64(-1)
	exportHandle, err := szEngine.ExportJSONEntityReport(ctx, flags)
	defer func() {
		err := szEngine.CloseExport(ctx, exportHandle)
		require.NoError(test, err)
	}()
	require.NoError(test, err)
	jsonEntityReport := ""
	for {
		jsonEntityReportFragment, err := szEngine.FetchNext(ctx, exportHandle)
		require.NoError(test, err)
		if len(jsonEntityReportFragment) == 0 {
			break
		}
		jsonEntityReport += jsonEntityReportFragment
	}
	require.NoError(test, err)
	assert.Greater(test, len(jsonEntityReport), 65536)
}
func TestSzengine_ExportJSONEntityReport_65536(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	aRecord := testfixtures.FixtureRecords["65536-periods"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.AddRecord(ctx, aRecord.DataSource, aRecord.ID, aRecord.JSON, flags)
	require.NoError(test, err)
	printActual(test, actual)
	defer func() { _, _ = szEngine.DeleteRecord(ctx, aRecord.DataSource, aRecord.ID, senzing.SzWithoutInfo) }()
	// TODO: Figure out correct flags.
	// flags := senzing.Flags(senzing.SZ_EXPORT_DEFAULT_FLAGS, senzing.SZ_EXPORT_INCLUDE_ALL_HAVING_RELATIONSHIPS, senzing.SZ_EXPORT_INCLUDE_ALL_HAVING_RELATIONSHIPS)
	flags = int64(-1)
	aHandle, err := szEngine.ExportJSONEntityReport(ctx, flags)
	defer func() {
		err := szEngine.CloseExport(ctx, aHandle)
		require.NoError(test, err)
	}()
	require.NoError(test, err)
	jsonEntityReport := ""
	for {
		jsonEntityReportFragment, err := szEngine.FetchNext(ctx, aHandle)
		require.NoError(test, err)
		if len(jsonEntityReportFragment) == 0 {
			break
		}
		jsonEntityReport += jsonEntityReportFragment
	}
	require.NoError(test, err)
	assert.Greater(test, len(jsonEntityReport), 65536)
}

// TODO: Implement TestSzengine_ExportJSONEntityReport_error
// func TestSzengine_ExportJSONEntityReport_error(test *testing.T) {}

func TestSzengine_ExportJSONEntityReportIterator(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()

	err := addRecords(ctx, records)
	require.NoError(test, err)
	expected := 1
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzExportIncludeAllEntities
	actualCount := 0
	for actual := range szEngine.ExportJSONEntityReportIterator(ctx, flags) {
		require.NoError(test, actual.Error)
		printActual(test, actual.Value)
		actualCount++
	}
	assert.Equal(test, expected, actualCount)
}

func TestSzengine_FetchNext(test *testing.T) {
	_ = test
	// Tested in:
	//  - TestSzengine_ExportJSONEntityReport
}

func TestSzengine_FetchNext_badExportHandle(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.FetchNext(ctx, badExportHandle)
	require.ErrorIs(test, err, szerror.ErrSzBase)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	entityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByEntityID(ctx, entityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByEntityID_badEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByEntityID(ctx, badEntityID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByRecordID(ctx, record.DataSource, record.ID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByRecordID_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByRecordID(ctx, badDataSourceCode, record.ID, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_FindInterestingEntitiesByRecordID_badRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindInterestingEntitiesByRecordID(ctx, record.DataSource, badRecordID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1, err := getEntityIDString(record1)
	require.NoError(test, err)
	entityID2, err := getEntityIDString(record2)
	require.NoError(test, err)
	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	maxDegrees := int64(2)
	buildOutDegree := int64(1)
	buildOutMaxEntities := int64(10)
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(ctx, entityIDs, maxDegrees, buildOutDegree, buildOutMaxEntities, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID_badEntityIDs(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	badEntityID1 := 0
	badEntityID2 := 1
	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + strconv.Itoa(badEntityID1) + `}, {"ENTITY_ID": ` + strconv.Itoa(badEntityID2) + `}]}`
	maxDegrees := int64(2)
	buildOutDegree := int64(1)
	buildOutMaxEntities := int64(10)
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(ctx, entityIDs, maxDegrees, buildOutDegree, buildOutMaxEntities, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID_badMaxDegrees(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1, err := getEntityIDString(record1)
	require.NoError(test, err)
	entityID2, err := getEntityIDString(record2)
	require.NoError(test, err)
	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	buildOutDegree := int64(1)
	buildOutMaxEntities := int64(10)
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(ctx, entityIDs, badMaxDegrees, buildOutDegree, buildOutMaxEntities, flags)
	require.ErrorIs(test, err, szerror.ErrSzBase)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID_badBuildOutDegree(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1, err := getEntityIDString(record1)
	require.NoError(test, err)
	entityID2, err := getEntityIDString(record2)
	require.NoError(test, err)
	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	maxDegrees := int64(2)
	buildOutMaxEntities := int64(10)
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(ctx, entityIDs, maxDegrees, badBuildOutDegree, buildOutMaxEntities, flags)
	require.ErrorIs(test, err, szerror.ErrSzBase)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByEntityID_badBuildOutMaxEntities(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	entityID1, err := getEntityIDString(record1)
	require.NoError(test, err)
	entityID2, err := getEntityIDString(record2)
	require.NoError(test, err)
	entityIDs := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
	maxDegrees := int64(2)
	buildOutDegree := int64(1)
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByEntityID(ctx, entityIDs, maxDegrees, buildOutDegree, badBuildOutMaxEntities, flags)
	require.ErrorIs(test, err, szerror.ErrSzBase)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordKeys := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.ID + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.ID + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_ID": "` + record3.ID + `"}]}`
	maxDegrees := int64(1)
	buildOutDegree := int64(2)
	buildOutMaxEntities := int64(10)
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByRecordID(ctx, recordKeys, maxDegrees, buildOutDegree, buildOutMaxEntities, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByRecordID_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordKeys := `{"RECORDS": [{"DATA_SOURCE": "` + badDataSourceCode + `", "RECORD_ID": "` + record1.ID + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.ID + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_ID": "` + record3.ID + `"}]}`
	maxDegrees := int64(1)
	buildOutDegree := int64(2)
	buildOutMaxEntities := int64(10)
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByRecordID(ctx, recordKeys, maxDegrees, buildOutDegree, buildOutMaxEntities, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_FindNetworkByRecordID_badRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	record3 := truthset.CustomerRecords["1003"]
	recordKeys := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + badRecordID + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.ID + `"}, {"DATA_SOURCE": "` + record3.DataSource + `", "RECORD_ID": "` + record3.ID + `"}]}`
	maxDegrees := int64(1)
	buildOutDegree := int64(2)
	buildOutMaxEntities := int64(10)
	flags := senzing.SzFindNetworkDefaultFlags
	actual, err := szEngine.FindNetworkByRecordID(ctx, recordKeys, maxDegrees, buildOutDegree, buildOutMaxEntities, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)
	maxDegrees := int64(1)
	exclusions := senzing.SzNoExclusions
	requiredDataSources := senzing.SzNoRequiredDatasources
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(ctx, startEntityID, endEntityID, maxDegrees, exclusions, requiredDataSources, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_badStartEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	badStartEntityID := badEntityID
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)
	maxDegrees := int64(1)
	exclusions := senzing.SzNoExclusions
	requiredDataSources := senzing.SzNoRequiredDatasources
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(ctx, badStartEntityID, endEntityID, maxDegrees, exclusions, requiredDataSources, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_badEndEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	badEndEntityID := badEntityID
	maxDegrees := int64(1)
	exclusions := senzing.SzNoExclusions
	requiredDataSources := senzing.SzNoRequiredDatasources
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(ctx, startEntityID, badEndEntityID, maxDegrees, exclusions, requiredDataSources, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_badMaxDegrees(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)
	exclusions := senzing.SzNoExclusions
	requiredDataSources := senzing.SzNoRequiredDatasources
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(ctx, startEntityID, endEntityID, badMaxDegrees, exclusions, requiredDataSources, flags)
	require.ErrorIs(test, err, szerror.ErrSzBase)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_badExclusions(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)
	maxDegrees := int64(1)
	requiredDataSources := senzing.SzNoRequiredDatasources
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(ctx, startEntityID, endEntityID, maxDegrees, badExclusions, requiredDataSources, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_badRequiredDataSource(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	startEntityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)
	maxDegrees := int64(1)
	exclusions := senzing.SzNoExclusions
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(ctx, startEntityID, endEntityID, maxDegrees, exclusions, badRequiredDataSources, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_excluding(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityID, err := getEntityID(startRecord)
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)
	maxDegrees := int64(1)
	startEntityIDString, err := getEntityIDStringForRecord("CUSTOMERS", "1001")
	require.NoError(test, err)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + startEntityIDString + `}]}`
	requiredDataSources := senzing.SzNoRequiredDatasources
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(ctx, startEntityID, endEntityID, maxDegrees, exclusions, requiredDataSources, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_excluding_badStartEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityID := badEntityID
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)
	maxDegrees := int64(1)
	startRecordEntityIDString, err := getEntityIDString(startRecord)
	require.NoError(test, err)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + startRecordEntityIDString + `}]}`
	requiredDataSources := senzing.SzNoRequiredDatasources
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(ctx, startEntityID, endEntityID, maxDegrees, exclusions, requiredDataSources, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_excludingAndIncluding(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityID, err := getEntityID(startRecord)
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)
	maxDegrees := int64(1)
	startRecordEntityIDString, err := getEntityIDString(startRecord)
	require.NoError(test, err)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + startRecordEntityIDString + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + startRecord.DataSource + `"]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(ctx, startEntityID, endEntityID, maxDegrees, exclusions, requiredDataSources, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_including(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	startRecord := truthset.CustomerRecords["1001"]
	startEntityID, err := getEntityID(startRecord)
	require.NoError(test, err)
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)
	maxDegrees := int64(1)
	exclusions := senzing.SzNoExclusions
	requiredDataSources := `{"DATA_SOURCES": ["` + startRecord.DataSource + `"]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(ctx, startEntityID, endEntityID, maxDegrees, exclusions, requiredDataSources, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByEntityID_including_badStartEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	startRecord := truthset.CustomerRecords["1001"]
	badStartEntityID := badEntityID
	endEntityID, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)
	maxDegrees := int64(1)
	exclusions := senzing.SzNoExclusions
	requiredDataSources := `{"DATA_SOURCES": ["` + startRecord.DataSource + `"]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByEntityID(ctx, badStartEntityID, endEntityID, maxDegrees, exclusions, requiredDataSources, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := senzing.SzNoExclusions
	requiredDataSources := senzing.SzNoRequiredDatasources
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(ctx, record1.DataSource, record1.ID, record2.DataSource, record2.ID, maxDegree, exclusions, requiredDataSources, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := senzing.SzNoExclusions
	requiredDataSources := senzing.SzNoRequiredDatasources
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(ctx, badDataSourceCode, record1.ID, record2.DataSource, record2.ID, maxDegree, exclusions, requiredDataSources, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_badRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := senzing.SzNoExclusions
	requiredDataSources := senzing.SzNoRequiredDatasources
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(ctx, record1.DataSource, badRecordID, record2.DataSource, record2.ID, maxDegree, exclusions, requiredDataSources, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_excluding(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.ID + `"}]}`
	requiredDataSources := senzing.SzNoRequiredDatasources
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(ctx, record1.DataSource, record1.ID, record2.DataSource, record2.ID, maxDegree, exclusions, requiredDataSources, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_excluding_badStartDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.ID + `"}]}`
	requiredDataSources := senzing.SzNoRequiredDatasources
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(ctx, badDataSourceCode, record1.ID, record2.DataSource, record2.ID, maxDegree, exclusions, requiredDataSources, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_excludingAndIncluding(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	exclusions := `{"RECORDS": [{ "DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.ID + `"}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(ctx, record1.DataSource, record1.ID, record2.DataSource, record2.ID, maxDegree, exclusions, requiredDataSources, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_including(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	record1EntityID, err := getEntityIDString(record1)
	require.NoError(test, err)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + record1EntityID + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(ctx, record1.DataSource, record1.ID, record2.DataSource, record2.ID, maxDegree, exclusions, requiredDataSources, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_FindPathByRecordID_including_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	maxDegree := int64(1)
	record1EntityID, err := getEntityIDString(record1)
	require.NoError(test, err)
	exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + record1EntityID + `}]}`
	requiredDataSources := `{"DATA_SOURCES": ["` + record1.DataSource + `"]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.FindPathByRecordID(ctx, badDataSourceCode, record1.ID, record2.DataSource, record2.ID, maxDegree, exclusions, requiredDataSources, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_GetActiveConfigID(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetActiveConfigID(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// TODO: Implement TestSzengine_GetActiveConfigID_error
// func TestSzengine_GetActiveConfigID_error(test *testing.T) {}

func TestSzengine_GetEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	entityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByEntityID(ctx, entityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetEntityByEntityID_badEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByEntityID(ctx, badEntityID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_GetEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByRecordID(ctx, record.DataSource, record.ID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetEntityByRecordID_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByRecordID(ctx, badDataSourceCode, record.ID, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_GetEntityByRecordID_badRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetEntityByRecordID(ctx, record.DataSource, badRecordID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_GetRecord(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetRecord(ctx, record.DataSource, record.ID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetRecord_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetRecord(ctx, badDataSourceCode, record.ID, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_GetRecord_badRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetRecord(ctx, record.DataSource, badRecordID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_GetRedoRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetRedoRecord(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// TODO: Implement TestSzengine_GetRedoRecord_error
// func TestSzengine_GetRedoRecord_error(test *testing.T) {}

func TestSzengine_GetStats(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	actual, err := szEngine.GetStats(ctx)
	require.NoError(test, err)
	printActual(test, actual)
}

// TODO: Implement TestSzengine_GetStats_error
// func TestSzengine_GetStats_error(test *testing.T) {}

func TestSzengine_GetVirtualEntityByRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + record1.ID + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.ID + `"}]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetVirtualEntityByRecordID(ctx, recordList, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_GetVirtualEntityByRecordID_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + badDataSourceCode + `", "RECORD_ID": "` + record1.ID + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.ID + `"}]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetVirtualEntityByRecordID(ctx, recordList, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_GetVirtualEntityByRecordID_badRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	recordList := `{"RECORDS": [{"DATA_SOURCE": "` + record1.DataSource + `", "RECORD_ID": "` + badRecordID + `"}, {"DATA_SOURCE": "` + record2.DataSource + `", "RECORD_ID": "` + record2.ID + `"}]}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.GetVirtualEntityByRecordID(ctx, recordList, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_HowEntityByEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	entityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	flags := senzing.SzNoFlags
	actual, err := szEngine.HowEntityByEntityID(ctx, entityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_HowEntityByEntityID_badEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzNoFlags
	actual, err := szEngine.HowEntityByEntityID(ctx, badEntityID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_PrimeEngine(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	err := szEngine.PrimeEngine(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzengine_PrimeEngine_error
// func TestSzengine_PrimeEngine_error(test *testing.T) {}

func TestSzengine_ProcessRedoRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	redoRecord, err := szEngine.GetRedoRecord(ctx)
	require.NoError(test, err)
	if len(redoRecord) > 0 {
		flags := senzing.SzWithoutInfo
		actual, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
		require.NoError(test, err)
		printActual(test, actual)
	}
}

func TestSzengine_ProcessRedoRecord_badRedoRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ProcessRedoRecord(ctx, badRedoRecord, flags)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
	printActual(test, actual)
}

func TestSzengine_ProcessRedoRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
		truthset.CustomerRecords["1004"],
		truthset.CustomerRecords["1005"],
		truthset.CustomerRecords["1009"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	redoRecord, err := szEngine.GetRedoRecord(ctx)
	require.NoError(test, err)
	if len(redoRecord) > 0 {
		flags := senzing.SzWithInfo
		actual, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
		require.NoError(test, err)
		printActual(test, actual)
	}
}

func TestSzengine_ProcessRedoRecord_withInfo_badRedoRecord(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzWithInfo
	actual, err := szEngine.ProcessRedoRecord(ctx, badRedoRecord, flags)
	require.ErrorIs(test, err, szerror.ErrSzConfiguration)
	printActual(test, actual)
}

func TestSzengine_ReevaluateEntity(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	entityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateEntity(ctx, entityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_ReevaluateEntity_badEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateEntity(ctx, badEntityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

// TODO: Implement TestSzengine_ReevaluateEntity_error
// func TestSzengine_ReevaluateEntity_error(test *testing.T) {}

func TestSzengine_ReevaluateEntity_withInfo(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	entityID, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateEntity(ctx, entityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_ReevaluateEntity_withInfo_badEntityID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateEntity(ctx, badEntityID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

// TODO: Implement TestSzengine_ReevaluateEntity_withInfo_error
// func TestSzengine_ReevaluateEntity_withInfo_error(test *testing.T) {}

func TestSzengine_ReevaluateRecord(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateRecord(ctx, record.DataSource, record.ID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateRecord(ctx, badDataSourceCode, record.ID, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord_badRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithoutInfo
	actual, err := szEngine.ReevaluateRecord(ctx, record.DataSource, badRecordID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}
func TestSzengine_ReevaluateRecord_withInfo(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateRecord(ctx, record.DataSource, record.ID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_ReevaluateRecord_withInfo_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzWithInfo
	actual, err := szEngine.ReevaluateRecord(ctx, badDataSourceCode, record.ID, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_SearchByAttributes(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	searchProfile := senzing.SzNoSearchProfile
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_SearchByAttributes_badAttributes(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	searchProfile := senzing.SzNoSearchProfile
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, badAttributes, searchProfile, flags)
	require.ErrorIs(test, err, szerror.ErrSzBase)
	printActual(test, actual)
}

func TestSzengine_SearchByAttributes_badSearchProfile(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, attributes, badSearchProfile, flags)
	require.ErrorIs(test, err, szerror.ErrSzBadInput)
	printActual(test, actual)
}

// TODO: Implement TestSzengine_SearchByAttributes_error
// func TestSzengine_SearchByAttributes_error(test *testing.T) {}

func TestSzengine_SearchByAttributes_withSearchProfile(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
		truthset.CustomerRecords["1003"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	searchProfile := "SEARCH"
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

// TODO: Implement TestSzengine_StreamExportCsvEntityReport
// func TestSzengine_StreamExportCsvEntityReport(test *testing.T) {}

// TODO: Implement TestSzengine_StreamExportJSONEntityReport
// func TestSzengine_StreamExportJSONEntityReport(test *testing.T) {}

func TestSzengine_SearchByAttributes_searchProfile(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "JOHNSON"}], "SSN_NUMBER": "053-39-3251"}`
	searchProfile := senzing.SzNoSearchProfile // TODO: Figure out the search profile
	flags := senzing.SzNoFlags
	actual, err := szEngine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_WhyEntities(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	entityID1, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	entityID2, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, entityID1, entityID2, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_WhyEntities_badEnitity1(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	entityID2, err := getEntityID(truthset.CustomerRecords["1002"])
	require.NoError(test, err)
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, badEntityID, entityID2, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_WhyEntities_badEnitity2(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	entityID1, err := getEntityID(truthset.CustomerRecords["1001"])
	require.NoError(test, err)
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyEntities(ctx, entityID1, badEntityID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_WhyRecordInEntity(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecordInEntity(ctx, record.DataSource, record.ID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_WhyRecordInEntity_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecordInEntity(ctx, badDataSourceCode, record.ID, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_WhyRecordInEntity_badRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record := truthset.CustomerRecords["1001"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecordInEntity(ctx, record.DataSource, badRecordID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

func TestSzengine_WhyRecords(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecords(ctx, record1.DataSource, record1.ID, record2.DataSource, record2.ID, flags)
	require.NoError(test, err)
	printActual(test, actual)
}

func TestSzengine_WhyRecords_badDataSourceCode(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecords(ctx, badDataSourceCode, record1.ID, record2.DataSource, record2.ID, flags)
	require.ErrorIs(test, err, szerror.ErrSzUnknownDataSource)
	printActual(test, actual)
}

func TestSzengine_WhyRecords_badRecordID(test *testing.T) {
	ctx := context.TODO()
	records := []record.Record{
		truthset.CustomerRecords["1001"],
		truthset.CustomerRecords["1002"],
	}
	defer func() { handleError(deleteRecords(ctx, records)) }()
	err := addRecords(ctx, records)
	require.NoError(test, err)
	szEngine := getTestObject(ctx, test)
	record1 := truthset.CustomerRecords["1001"]
	record2 := truthset.CustomerRecords["1002"]
	flags := senzing.SzNoFlags
	actual, err := szEngine.WhyRecords(ctx, record1.DataSource, record1.ID, record2.DataSource, badRecordID, flags)
	require.ErrorIs(test, err, szerror.ErrSzNotFound)
	printActual(test, actual)
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func TestSzengine_SetLogLevel_badLogLevelName(test *testing.T) {
	ctx := context.TODO()
	szConfig := getTestObject(ctx, test)
	_ = szConfig.SetLogLevel(ctx, badLogLevelName)
}

func TestSzengine_SetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
}

func TestSzengine_GetObserverOrigin(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
	actual := szEngine.GetObserverOrigin(ctx)
	assert.Equal(test, origin, actual)
	printActual(test, actual)
}

func TestSzengine_UnregisterObserver(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	err := szEngine.UnregisterObserver(ctx, observerSingleton)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func TestSzengine_AsInterface(test *testing.T) {
	expected := int64(4)
	ctx := context.TODO()
	szEngine := getSzEngineAsInterface(ctx)
	actual, err := szEngine.CountRedoRecords(ctx)
	require.NoError(test, err)
	printActual(test, actual)
	assert.Equal(test, expected, actual)
}

func TestSzengine_Initialize(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	settings, err := getSettings()
	require.NoError(test, err)
	configID := senzing.SzInitializeWithDefaultConfiguration
	err = szEngine.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	require.NoError(test, err)
}

// TODO: Implement TestSzengine_Initialize_error
// func TestSzengine_Initialize_error(test *testing.T) {}

func TestSzengine_Initialize_withConfigID(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	settings, err := getSettings()
	require.NoError(test, err)
	configID := getDefaultConfigID()
	err = szEngine.Initialize(ctx, instanceName, settings, configID, verboseLogging)
	require.NoError(test, err)
}

// TODO: Implement TestSzengine_Initialize_withConfigID_error
// func TestSzengine_Initialize_withConfigID_error(test *testing.T) {}

func TestSzengine_Reinitialize(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	configID, err := szEngine.GetActiveConfigID(ctx)
	require.NoError(test, err)
	err = szEngine.Reinitialize(ctx, configID)
	require.NoError(test, err)
	printActual(test, configID)
}

// TODO: Implement TestSzengine_Reinitialize_badConfigID
// func TestSzengine_Reinitialize_badConfigID(test *testing.T) {}

func TestSzengine_Destroy(test *testing.T) {
	ctx := context.TODO()
	szEngine := getTestObject(ctx, test)
	err := szEngine.Destroy(ctx)
	require.NoError(test, err)
}

// TODO: Implement TestSzengine_Destroy_error
// func TestSzengine_Destroy_error(test *testing.T) {}

func TestSzengine_Destroy_withObserver(test *testing.T) {
	ctx := context.TODO()
	szEngineSingleton = nil
	szEngine := getTestObject(ctx, test)
	err := szEngine.Destroy(ctx)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func addRecords(ctx context.Context, records []record.Record) error {
	var err error
	szEngine, err := getSzEngine(ctx)
	if err != nil {
		return err
	}
	flags := senzing.SzWithoutInfo
	for _, record := range records {
		_, err = szEngine.AddRecord(ctx, record.DataSource, record.ID, record.JSON, flags)
		if err != nil {
			return err
		}
	}
	return err
}

func deleteRecords(ctx context.Context, records []record.Record) error {
	var err error
	szEngine, err := getSzEngine(ctx)
	if err != nil {
		return err
	}
	flags := senzing.SzWithoutInfo
	for _, record := range records {
		_, err = szEngine.DeleteRecord(ctx, record.DataSource, record.ID, flags)
		if err != nil {
			return err
		}
	}
	return err
}

func getDefaultConfigID() int64 {
	return defaultConfigID
}

func getEntityID(record record.Record) (int64, error) {
	return getEntityIDForRecord(record.DataSource, record.ID)
}

func getEntityIDForRecord(datasource string, id string) (int64, error) {
	var result int64
	var err error
	ctx := context.TODO()
	// var result int64
	szEngine, err := getSzEngine(ctx)
	if err != nil {
		return result, err
	}
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

func getEntityIDString(record record.Record) (string, error) {
	var result string
	var err error
	entityID, err := getEntityID(record)
	if err != nil {
		return result, err
	}
	result = strconv.FormatInt(entityID, baseTen)
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

func getSzConfig(ctx context.Context) (senzing.SzConfig, error) {
	var err error
	if szConfigSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			return szConfigSingleton, fmt.Errorf("getSettings() Error: %w", err)
		}
		grpcConnection := getGrpcConnection()
		szConfigSingleton = &szconfig.Szconfig{
			GrpcClient: szconfigpb.NewSzConfigClient(grpcConnection),
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

func getSzConfigManager(ctx context.Context) (senzing.SzConfigManager, error) {
	var err error
	if szConfigManagerSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			return szConfigManagerSingleton, fmt.Errorf("getSettings() Error: %w", err)
		}
		grpcConnection := getGrpcConnection()
		szConfigManagerSingleton = &szconfigmanager.Szconfigmanager{
			GrpcClient: szconfigmanagerpb.NewSzConfigManagerClient(grpcConnection),
		}
		err = szConfigManagerSingleton.SetLogLevel(ctx, logLevel)
		if err != nil {
			return szConfigManagerSingleton, fmt.Errorf("SetLogLevel() Error: %w", err)
		}
		if logLevel == "TRACE" {
			szConfigManagerSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szConfigManagerSingleton.RegisterObserver(ctx, observerSingleton)
			if err != nil {
				return szConfigManagerSingleton, fmt.Errorf("RegisterObserver() Error: %w", err)
			}
			err = szConfigManagerSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			if err != nil {
				return szConfigManagerSingleton, fmt.Errorf("SetLogLevel() - 2 Error: %w", err)
			}
		}
		err = szConfigManagerSingleton.Initialize(ctx, instanceName, settings, verboseLogging)
		if err != nil {
			return szConfigManagerSingleton, fmt.Errorf("Initialize() Error: %w", err)
		}
	}
	return szConfigManagerSingleton, err
}

func getSzDiagnostic(ctx context.Context) (senzing.SzDiagnostic, error) {
	var err error
	if szDiagnosticSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			return szDiagnosticSingleton, fmt.Errorf("getSettings() Error: %w", err)
		}
		grpcConnection := getGrpcConnection()
		szDiagnosticSingleton = &szdiagnostic.Szdiagnostic{
			GrpcClient: szdiagnosticpb.NewSzDiagnosticClient(grpcConnection),
		}
		err = szDiagnosticSingleton.SetLogLevel(ctx, logLevel)
		if err != nil {
			return szDiagnosticSingleton, fmt.Errorf("SetLogLevel() Error: %w", err)
		}
		if logLevel == "TRACE" {
			szDiagnosticSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szDiagnosticSingleton.RegisterObserver(ctx, observerSingleton)
			if err != nil {
				return szDiagnosticSingleton, fmt.Errorf("RegisterObserver() Error: %w", err)
			}
			err = szDiagnosticSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			if err != nil {
				return szDiagnosticSingleton, fmt.Errorf("SetLogLevel() - 2 Error: %w", err)
			}
		}
		err = szDiagnosticSingleton.Initialize(ctx, instanceName, settings, getDefaultConfigID(), verboseLogging)
		if err != nil {
			return szDiagnosticSingleton, fmt.Errorf("Initialize() Error: %w", err)
		}
	}
	return szDiagnosticSingleton, err
}

func getSzEngine(ctx context.Context) (*Szengine, error) {
	var err error
	if szEngineSingleton == nil {
		settings, err := getSettings()
		if err != nil {
			return szEngineSingleton, fmt.Errorf("getSettings() Error: %w", err)
		}
		grpcConnection := getGrpcConnection()
		szEngineSingleton = &Szengine{
			GrpcClient: szpb.NewSzEngineClient(grpcConnection),
		}
		err = szEngineSingleton.SetLogLevel(ctx, logLevel)
		if err != nil {
			return szEngineSingleton, fmt.Errorf("SetLogLevel() Error: %w", err)
		}
		if logLevel == "TRACE" {
			szEngineSingleton.SetObserverOrigin(ctx, observerOrigin)
			err = szEngineSingleton.RegisterObserver(ctx, observerSingleton)
			if err != nil {
				return szEngineSingleton, fmt.Errorf("RegisterObserver() Error: %w", err)
			}
			err = szEngineSingleton.SetLogLevel(ctx, logLevel) // Duplicated for coverage testing
			if err != nil {
				return szEngineSingleton, fmt.Errorf("SetLogLevel() - 2 Error: %w", err)
			}
		}
		err = szEngineSingleton.Initialize(ctx, instanceName, settings, getDefaultConfigID(), verboseLogging)
		if err != nil {
			return szEngineSingleton, fmt.Errorf("Initialize() Error: %w", err)
		}
	}
	return szEngineSingleton, err
}

func getSzEngineAsInterface(ctx context.Context) senzing.SzEngine {
	result, err := getSzEngine(ctx)
	handleError(err)
	return result
}

func getTestObject(ctx context.Context, test *testing.T) *Szengine {
	result, err := getSzEngine(ctx)
	require.NoError(test, err)
	return result
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func handleErrorWithString(aString string, err error) {
	_ = aString
	handleError(err)
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

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		if errors.Is(err, szerror.ErrSzUnrecoverable) {
			fmt.Printf("\nUnrecoverable error detected. \n\n")
		}
		if errors.Is(err, szerror.ErrSzRetryable) {
			fmt.Printf("\nRetryable error detected. \n\n")
		}
		if errors.Is(err, szerror.ErrSzBadInput) {
			fmt.Printf("\nBad user input error detected. \n\n")
		}
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
	var err error
	err = setupSenzingConfiguration()
	if err != nil {
		return fmt.Errorf("Failed to set up Senzing configuration. Error: %w", err)
	}
	err = setupPurgeRepository()
	if err != nil {
		return fmt.Errorf("failed to purge repository. Error: %w", err)
	}
	return err
}

func setupSenzingConfiguration() error {
	ctx := context.TODO()
	now := time.Now()

	// Create sz objects.

	szConfig, err := getSzConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to getSzConfig(). Error: %w", err)
	}
	szConfigManager, err := getSzConfigManager(ctx)
	if err != nil {
		return fmt.Errorf("failed to getSzConfigManager(). Error: %w", err)
	}

	// Create an in memory Senzing configuration.

	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to szConfig.CreateConfig(). Error: %w", err)
	}

	// Add data sources to in-memory Senzing configuration.

	dataSourceCodes := []string{"CUSTOMERS", "REFERENCE", "WATCHLIST"}
	for _, dataSourceCode := range dataSourceCodes {
		_, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return fmt.Errorf("failed to szConfig.AddDataSource(). Error: %w", err)
		}
	}

	// Create a string representation of the in-memory configuration.

	configDefinition, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return fmt.Errorf("failed to szConfig.ExportConfig(). Error: %w", err)
	}

	// Close szConfig in-memory object.

	err = szConfig.CloseConfig(ctx, configHandle)
	if err != nil {
		return fmt.Errorf("failed to szConfig.CloseConfig(). Error: %w", err)
	}

	// Persist the Senzing configuration to the Senzing repository as default.

	configComment := fmt.Sprintf("Created by szengine_test at %s", now.UTC())
	configID, err := szConfigManager.AddConfig(ctx, configDefinition, configComment)
	if err != nil {
		return fmt.Errorf("failed to szConfigManager.AddConfig(). Error: %w", err)
	}

	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return fmt.Errorf("failed to szConfigManager.SetDefaultConfigID(). Error: %w", err)
	}

	return err
}

func setupPurgeRepository() error {
	ctx := context.TODO()
	szDiagnostic, err := getSzDiagnostic(ctx)
	if err != nil {
		return fmt.Errorf("failed to get SzDiagnostic. Error: %w", err)
	}
	err = szDiagnostic.PurgeRepository(ctx)
	if err != nil {
		return fmt.Errorf("failed to szDiagnostic.PurgeRepository(). Error: %w", err)
	}
	return err
}

func teardown() error {
	var err error
	return err
}
