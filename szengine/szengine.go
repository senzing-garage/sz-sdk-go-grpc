/*
The [Szengine] implementation of the [senzing.SzEngine] interface
that communicates a gRPC server.
*/
package szengine

import (
	"context"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szengine"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
)

type Szengine struct {
	GrpcClient     szpb.SzEngineClient
	isTrace        bool // Performance optimization
	logger         logging.Logging
	observerOrigin string
	observers      subject.Subject
}

const (
	baseCallerSkip = 4
	baseTen        = 10
)

// ----------------------------------------------------------------------------
// sz-sdk-go.SzEngine interface methods
// ----------------------------------------------------------------------------

/*
Method AddRecord loads a record into the repository and performs entity resolution.

If a record already exists with the same data source code and record ID, it will be replaced.

If the record definition contains DATA_SOURCE and RECORD_ID JSON keys,
the values must match the dataSourceCode and recordID parameters.

The data source code must be registered in the active configuration.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - recordDefinition: A JSON document containing the record to be added to the Senzing repository.
  - flags: Flags used to control information returned.

Output
  - A JSON document containing metadata as specified by the flags.
*/
func (client *Szengine) AddRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	recordDefinition string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(1, dataSourceCode, recordID, recordDefinition, flags)

		entryTime := time.Now()

		defer func() {
			client.traceExit(2, dataSourceCode, recordID, recordDefinition, flags, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.addRecord(ctx, dataSourceCode, recordID, recordDefinition, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
				"flags":          strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8001, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method CloseExportReport closes an export report.

Used in conjunction with ExportJsonEntityReport(), ExportCsvEntityReport(), and FetchNext().

Input
  - ctx: A context to control lifecycle.
  - exportHandle: A handle created by [Szengine.ExportJSONEntityReport] or [Szengine.ExportCsvEntityReport]
    that is to be closed.
*/
func (client *Szengine) CloseExportReport(ctx context.Context, exportHandle uintptr) error {
	var err error

	if client.isTrace {
		client.traceEntry(5, exportHandle)

		entryTime := time.Now()

		defer func() { client.traceExit(6, exportHandle, err, time.Since(entryTime)) }()
	}

	err = client.closeExportReport(ctx, exportHandle)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8002, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method CountRedoRecords gets the number of redo records pending processing.

WARNING: When there is a large number of redo records, this is an expensive call.
Hint: If processing redo records, use result of [Szengine.GetRedoRecord] to manage looping.

Input
  - ctx: A context to control lifecycle.

Output
  - The number of redo records in Senzing's redo queue.
*/
func (client *Szengine) CountRedoRecords(ctx context.Context) (int64, error) {
	var (
		err    error
		result int64
	)

	if client.isTrace {
		client.traceEntry(7)

		entryTime := time.Now()

		defer func() { client.traceExit(8, result, err, time.Since(entryTime)) }()
	}

	result, err = client.countRedoRecords(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8003, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method DeleteRecord deletes a record from the repository and performs entity resolution.

The data source code must be registered in the active configuration.

Is idempotent.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document containing metadata as specified by the flags.
*/
func (client *Szengine) DeleteRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(9, dataSourceCode, recordID, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(10, dataSourceCode, recordID, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.deleteRecord(ctx, dataSourceCode, recordID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
				"flags":          strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8004, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method Destroy is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szengine) Destroy(ctx context.Context) error {
	var err error

	if client.isTrace {
		client.traceEntry(11)

		entryTime := time.Now()

		defer func() { client.traceExit(12, err, time.Since(entryTime)) }()
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8005, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method ExportCsvEntityReport initiates an export report of entity data in CSV format.

Used in conjunction with FetchNext() and CloseEntityReport().

The first FetchNext() call, after calling this method, returns the CSV header.

Subsequent FetchNext() calls return exported entity data in CSV format.

Use with large repositories is not advised.

Input
  - ctx: A context to control lifecycle.
  - csvColumnList: Use `*` to request all columns, an empty string to request "standard" columns,
    or a comma-separated list of column names for customized columns.
  - flags: Flags used to control information returned.

Output
  - exportHandle: A handle that identifies the document to be scrolled through using [Szengine.FetchNext].
*/
func (client *Szengine) ExportCsvEntityReport(ctx context.Context, csvColumnList string, flags int64) (uintptr, error) {
	var (
		err    error
		result uintptr
	)

	if client.isTrace {
		client.traceEntry(13, csvColumnList, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(14, csvColumnList, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.exportCsvEntityReport(ctx, csvColumnList, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"flags": strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8006, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method ExportCsvEntityReportIterator creates an Iterator that can be used in a for-loop
to scroll through a CSV document of exported entities.

It is a convenience method for the [Szenzine.ExportCsvEntityReport], [Szengine.FetchNext], [Szengine.CloseExportReport]
lifecycle of a list of entities to export.

Input
  - ctx: A context to control lifecycle.
  - csvColumnList: Use `*` to request all columns, an empty string to request "standard" columns,
    or a comma-separated list of column names for customized columns.
  - flags: Flags used to control information returned.

Output
  - A channel of strings that can be iterated over.
*/
func (client *Szengine) ExportCsvEntityReportIterator(
	ctx context.Context,
	csvColumnList string,
	flags int64,
) chan senzing.StringFragment {
	stringFragmentChannel := make(chan senzing.StringFragment)

	go func() {
		defer close(stringFragmentChannel)

		var err error

		if client.isTrace {
			client.traceEntry(15, csvColumnList, flags)

			entryTime := time.Now()

			defer func() { client.traceExit(16, csvColumnList, flags, err, time.Since(entryTime)) }()
		}

		request := &szpb.StreamExportCsvEntityReportRequest{
			CsvColumnList: csvColumnList,
			Flags:         flags,
		}

		stream, err := client.GrpcClient.StreamExportCsvEntityReport(ctx, request)
		if err != nil {
			stringFragmentChannel <- senzing.StringFragment{
				Error: helper.ConvertGrpcError(err),
			}

			return
		}

	forLoop:
		for {
			select {
			case <-ctx.Done():
				stringFragmentChannel <- senzing.StringFragment{
					Error: helper.ConvertGrpcError(ctx.Err()),
				}

				break forLoop
			default:
				response, err := stream.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break forLoop
					}

					stringFragmentChannel <- senzing.StringFragment{
						Error: helper.ConvertGrpcError(err),
					}

					break forLoop
				}

				stringFragmentChannel <- senzing.StringFragment{
					Value: response.GetResult(),
				}
			}
		}

		if client.observers != nil {
			go func() {
				details := map[string]string{
					"flags": strconv.FormatInt(flags, baseTen),
				}
				notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8007, err, details)
			}()
		}
	}()

	return stringFragmentChannel
}

/*
Method ExportJSONEntityReport initiates an export report of entity data in JSON Lines format.

Used in conjunction with FetchNext)() and CloseEntityReport().

Each fetchNext() call returns exported entity data as a JSON object.

Use with large repositories is not advised.

Input
  - ctx: A context to control lifecycle.
  - flags: Flags used to control information returned.

Output
  - A handle that identifies the document to be scrolled through using [Szengine.FetchNext].
*/
func (client *Szengine) ExportJSONEntityReport(ctx context.Context, flags int64) (uintptr, error) {
	var (
		err    error
		result uintptr
	)

	if client.isTrace {
		client.traceEntry(17, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(18, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.exportJSONEntityReport(ctx, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"flags": strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8008, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method ExportJSONEntityReportIterator creates an Iterator that can be used in a for-loop
to scroll through a JSON document of exported entities.

It is a convenience method for the [Szengine.ExportJSONEntityReport], [Szengine.FetchNext], [Szengine.CloseExportReport]
lifecycle of a list of entities to export.

Input
  - ctx: A context to control lifecycle.
  - flags: Flags used to control information returned.

Output
  - A channel of strings that can be iterated over.
*/
func (client *Szengine) ExportJSONEntityReportIterator(ctx context.Context, flags int64) chan senzing.StringFragment {
	stringFragmentChannel := make(chan senzing.StringFragment)

	go func() {
		defer close(stringFragmentChannel)

		var err error

		if client.isTrace {
			client.traceEntry(19, flags)

			entryTime := time.Now()

			defer func() { client.traceExit(20, flags, err, time.Since(entryTime)) }()
		}

		request := &szpb.StreamExportJsonEntityReportRequest{
			Flags: flags,
		}

		stream, err := client.GrpcClient.StreamExportJsonEntityReport(ctx, request)
		if err != nil {
			stringFragmentChannel <- senzing.StringFragment{
				Error: helper.ConvertGrpcError(err),
			}

			return
		}

	forLoop:
		for {
			select {
			case <-ctx.Done():
				stringFragmentChannel <- senzing.StringFragment{
					Error: helper.ConvertGrpcError(ctx.Err()),
				}

				break forLoop
			default:
				response, err := stream.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break forLoop
					}

					stringFragmentChannel <- senzing.StringFragment{
						Error: helper.ConvertGrpcError(err),
					}

					break forLoop
				}

				stringFragmentChannel <- senzing.StringFragment{
					Value: response.GetResult(),
				}
			}
		}

		if client.observers != nil {
			go func() {
				details := map[string]string{}
				notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8009, err, details)
			}()
		}
	}()

	return stringFragmentChannel
}

/*
Method FetchNext fetches the next line of entity data from an open export report.

Used in conjunction with ExportJsonEntityReport(), ExportCsvEntityReport(), and CloseEntityReport().

If the export handle was obtained from ExportCsvEntityReport(), this returns the CSV header on the first call and
exported entity data in CSV format on subsequent calls.

If the export handle was obtained from ExportJsonEntityReport(), this returns exported entity data as a JSON object.

When empty string is returned, the export report is complete
and the caller should invoke closeExportReport() to free resources.

Input
  - ctx: A context to control lifecycle.
  - exportHandle: A handle created by [Szengine.ExportJSONEntityReport] or [Szengine.ExportCsvEntityReport].

Output
  - The next chunk of exported data. An empty string signifies end of data.
*/
func (client *Szengine) FetchNext(ctx context.Context, exportHandle uintptr) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(21, exportHandle)

		entryTime := time.Now()

		defer func() { client.traceExit(22, exportHandle, result, err, time.Since(entryTime)) }()
	}

	result, err = client.fetchNext(ctx, exportHandle)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8010, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method FindInterestingEntitiesByEntityID is an experimental method.

Contact Senzing support.

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) FindInterestingEntitiesByEntityID(
	ctx context.Context,
	entityID int64,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(23, entityID, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(24, entityID, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.findInterestingEntitiesByEntityID(ctx, entityID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityID": formatEntityID(entityID),
				"flags":    strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8011, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method FindInterestingEntitiesByRecordID is an experimental method.

Contact Senzing support.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) FindInterestingEntitiesByRecordID(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(25, dataSourceCode, recordID, flags)

		entryTime := time.Now()

		defer func() {
			client.traceExit(26, dataSourceCode, recordID, flags, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.findInterestingEntitiesByRecordID(ctx, dataSourceCode, recordID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
				"flags":          strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8012, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method FindNetworkByEntityID retrieves a network of relationships among entities, specified by entity IDs.

Warning: Entity networks may be very large due to the volume of inter-related data in the repository.
The parameters of this method can be used to limit the information returned.

Input
  - ctx: A context to control lifecycle.
  - entityIDs: A JSON document listing entities.
    Example: `{"ENTITIES": [{"ENTITY_ID": 1}, {"ENTITY_ID": 2}, {"ENTITY_ID": 3}]}`
  - maxDegrees: The maximum number of degrees in paths between entityIDs.
  - buildOutDegrees: The number of degrees of relationships to show around each search entity. Zero (0)
    prevents buildout.
  - buildOutMaxEntities: The maximum number of entities to build out in the returned network.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) FindNetworkByEntityID(
	ctx context.Context,
	entityIDs string,
	maxDegrees int64,
	buildOutDegrees int64,
	buildOutMaxEntities int64,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(27, entityIDs, maxDegrees, buildOutDegrees, buildOutMaxEntities, flags)

		entryTime := time.Now()

		defer func() {
			client.traceExit(
				28,
				entityIDs,
				maxDegrees,
				buildOutDegrees,
				buildOutMaxEntities,
				flags,
				result,
				err,
				time.Since(entryTime),
			)
		}()
	}

	result, err = client.findNetworkByEntityID(
		ctx,
		entityIDs,
		maxDegrees,
		buildOutDegrees,
		buildOutMaxEntities,
		flags,
	)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityIDs": entityIDs,
				"flags":     strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8013, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method FindNetworkByRecordID retrieves a network of relationships among entities, specified by record IDs.

Warning: Entity networks may be very large due to the volume of inter-related data in the repository.
The parameters of this method can be used to limit the information returned.

Input
  - ctx: A context to control lifecycle.
  - recordKeys: A JSON document listing records.
    Example: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}]}`
  - maxDegrees: The maximum number of degrees in paths between entities identified by the recordKeys.
  - buildOutDegrees: The number of degrees of relationships to show around each search entity.
    Zero (0) prevents buildout.
  - buildOutMaxEntities: The maximum number of entities to build out in the returned network.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) FindNetworkByRecordID(
	ctx context.Context,
	recordKeys string,
	maxDegrees int64,
	buildOutDegrees int64,
	buildOutMaxEntities int64,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(29, recordKeys, maxDegrees, buildOutDegrees, buildOutMaxEntities, flags)

		entryTime := time.Now()

		defer func() {
			client.traceExit(
				30,
				recordKeys,
				maxDegrees,
				buildOutDegrees,
				buildOutMaxEntities,
				flags,
				result,
				err,
				time.Since(entryTime),
			)
		}()
	}

	result, err = client.findNetworkByRecordID(
		ctx,
		recordKeys,
		maxDegrees,
		buildOutDegrees,
		buildOutMaxEntities,
		flags,
	)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"recordKeys": recordKeys,
				"flags":      strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8014, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method FindPathByEntityID searches for the shortest relationship path between two entities, specified by entity IDs.

The returned path is the shortest path among the paths that satisfy the parameters.

Input
  - ctx: A context to control lifecycle.
  - startEntityID: The entity ID for the starting entity of the search path.
  - endEntityID: The entity ID for the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - avoidEntityIDs: A JSON document listing entities that should be avoided on the path.
    An empty string disables this capability.
    Example: `{"ENTITIES": [{"ENTITY_ID": 1}, {"ENTITY_ID": 2}, {"ENTITY_ID": 3}]}`
  - requiredDataSources: A JSON document listing data sources that should be included on the path.
    An empty string disables this capability.
    Example: `{"DATA_SOURCES": ["MY_DATASOURCE_1", "MY_DATASOURCE_2", "MY_DATASOURCE_3"]}`
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) FindPathByEntityID(
	ctx context.Context,
	startEntityID int64,
	endEntityID int64,
	maxDegrees int64,
	avoidEntityIDs string,
	requiredDataSources string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(31, startEntityID, endEntityID, maxDegrees, avoidEntityIDs, requiredDataSources, flags)

		entryTime := time.Now()

		defer func() {
			client.traceExit(32, startEntityID, endEntityID, maxDegrees, avoidEntityIDs, requiredDataSources,
				flags, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.findPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startEntityID":       formatEntityID(startEntityID),
				"endEntityID":         formatEntityID(endEntityID),
				"avoidEntityIDs":      avoidEntityIDs,
				"requiredDataSources": requiredDataSources,
				"flags":               strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8015, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method FindPathByRecordID searches for the shortest relationship path between two entities, specifiec by record IDs.

The returned path is the shortest path among the paths that satisfy the parameters.

Input
  - ctx: A context to control lifecycle.
  - startDataSourceCode: Identifies the provenance of the record for the starting
    entity of the search path.
  - startRecordID: The unique identifier within the records of the same data source
    for the starting entity of the search path.
  - endDataSourceCode: Identifies the provenance of the record for the ending entity
    of the search path.
  - endRecordID: The unique identifier within the records of the same data source for
    the ending entity of the search path.
  - maxDegrees: The maximum number of degrees in paths between search entities.
  - avoidRecordKeys: A JSON document listing entities that should be avoided on the path.
    An empty string disables this capability.
    Example: `{"RECORDS": [
    {"DATA_SOURCE": "MY_DATASOURCE", "RECORD_ID": "1"},
    {"DATA_SOURCE": "MY_DATASOURCE", "RECORD_ID": "2"},
    {"DATA_SOURCE": "MY_DATASOURCE", "RECORD_ID": "3"}
    ]}`
  - requiredDataSources: A JSON document listing data sources that should be included on the path.
    An empty string disables this capability.
    Example: `{"DATA_SOURCES": ["MY_DATASOURCE_1", "MY_DATASOURCE_2", "MY_DATASOURCE_3"]}`
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) FindPathByRecordID(
	ctx context.Context,
	startDataSourceCode string,
	startRecordID string,
	endDataSourceCode string,
	endRecordID string,
	maxDegrees int64,
	avoidRecordKeys string,
	requiredDataSources string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(33, startDataSourceCode, startRecordID, endDataSourceCode, endRecordID, maxDegrees,
			avoidRecordKeys, requiredDataSources, flags)

		entryTime := time.Now()

		defer func() {
			client.traceExit(34, startDataSourceCode, startRecordID, endDataSourceCode, endRecordID, maxDegrees,
				avoidRecordKeys, requiredDataSources, flags, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.findPathByRecordID(
		ctx,
		startDataSourceCode,
		startRecordID,
		endDataSourceCode,
		endRecordID,
		maxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"startDataSourceCode": startDataSourceCode,
				"startRecordID":       startRecordID,
				"endDataSourceCode":   endDataSourceCode,
				"endRecordID":         endRecordID,
				"avoidRecordKeys":     avoidRecordKeys,
				"requiredDataSources": requiredDataSources,
				"flags":               strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8016, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method GetActiveConfigID gets the currently active configuration ID.

May not be the default configuration ID.

Input
  - ctx: A context to control lifecycle.

Output
  - configID: The Senzing configuration JSON document identifier that is currently in use by the Senzing engine.
*/
func (client *Szengine) GetActiveConfigID(ctx context.Context) (int64, error) {
	var (
		err    error
		result int64
	)

	if client.isTrace {
		client.traceEntry(35)

		entryTime := time.Now()

		defer func() { client.traceExit(36, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getActiveConfigID(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8017, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method GetEntityByEntityID retrieves information about an entity, specified by entity ID.

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.
  - flags: Flags used to control information returned.

Output

  - A JSON document.
*/
func (client *Szengine) GetEntityByEntityID(ctx context.Context, entityID int64, flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(37, entityID, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(38, entityID, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getEntityByEntityID(ctx, entityID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityID": formatEntityID(entityID),
				"flags":    strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8018, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method GetEntityByRecordID retrieves information about an entity, specified by record ID.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) GetEntityByRecordID(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(39, dataSourceCode, recordID, flags)

		entryTime := time.Now()

		defer func() {
			client.traceExit(40, dataSourceCode, recordID, flags, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.getEntityByRecordID(ctx, dataSourceCode, recordID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
				"flags":          strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8019, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method GetRecord retrieves information about a record.

The information contains the original record data that was loaded and may contain other information
depending on the flags parameter.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) GetRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(45, dataSourceCode, recordID, flags)

		entryTime := time.Now()

		defer func() {
			client.traceExit(46, dataSourceCode, recordID, flags, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.getRecord(ctx, dataSourceCode, recordID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
				"flags":          strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8020, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method GetRecordPreview describes the features resulting from the hypothetical load of a record.

Used to preview the features for a record that has not been loaded.

Input
  - ctx: A context to control lifecycle.
  - recordDefinition: A JSON document containing the record to be tested against the Senzing repository.
  - flags: Flags used to control information returned.

Output
  - A JSON document containing metadata as specified by the flags.
*/
func (client *Szengine) GetRecordPreview(ctx context.Context, recordDefinition string, flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(77, recordDefinition, flags)

		entryTime := time.Now()

		defer func() {
			client.traceExit(78, recordDefinition, flags, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.getRecordPreview(ctx, recordDefinition, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"flags": strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8035, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method GetRedoRecord retrieves and removes a pending redo record.

An empty value will be returned if there are no pending redo records.

Use processRedoRecord() to process the result of this function.

Once a redo record is retrieved, it is no longer tracked by Senzing.

The redo record may be stored externally for later processing.

See also countRedoRecords(), processRedoRecord().

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document. If no redo records exist, an empty string is returned.
*/
func (client *Szengine) GetRedoRecord(ctx context.Context) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(47)

		entryTime := time.Now()

		defer func() { client.traceExit(48, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getRedoRecord(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8021, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method GetStats gets and resets the internal engine workload statistics for the current operating system process.

The output is helpful when interacting with Senzing support.

Best practice to periodically log the results.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document.
*/
func (client *Szengine) GetStats(ctx context.Context) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(49)

		entryTime := time.Now()

		defer func() { client.traceExit(50, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getStats(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8022, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method GetVirtualEntityByRecordID describes how an entity would look if composed of a given set of records.

Virtual entities do not have relationships.

Input
  - ctx: A context to control lifecycle.
  - recordKeys: A JSON document listing records to include in the hypothetical entity.
    Example: `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}]}`
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) GetVirtualEntityByRecordID(
	ctx context.Context,
	recordKeys string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(51, recordKeys, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(52, recordKeys, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getVirtualEntityByRecordID(ctx, recordKeys, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"recordKeys": recordKeys,
				"flags":      strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8023, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method HowEntityByEntityID explains how an entity was constructed from its records.

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) HowEntityByEntityID(ctx context.Context, entityID int64, flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(53, entityID, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(54, entityID, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.howEntityByEntityID(ctx, entityID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityID": formatEntityID(entityID),
				"flags":    strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8024, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method PrimeEngine pre-loads engine resources.

Explicitly calling this method ensures the performance cost is incurred at a predictable time rather than unexpectedly
with the first call requiring the resource.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szengine) PrimeEngine(ctx context.Context) error {
	var err error

	if client.isTrace {
		client.traceEntry(57)

		entryTime := time.Now()

		defer func() { client.traceExit(58, err, time.Since(entryTime)) }()
	}

	err = client.primeEngine(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8026, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method ProcessRedoRecord processes the provided redo record.

This operation performs entity resolution.

Calling processRedoRecord() has the potential to create more redo records in certain situations.

See also getRedoRecord() countRedoRecords().

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document.
*/
func (client *Szengine) ProcessRedoRecord(ctx context.Context, redoRecord string, flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(59, redoRecord, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(60, redoRecord, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.processRedoRecord(ctx, redoRecord, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"flags": strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8027, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method ReevaluateEntity reevaluates an entity by entity ID.

This operation performs entity resolution.

If the entity is not found, then no changes are made.

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.
  - flags: Flags used to control information returned.
*/
func (client *Szengine) ReevaluateEntity(ctx context.Context, entityID int64, flags int64) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(61, entityID, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(62, entityID, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.reevaluateEntity(ctx, entityID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityID": formatEntityID(entityID),
				"flags":    strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8028, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method ReevaluateRecord reevaluates an entity by record ID.

This operation performs entity resolution.

If the record is not found, then no changes are made.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.
*/
func (client *Szengine) ReevaluateRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(63, dataSourceCode, recordID, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(64, dataSourceCode, recordID, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.reevaluateRecord(ctx, dataSourceCode, recordID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
				"flags":          strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8029, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method SearchByAttributes searches for entities that match or relate to the provided attributes.

The default search profile is SEARCH. Alternatively, INGEST may be used.

Input
  - ctx: A context to control lifecycle.
  - attributes: A JSON document containing the attributes desired in the result set.
    Example: `{"NAME_FULL": "BOB SMITH", "EMAIL_ADDRESS": "bsmith@work.com"}`
  - searchProfile: The name of the search profile to use in the search.
    An empty string will use the default search profile.
    Example: "SEARCH"
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) SearchByAttributes(
	ctx context.Context,
	attributes string,
	searchProfile string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(69, attributes, searchProfile, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(70, attributes, searchProfile, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.searchByAttributes(ctx, attributes, searchProfile, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"attributes":    attributes,
				"searchProfile": searchProfile,
				"flags":         strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8031, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method WhyEntities describes the ways two entities relate to each other.

Input
  - ctx: A context to control lifecycle.
  - entityID1: The first of two entity IDs.
  - entityID2: The second of two entity IDs.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) WhyEntities(
	ctx context.Context,
	entityID1 int64,
	entityID2 int64,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(71, entityID1, entityID2, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(72, entityID1, entityID2, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.whyEntities(ctx, entityID1, entityID2, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"entityID1": formatEntityID(entityID1),
				"entityID2": formatEntityID(entityID2),
				"flags":     strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8032, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method WhyRecordInEntity describes the ways a record relates to the rest of its respective entity.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Identifies the provenance of the data.
  - recordID: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) WhyRecordInEntity(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(73, dataSourceCode, recordID, flags)

		entryTime := time.Now()

		defer func() { client.traceExit(74, dataSourceCode, recordID, flags, result, err, time.Since(entryTime)) }()
	}

	result, err = client.whyRecordInEntity(ctx, dataSourceCode, recordID, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"recordID":       recordID,
				"flags":          strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8033, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method WhyRecords describes the ways two records relate to each other.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode1: Identifies the provenance of the data.
  - recordID1: The unique identifier within the records of the same data source.
  - dataSourceCode2: Identifies the provenance of the data.
  - recordID2: The unique identifier within the records of the same data source.
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) WhyRecords(
	ctx context.Context,
	dataSourceCode1 string,
	recordID1 string,
	dataSourceCode2 string,
	recordID2 string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(75, dataSourceCode1, recordID1, dataSourceCode2, recordID2, flags)

		entryTime := time.Now()

		defer func() {
			client.traceExit(
				76,
				dataSourceCode1,
				recordID1,
				dataSourceCode2,
				recordID2,
				flags,
				result,
				err,
				time.Since(entryTime),
			)
		}()
	}

	result, err = client.whyRecords(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode1": dataSourceCode1,
				"recordID1":       recordID1,
				"dataSourceCode2": dataSourceCode2,
				"recordID2":       recordID2,
				"flags":           strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8034, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method WhySearch describes the ways a set of search attributes relate to an entity.

The default search profile is SEARCH. Alternatively, INGEST may be used.

Input
  - ctx: A context to control lifecycle.
  - attributes: A JSON document containing the attributes desired in the result set.
    Example: `{"NAME_FULL": "BOB SMITH", "EMAIL_ADDRESS": "bsmith@work.com"}`
  - entityID:
  - searchProfile: The name of the search profile to use in the search.
    An empty string will use the default search profile.
    Example: "SEARCH"
  - flags: Flags used to control information returned.

Output
  - A JSON document.
*/
func (client *Szengine) WhySearch(
	ctx context.Context,
	attributes string,
	entityID int64,
	searchProfile string,
	flags int64,
) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(69, attributes, entityID, searchProfile, flags)

		entryTime := time.Now()

		defer func() {
			client.traceExit(70, attributes, entityID, searchProfile, flags, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.whySearch(ctx, attributes, entityID, searchProfile, flags)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"attributes":    attributes,
				"entityID":      formatEntityID(entityID),
				"searchProfile": searchProfile,
				"flags":         strconv.FormatInt(flags, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8031, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

/*
Method GetObserverOrigin returns the "origin" value of past Observer messages.

Input
  - ctx: A context to control lifecycle.

Output
  - The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szengine) GetObserverOrigin(ctx context.Context) string {
	_ = ctx

	return client.observerOrigin
}

/*
Method Initialize is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - configID: The configuration ID used for the initialization.  0 for current default configuration.
  - verboseLogging: A flag to enable deeper logging of the Sz processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szengine) Initialize(
	ctx context.Context,
	instanceName string,
	settings string,
	configID int64,
	verboseLogging int64,
) error {
	var err error

	if client.isTrace {
		client.traceEntry(55, instanceName, settings, configID, verboseLogging)

		entryTime := time.Now()

		defer func() {
			client.traceExit(56, instanceName, settings, configID, verboseLogging, err, time.Since(entryTime))
		}()
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configID":       strconv.FormatInt(configID, baseTen),
				"instanceName":   instanceName,
				"settings":       settings,
				"verboseLogging": strconv.FormatInt(verboseLogging, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8025, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method RegisterObserver adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szengine) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if client.isTrace {
		client.traceEntry(703, observer.GetObserverID(ctx))

		entryTime := time.Now()

		defer func() { client.traceExit(704, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}

	if client.observers == nil {
		client.observers = &subject.SimpleSubject{}
	}

	err = client.observers.RegisterObserver(ctx, observer)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"observerID": observer.GetObserverID(ctx),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8702, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method Reinitialize re-initializes the Senzing engine with a specific Senzing configuration JSON document identifier.

Input
  - ctx: A context to control lifecycle.
  - configID: The Senzing configuration JSON document identifier used for the initialization.
*/
func (client *Szengine) Reinitialize(ctx context.Context, configID int64) error {
	var err error

	if client.isTrace {
		client.traceEntry(65, configID)

		entryTime := time.Now()

		defer func() { client.traceExit(66, configID, err, time.Since(entryTime)) }()
	}

	err = client.reinitialize(ctx, configID)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configID": strconv.FormatInt(configID, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8030, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method SetLogLevel sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevelName: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *Szengine) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error

	if client.isTrace {
		client.traceEntry(705, logLevelName)

		entryTime := time.Now()

		defer func() { client.traceExit(706, logLevelName, err, time.Since(entryTime)) }()
	}

	if !logging.IsValidLogLevelName(logLevelName) {
		return wraperror.Errorf(szerror.ErrSzSdk, "invalid error level: %s", logLevelName)
	}

	err = client.getLogger().SetLogLevel(logLevelName)
	client.isTrace = (logLevelName == logging.LevelTraceName)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8703, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method SetObserverOrigin sets the "origin" value in future Observer messages.

Input
  - ctx: A context to control lifecycle.
  - origin: The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szengine) SetObserverOrigin(ctx context.Context, origin string) {
	_ = ctx
	client.observerOrigin = origin
}

/*
Method UnregisterObserver removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szengine) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if client.isTrace {
		client.traceEntry(707, observer.GetObserverID(ctx))

		entryTime := time.Now()

		defer func() { client.traceExit(708, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}

	if client.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8704, err, details)
		err = client.observers.UnregisterObserver(ctx, observer)

		if !client.observers.HasObservers(ctx) {
			client.observers = nil
		}
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private methods for gRPC request/response
// ----------------------------------------------------------------------------

func (client *Szengine) addRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	recordDefinition string,
	flags int64,
) (string, error) {
	request := &szpb.AddRecordRequest{
		DataSourceCode:   dataSourceCode,
		Flags:            flags,
		RecordDefinition: recordDefinition,
		RecordId:         recordID,
	}
	response, err := client.GrpcClient.AddRecord(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) closeExportReport(ctx context.Context, exportHandle uintptr) error {
	request := &szpb.CloseExportReportRequest{
		ExportHandle: int64(exportHandle),
	}
	_, err := client.GrpcClient.CloseExportReport(ctx, request)

	return helper.ConvertGrpcError(err)
}

func (client *Szengine) countRedoRecords(ctx context.Context) (int64, error) {
	request := &szpb.CountRedoRecordsRequest{}
	response, err := client.GrpcClient.CountRedoRecords(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) deleteRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	request := &szpb.DeleteRecordRequest{
		DataSourceCode: dataSourceCode,
		Flags:          flags,
		RecordId:       recordID,
	}
	response, err := client.GrpcClient.DeleteRecord(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) exportCsvEntityReport(ctx context.Context, csvColumnList string, flags int64) (uintptr, error) {
	request := &szpb.ExportCsvEntityReportRequest{
		CsvColumnList: csvColumnList,
		Flags:         flags,
	}
	response, err := client.GrpcClient.ExportCsvEntityReport(ctx, request)
	result := uintptr(response.GetResult())

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) exportJSONEntityReport(ctx context.Context, flags int64) (uintptr, error) {
	request := &szpb.ExportJsonEntityReportRequest{
		Flags: flags,
	}
	response, err := client.GrpcClient.ExportJsonEntityReport(ctx, request)
	result := (uintptr)(response.GetResult())

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) fetchNext(ctx context.Context, exportHandle uintptr) (string, error) {
	request := &szpb.FetchNextRequest{
		ExportHandle: int64(exportHandle),
	}
	response, err := client.GrpcClient.FetchNext(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) findInterestingEntitiesByEntityID(
	ctx context.Context,
	entityID int64,
	flags int64,
) (string, error) {
	request := &szpb.FindInterestingEntitiesByEntityIdRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	response, err := client.GrpcClient.FindInterestingEntitiesByEntityId(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) findInterestingEntitiesByRecordID(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	request := &szpb.FindInterestingEntitiesByRecordIdRequest{
		DataSourceCode: dataSourceCode,
		Flags:          flags,
		RecordId:       recordID,
	}
	response, err := client.GrpcClient.FindInterestingEntitiesByRecordId(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) findNetworkByEntityID(
	ctx context.Context,
	entityIDs string,
	maxDegrees int64,
	buildOutDegree int64,
	buildOutMaxEntities int64,
	flags int64,
) (string, error) {
	request := &szpb.FindNetworkByEntityIdRequest{
		BuildOutDegrees:     buildOutDegree,
		BuildOutMaxEntities: buildOutMaxEntities,
		EntityIds:           entityIDs,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
	}
	response, err := client.GrpcClient.FindNetworkByEntityId(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) findNetworkByRecordID(
	ctx context.Context,
	recordKeys string,
	maxDegrees int64,
	buildOutDegree int64,
	buildOutMaxEntities int64,
	flags int64,
) (string, error) {
	request := &szpb.FindNetworkByRecordIdRequest{
		BuildOutDegrees:     buildOutDegree,
		BuildOutMaxEntities: buildOutMaxEntities,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
		RecordKeys:          recordKeys,
	}
	response, err := client.GrpcClient.FindNetworkByRecordId(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) findPathByEntityID(
	ctx context.Context,
	startEntityID int64,
	endEntityID int64,
	maxDegrees int64,
	avoidEntityIDs string,
	requiredDataSources string,
	flags int64,
) (string, error) {
	request := &szpb.FindPathByEntityIdRequest{
		AvoidEntityIds:      avoidEntityIDs,
		EndEntityId:         endEntityID,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
		RequiredDataSources: requiredDataSources,
		StartEntityId:       startEntityID,
	}
	response, err := client.GrpcClient.FindPathByEntityId(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) findPathByRecordID(
	ctx context.Context,
	startDataSourceCode string,
	startRecordID string,
	endDataSourceCode string,
	endRecordID string,
	maxDegrees int64,
	avoidRecordKeys string,
	requiredDataSources string,
	flags int64,
) (string, error) {
	request := &szpb.FindPathByRecordIdRequest{
		AvoidRecordKeys:     avoidRecordKeys,
		EndDataSourceCode:   endDataSourceCode,
		EndRecordId:         endRecordID,
		Flags:               flags,
		MaxDegrees:          maxDegrees,
		RequiredDataSources: requiredDataSources,
		StartDataSourceCode: startDataSourceCode,
		StartRecordId:       startRecordID,
	}
	response, err := client.GrpcClient.FindPathByRecordId(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) getActiveConfigID(ctx context.Context) (int64, error) {
	request := &szpb.GetActiveConfigIdRequest{}
	response, err := client.GrpcClient.GetActiveConfigId(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) getEntityByEntityID(ctx context.Context, entityID int64, flags int64) (string, error) {
	request := &szpb.GetEntityByEntityIdRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	response, err := client.GrpcClient.GetEntityByEntityId(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) getEntityByRecordID(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	request := &szpb.GetEntityByRecordIdRequest{
		DataSourceCode: dataSourceCode,
		Flags:          flags,
		RecordId:       recordID,
	}
	response, err := client.GrpcClient.GetEntityByRecordId(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) getRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	request := &szpb.GetRecordRequest{
		DataSourceCode: dataSourceCode,
		Flags:          flags,
		RecordId:       recordID,
	}
	response, err := client.GrpcClient.GetRecord(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) getRedoRecord(ctx context.Context) (string, error) {
	request := &szpb.GetRedoRecordRequest{}
	response, err := client.GrpcClient.GetRedoRecord(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) getStats(ctx context.Context) (string, error) {
	request := &szpb.GetStatsRequest{}
	response, err := client.GrpcClient.GetStats(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) getVirtualEntityByRecordID(
	ctx context.Context,
	recordKeys string,
	flags int64,
) (string, error) {
	request := &szpb.GetVirtualEntityByRecordIdRequest{
		Flags:      flags,
		RecordKeys: recordKeys,
	}
	response, err := client.GrpcClient.GetVirtualEntityByRecordId(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) howEntityByEntityID(ctx context.Context, entityID int64, flags int64) (string, error) {
	request := &szpb.HowEntityByEntityIdRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	response, err := client.GrpcClient.HowEntityByEntityId(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) getRecordPreview(ctx context.Context, recordDefinition string, flags int64) (string, error) {
	request := &szpb.GetRecordPreviewRequest{
		Flags:            flags,
		RecordDefinition: recordDefinition,
	}
	response, err := client.GrpcClient.GetRecordPreview(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) primeEngine(ctx context.Context) error {
	request := &szpb.PrimeEngineRequest{}
	_, err := client.GrpcClient.PrimeEngine(ctx, request)

	return helper.ConvertGrpcError(err)
}

func (client *Szengine) processRedoRecord(ctx context.Context, redoRecord string, flags int64) (string, error) {
	request := &szpb.ProcessRedoRecordRequest{
		Flags:      flags,
		RedoRecord: redoRecord,
	}
	response, err := client.GrpcClient.ProcessRedoRecord(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) reevaluateEntity(ctx context.Context, entityID int64, flags int64) (string, error) {
	request := &szpb.ReevaluateEntityRequest{
		EntityId: entityID,
		Flags:    flags,
	}
	response, err := client.GrpcClient.ReevaluateEntity(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) reevaluateRecord(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	request := &szpb.ReevaluateRecordRequest{
		DataSourceCode: dataSourceCode,
		Flags:          flags,
		RecordId:       recordID,
	}
	response, err := client.GrpcClient.ReevaluateRecord(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) reinitialize(ctx context.Context, configID int64) error {
	request := &szpb.ReinitializeRequest{
		ConfigId: configID,
	}
	_, err := client.GrpcClient.Reinitialize(ctx, request)

	return helper.ConvertGrpcError(err)
}

func (client *Szengine) searchByAttributes(
	ctx context.Context,
	attributes string,
	searchProfile string,
	flags int64,
) (string, error) {
	request := &szpb.SearchByAttributesRequest{
		Attributes:    attributes,
		Flags:         flags,
		SearchProfile: searchProfile,
	}
	response, err := client.GrpcClient.SearchByAttributes(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) whyEntities(
	ctx context.Context,
	entityID1 int64,
	entityID2 int64,
	flags int64,
) (string, error) {
	request := &szpb.WhyEntitiesRequest{
		EntityId_1: entityID1,
		EntityId_2: entityID2,
		Flags:      flags,
	}
	response, err := client.GrpcClient.WhyEntities(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) whyRecordInEntity(
	ctx context.Context,
	dataSourceCode string,
	recordID string,
	flags int64,
) (string, error) {
	request := &szpb.WhyRecordInEntityRequest{
		DataSourceCode: dataSourceCode,
		Flags:          flags,
		RecordId:       recordID,
	}
	response, err := client.GrpcClient.WhyRecordInEntity(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) whyRecords(
	ctx context.Context,
	dataSourceCode1 string,
	recordID1 string,
	dataSourceCode2 string,
	recordID2 string,
	flags int64,
) (string, error) {
	request := &szpb.WhyRecordsRequest{
		DataSourceCode_1: dataSourceCode1,
		DataSourceCode_2: dataSourceCode2,
		RecordId_1:       recordID1,
		RecordId_2:       recordID2,
		Flags:            flags,
	}
	response, err := client.GrpcClient.WhyRecords(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szengine) whySearch(
	ctx context.Context,
	attributes string,
	entityID int64,
	searchProfile string,
	flags int64,
) (string, error) {
	request := &szpb.WhySearchRequest{
		Attributes:    attributes,
		EntityId:      entityID,
		Flags:         flags,
		SearchProfile: searchProfile,
	}
	response, err := client.GrpcClient.WhySearch(ctx, request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szengine) getLogger() logging.Logging {
	if client.logger == nil {
		client.logger = helper.GetLogger(ComponentID, szengine.IDMessages, baseCallerSkip)
	}

	return client.logger
}

// Trace method entry.
func (client *Szengine) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *Szengine) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

func formatEntityID(entityID int64) string {
	return strconv.FormatInt(entityID, baseTen)
}
