/*
 *
 */

// Package main implements a client for the service.
package g2diagnosticclient

import (
	pb "github.com/senzing/g2-sdk-go-grpc/protobuf/g2diagnostic"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const MessageIdFormat = "senzing-6025%04d"

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2diagnosticClient struct {
	G2DiagnosticGrpcClient pb.G2DiagnosticClient
}

// ----------------------------------------------------------------------------
// Interfaces
// ----------------------------------------------------------------------------

//type G2GiagnosticClient interface {
//	CheckDBPerf(ctx context.Context, secondsToRun int) (string, error)
//	ClearLastException(ctx context.Context) error
//	CloseEntityListBySize(ctx context.Context, entityListBySizeHandle interface{}) error
//	Destroy(ctx context.Context) error
//	FetchNextEntityBySize(ctx context.Context, entityListBySizeHandle interface{}) (string, error)
//	FindEntitiesByFeatureIDs(ctx context.Context, features string) (string, error)
//	GetAvailableMemory(ctx context.Context) (int64, error)
//	GetDataSourceCounts(ctx context.Context) (string, error)
//	GetDBInfo(ctx context.Context) (string, error)
//	GetEntityDetails(ctx context.Context, entityID int64, includeInternalFeatures int) (string, error)
//	GetEntityListBySize(ctx context.Context, entitySize int) (interface{}, error)
//	GetEntityResume(ctx context.Context, entityID int64) (string, error)
//	GetEntitySizeBreakdown(ctx context.Context, minimumEntitySize int, includeInternalFeatures int) (string, error)
//	GetFeature(ctx context.Context, libFeatID int64) (string, error)
//	GetGenericFeatures(ctx context.Context, featureType string, maximumEstimatedCount int) (string, error)
//	GetLastException(ctx context.Context) (string, error)
//	GetLastExceptionCode(ctx context.Context) (string, error)
//	GetLogicalCores(ctx context.Context) (int, error)
//	GetMappingStatistics(ctx context.Context, includeInternalFeatures int) (string, error)
//	GetPhysicalCores(ctx context.Context) (int, error)
//	GetRelationshipDetails(ctx context.Context, relationshipID int64, includeInternalFeatures int) (string, error)
//	GetResolutionStatistics(ctx context.Context) (string, error)
//	GetTotalSystemMemory(ctx context.Context) (int64, error)
//	Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error
//	InitWithConfigID(ctx context.Context, moduleName string, iniParams string, initConfigID int64, verboseLogging int) error
//	Reinit(ctx context.Context, initConfigID int64) error
//}
