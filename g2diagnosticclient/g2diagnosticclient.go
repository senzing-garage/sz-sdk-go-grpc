/*
 *
 */

// Package main implements a client for the service.
package g2diagnosticclient

import (
	"context"
	"fmt"

	pb "github.com/senzing/g2-sdk-go-grpc/protobuf/g2diagnostic"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

// CheckDBPerf performs inserts to determine rate of insertion.
func (client *G2diagnosticClient) CheckDBPerf(ctx context.Context, secondsToRun int) (string, error) {
	request := pb.CheckDBPerfRequest{
		SecondsToRun: int32(secondsToRun),
	}
	response, err := client.G2DiagnosticGrpcClient.CheckDBPerf(ctx, &request)
	result := response.GetResult()
	return result, err
}

// TODO: Document.
func (client *G2diagnosticClient) CloseEntityListBySize(ctx context.Context, entityListBySizeHandle interface{}) error {
	request := pb.CloseEntityListBySizeRequest{
		EntityListBySizeHandle: fmt.Sprintf("%v", entityListBySizeHandle),
	}
	_, err := client.G2DiagnosticGrpcClient.CloseEntityListBySize(ctx, &request)
	return err
}

// TODO: Document.
func (client *G2diagnosticClient) Destroy(ctx context.Context) error {
	request := pb.DestroyRequest{}
	_, err := client.G2DiagnosticGrpcClient.Destroy(ctx, &request)
	return err
}

// TODO: Document.
func (client *G2diagnosticClient) FetchNextEntityBySize(ctx context.Context, entityListBySizeHandle interface{}) (string, error) {
	request := pb.FetchNextEntityBySizeRequest{
		EntityListBySizeHandle: fmt.Sprintf("%v", entityListBySizeHandle),
	}
	response, err := client.G2DiagnosticGrpcClient.FetchNextEntityBySize(ctx, &request)
	result := response.GetResult()
	return result, err
}

// TODO: Document.
func (client *G2diagnosticClient) FindEntitiesByFeatureIDs(ctx context.Context, features string) (string, error) {
	request := pb.FindEntitiesByFeatureIDsRequest{
		Features: features,
	}
	response, err := client.G2DiagnosticGrpcClient.FindEntitiesByFeatureIDs(ctx, &request)
	result := response.GetResult()
	return result, err
}

// GetAvailableMemory returns the available memory, in bytes, on the host system.
func (client *G2diagnosticClient) GetAvailableMemory(ctx context.Context) (int64, error) {
	request := pb.GetAvailableMemoryRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetAvailableMemory(ctx, &request)
	result := response.GetResult()
	return result, err
}

// TODO: Document.
func (client *G2diagnosticClient) GetDataSourceCounts(ctx context.Context) (string, error) {
	request := pb.GetDataSourceCountsRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetDataSourceCounts(ctx, &request)
	result := response.GetResult()
	return result, err
}

// GetDBInfo returns information about the database connection.
func (client *G2diagnosticClient) GetDBInfo(ctx context.Context) (string, error) {
	request := pb.GetDBInfoRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetDBInfo(ctx, &request)
	result := response.GetResult()
	return result, err
}

// TODO: Document.
func (client *G2diagnosticClient) GetEntityDetails(ctx context.Context, entityID int64, includeInternalFeatures int) (string, error) {
	request := pb.GetEntityDetailsRequest{
		EntityID:                entityID,
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.G2DiagnosticGrpcClient.GetEntityDetails(ctx, &request)
	result := response.GetResult()
	return result, err
}

// TODO: Document.
func (client *G2diagnosticClient) GetEntityListBySize(ctx context.Context, entitySize int) (interface{}, error) {
	request := pb.GetEntityListBySizeRequest{
		EntitySize: int32(entitySize),
	}
	response, err := client.G2DiagnosticGrpcClient.GetEntityListBySize(ctx, &request)
	result := response.GetResult()
	return result, err
}

// TODO: Document.
func (client *G2diagnosticClient) GetEntityResume(ctx context.Context, entityID int64) (string, error) {
	request := pb.GetEntityResumeRequest{
		EntityID: entityID,
	}
	response, err := client.G2DiagnosticGrpcClient.GetEntityResume(ctx, &request)
	result := response.GetResult()
	return result, err
}

// TODO: Document.
func (client *G2diagnosticClient) GetEntitySizeBreakdown(ctx context.Context, minimumEntitySize int, includeInternalFeatures int) (string, error) {
	request := pb.GetEntitySizeBreakdownRequest{
		MinimumEntitySize:       int32(minimumEntitySize),
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.G2DiagnosticGrpcClient.GetEntitySizeBreakdown(ctx, &request)
	result := response.GetResult()
	return result, err
}

// TODO: Document.
func (client *G2diagnosticClient) GetFeature(ctx context.Context, libFeatID int64) (string, error) {
	request := pb.GetFeatureRequest{
		LibFeatID: libFeatID,
	}
	response, err := client.G2DiagnosticGrpcClient.GetFeature(ctx, &request)
	result := response.GetResult()
	return result, err
}

// TODO: Document.
func (client *G2diagnosticClient) GetGenericFeatures(ctx context.Context, featureType string, maximumEstimatedCount int) (string, error) {
	request := pb.GetGenericFeaturesRequest{
		FeatureType:           featureType,
		MaximumEstimatedCount: int32(maximumEstimatedCount),
	}
	response, err := client.G2DiagnosticGrpcClient.GetGenericFeatures(ctx, &request)
	result := response.GetResult()
	return result, err
}

// GetLogicalCores returns the number of logical cores on the host system.
func (client *G2diagnosticClient) GetLogicalCores(ctx context.Context) (int, error) {
	request := pb.GetLogicalCoresRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetLogicalCores(ctx, &request)
	result := int(response.GetResult())
	return result, err
}

// TODO: Document.
func (client *G2diagnosticClient) GetMappingStatistics(ctx context.Context, includeInternalFeatures int) (string, error) {
	request := pb.GetMappingStatisticsRequest{
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.G2DiagnosticGrpcClient.GetMappingStatistics(ctx, &request)
	result := response.GetResult()
	return result, err
}

// GetPhysicalCores returns the number of physical cores on the host system.
func (client *G2diagnosticClient) GetPhysicalCores(ctx context.Context) (int, error) {
	request := pb.GetPhysicalCoresRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetPhysicalCores(ctx, &request)
	result := int(response.GetResult())
	return result, err
}

// TODO: Document.
func (client *G2diagnosticClient) GetRelationshipDetails(ctx context.Context, relationshipID int64, includeInternalFeatures int) (string, error) {
	request := pb.GetRelationshipDetailsRequest{
		RelationshipID:          relationshipID,
		IncludeInternalFeatures: int32(includeInternalFeatures),
	}
	response, err := client.G2DiagnosticGrpcClient.GetRelationshipDetails(ctx, &request)
	result := response.GetResult()
	return result, err
}

// TODO: Document.
func (client *G2diagnosticClient) GetResolutionStatistics(ctx context.Context) (string, error) {
	request := pb.GetResolutionStatisticsRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetResolutionStatistics(ctx, &request)
	result := response.GetResult()
	return result, err
}

// GetTotalSystemMemory returns the total memory, in bytes, on the host system.
func (client *G2diagnosticClient) GetTotalSystemMemory(ctx context.Context) (int64, error) {
	request := pb.GetTotalSystemMemoryRequest{}
	response, err := client.G2DiagnosticGrpcClient.GetTotalSystemMemory(ctx, &request)
	result := response.GetResult()
	return result, err
}

// Init initializes the Senzing G2diagnosis.
func (client *G2diagnosticClient) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	request := pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	_, err := client.G2DiagnosticGrpcClient.Init(ctx, &request)
	return err
}

// TODO: Document.
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

// TODO: Document.
func (client *G2diagnosticClient) Reinit(ctx context.Context, initConfigID int64) error {
	request := pb.ReinitRequest{
		InitConfigID: initConfigID,
	}
	_, err := client.G2DiagnosticGrpcClient.Reinit(ctx, &request)
	return err
}
