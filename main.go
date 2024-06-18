package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go-grpc/szengine"
	"github.com/senzing-garage/sz-sdk-go-grpc/szproduct"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szconfigmanagerpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
	szenginepb "github.com/senzing-garage/sz-sdk-proto/go/szengine"
	szproductpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const MessageIDTemplate = "senzing-9999%04d"

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var Messages = map[int]string{
	1:    "%s",
	2:    "WithInfo: %s",
	2001: "Testing %s.",
	2002: "Physical cores: %d.",
	2003: "withInfo",
	2004: "License",
	2999: "Cannot retrieve last error message.",
}

// Values updated via "go install -ldflags" parameters.

var (
	buildIteration = "0"
	buildVersion   = "0.0.0"
	grpcAddress    = "localhost:8261"
	grpcConnection *grpc.ClientConn
	logger         logging.Logging
	programName    = "unknown"
)

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func main() {
	var err error
	ctx := context.TODO()

	// Configure the "log" standard library.

	log.SetFlags(0)
	logger, err = getLogger(ctx)
	if err != nil {
		failOnError(5000, err)
	}

	// Test logger.

	programmMetadataMap := map[string]interface{}{
		"ProgramName":    programName,
		"BuildVersion":   buildVersion,
		"BuildIteration": buildIteration,
	}

	fmt.Printf("\n-------------------------------------------------------------------------------\n\n")
	logger.Log(2001, "Just a test of logging", programmMetadataMap)

	// Create observers.

	// observer1 := &observer.ObserverNull{
	// 	Id: "Observer 1",
	// }
	// observer2 := &observer.ObserverNull{
	// 	Id: "Observer 2",
	// }

	// grpcConnection, err := grpc.Dial("localhost:8261", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	fmt.Printf("Did not connect: %v\n", err)
	// }

	// observer3 := &observer.ObserverGrpc{
	// 	Id:         "Observer 3",
	// 	GrpcClient: observerpb.NewObserverClient(grpcConnection),
	// }

	// Get Senzing objects for installing a Senzing Engine configuration.

	szConfig, err := getSzConfig(ctx)
	if err != nil {
		failOnError(5001, err)
	}
	// err = szConfig.RegisterObserver(ctx, observer1)
	// if err != nil {
	// 	panic(err)
	// }
	// err = szConfig.RegisterObserver(ctx, observer2)
	// if err != nil {
	// 	panic(err)
	// }
	// err = szConfig.RegisterObserver(ctx, observer3)
	// if err != nil {
	// 	panic(err)
	// }
	// szConfig.SetObserverOrigin(ctx, "s-sdk-go-grpc main.go")

	szConfigManager, err := getSzConfigManager(ctx)
	if err != nil {
		failOnError(5005, err)
	}
	// err = szConfigManager.RegisterObserver(ctx, observer1)
	// if err != nil {
	// 	panic(err)
	// }

	// Persist the Senzing configuration to the Senzing repository.

	err = demonstrateConfigFunctions(ctx, szConfig, szConfigManager)
	if err != nil {
		failOnError(5008, err)
	}

	// Now that a Senzing configuration is installed, get the remainder of the Senzing objects.

	szEngine, err := getSzEngine(ctx)
	if err != nil {
		failOnError(5010, err)
	}
	// err = szEngine.RegisterObserver(ctx, observer1)
	// if err != nil {
	// 	panic(err)
	// }

	szProduct, err := getSzProduct(ctx)
	if err != nil {
		failOnError(5011, err)
	}
	// err = szProduct.RegisterObserver(ctx, observer1)
	// if err != nil {
	// 	panic(err)
	// }

	// Demonstrate tests.

	err = demonstrateAdditionalFunctions(ctx, szEngine, szProduct)
	if err != nil {
		failOnError(5015, err)
	}

	fmt.Printf("\n-------------------------------------------------------------------------------\n\n")
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func getGrpcConnection() *grpc.ClientConn {
	var err error
	if grpcConnection == nil {
		grpcConnection, err = grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Did not connect: %v\n", err)
		}
		//		defer grpcConnection.Close()
	}
	return grpcConnection
}

func getSzConfig(ctx context.Context) (senzing.SzConfig, error) {
	_ = ctx
	var err error
	grpcConnection := getGrpcConnection()
	result := &szconfig.Szconfig{
		GrpcClient: szconfigpb.NewSzConfigClient(grpcConnection),
	}
	return result, err
}

func getSzConfigManager(ctx context.Context) (senzing.SzConfigManager, error) {
	_ = ctx
	var err error
	grpcConnection := getGrpcConnection()
	result := &szconfigmanager.Szconfigmanager{
		GrpcClient: szconfigmanagerpb.NewSzConfigManagerClient(grpcConnection),
	}
	return result, err
}

func getSzEngine(ctx context.Context) (senzing.SzEngine, error) {
	_ = ctx
	var err error
	grpcConnection := getGrpcConnection()
	result := &szengine.Szengine{
		GrpcClient: szenginepb.NewSzEngineClient(grpcConnection),
	}
	return result, err
}

func getSzProduct(ctx context.Context) (senzing.SzProduct, error) {
	_ = ctx
	var err error
	grpcConnection := getGrpcConnection()
	result := &szproduct.Szproduct{
		GrpcClient: szproductpb.NewSzProductClient(grpcConnection),
	}
	return result, err
}

func getLogger(ctx context.Context) (logging.Logging, error) {
	_ = ctx
	logger, err := logging.NewSenzingLogger(9999, Messages)
	if err != nil {
		fmt.Println(err)
	}

	return logger, err
}

func demonstrateConfigFunctions(ctx context.Context, szConfig senzing.SzConfig, szConfigManager senzing.SzConfigManager) error {
	now := time.Now()

	// Using SzConfig: Create a default configuration in memory.

	configHandle, err := szConfig.CreateConfig(ctx)
	if err != nil {
		return logger.NewError(5100, err)
	}

	// Using SzConfig: Add data source to in-memory configuration.

	for dataSourceCode := range truthset.TruthsetDataSources {
		_, err := szConfig.AddDataSource(ctx, configHandle, dataSourceCode)
		if err != nil {
			return logger.NewError(5101, err)
		}
	}

	// Using SzConfig: Persist configuration to a string.

	configStr, err := szConfig.ExportConfig(ctx, configHandle)
	if err != nil {
		return logger.NewError(5102, err)
	}

	// Using SzConfigManager: Persist configuration string to database.

	configComments := fmt.Sprintf("Created by szmain.go at %s", now.UTC())
	configID, err := szConfigManager.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return logger.NewError(5103, err)
	}

	// Using SzConfigManager: Set new configuration as the default.

	err = szConfigManager.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return logger.NewError(5104, err)
	}

	return err
}

func demonstrateAddRecord(ctx context.Context, szEngine senzing.SzEngine) (string, error) {
	dataSourceCode := "TEST"
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(1000000000))
	if err != nil {
		panic(err)
	}
	recordID := randomNumber.String()
	recordDefinition := fmt.Sprintf(
		"%s%s%s",
		`{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "`,
		recordID,
		`", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "SEAMAN", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`)
	var flags int64 = senzing.SzWithInfo

	// Using SzEngine: Add record and return "withInfo".

	return szEngine.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, flags)
}

func demonstrateAdditionalFunctions(ctx context.Context, szEngine senzing.SzEngine, szProduct senzing.SzProduct) error {

	// Using SzEngine: Add records with information returned.

	withInfo, err := demonstrateAddRecord(ctx, szEngine)
	if err != nil {
		failOnError(5302, err)
	}
	logger.Log(2003, withInfo)

	// Using SzProduct: Show license metadata.

	license, err := szProduct.GetLicense(ctx)
	if err != nil {
		failOnError(5303, err)
	}
	logger.Log(2004, license)

	return err
}

func failOnError(msgID int, err error) {
	logger.Log(msgID, err)
	panic(err.Error())
}
