/*
 *
 */

// Package main implements a client for the service.
package g2diagnosticclient

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	pb "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (client *G2diagnosticClient) getLogger() messagelogger.MessageLoggerInterface {
	if client.logger == nil {
		client.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return client.logger
}

func (client *G2diagnosticClient) notify(ctx context.Context, messageId int, err error, details map[string]string) {
	now := time.Now()
	details["subjectId"] = strconv.Itoa(ProductId)
	details["messageId"] = strconv.Itoa(messageId)
	details["messageTime"] = strconv.FormatInt(now.UnixNano(), 10)
	if err != nil {
		details["error"] = err.Error()
	}
	message, err := json.Marshal(details)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		client.observers.NotifyObservers(ctx, string(message))
	}
}

// Trace method entry.
func (client *G2diagnosticClient) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *G2diagnosticClient) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The CheckDBPerf method performs inserts to determine rate of insertion.

Input
  - ctx: A context to control lifecycle.
  - secondsToRun: Duration of the test in seconds.

Output

  - A string containing a JSON document.
    Example: `{"numRecordsInserted":0,"insertTime":0}`
*/
func (client *G2diagnosticClient) CheckDBPerf(ctx context.Context, secondsToRun int) (string, error) {
	if client.isTrace {
		client.traceEntry(1, secondsToRun)
	}
	entryTime := time.Now()
	request := pb.CheckDBPerfRequest{
		SecondsToRun: int32(secondsToRun),
	}
	response, err := client.GrpcClient.CheckDBPerf(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8001, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(2, secondsToRun, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The CloseEntityListBySize method closes the list created by GetEntityListBySize().
It is part of the GetEntityListBySize(), FetchNextEntityBySize(), CloseEntityListBySize()
lifecycle of a list of sized entities.
The entityListBySizeHandle is created by the GetEntityListBySize() method.

Input
  - ctx: A context to control lifecycle.
  - entityListBySizeHandle: A handle created by GetEntityListBySize().
*/
func (client *G2diagnosticClient) CloseEntityListBySize(ctx context.Context, entityListBySizeHandle uintptr) error {
	if client.isTrace {
		client.traceEntry(5)
	}
	entryTime := time.Now()
	request := pb.CloseEntityListBySizeRequest{
		EntityListBySizeHandle: fmt.Sprintf("%v", entityListBySizeHandle),
	}
	_, err := client.GrpcClient.CloseEntityListBySize(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8002, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(6, err, time.Since(entryTime))
	}
	return err
}

/*
The Destroy method will destroy and perform cleanup for the Senzing G2Diagnostic object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2diagnosticClient) Destroy(ctx context.Context) error {
	if client.isTrace {
		client.traceEntry(7)
	}
	entryTime := time.Now()
	request := pb.DestroyRequest{}
	_, err := client.GrpcClient.Destroy(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8003, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(8, err, time.Since(entryTime))
	}
	return err
}

/*
The FetchNextEntityBySize method gets the next section of the list created by GetEntityListBySize().
It is part of the GetEntityListBySize(), FetchNextEntityBySize(), CloseEntityListBySize()
lifecycle of a list of sized entities.
The entityListBySizeHandle is created by the GetEntityListBySize() method.

Input
  - ctx: A context to control lifecycle.
  - entityListBySizeHandle: A handle created by GetEntityListBySize().

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnosticClient) FetchNextEntityBySize(ctx context.Context, entityListBySizeHandle uintptr) (string, error) {
	if client.isTrace {
		client.traceEntry(9)
	}
	entryTime := time.Now()
	request := pb.FetchNextEntityBySizeRequest{
		EntityListBySizeHandle: fmt.Sprintf("%v", entityListBySizeHandle),
	}
	response, err := client.GrpcClient.FetchNextEntityBySize(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8004, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(10, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The FindEntitiesByFeatureIDs method finds entities having any of the lib feat id specified in the "features" JSON document.
The "features" also contains an entity id.
This entity is ignored in the returned values.

Input
  - ctx: A context to control lifecycle.
  - features: A JSON document having the format: `{"ENTITY_ID":<entity id>,"LIB_FEAT_IDS":[<id1>,<id2>,...<idn>]}` where ENTITY_ID specifies the entity to ignore in the returns and <id#> are the lib feat ids used to query for entities.

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnosticClient) FindEntitiesByFeatureIDs(ctx context.Context, features string) (string, error) {
	if client.isTrace {
		client.traceEntry(11, features)
	}
	entryTime := time.Now()
	request := pb.FindEntitiesByFeatureIDsRequest{
		Features: features,
	}
	response, err := client.GrpcClient.FindEntitiesByFeatureIDs(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8005, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(12, features, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The GetAvailableMemory method returns the available memory, in bytes, on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of bytes of available memory.
*/
func (client *G2diagnosticClient) GetAvailableMemory(ctx context.Context) (int64, error) {
	if client.isTrace {
		client.traceEntry(13)
	}
	entryTime := time.Now()
	request := pb.GetAvailableMemoryRequest{}
	response, err := client.GrpcClient.GetAvailableMemory(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8006, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(14, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The GetDataSourceCounts method returns information about data sources.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document enumerating data sources.
    See the example output.
*/
func (client *G2diagnosticClient) GetDataSourceCounts(ctx context.Context) (string, error) {
	if client.isTrace {
		client.traceEntry(15)
	}
	entryTime := time.Now()
	request := pb.GetDataSourceCountsRequest{}
	response, err := client.GrpcClient.GetDataSourceCounts(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8007, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(16, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The GetDBInfo method returns information about the database connection.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document enumerating data sources.
    Example: `{"Hybrid Mode":false,"Database Details":[{"Name":"0.0.0.0","Type":"postgresql"}]}`
*/
func (client *G2diagnosticClient) GetDBInfo(ctx context.Context) (string, error) {
	if client.isTrace {
		client.traceEntry(17)
	}
	entryTime := time.Now()
	request := pb.GetDBInfoRequest{}
	response, err := client.GrpcClient.GetDBInfo(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8008, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(18, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The GetEntityDetails method returns information about the database connection.

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.
  - includeInternalFeatures: FIXME:

Output
  - A JSON document enumerating FIXME:.
    See the example output.
*/
func (client *G2diagnosticClient) GetEntityDetails(ctx context.Context, entityID int64, includeInternalFeatures int) (string, error) {
	if client.isTrace {
		client.traceEntry(19, entityID, includeInternalFeatures)
	}
	entryTime := time.Now()
	request := pb.GetEntityDetailsRequest{
		EntityID:                entityID,
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.GrpcClient.GetEntityDetails(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8009, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(20, entityID, includeInternalFeatures, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The GetEntityListBySize method gets the next section of the list created by GetEntityListBySize().
It is part of the GetEntityListBySize(), FetchNextEntityBySize(), CloseEntityListBySize()
lifecycle of a list of sized entities.
The entityListBySizeHandle is used by the FetchNextEntityBySize() and CloseEntityListBySize() methods.

Input
  - ctx: A context to control lifecycle.
  - entitySize: FIXME:

Output
  - A handle to an entity list to be used with FetchNextEntityBySize() and CloseEntityListBySize().
*/
func (client *G2diagnosticClient) GetEntityListBySize(ctx context.Context, entitySize int) (uintptr, error) {
	if client.isTrace {
		client.traceEntry(21, entitySize)
	}
	entryTime := time.Now()
	request := pb.GetEntityListBySizeRequest{
		EntitySize: int32(entitySize),
	}
	response, err := client.GrpcClient.GetEntityListBySize(ctx, &request)
	if err != nil {
		return 0, err
	}
	result := response.GetResult()
	result_int, err := strconv.Atoi(result)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8010, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(22, entitySize, (uintptr)(result_int), err, time.Since(entryTime))
	}
	return uintptr(result_int), err
}

/*
The GetEntityResume method FIXME:

Input
  - ctx: A context to control lifecycle.
  - entityID: The unique identifier of an entity.

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnosticClient) GetEntityResume(ctx context.Context, entityID int64) (string, error) {
	if client.isTrace {
		client.traceEntry(23, entityID)
	}
	entryTime := time.Now()
	request := pb.GetEntityResumeRequest{
		EntityID: entityID,
	}
	response, err := client.GrpcClient.GetEntityResume(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8011, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(24, entityID, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The GetEntitySizeBreakdown method FIXME:

Input
  - ctx: A context to control lifecycle.
  - minimumEntitySize: FIXME:
  - includeInternalFeatures: FIXME:

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnosticClient) GetEntitySizeBreakdown(ctx context.Context, minimumEntitySize int, includeInternalFeatures int) (string, error) {
	if client.isTrace {
		client.traceEntry(25, minimumEntitySize, includeInternalFeatures)
	}
	entryTime := time.Now()
	request := pb.GetEntitySizeBreakdownRequest{
		MinimumEntitySize:       int32(minimumEntitySize),
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.GrpcClient.GetEntitySizeBreakdown(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8012, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(26, minimumEntitySize, includeInternalFeatures, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The GetFeature method retrieves a stored feature.

Input
  - ctx: A context to control lifecycle.
  - libFeatID: The identifier of the feature requested in the search.

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnosticClient) GetFeature(ctx context.Context, libFeatID int64) (string, error) {
	if client.isTrace {
		client.traceEntry(27, libFeatID)
	}
	entryTime := time.Now()
	request := pb.GetFeatureRequest{
		LibFeatID: libFeatID,
	}
	response, err := client.GrpcClient.GetFeature(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8013, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(28, libFeatID, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The GetGenericFeatures method retrieves a stored feature.

Input
  - ctx: A context to control lifecycle.
  - featureType: FIXME:
  - maximumEstimatedCount: FIXME:

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnosticClient) GetGenericFeatures(ctx context.Context, featureType string, maximumEstimatedCount int) (string, error) {
	if client.isTrace {
		client.traceEntry(29, featureType, maximumEstimatedCount)
	}
	entryTime := time.Now()
	request := pb.GetGenericFeaturesRequest{
		FeatureType:           featureType,
		MaximumEstimatedCount: int32(maximumEstimatedCount),
	}
	response, err := client.GrpcClient.GetGenericFeatures(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8014, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(30, featureType, maximumEstimatedCount, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The GetLogicalCores method returns the number of logical cores on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of logical cores.
*/
func (client *G2diagnosticClient) GetLogicalCores(ctx context.Context) (int, error) {
	if client.isTrace {
		client.traceEntry(35)
	}
	entryTime := time.Now()
	request := pb.GetLogicalCoresRequest{}
	response, err := client.GrpcClient.GetLogicalCores(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8015, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(36, response.GetResult(), err, time.Since(entryTime))
	}
	return int(response.GetResult()), err
}

/*
The GetMappingStatistics method FIXME:

Input
  - ctx: A context to control lifecycle.
  - includeInternalFeatures: FIXME:

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnosticClient) GetMappingStatistics(ctx context.Context, includeInternalFeatures int) (string, error) {
	if client.isTrace {
		client.traceEntry(37, includeInternalFeatures)
	}
	entryTime := time.Now()
	request := pb.GetMappingStatisticsRequest{
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.GrpcClient.GetMappingStatistics(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8016, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(38, includeInternalFeatures, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The GetPhysicalCores method returns the number of physical cores on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of physical cores.
*/
func (client *G2diagnosticClient) GetPhysicalCores(ctx context.Context) (int, error) {
	if client.isTrace {
		client.traceEntry(39)
	}
	entryTime := time.Now()
	request := pb.GetPhysicalCoresRequest{}
	response, err := client.GrpcClient.GetPhysicalCores(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8017, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(40, response.GetResult(), err, time.Since(entryTime))
	}
	return int(response.GetResult()), err
}

/*
The GetRelationshipDetails method FIXME:

Input
  - ctx: A context to control lifecycle.
  - relationshipID: FIXME:
  - includeInternalFeatures: FIXME:

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnosticClient) GetRelationshipDetails(ctx context.Context, relationshipID int64, includeInternalFeatures int) (string, error) {
	if client.isTrace {
		client.traceEntry(41, relationshipID, includeInternalFeatures)
	}
	entryTime := time.Now()
	request := pb.GetRelationshipDetailsRequest{
		RelationshipID:          relationshipID,
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.GrpcClient.GetRelationshipDetails(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8018, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(42, relationshipID, includeInternalFeatures, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The GetResolutionStatistics method FIXME:

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnosticClient) GetResolutionStatistics(ctx context.Context) (string, error) {
	if client.isTrace {
		client.traceEntry(43)
	}
	entryTime := time.Now()
	request := pb.GetResolutionStatisticsRequest{}
	response, err := client.GrpcClient.GetResolutionStatistics(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8019, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(44, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The GetTotalSystemMemory method returns the total memory, in bytes, on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of bytes of memory.
*/
func (client *G2diagnosticClient) GetTotalSystemMemory(ctx context.Context) (int64, error) {
	if client.isTrace {
		client.traceEntry(57)
	}
	entryTime := time.Now()
	request := pb.GetTotalSystemMemoryRequest{}
	response, err := client.GrpcClient.GetTotalSystemMemory(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8020, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(46, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The Init method initializes the Senzing G2Diagnosis object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - moduleName: A name for the auditing node, to help identify it within system logs.
  - iniParams: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *G2diagnosticClient) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	if client.isTrace {
		client.traceEntry(47, moduleName, iniParams, verboseLogging)
	}
	entryTime := time.Now()
	request := pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	_, err := client.GrpcClient.Init(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"iniParams":      iniParams,
				"moduleName":     moduleName,
				"verboseLogging": strconv.Itoa(verboseLogging),
			}
			client.notify(ctx, 8021, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(48, moduleName, iniParams, verboseLogging, err, time.Since(entryTime))
	}
	return err
}

/*
The InitWithConfigID method initializes the Senzing G2Diagnosis object with a non-default configuration ID.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - moduleName: A name for the auditing node, to help identify it within system logs.
  - iniParams: A JSON string containing configuration parameters.
  - initConfigID: The configuration ID used for the initialization.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *G2diagnosticClient) InitWithConfigID(ctx context.Context, moduleName string, iniParams string, initConfigID int64, verboseLogging int) error {
	if client.isTrace {
		client.traceEntry(49, moduleName, iniParams, initConfigID, verboseLogging)
	}
	entryTime := time.Now()
	request := pb.InitWithConfigIDRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		InitConfigID:   initConfigID,
		VerboseLogging: int32(verboseLogging),
	}
	_, err := client.GrpcClient.InitWithConfigID(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"iniParams":      iniParams,
				"initConfigID":   strconv.FormatInt(initConfigID, 10),
				"moduleName":     moduleName,
				"verboseLogging": strconv.Itoa(verboseLogging),
			}
			client.notify(ctx, 8022, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(50, moduleName, iniParams, initConfigID, verboseLogging, err, time.Since(entryTime))
	}
	return err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2diagnosticClient) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	if client.observers == nil {
		client.observers = &subject.SubjectImpl{}
	}
	return client.observers.RegisterObserver(ctx, observer)
}

/*
The Reinit method re-initializes the Senzing G2Diagnosis object.

Input
  - ctx: A context to control lifecycle.
  - initConfigID: The configuration ID used for the initialization.
*/
func (client *G2diagnosticClient) Reinit(ctx context.Context, initConfigID int64) error {
	if client.isTrace {
		client.traceEntry(51, initConfigID)
	}
	entryTime := time.Now()
	request := pb.ReinitRequest{
		InitConfigID: initConfigID,
	}
	_, err := client.GrpcClient.Reinit(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"initConfigID": strconv.FormatInt(initConfigID, 10),
			}
			client.notify(ctx, 8023, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(52, initConfigID, err, time.Since(entryTime))
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *G2diagnosticClient) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if client.isTrace {
		client.traceEntry(53, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	if client.isTrace {
		defer client.traceExit(54, logLevel, err, time.Since(entryTime))
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2diagnosticClient) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	err := client.observers.UnregisterObserver(ctx, observer)
	if err != nil {
		return err
	}
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	return err
}
