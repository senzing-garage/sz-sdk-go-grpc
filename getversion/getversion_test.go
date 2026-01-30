package getversion_test

import (
	"context"
	"testing"

	"github.com/senzing-garage/sz-sdk-go-grpc/getversion"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var (
	grpcAddress    = "0.0.0.0:8261"
	grpcConnection *grpc.ClientConn
)

// ----------------------------------------------------------------------------
// Interface methods - test
// ----------------------------------------------------------------------------

func TestVersion_GetSenzingVersion(test *testing.T) {
	ctx := context.Background()
	grpcConnection := getGrpcConnection(ctx)
	x := getversion.GetSenzingVersion(ctx, grpcConnection)
	assert.NotEmpty(test, x)
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getGrpcConnection(ctx context.Context) *grpc.ClientConn {
	if grpcConnection == nil {
		transportCredentials, err := helper.GetGrpcTransportCredentials(ctx)
		panicOnError(err)

		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(transportCredentials),
		}

		grpcConnection, err = grpc.NewClient(grpcAddress, dialOptions...)
		panicOnError(err)
	}

	return grpcConnection
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
