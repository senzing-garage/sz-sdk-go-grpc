/*
 *
 */

// Package main implements a client for the service.
package g2diagnosticclient

import (
	"context"
	"fmt"

	pb "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
)

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
	request := pb.CheckDBPerfRequest{
		SecondsToRun: int32(secondsToRun),
	}
	response, err := client.G2DiagnosticGrpcClient.CheckDBPerf(ctx, &request)
	result := response.GetResult()
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
func (client *G2diagnosticClient) CloseEntityListBySize(ctx context.Context, entityListBySizeHandle interface{}) error {
	request := pb.CloseEntityListBySizeRequest{
		EntityListBySizeHandle: fmt.Sprintf("%v", entityListBySizeHandle),
	}
	_, err := client.G2DiagnosticGrpcClient.CloseEntityListBySize(ctx, &request)
	return err
}

/*
The Destroy method will destroy and perform cleanup for the Senzing G2Diagnostic object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2diagnosticClient) Destroy(ctx context.Context) error {
	request := pb.DestroyRequest{}
	_, err := client.G2DiagnosticGrpcClient.Destroy(ctx, &request)
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
func (client *G2diagnosticClient) FetchNextEntityBySize(ctx context.Context, entityListBySizeHandle interface{}) (string, error) {
	request := pb.FetchNextEntityBySizeRequest{
		EntityListBySizeHandle: fmt.Sprintf("%v", entityListBySizeHandle),
	}
	response, err := client.G2DiagnosticGrpcClient.FetchNextEntityBySize(ctx, &request)
	result := response.GetResult()
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
func (client *G2diagnosticClient) FindEntitiesByFeatureIDs(ctx context.Context, features string) (string, error) {
	request := pb.FindEntitiesByFeatureIDsRequest{
		Features: features,
	}
	response, err := client.G2DiagnosticGrpcClient.FindEntitiesByFeatureIDs(ctx, &request)
	result := response.GetResult()
	return result, err
}

/*
The GetAvailableMemory method returns the available memory, in bytes, on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of bytes of available memory.
*/
func (client *G2diagnosticClient) GetAvailableMemory(ctx context.Context) (int64, error) {
	request := pb.GetAvailableMemoryRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetAvailableMemory(ctx, &request)
	result := response.GetResult()
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
func (client *G2diagnosticClient) GetDataSourceCounts(ctx context.Context) (string, error) {
	request := pb.GetDataSourceCountsRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetDataSourceCounts(ctx, &request)
	result := response.GetResult()
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
func (client *G2diagnosticClient) GetDBInfo(ctx context.Context) (string, error) {
	request := pb.GetDBInfoRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetDBInfo(ctx, &request)
	result := response.GetResult()
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
func (client *G2diagnosticClient) GetEntityDetails(ctx context.Context, entityID int64, includeInternalFeatures int) (string, error) {
	request := pb.GetEntityDetailsRequest{
		EntityID:                entityID,
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.G2DiagnosticGrpcClient.GetEntityDetails(ctx, &request)
	result := response.GetResult()
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
func (client *G2diagnosticClient) GetEntityListBySize(ctx context.Context, entitySize int) (interface{}, error) {
	request := pb.GetEntityListBySizeRequest{
		EntitySize: int32(entitySize),
	}
	response, err := client.G2DiagnosticGrpcClient.GetEntityListBySize(ctx, &request)
	result := response.GetResult()
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
func (client *G2diagnosticClient) GetEntityResume(ctx context.Context, entityID int64) (string, error) {
	request := pb.GetEntityResumeRequest{
		EntityID: entityID,
	}
	response, err := client.G2DiagnosticGrpcClient.GetEntityResume(ctx, &request)
	result := response.GetResult()
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
func (client *G2diagnosticClient) GetEntitySizeBreakdown(ctx context.Context, minimumEntitySize int, includeInternalFeatures int) (string, error) {
	request := pb.GetEntitySizeBreakdownRequest{
		MinimumEntitySize:       int32(minimumEntitySize),
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.G2DiagnosticGrpcClient.GetEntitySizeBreakdown(ctx, &request)
	result := response.GetResult()
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
func (client *G2diagnosticClient) GetFeature(ctx context.Context, libFeatID int64) (string, error) {
	request := pb.GetFeatureRequest{
		LibFeatID: libFeatID,
	}
	response, err := client.G2DiagnosticGrpcClient.GetFeature(ctx, &request)
	result := response.GetResult()
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
func (client *G2diagnosticClient) GetGenericFeatures(ctx context.Context, featureType string, maximumEstimatedCount int) (string, error) {
	request := pb.GetGenericFeaturesRequest{
		FeatureType:           featureType,
		MaximumEstimatedCount: int32(maximumEstimatedCount),
	}
	response, err := client.G2DiagnosticGrpcClient.GetGenericFeatures(ctx, &request)
	result := response.GetResult()
	return result, err
}

/*
The GetLogicalCores method returns the number of logical cores on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of logical cores.
*/
func (client *G2diagnosticClient) GetLogicalCores(ctx context.Context) (int, error) {
	request := pb.GetLogicalCoresRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetLogicalCores(ctx, &request)
	result := int(response.GetResult())
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
func (client *G2diagnosticClient) GetMappingStatistics(ctx context.Context, includeInternalFeatures int) (string, error) {
	request := pb.GetMappingStatisticsRequest{
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.G2DiagnosticGrpcClient.GetMappingStatistics(ctx, &request)
	result := response.GetResult()
	return result, err
}

/*
The GetPhysicalCores method returns the number of physical cores on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of physical cores.
*/
func (client *G2diagnosticClient) GetPhysicalCores(ctx context.Context) (int, error) {
	request := pb.GetPhysicalCoresRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetPhysicalCores(ctx, &request)
	result := int(response.GetResult())
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
func (client *G2diagnosticClient) GetRelationshipDetails(ctx context.Context, relationshipID int64, includeInternalFeatures int) (string, error) {
	request := pb.GetRelationshipDetailsRequest{
		RelationshipID:          relationshipID,
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.G2DiagnosticGrpcClient.GetRelationshipDetails(ctx, &request)
	result := response.GetResult()
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
func (client *G2diagnosticClient) GetResolutionStatistics(ctx context.Context) (string, error) {
	request := pb.GetResolutionStatisticsRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetResolutionStatistics(ctx, &request)
	result := response.GetResult()
	return result, err
}

/*
The GetTotalSystemMemory method returns the total memory, in bytes, on the host system.

Input
  - ctx: A context to control lifecycle.

Output
  - Number of bytes of memory.
*/
func (client *G2diagnosticClient) GetTotalSystemMemory(ctx context.Context) (int64, error) {
	request := pb.GetTotalSystemMemoryRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetTotalSystemMemory(ctx, &request)
	result := response.GetResult()
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
func (client *G2diagnosticClient) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	request := pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	_, err := client.G2DiagnosticGrpcClient.Init(ctx, &request)
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
	request := pb.InitWithConfigIDRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		InitConfigID:   initConfigID,
		VerboseLogging: int32(verboseLogging),
	}
	_, err := client.G2DiagnosticGrpcClient.InitWithConfigID(ctx, &request)
	return err
}

/*
The Reinit method re-initializes the Senzing G2Diagnosis object.

Input
  - ctx: A context to control lifecycle.
  - initConfigID: The configuration ID used for the initialization.
*/
func (client *G2diagnosticClient) Reinit(ctx context.Context, initConfigID int64) error {
	request := pb.ReinitRequest{
		InitConfigID: initConfigID,
	}
	_, err := client.G2DiagnosticGrpcClient.Reinit(ctx, &request)
	return err
}
