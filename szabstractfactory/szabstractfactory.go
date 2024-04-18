package szabstractfactory

import (
	"context"

	"github.com/senzing-garage/sz-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-grpc/szdiagnostic"
	"github.com/senzing-garage/sz-sdk-go-grpc/szengine"
	"github.com/senzing-garage/sz-sdk-go-grpc/szproduct"
	"github.com/senzing-garage/sz-sdk-go/sz"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	szdiagnosticpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
	szenginepb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
	szproductpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

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
  - An sz.SzConfig object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzConfig(ctx context.Context) (sz.SzConfig, error) {
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
  - An sz.CreateConfigManager object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzConfigManager(ctx context.Context) (sz.SzConfigManager, error) {
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
  - An sz.SzDiagnostic object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzDiagnostic(ctx context.Context) (sz.SzDiagnostic, error) {
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
  - An sz.SzEngine object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzEngine(ctx context.Context) (sz.SzEngine, error) {
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
  - An sz.SzProduct object.
    See the example output.
*/
func (factory *Szabstractfactory) CreateSzProduct(ctx context.Context) (sz.SzProduct, error) {
	result := &szproduct.Szproduct{
		GrpcClient: szproductpb.NewSzProductClient(factory.GrpcConnection),
	}
	return result, nil
}