/*
 *
 */

// Package main implements a client for the service.
package g2diagnostic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/senzing/g2-sdk-go-grpc/helper"
	g2diagnosticapi "github.com/senzing/g2-sdk-go/g2diagnostic"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-observing/notifier"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2diagnostic struct {
	GrpcClient g2pb.G2DiagnosticClient
	isTrace    bool // Performance optimization
	logger     logging.LoggingInterface
	observers  subject.Subject
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *G2diagnostic) getLogger() logging.LoggingInterface {
	var err error = nil
	if client.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		client.logger, err = logging.NewSenzingSdkLogger(ProductId, g2diagnosticapi.IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return client.logger
}

// Trace method entry.
func (client *G2diagnostic) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *G2diagnostic) traceExit(errorNumber int, details ...interface{}) {
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
func (client *G2diagnostic) CheckDBPerf(ctx context.Context, secondsToRun int) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(1, secondsToRun)
		defer func() { client.traceExit(2, secondsToRun, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.CheckDBPerfRequest{
		SecondsToRun: int32(secondsToRun),
	}
	response, err := client.GrpcClient.CheckDBPerf(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8001, err, details)
		}()
	}
	return result, err
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
func (client *G2diagnostic) CloseEntityListBySize(ctx context.Context, entityListBySizeHandle uintptr) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(5)
		defer func() { client.traceExit(6, err, time.Since(entryTime)) }()
	}
	request := g2pb.CloseEntityListBySizeRequest{
		EntityListBySizeHandle: fmt.Sprintf("%v", entityListBySizeHandle),
	}
	_, err = client.GrpcClient.CloseEntityListBySize(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8002, err, details)
		}()
	}
	return err
}

/*
The Destroy method will destroy and perform cleanup for the Senzing G2Diagnostic object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2diagnostic) Destroy(ctx context.Context) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(7)
		defer func() { client.traceExit(8, err, time.Since(entryTime)) }()
	}
	request := g2pb.DestroyRequest{}
	_, err = client.GrpcClient.Destroy(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8003, err, details)
		}()
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
func (client *G2diagnostic) FetchNextEntityBySize(ctx context.Context, entityListBySizeHandle uintptr) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(9)
		defer func() { client.traceExit(10, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.FetchNextEntityBySizeRequest{
		EntityListBySizeHandle: fmt.Sprintf("%v", entityListBySizeHandle),
	}
	response, err := client.GrpcClient.FetchNextEntityBySize(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8004, err, details)
		}()
	}
	return result, err
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
func (client *G2diagnostic) FindEntitiesByFeatureIDs(ctx context.Context, features string) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(11, features)
		defer func() { client.traceExit(12, features, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.FindEntitiesByFeatureIDsRequest{
		Features: features,
	}
	response, err := client.GrpcClient.FindEntitiesByFeatureIDs(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8005, err, details)
		}()
	}
	return result, err
}

/*
The GetAvailableMemory method returns the available memory, in bytes, on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of bytes of available memory.
*/
func (client *G2diagnostic) GetAvailableMemory(ctx context.Context) (int64, error) {
	var err error = nil
	var result int64 = 0
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(13)
		defer func() { client.traceExit(14, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetAvailableMemoryRequest{}
	response, err := client.GrpcClient.GetAvailableMemory(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8006, err, details)
		}()
	}
	return result, err
}

/*
The GetDataSourceCounts method returns information about data sources.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document enumerating data sources.
    See the example output.
*/
func (client *G2diagnostic) GetDataSourceCounts(ctx context.Context) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(15)
		defer func() { client.traceExit(16, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetDataSourceCountsRequest{}
	response, err := client.GrpcClient.GetDataSourceCounts(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8007, err, details)
		}()
	}
	return result, err
}

/*
The GetDBInfo method returns information about the database connection.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document enumerating data sources.
    Example: `{"Hybrid Mode":false,"Database Details":[{"Name":"0.0.0.0","Type":"postgresql"}]}`
*/
func (client *G2diagnostic) GetDBInfo(ctx context.Context) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(17)
		defer func() { client.traceExit(18, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetDBInfoRequest{}
	response, err := client.GrpcClient.GetDBInfo(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8008, err, details)
		}()
	}
	return result, err
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
func (client *G2diagnostic) GetEntityDetails(ctx context.Context, entityID int64, includeInternalFeatures int) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(19, entityID, includeInternalFeatures)
		defer func() { client.traceExit(20, entityID, includeInternalFeatures, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetEntityDetailsRequest{
		EntityID:                entityID,
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.GrpcClient.GetEntityDetails(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8009, err, details)
		}()
	}
	return result, err
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
func (client *G2diagnostic) GetEntityListBySize(ctx context.Context, entitySize int) (uintptr, error) {
	var err error = nil
	var result uintptr = 0
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(21, entitySize)
		defer func() { client.traceExit(22, entitySize, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetEntityListBySizeRequest{
		EntitySize: int32(entitySize),
	}
	response, err := client.GrpcClient.GetEntityListBySize(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if err != nil {
		return 0, err
	}
	resultInt, err := strconv.Atoi(response.GetResult())
	if err != nil {
		return 0, err
	}
	result = uintptr(resultInt)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8010, err, details)
		}()
	}
	return result, err
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
func (client *G2diagnostic) GetEntityResume(ctx context.Context, entityID int64) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(23, entityID)
		defer func() { client.traceExit(24, entityID, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetEntityResumeRequest{
		EntityID: entityID,
	}
	response, err := client.GrpcClient.GetEntityResume(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8011, err, details)
		}()
	}
	return result, err
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
func (client *G2diagnostic) GetEntitySizeBreakdown(ctx context.Context, minimumEntitySize int, includeInternalFeatures int) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(25, minimumEntitySize, includeInternalFeatures)
		defer func() {
			client.traceExit(26, minimumEntitySize, includeInternalFeatures, result, err, time.Since(entryTime))
		}()
	}
	request := g2pb.GetEntitySizeBreakdownRequest{
		MinimumEntitySize:       int32(minimumEntitySize),
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.GrpcClient.GetEntitySizeBreakdown(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8012, err, details)
		}()
	}
	return result, err
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
func (client *G2diagnostic) GetFeature(ctx context.Context, libFeatID int64) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(27, libFeatID)
		defer func() { client.traceExit(28, libFeatID, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetFeatureRequest{
		LibFeatID: libFeatID,
	}
	response, err := client.GrpcClient.GetFeature(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8013, err, details)
		}()
	}
	return result, err
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
func (client *G2diagnostic) GetGenericFeatures(ctx context.Context, featureType string, maximumEstimatedCount int) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(29, featureType, maximumEstimatedCount)
		defer func() { client.traceExit(30, featureType, maximumEstimatedCount, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetGenericFeaturesRequest{
		FeatureType:           featureType,
		MaximumEstimatedCount: int32(maximumEstimatedCount),
	}
	response, err := client.GrpcClient.GetGenericFeatures(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8014, err, details)
		}()
	}
	return result, err
}

/*
The GetLogicalCores method returns the number of logical cores on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of logical cores.
*/
func (client *G2diagnostic) GetLogicalCores(ctx context.Context) (int, error) {
	var err error = nil
	var result int = 0
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(35)
		defer func() { client.traceExit(36, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetLogicalCoresRequest{}
	response, err := client.GrpcClient.GetLogicalCores(ctx, &request)
	result = int(response.GetResult())
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8015, err, details)
		}()
	}
	return result, err
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
func (client *G2diagnostic) GetMappingStatistics(ctx context.Context, includeInternalFeatures int) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(37, includeInternalFeatures)
		defer func() { client.traceExit(38, includeInternalFeatures, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetMappingStatisticsRequest{
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.GrpcClient.GetMappingStatistics(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8016, err, details)
		}()
	}
	return result, err
}

/*
The GetPhysicalCores method returns the number of physical cores on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of physical cores.
*/
func (client *G2diagnostic) GetPhysicalCores(ctx context.Context) (int, error) {
	var err error = nil
	var result int = 0
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(39)
		defer func() { client.traceExit(40, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetPhysicalCoresRequest{}
	response, err := client.GrpcClient.GetPhysicalCores(ctx, &request)
	result = int(response.GetResult())
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8017, err, details)
		}()
	}
	return result, err
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
func (client *G2diagnostic) GetRelationshipDetails(ctx context.Context, relationshipID int64, includeInternalFeatures int) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(41, relationshipID, includeInternalFeatures)
		defer func() {
			client.traceExit(42, relationshipID, includeInternalFeatures, result, err, time.Since(entryTime))
		}()
	}
	request := g2pb.GetRelationshipDetailsRequest{
		RelationshipID:          relationshipID,
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.GrpcClient.GetRelationshipDetails(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8018, err, details)
		}()
	}
	return result, err
}

/*
The GetResolutionStatistics method FIXME:

Input
  - ctx: A context to control lifecycle.

Output
  - A string containing a JSON document.
    See the example output.
*/
func (client *G2diagnostic) GetResolutionStatistics(ctx context.Context) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(43)
		defer func() { client.traceExit(44, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetResolutionStatisticsRequest{}
	response, err := client.GrpcClient.GetResolutionStatistics(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8019, err, details)
		}()
	}
	return result, err
}

/*
The GetSdkId method returns the identifier of this particular Software Development Kit (SDK).
It is handy when working with multiple implementations of the same G2diagnosticInterface.
For this implementation, "grpc" is returned.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2diagnostic) GetSdkId(ctx context.Context) string {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(59)
		defer func() { client.traceExit(60, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8024, err, details)
		}()
	}
	return "grpc"
}

/*
The GetTotalSystemMemory method returns the total memory, in bytes, on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of bytes of memory.
*/
func (client *G2diagnostic) GetTotalSystemMemory(ctx context.Context) (int64, error) {
	var err error = nil
	var result int64 = 0
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(57)
		defer func() { client.traceExit(46, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetTotalSystemMemoryRequest{}
	response, err := client.GrpcClient.GetTotalSystemMemory(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8020, err, details)
		}()
	}
	return result, err
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
func (client *G2diagnostic) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(47, moduleName, iniParams, verboseLogging)
		defer func() { client.traceExit(48, moduleName, iniParams, verboseLogging, err, time.Since(entryTime)) }()
	}
	request := g2pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	_, err = client.GrpcClient.Init(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"iniParams":      iniParams,
				"moduleName":     moduleName,
				"verboseLogging": strconv.Itoa(verboseLogging),
			}
			notifier.Notify(ctx, client.observers, ProductId, 8021, err, details)
		}()
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
func (client *G2diagnostic) InitWithConfigID(ctx context.Context, moduleName string, iniParams string, initConfigID int64, verboseLogging int) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(49, moduleName, iniParams, initConfigID, verboseLogging)
		defer func() {
			client.traceExit(50, moduleName, iniParams, initConfigID, verboseLogging, err, time.Since(entryTime))
		}()
	}
	request := g2pb.InitWithConfigIDRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		InitConfigID:   initConfigID,
		VerboseLogging: int32(verboseLogging),
	}
	_, err = client.GrpcClient.InitWithConfigID(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"iniParams":      iniParams,
				"initConfigID":   strconv.FormatInt(initConfigID, 10),
				"moduleName":     moduleName,
				"verboseLogging": strconv.Itoa(verboseLogging),
			}
			notifier.Notify(ctx, client.observers, ProductId, 8022, err, details)
		}()
	}
	return err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2diagnostic) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(55, observer.GetObserverId(ctx))
		defer func() { client.traceExit(56, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	if client.observers == nil {
		client.observers = &subject.SubjectImpl{}
	}
	err = client.observers.RegisterObserver(ctx, observer)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"observerID": observer.GetObserverId(ctx),
			}
			notifier.Notify(ctx, client.observers, ProductId, 8025, err, details)
		}()
	}
	return err
}

/*
The Reinit method re-initializes the Senzing G2Diagnosis object.

Input
  - ctx: A context to control lifecycle.
  - initConfigID: The configuration ID used for the initialization.
*/
func (client *G2diagnostic) Reinit(ctx context.Context, initConfigID int64) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(51, initConfigID)
		defer func() { client.traceExit(52, initConfigID, err, time.Since(entryTime)) }()
	}
	request := g2pb.ReinitRequest{
		InitConfigID: initConfigID,
	}
	_, err = client.GrpcClient.Reinit(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"initConfigID": strconv.FormatInt(initConfigID, 10),
			}
			notifier.Notify(ctx, client.observers, ProductId, 8023, err, details)
		}()
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevelName: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *G2diagnostic) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(53, logLevelName)
		defer func() { client.traceExit(54, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	client.getLogger().SetLogLevel(logLevelName)
	client.isTrace = (logLevelName == logging.LevelTraceName)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"logLevel": logLevelName,
			}
			notifier.Notify(ctx, client.observers, ProductId, 8026, err, details)
		}()
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2diagnostic) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(57, observer.GetObserverId(ctx))
		defer client.traceExit(58, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	if client.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, client.observers, ProductId, 8027, err, details)
	}
	err = client.observers.UnregisterObserver(ctx, observer)
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	return err
}
