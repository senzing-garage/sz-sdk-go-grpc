/*
 *
 */

// Package g2config implements a client for the service.
package g2config

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	g2configapi "github.com/senzing/g2-sdk-go/g2config"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2config"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2config struct {
	GrpcClient g2pb.G2ConfigClient
	isTrace    bool
	logger     messagelogger.MessageLoggerInterface
	observers  subject.Subject
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (client *G2config) getLogger() messagelogger.MessageLoggerInterface {
	if client.logger == nil {
		client.logger, _ = messagelogger.NewSenzingApiLogger(ProductId, g2configapi.IdMessages, g2configapi.IdStatuses, messagelogger.LevelInfo)
	}
	return client.logger
}

// Notify registered observers.
func (client *G2config) notify(ctx context.Context, messageId int, err error, details map[string]string) {
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
func (client *G2config) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *G2config) traceExit(errorNumber int, details ...interface{}) {
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
func (client *G2config) AddDataSource(ctx context.Context, configHandle uintptr, inputJson string) (string, error) {
	if client.isTrace {
		client.traceEntry(1, configHandle, inputJson)
	}
	entryTime := time.Now()
	request := g2pb.AddDataSourceRequest{
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
func (client *G2config) Close(ctx context.Context, configHandle uintptr) error {
	if client.isTrace {
		client.traceEntry(5, configHandle)
	}
	entryTime := time.Now()
	request := g2pb.CloseRequest{
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
func (client *G2config) Create(ctx context.Context) (uintptr, error) {
	if client.isTrace {
		client.traceEntry(7)
	}
	entryTime := time.Now()
	request := g2pb.CreateRequest{}
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
func (client *G2config) DeleteDataSource(ctx context.Context, configHandle uintptr, inputJson string) error {
	if client.isTrace {
		client.traceEntry(9, configHandle, inputJson)
	}
	entryTime := time.Now()
	request := g2pb.DeleteDataSourceRequest{
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
func (client *G2config) Destroy(ctx context.Context) error {
	if client.isTrace {
		client.traceEntry(11)
	}
	entryTime := time.Now()
	request := g2pb.DestroyRequest{}
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
The GetSdkId method returns the identifier of this particular Software Development Kit (SDK).
It is handy when working with multiple implementations of the same G2configInterface.
For this implementation, "grpc" is returned.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2config) GetSdkId(ctx context.Context) string {
	if client.isTrace {
		client.traceEntry(31)
	}
	entryTime := time.Now()
	var err error = nil
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			client.notify(ctx, 8010, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(32, err, time.Since(entryTime))
	}
	return "grpc"
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
func (client *G2config) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	if client.isTrace {
		client.traceEntry(17, moduleName, iniParams, verboseLogging)
	}
	entryTime := time.Now()
	request := g2pb.InitRequest{
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
func (client *G2config) ListDataSources(ctx context.Context, configHandle uintptr) (string, error) {
	if client.isTrace {
		client.traceEntry(19, configHandle)
	}
	entryTime := time.Now()
	request := g2pb.ListDataSourcesRequest{
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
func (client *G2config) Load(ctx context.Context, configHandle uintptr, jsonConfig string) error {
	if client.isTrace {
		client.traceEntry(21, configHandle, jsonConfig)
	}
	entryTime := time.Now()
	request := g2pb.LoadRequest{
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
func (client *G2config) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	if client.isTrace {
		client.traceEntry(27, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	if client.observers == nil {
		client.observers = &subject.SubjectImpl{}
	}
	err := client.observers.RegisterObserver(ctx, observer)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"observerID": observer.GetObserverId(ctx),
			}
			client.notify(ctx, 8011, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(28, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
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
func (client *G2config) Save(ctx context.Context, configHandle uintptr) (string, error) {
	if client.isTrace {
		client.traceEntry(23, configHandle)
	}
	entryTime := time.Now()
	request := g2pb.SaveRequest{
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
func (client *G2config) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	if client.isTrace {
		client.traceEntry(25, logLevel)
	}
	entryTime := time.Now()
	var err error = nil
	client.getLogger().SetLogLevel(messagelogger.Level(logLevel))
	client.isTrace = (client.getLogger().GetLogLevel() == messagelogger.LevelTrace)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"logLevel": logger.LevelToTextMap[logLevel],
			}
			client.notify(ctx, 8012, err, details)
		}()
	}
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
func (client *G2config) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	if client.isTrace {
		client.traceEntry(29, observer.GetObserverId(ctx))
	}
	entryTime := time.Now()
	var err error = nil
	if client.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverId(ctx),
		}
		client.notify(ctx, 8013, err, details)
	}
	err = client.observers.UnregisterObserver(ctx, observer)
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	if client.isTrace {
		defer client.traceExit(30, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}
