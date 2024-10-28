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
The CreateSzConfig method returns an SzConfig object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzConfig object.
*/
func (factory *Szabstractfactory) CreateSzConfig(ctx context.Context) (senzing.SzConfig, error) {
	_ = ctx
	result := &szconfig.Szconfig{
		GrpcClient: szconfigpb.NewSzConfigClient(factory.GrpcConnection),
	}
	return result, nil
}

/*
The CreateSzConfigManager method returns an SzConfigManager object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzConfigManager object.
*/
func (factory *Szabstractfactory) CreateSzConfigManager(ctx context.Context) (senzing.SzConfigManager, error) {
	_ = ctx
	result := &szconfigmanager.Szconfigmanager{
		GrpcClient: szconfigmanagerpb.NewSzConfigManagerClient(factory.GrpcConnection),
	}
	return result, nil
}

/*
The CreateSzDiagnostic method returns an SzDiagnostic object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzDiagnostic object.
*/
func (factory *Szabstractfactory) CreateSzDiagnostic(ctx context.Context) (senzing.SzDiagnostic, error) {
	_ = ctx
	result := &szdiagnostic.Szdiagnostic{
		GrpcClient: szdiagnosticpb.NewSzDiagnosticClient(factory.GrpcConnection),
	}
	return result, nil
}

/*
The CreateSzEngine method returns an SzEngine object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzEngine object.
*/
func (factory *Szabstractfactory) CreateSzEngine(ctx context.Context) (senzing.SzEngine, error) {
	_ = ctx
	result := &szengine.Szengine{
		GrpcClient: szenginepb.NewSzEngineClient(factory.GrpcConnection),
	}
	return result, nil
}

/*
The CreateSzProduct method returns an SzProduct object
implemented to use the Senzing native C binary, libSz.so.

Input
  - ctx: A context to control lifecycle.

Output
  - An SzProduct object.
*/
func (factory *Szabstractfactory) CreateSzProduct(ctx context.Context) (senzing.SzProduct, error) {
	_ = ctx
	result := &szproduct.Szproduct{
		GrpcClient: szproductpb.NewSzProductClient(factory.GrpcConnection),
	}
	return result, nil
}

/*
Method Destroy will destroy and perform cleanup for the Senzing objects created by the AbstractFactory.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (factory *Szabstractfactory) Destroy(ctx context.Context) error {
	var err error
	return err
}

/*
Method Reinitialize re-initializes the Senzing objects created by the AbstractFactory with a specific Senzing configuration JSON document identifier.

Input
  - ctx: A context to control lifecycle.
  - configID: The Senzing configuration JSON document identifier used for the initialization.
*/
func (factory *Szabstractfactory) Reinitialize(ctx context.Context, configID int64) error {
	var err error
	factory.ConfigID = configID
	if factory.isSzdiagnosticInitialized {
		szDiagnostic := &szdiagnostic.Szdiagnostic{}
		err = szDiagnostic.Reinitialize(ctx, configID)
		if err != nil {
			return err
		}
	}
	if factory.isSzengineInitialized {
		szEngine := &szengine.Szengine{}
		err = szEngine.Reinitialize(ctx, configID)
		if err != nil {
			return err
		}
	}
	return err
}
