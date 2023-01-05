/*
 *
 */

// Package main implements a client for the service.
package main

import (
	"context"
	"flag"
	"time"

	"github.com/senzing/g2-sdk-go-grpc/g2diagnosticclient"
	pb "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
	"github.com/senzing/go-helpers/g2engineconfigurationjson"
	"github.com/senzing/go-logging/messagelogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2config package found messages having the format "senzing-6001xxxx".
const ProductId = 9999

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2config package.
var IdMessages = map[int]string{
	1:    "Enter AddDataSource(%v, %s).",
	2:    "Exit  AddDataSource(%v, %s) returned (%s, %v).",
	4001: "Call to G2Config_addDataSource(%v, %s) failed. Return code: %d",
}

// Status strings for specific g2config messages.
var IdStatuses = map[int]string{}

// ----------------------------------------------------------------------------
// Interfaces
// ----------------------------------------------------------------------------

var (
	grpcAddress = flag.String("addr", "localhost:8258", "the address to connect to")
)

func main() {

	// Configure the "log" standard library.

	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)

	// Create a context.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Quick-and-dirty command line parameters. (Replace with Viper)

	flag.Parse()

	// Set up a connection to the gRPC server.

	grpcConnection, err := grpc.Dial(*grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log(5000, "Did not connect: %v", err)
	}
	defer grpcConnection.Close()

	// Set up a client to the G2diagnosis gRPC server.

	g2diagnosticClient := g2diagnosticclient.G2diagnosticClient{
		GrpcClient: pb.NewG2DiagnosticClient(grpcConnection),
	}

	// Create request parameters.

	moduleName := "Test module name"
	verboseLogging := 0
	iniParams, jsonErr := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
	if jsonErr != nil {
		logger.Log(5001, "Could not build Configuration JSON: %v", jsonErr)
	}

	// Call.

	err = g2diagnosticClient.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		logger.Log(5002, "Could not Init: %v", err)
	}

	// Call.

	var result string
	result, err = g2diagnosticClient.CheckDBPerf(ctx, 10)
	if err != nil {
		logger.Log(5003, "Could not CheckDBPerf: %v", err)
	}
	logger.Log(2001, result)
}
