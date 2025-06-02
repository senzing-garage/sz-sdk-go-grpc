package szabstractfactory

import (
	"context"

	"github.com/senzing-garage/go-helpers/wraperror"
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

/*
Szabstractfactory is an implementation of the [senzing.SzAbstractFactory] interface.

[senzing.SzAbstractFactory]: https://pkg.go.dev/github.com/senzing-garage/sz-sdk-go/senzing#SzAbstractFactory
*/
type Szabstractfactory struct {
	GrpcConnection *grpc.ClientConn
}

// ----------------------------------------------------------------------------
// senzing.SzAbstractFactory interface methods
// ----------------------------------------------------------------------------

/*
Method CreateConfigManager returns an SzConfigManager object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzConfigManager object.
*/
func (factory *Szabstractfactory) CreateConfigManager(ctx context.Context) (senzing.SzConfigManager, error) {
	var err error

	_ = ctx
	result := &szconfigmanager.Szconfigmanager{
		GrpcClient:         szconfigmanagerpb.NewSzConfigManagerClient(factory.GrpcConnection),
		GrpcClientSzConfig: szconfigpb.NewSzConfigClient(factory.GrpcConnection),
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method CreateDiagnostic returns an SzDiagnostic object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzDiagnostic object.
*/
func (factory *Szabstractfactory) CreateDiagnostic(ctx context.Context) (senzing.SzDiagnostic, error) {
	var err error

	_ = ctx
	result := &szdiagnostic.Szdiagnostic{
		GrpcClient: szdiagnosticpb.NewSzDiagnosticClient(factory.GrpcConnection),
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method CreateEngine returns an SzEngine object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzEngine object.
*/
func (factory *Szabstractfactory) CreateEngine(ctx context.Context) (senzing.SzEngine, error) {
	var err error

	_ = ctx
	result := &szengine.Szengine{
		GrpcClient: szenginepb.NewSzEngineClient(factory.GrpcConnection),
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method CreateProduct returns an SzProduct object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzProduct object.
*/
func (factory *Szabstractfactory) CreateProduct(ctx context.Context) (senzing.SzProduct, error) {
	var err error

	_ = ctx
	result := &szproduct.Szproduct{
		GrpcClient: szproductpb.NewSzProductClient(factory.GrpcConnection),
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method Destroy will destroy and perform cleanup for the Senzing objects created by the AbstractFactory.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (factory *Szabstractfactory) Destroy(ctx context.Context) error {
	var err error

	szConfigmanager := &szconfigmanager.Szconfigmanager{
		GrpcClient: szconfigmanagerpb.NewSzConfigManagerClient(factory.GrpcConnection),
	}

	err = szConfigmanager.Destroy(ctx)
	if err != nil {
		return wraperror.Errorf(err, "szConfigmanager")
	}

	szDiagnostic := &szdiagnostic.Szdiagnostic{
		GrpcClient: szdiagnosticpb.NewSzDiagnosticClient(factory.GrpcConnection),
	}

	err = szDiagnostic.Destroy(ctx)
	if err != nil {
		return wraperror.Errorf(err, "szDiagnostic")
	}

	szEngine := &szengine.Szengine{
		GrpcClient: szenginepb.NewSzEngineClient(factory.GrpcConnection),
	}

	err = szEngine.Destroy(ctx)
	if err != nil {
		return wraperror.Errorf(err, "szEngine")
	}

	szProduct := &szproduct.Szproduct{
		GrpcClient: szproductpb.NewSzProductClient(factory.GrpcConnection),
	}

	err = szProduct.Destroy(ctx)
	if err != nil {
		return wraperror.Errorf(err, "szProduct")
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method Reinitialize re-initializes the Senzing objects created by the AbstractFactory
with a specific Senzing configuration JSON document identifier.

Input
  - ctx: A context to control lifecycle.
  - configID: The Senzing configuration JSON document identifier used for the initialization.
*/
func (factory *Szabstractfactory) Reinitialize(ctx context.Context, configID int64) error {
	var err error

	szDiagnostic := &szdiagnostic.Szdiagnostic{
		GrpcClient: szdiagnosticpb.NewSzDiagnosticClient(factory.GrpcConnection),
	}

	err = szDiagnostic.Reinitialize(ctx, configID)
	if err != nil {
		return wraperror.Errorf(err, "szDiagnostic")
	}

	szEngine := &szengine.Szengine{
		GrpcClient: szenginepb.NewSzEngineClient(factory.GrpcConnection),
	}

	err = szEngine.Reinitialize(ctx, configID)
	if err != nil {
		return wraperror.Errorf(err, "Reinitialize")
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}
