/*
 *
 */

// Package g2configclient implements a client for the service.
package g2configclient

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	pb "github.com/senzing/g2-sdk-proto/go/g2config"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (client *G2configClient) getLogger() messagelogger.MessageLoggerInterface {
	if client.logger == nil {
		client.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, messagelogger.LevelInfo)
	}
	return client.logger
}

func (client *G2configClient) notify(ctx context.Context, messageId int, err error, details map[string]string) {
	now := time.Now()
	details["subjectId"] = strconv.Itoa(ProductId)
	details["messageId"] = strconv.Itoa(messageId)
	details["messageTime"] = strconv.FormatInt(now.UnixNano(), 10)
	if err != nil {
		details["error"] = err.Error()
	}
	message, err := json.Marshal(details)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		client.observers.NotifyObservers(ctx, string(message))
	}
}

// Trace method entry.
func (client *G2configClient) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *G2configClient) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The AddDataSource method adds a data source to an existing in-memory configuration.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
  - inputJson: A JSON document in the format `{"DSRC_CODE": "NAME_OF_DATASOURCE"}`.

Output
  - A string containing a JSON document listing the newly created data source.
    See the example output.
*/
func (client *G2configClient) AddDataSource(ctx context.Context, configHandle uintptr, inputJson string) (string, error) {
	if client.isTrace {
		client.traceEntry(1, configHandle, inputJson)
	}
	entryTime := time.Now()
	request := pb.AddDataSourceRequest{
		ConfigHandle: int64(configHandle),
		InputJson:    inputJson,
	}
	response, err := client.GrpcClient.AddDataSource(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"inputJson": inputJson,
				"return":    response.GetResult(),
			}
			client.notify(ctx, 8001, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(2, configHandle, inputJson, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The Close method cleans up the Senzing G2Config object pointed to by the handle.
The handle was created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
*/
func (client *G2configClient) Close(ctx context.Context, configHandle uintptr) error {
	if client.isTrace {
		client.traceEntry(5, configHandle)
	}
	entryTime := time.Now()
	request := pb.CloseRequest{
		ConfigHandle: int64(configHandle),
	}
	_, err := client.GrpcClient.Close(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8002, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(6, configHandle, err, time.Since(entryTime))
	}
	return err
}

/*
The Create method creates an in-memory Senzing configuration from the g2config.json
template configuration file located in the PIPELINE.RESOURCEPATH path.
A handle is returned to identify the in-memory configuration.
The handle is used by the AddDataSource(), ListDataSources(), DeleteDataSource(), Load(), and Save() methods.
The handle is terminated by the Close() method.

Input
  - ctx: A context to control lifecycle.

Output
  - A Pointer to an in-memory Senzing configuration.
*/
func (client *G2configClient) Create(ctx context.Context) (uintptr, error) {
	if client.isTrace {
		client.traceEntry(7)
	}
	entryTime := time.Now()
	request := pb.CreateRequest{}
	response, err := client.GrpcClient.Create(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8003, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(8, (uintptr)(response.GetResult()), err, time.Since(entryTime))
	}
	return uintptr(response.GetResult()), err
}

/*
The DeleteDataSource method removes a data source from an existing configuration.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
  - inputJson: A JSON document in the format `{"DSRC_CODE": "NAME_OF_DATASOURCE"}`.
*/
func (client *G2configClient) DeleteDataSource(ctx context.Context, configHandle uintptr, inputJson string) error {
	if client.isTrace {
		client.traceEntry(9, configHandle, inputJson)
	}
	entryTime := time.Now()
	request := pb.DeleteDataSourceRequest{
		ConfigHandle: int64(configHandle),
		InputJson:    inputJson,
	}
	_, err := client.GrpcClient.DeleteDataSource(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"inputJson": inputJson,
			}
			client.notify(ctx, 8004, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(10, configHandle, inputJson, err, time.Since(entryTime))
	}
	return err
}

/*
The Destroy method will destroy and perform cleanup for the Senzing G2Config object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2configClient) Destroy(ctx context.Context) error {
	if client.isTrace {
		client.traceEntry(11)
	}
	entryTime := time.Now()
	request := pb.DestroyRequest{}
	_, err := client.GrpcClient.Destroy(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8005, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(12, err, time.Since(entryTime))
	}
	return err
}

/*
The Init method initializes the Senzing G2Config object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - moduleName: A name for the auditing node, to help identify it within system logs.
  - iniParams: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *G2configClient) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	if client.isTrace {
		client.traceEntry(17, moduleName, iniParams, verboseLogging)
	}
	entryTime := time.Now()
	request := pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	_, err := client.GrpcClient.Init(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"iniParams":      iniParams,
				"moduleName":     moduleName,
				"verboseLogging": strconv.Itoa(verboseLogging),
			}
			client.notify(ctx, 8006, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(18, moduleName, iniParams, verboseLogging, err, time.Since(entryTime))
	}
	return err
}

/*
The ListDataSources method returns a JSON document of data sources.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.

Output
  - A string containing a JSON document listing all of the data sources.
    See the example output.
*/
func (client *G2configClient) ListDataSources(ctx context.Context, configHandle uintptr) (string, error) {
	if client.isTrace {
		client.traceEntry(19, configHandle)
	}
	entryTime := time.Now()
	request := pb.ListDataSourcesRequest{
		ConfigHandle: int64(configHandle),
	}
	response, err := client.GrpcClient.ListDataSources(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8007, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(20, configHandle, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The Load method initializes the Senzing G2Config object from a JSON string.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
  - jsonConfig: A JSON document containing the Senzing configuration.
*/
func (client *G2configClient) Load(ctx context.Context, configHandle uintptr, jsonConfig string) error {
	if client.isTrace {
		client.traceEntry(21, configHandle, jsonConfig)
	}
	entryTime := time.Now()
	request := pb.LoadRequest{
		ConfigHandle: int64(configHandle),
		JsonConfig:   jsonConfig,
	}
	_, err := client.GrpcClient.Load(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8008, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(22, configHandle, jsonConfig, err, time.Since(entryTime))
	}
	return err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2configClient) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	if client.observers == nil {
		client.observers = &subject.SubjectImpl{}
	}
	return client.observers.RegisterObserver(ctx, observer)
}

/*
The Save method creates a JSON string representation of the Senzing G2Config object.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.

Output
  - A string containing a JSON Document representation of the Senzing G2Config object.
    See the example output.
*/
func (client *G2configClient) Save(ctx context.Context, configHandle uintptr) (string, error) {
	if client.isTrace {
		client.traceEntry(23, configHandle)
	}
	entryTime := time.Now()
	request := pb.SaveRequest{
		ConfigHandle: int64(configHandle),
	}
	response, err := client.GrpcClient.Save(ctx, &request)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8009, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(24, configHandle, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *G2configClient) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if client.isTrace {
		client.traceEntry(25, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	if client.isTrace {
		defer client.traceExit(26, logLevel, err, time.Since(entryTime))
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2configClient) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	err := client.observers.UnregisterObserver(ctx, observer)
	if err != nil {
		return err
	}
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	return err
}
