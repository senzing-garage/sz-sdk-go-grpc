package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/senzing/g2-sdk-go-grpc/g2configclient"
	"github.com/senzing/g2-sdk-go-grpc/g2configmgrclient"
	"github.com/senzing/g2-sdk-go-grpc/g2diagnosticclient"
	"github.com/senzing/g2-sdk-go-grpc/g2engineclient"
	"github.com/senzing/g2-sdk-go-grpc/g2productclient"
	"github.com/senzing/g2-sdk-go/g2config"
	"github.com/senzing/g2-sdk-go/g2configmgr"
	"github.com/senzing/g2-sdk-go/g2diagnostic"
	"github.com/senzing/g2-sdk-go/g2engine"
	"github.com/senzing/g2-sdk-go/g2product"
	pbg2config "github.com/senzing/g2-sdk-proto/go/g2config"
	pbg2configmgr "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	pbg2diagnostic "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
	pbg2engine "github.com/senzing/g2-sdk-proto/go/g2engine"
	pbg2product "github.com/senzing/g2-sdk-proto/go/g2product"
	"github.com/senzing/go-common/truthset"
	"github.com/senzing/go-logging/messageformat"
	"github.com/senzing/go-logging/messageid"
	"github.com/senzing/go-logging/messagelevel"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-logging/messagestatus"
	"github.com/senzing/go-logging/messagetext"
	"github.com/senzing/go-observing/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const MessageIdTemplate = "senzing-9999%04d"

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

var programName string = "unknown"
var buildVersion string = "0.0.0"
var buildIteration string = "0"
var logger messagelogger.MessageLoggerInterface = nil

var (
	grpcAddress    = "localhost:8258"
	grpcConnection *grpc.ClientConn
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func getGrpcConnection() *grpc.ClientConn {
	var err error = nil
	if grpcConnection == nil {
		grpcConnection, err = grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Did not connect: %v\n", err)
		}
		//		defer grpcConnection.Close()
	}
	return grpcConnection
}

func getG2config(ctx context.Context) (g2config.G2config, error) {
	var err error = nil
	grpcConnection := getGrpcConnection()
	result := &g2configclient.G2configClient{
		GrpcClient: pbg2config.NewG2ConfigClient(grpcConnection),
	}
	return result, err
}

func getG2configmgr(ctx context.Context) (g2configmgr.G2configmgr, error) {
	var err error = nil
	grpcConnection := getGrpcConnection()
	result := &g2configmgrclient.G2configmgrClient{
		GrpcClient: pbg2configmgr.NewG2ConfigMgrClient(grpcConnection),
	}
	return result, err
}

func getG2diagnostic(ctx context.Context) (g2diagnostic.G2diagnostic, error) {
	var err error = nil
	grpcConnection := getGrpcConnection()
	result := &g2diagnosticclient.G2diagnosticClient{
		GrpcClient: pbg2diagnostic.NewG2DiagnosticClient(grpcConnection),
	}
	return result, err
}

func getG2engine(ctx context.Context) (g2engine.G2engine, error) {
	var err error = nil
	grpcConnection := getGrpcConnection()
	result := &g2engineclient.G2engineClient{
		GrpcClient: pbg2engine.NewG2EngineClient(grpcConnection),
	}
	return result, err
}

func getG2product(ctx context.Context) (g2product.G2product, error) {
	var err error = nil
	grpcConnection := getGrpcConnection()
	result := &g2productclient.G2productClient{
		GrpcClient: pbg2product.NewG2ProductClient(grpcConnection),
	}
	return result, err
}

func getLogger(ctx context.Context) (messagelogger.MessageLoggerInterface, error) {
	messageFormat := &messageformat.MessageFormatJson{}
	messageIdTemplate := &messageid.MessageIdTemplated{
		MessageIdTemplate: MessageIdTemplate,
	}
	messageLevel := &messagelevel.MessageLevelByIdRange{
		IdLevelRanges: messagelevel.IdLevelRanges,
	}
	messageStatus := &messagestatus.MessageStatusByIdRange{
		IdStatusRanges: messagestatus.IdLevelRangesAsString,
	}
	messageText := &messagetext.MessageTextTemplated{
		IdMessages: Messages,
	}
	return messagelogger.New(messageFormat, messageIdTemplate, messageLevel, messageStatus, messageText, messagelogger.LevelInfo)
}

func demonstrateConfigFunctions(ctx context.Context, g2Config g2config.G2config, g2Configmgr g2configmgr.G2configmgr) error {
	now := time.Now()

	// Using G2Config: Create a default configuration in memory

	configHandle, err := g2Config.Create(ctx)
	if err != nil {
		return logger.Error(5100, err)
	}

	// Using G2Config: Add data source to in-memory configuration.

	for _, testDataSource := range truthset.TruthsetDataSources {
		_, err := g2Config.AddDataSource(ctx, configHandle, testDataSource.Json)
		if err != nil {
			return logger.Error(5101, err)
		}
	}

	// Using G2Config: Persist configuration to a string.

	configStr, err := g2Config.Save(ctx, configHandle)
	if err != nil {
		return logger.Error(5102, err)
	}

	// Using G2Configmgr: Persist configuration string to database.

	configComments := fmt.Sprintf("Created by g2diagnostic_test at %s", now.UTC())
	configID, err := g2Configmgr.AddConfig(ctx, configStr, configComments)
	if err != nil {
		return logger.Error(5103, err)
	}

	// Using G2Configmgr: Set new configuration as the default.

	err = g2Configmgr.SetDefaultConfigID(ctx, configID)
	if err != nil {
		return logger.Error(5104, err)
	}

	return err
}

func demonstrateAddRecord(ctx context.Context, g2Engine g2engine.G2engine) (string, error) {
	dataSourceCode := "TEST"
	recordID := strconv.Itoa(rand.Intn(1000000000))
	jsonData := fmt.Sprintf(
		"%s%s%s",
		`{"SOCIAL_HANDLE": "flavorh", "DATE_OF_BIRTH": "4/8/1983", "ADDR_STATE": "LA", "ADDR_POSTAL_CODE": "71232", "SSN_NUMBER": "053-39-3251", "ENTITY_TYPE": "TEST", "GENDER": "F", "srccode": "MDMPER", "CC_ACCOUNT_NUMBER": "5534202208773608", "RECORD_ID": "`,
		recordID,
		`", "DSRC_ACTION": "A", "ADDR_CITY": "Delhi", "DRIVERS_LICENSE_STATE": "DE", "PHONE_NUMBER": "225-671-0796", "NAME_LAST": "SEAMAN", "entityid": "284430058", "ADDR_LINE1": "772 Armstrong RD"}`)
	loadID := dataSourceCode
	var flags int64 = 0

	// Using G2Engine: Add record and return "withInfo".

	return g2Engine.AddRecordWithInfo(ctx, dataSourceCode, recordID, jsonData, loadID, flags)
}

func demonstrateAdditionalFunctions(ctx context.Context, g2Diagnostic g2diagnostic.G2diagnostic, g2Engine g2engine.G2engine, g2Product g2product.G2product) error {
	// Using G2Diagnostic: Check physical cores.

	actual, err := g2Diagnostic.GetPhysicalCores(ctx)
	if err != nil {
		logger.Log(5300, err)
	}
	logger.Log(2002, actual)

	// Using G2Engine: Add records with information returned.

	withInfo, err := demonstrateAddRecord(ctx, g2Engine)
	if err != nil {
		logger.Log(5302, err)
	}
	logger.Log(2003, withInfo)

	// Using G2Product: Show license metadata.

	license, err := g2Product.License(ctx)
	if err != nil {
		logger.Log(5303, err)
	}
	logger.Log(2004, license)

	return err
}

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func main() {
	ctx := context.TODO()
	var err error = nil

	// Randomize random number generator.

	rand.Seed(time.Now().UnixNano())

	// Configure the "log" standard library.

	log.SetFlags(0)
	logger, err = getLogger(ctx)
	if err != nil {
		logger.Log(5000, err)
	}

	// Test logger.

	programmMetadataMap := map[string]interface{}{
		"ProgramName":    programName,
		"BuildVersion":   buildVersion,
		"BuildIteration": buildIteration,
	}

	fmt.Printf("\n-------------------------------------------------------------------------------\n\n")
	logger.Log(2001, "Just a test of logging", programmMetadataMap)

	// Create 2 observers.

	observer1 := &observer.ObserverNull{
		Id: "Observer 1",
	}
	observer2 := &observer.ObserverNull{
		Id: "Observer 2",
	}

	// Get Senzing objects for installing a Senzing Engine configuration.

	g2Config, err := getG2config(ctx)
	if err != nil {
		logger.Log(5001, err)
	}
	g2Config.RegisterObserver(ctx, observer1)
	g2Config.RegisterObserver(ctx, observer2)

	g2Configmgr, err := getG2configmgr(ctx)
	if err != nil {
		logger.Log(5002, err)
	}
	g2Configmgr.RegisterObserver(ctx, observer1)

	// Persist the Senzing configuration to the Senzing repository.

	err = demonstrateConfigFunctions(ctx, g2Config, g2Configmgr)
	if err != nil {
		logger.Log(5003, err)
	}

	// Now that a Senzing configuration is installed, get the remainder of the Senzing objects.

	g2Diagnostic, err := getG2diagnostic(ctx)
	if err != nil {
		logger.Log(5004, err)
	}
	g2Diagnostic.RegisterObserver(ctx, observer1)

	g2Engine, err := getG2engine(ctx)
	if err != nil {
		logger.Log(5005, err)
	}
	g2Engine.RegisterObserver(ctx, observer1)

	g2Product, err := getG2product(ctx)
	if err != nil {
		logger.Log(5006, err)
	}
	g2Product.RegisterObserver(ctx, observer1)

	// Demonstrate tests.

	err = demonstrateAdditionalFunctions(ctx, g2Diagnostic, g2Engine, g2Product)
	if err != nil {
		logger.Log(5007, err)
	}

	fmt.Printf("\n-------------------------------------------------------------------------------\n\n")
}
