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
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"google.golang.org/grpc"
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
	grpcAddress    = "0.0.0.0:8261"
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
	failOnError(5000, err)

	// Test logger.

	programmMetadataMap := map[string]interface{}{
		"ProgramName":    programName,
		"BuildVersion":   buildVersion,
		"BuildIteration": buildIteration,
	}

	outputf("\n-------------------------------------------------------------------------------\n\n")
	logger.Log(2001, "Just a test of logging", programmMetadataMap)

	szAbstractFactory := getSzAbstractFactory(ctx)

	// Persist the Senzing configuration to the Senzing repository.

	err = demonstrateConfigFunctions(ctx, szAbstractFactory)
	failOnError(5002, err)

	// Demonstrate tests.

	err = demonstrateAdditionalFunctions(ctx, szAbstractFactory)
	failOnError(5003, err)

	outputf("\n-------------------------------------------------------------------------------\n\n")
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func demonstrateAdditionalFunctions(ctx context.Context, szAbstractFactory senzing.SzAbstractFactory) error {

	// Create Senzing objects.

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	failOnError(5301, err)

	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	failOnError(5302, err)

	// Using SzEngine: Add records with information returned.

	withInfo, err := demonstrateAddRecord(ctx, szEngine)
	failOnError(5303, err)
	logger.Log(2003, withInfo)

	// Using SzProduct: Show license metadata.

	license, err := szProduct.GetLicense(ctx)
	failOnError(5304, err)
	logger.Log(2004, license)

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
		`", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "SEAMAN", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`,
	)

	var flags = senzing.SzWithInfo

	// Using SzEngine: Add record and return "withInfo".

	return szEngine.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, flags)
}

func demonstrateConfigFunctions(ctx context.Context, szAbstractFactory senzing.SzAbstractFactory) error {
	now := time.Now()

	// Create Senzing objects.

	szConfigManager, err := szAbstractFactory.CreateConfigManager(ctx)
	if err != nil {
		return logger.NewError(5100, err)
	}

	szConfig, err := szConfigManager.CreateConfigFromTemplate(ctx)
	if err != nil {
		return logger.NewError(5101, err)
	}

	// Using SzConfig: Add data source to in-memory configuration.

	for dataSourceCode := range truthset.TruthsetDataSources {
		_, err := szConfig.AddDataSource(ctx, dataSourceCode)
		if err != nil {
			return logger.NewError(5102, err)
		}
	}

	// Using SzConfig: Persist configuration to a string.

	configComments := fmt.Sprintf("Created by szmain.go at %s", now.UTC())

	configDefinition, err := szConfig.Export(ctx)
	if err != nil {
		return logger.NewError(5103, err)
	}

	// Using SzConfigManager: Persist configuration string to database.

	_, err = szConfigManager.SetDefaultConfig(ctx, configDefinition, configComments)
	if err != nil {
		return logger.NewError(5104, err)
	}

	return err
}

func failOnError(msgID int, err error) {
	if err != nil {
		logger.Log(msgID, err)
		panic(err.Error())
	}
}

func getGrpcConnection() *grpc.ClientConn {
	if grpcConnection == nil {
		transportCredentials, err := helper.GetGrpcTransportCredentials()
		if err != nil {
			panic(err)
		}

		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(transportCredentials),
		}

		grpcConnection, err = grpc.NewClient(grpcAddress, dialOptions...)
		if err != nil {
			panic(err)
		}
	}

	return grpcConnection
}

func getLogger(ctx context.Context) (logging.Logging, error) {
	_ = ctx
	loggerOptions := []interface{}{
		logging.OptionMessageFields{Value: []string{"id", "text", "reason", "errors", "details"}},
	}

	logger, err := logging.NewSenzingLogger(9999, Messages, loggerOptions...)
	if err != nil {
		outputln(err)
	}

	return logger, err
}

func getSzAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	_ = ctx

	result := &szabstractfactory.Szabstractfactory{
		GrpcConnection: getGrpcConnection(),
	}

	return result
}

func outputf(format string, message ...any) {
	fmt.Printf(format, message...) //nolint
}

func outputln(message ...any) {
	fmt.Println(message...) //nolint
}
