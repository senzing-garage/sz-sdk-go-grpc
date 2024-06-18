package szabstractfactory

import (
	"context"

	"github.com/senzing-garage/sz-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-grpc/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go-grpc/szengine"
	"github.com/senzing-garage/sz-sdk-go-grpc/szproduct"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	szdiagnosticpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
	szenginepb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
	szproductpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
	"google.golang.org/grpc"
)

// Szconfig is the default implementation of the Szconfig interface.
type Szabstractfactory struct {
	GrpcConnection *grpc.ClientConn
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
TODO: Write description for CreateSzConfig
The CreateSzConfig method...

Input
  - ctx: A context to control lifecycle.

Output
  - An senzing.SzConfig object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzConfig(ctx context.Context) (senzing.SzConfig, error) {
	result := &szconfig.Szconfig{
		GrpcClient: szconfigpb.NewSzConfigClient(factory.GrpcConnection),
	}
	return result, nil
}

/*
TODO: Write description for CreateSzConfigManager
The CreateSzConfigManager method...

Input
  - ctx: A context to control lifecycle.

Output
  - An senzing.CreateConfigManager object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzConfigManager(ctx context.Context) (senzing.SzConfigManager, error) {
	result := &szconfigmanager.Szconfigmanager{
		GrpcClient: szconfigmanagerpb.NewSzConfigManagerClient(factory.GrpcConnection),
	}
	return result, nil
}

/*
TODO: Write description for CreateSzDiagnostic
The CreateSzDiagnostic method...

Input
  - ctx: A context to control lifecycle.

Output
  - An senzing.SzDiagnostic object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzDiagnostic(ctx context.Context) (senzing.SzDiagnostic, error) {
	result := &szdiagnostic.Szdiagnostic{
		GrpcClient: szdiagnosticpb.NewSzDiagnosticClient(factory.GrpcConnection),
	}
	return result, nil
}

/*
TODO: Write description for CreateSzEngine
The CreateSzEngine method...

Input
  - ctx: A context to control lifecycle.

Output
  - An senzing.SzEngine object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzEngine(ctx context.Context) (senzing.SzEngine, error) {
	result := &szengine.Szengine{
		GrpcClient: szenginepb.NewSzEngineClient(factory.GrpcConnection),
	}
	return result, nil
}

/*
TODO: Write description for CreateSzProduct
The CreateSzProduct method...

Input
  - ctx: A context to control lifecycle.

Output
  - An senzing.SzProduct object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzProduct(ctx context.Context) (senzing.SzProduct, error) {
	result := &szproduct.Szproduct{
		GrpcClient: szproductpb.NewSzProductClient(factory.GrpcConnection),
	}
	return result, nil
}
