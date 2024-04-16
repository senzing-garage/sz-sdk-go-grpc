/*
 *
 */

// Package szconfig implements a client for the service.
package szconfig

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	szconfigapi "github.com/senzing-garage/sz-sdk-go/szconfig"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Szconfig struct {
	GrpcClient     szpb.SzConfigClient
	isTrace        bool // Performance optimization
	logger         logging.LoggingInterface
	observerOrigin string
	observers      subject.Subject
}

// ----------------------------------------------------------------------------
// sz-sdk-go.SzConfig interface methods
// ----------------------------------------------------------------------------

/*
The AddDataSource method adds a data source to an existing in-memory configuration.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
  - dataSourceCode: A JSON document in the format `{"DSRC_CODE": "NAME_OF_DATASOURCE"}`.

Output
  - A string containing a JSON document listing the newly created data source.
    See the example output.
*/
func (client *Szconfig) AddDataSource(ctx context.Context, configHandle uintptr, dataSourceCode string) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(1, configHandle, dataSourceCode)
		defer func() { client.traceExit(2, configHandle, dataSourceCode, result, err, time.Since(entryTime)) }()
	}
	request := szpb.AddDataSourceRequest{
		ConfigHandle:   int64(configHandle),
		DataSourceCode: dataSourceCode,
	}
	response, err := client.GrpcClient.AddDataSource(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"return":         result,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8001, err, details)
		}()
	}
	return result, err
}

/*
The CloseConfig method cleans up the Senzing G2Config object pointed to by the handle.
The handle was created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
*/
func (client *Szconfig) CloseConfig(ctx context.Context, configHandle uintptr) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(5, configHandle)
		defer func() { client.traceExit(6, configHandle, err, time.Since(entryTime)) }()
	}
	request := szpb.CloseConfigRequest{
		ConfigHandle: int64(configHandle),
	}
	_, err = client.GrpcClient.CloseConfig(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8002, err, details)
		}()
	}
	return err
}

/*
The CreateConfig method creates an in-memory Senzing configuration from the g2config.json
template configuration file located in the PIPELINE.RESOURCEPATH path.
A handle is returned to identify the in-memory configuration.
The handle is used by the AddDataSource(), ListDataSources(), DeleteDataSource(), Load(), and Save() methods.
The handle is terminated by the Close() method.

Input
  - ctx: A context to control lifecycle.

Output
  - A Pointer to an in-memory Senzing configuration.
*/
func (client *Szconfig) CreateConfig(ctx context.Context) (uintptr, error) {
	var err error = nil
	var result uintptr = 0
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(7)
		defer func() { client.traceExit(8, result, err, time.Since(entryTime)) }()
	}
	request := szpb.CreateConfigRequest{}
	response, err := client.GrpcClient.CreateConfig(ctx, &request)
	result = (uintptr)(response.GetResult())
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8003, err, details)
		}()
	}
	return result, err
}

/*
The DeleteDataSource method removes a data source from an existing configuration.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
  - dataSourceCode: A JSON document in the format `{"DSRC_CODE": "NAME_OF_DATASOURCE"}`.
*/
func (client *Szconfig) DeleteDataSource(ctx context.Context, configHandle uintptr, dataSourceCode string) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(9, configHandle, dataSourceCode)
		defer func() { client.traceExit(10, configHandle, dataSourceCode, err, time.Since(entryTime)) }()
	}
	request := szpb.DeleteDataSourceRequest{
		ConfigHandle:   int64(configHandle),
		DataSourceCode: dataSourceCode,
	}
	_, err = client.GrpcClient.DeleteDataSource(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"inputJson": dataSourceCode,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8004, err, details)
		}()
	}
	return err
}

/*
The Destroy method is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szconfig) Destroy(ctx context.Context) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(11)
		defer func() { client.traceExit(12, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8005, err, details)
		}()
	}
	return err
}

/*
The ExportConfig method creates a JSON string representation of the Senzing G2Config object.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.

Output
  - A string containing a JSON Document representation of the Senzing G2Config object.
    See the example output.
*/
func (client *Szconfig) ExportConfig(ctx context.Context, configHandle uintptr) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(23, configHandle)
		defer func() { client.traceExit(24, configHandle, result, err, time.Since(entryTime)) }()
	}
	request := szpb.ExportConfigRequest{
		ConfigHandle: int64(configHandle),
	}
	response, err := client.GrpcClient.ExportConfig(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8009, err, details)
		}()
	}
	return result, err
}

/*
The GetDataSources method returns a JSON document of data sources.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.

Output
  - A string containing a JSON document listing all of the data sources.
    See the example output.
*/
func (client *Szconfig) GetDataSources(ctx context.Context, configHandle uintptr) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(19, configHandle)
		defer func() { client.traceExit(20, configHandle, result, err, time.Since(entryTime)) }()
	}
	request := szpb.GetDataSourcesRequest{
		ConfigHandle: int64(configHandle),
	}
	response, err := client.GrpcClient.GetDataSources(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8007, err, details)
		}()
	}
	return result, err
}

/*
The ImportConfig method initializes the Senzing G2Config object from a JSON string.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
  - configDefinition: A JSON document containing the Senzing configuration.
*/
func (client *Szconfig) ImportConfig(ctx context.Context, configDefinition string) (uintptr, error) {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(21, configDefinition)
		defer func() { client.traceExit(22, configDefinition, err, time.Since(entryTime)) }()
	}
	request := szpb.ImportConfigRequest{
		ConfigDefinition: configDefinition,
	}
	response, err := client.GrpcClient.ImportConfig(ctx, &request)
	result := (uintptr)(response.GetResult())
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8008, err, details)
		}()
	}
	return result, err
}

/*
The Initialize method is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szconfig) Initialize(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(17, instanceName, settings, verboseLogging)
		defer func() { client.traceExit(18, instanceName, settings, verboseLogging, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"settings":       settings,
				"instanceName":   instanceName,
				"verboseLogging": strconv.FormatInt(verboseLogging, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8006, err, details)
		}()
	}
	return err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szconfig) getLogger() logging.LoggingInterface {
	var err error = nil
	if client.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		client.logger, err = logging.NewSenzingSdkLogger(ComponentId, szconfigapi.IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return client.logger
}

// Trace method entry.
func (client *Szconfig) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *Szconfig) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevelName: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *Szconfig) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(25, logLevelName)
		defer func() { client.traceExit(26, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	err = client.getLogger().SetLogLevel(logLevelName)
	client.isTrace = (logLevelName == logging.LevelTraceName)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"logLevel": logLevelName,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8012, err, details)
		}()
	}
	return err
}

// --- Observer ---------------------------------------------------------------

/*
The GetObserverOrigin method returns the "origin" value of past Observer messages.

Input
  - ctx: A context to control lifecycle.

Output
  - The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szconfig) GetObserverOrigin(ctx context.Context) string {
	return client.observerOrigin
}

/*
The GetSdkId method returns the identifier of this particular Software Development Kit (SDK).
It is handy when working with multiple implementations of the same G2configInterface.
For this implementation, "grpc" is returned.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szconfig) GetSdkId(ctx context.Context) string {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(31)
		defer func() { client.traceExit(32, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8010, err, details)
		}()
	}
	return "grpc"
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szconfig) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(27, observer.GetObserverId(ctx))
		defer func() { client.traceExit(28, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	if client.observers == nil {
		client.observers = &subject.SubjectImpl{}
	}
	err = client.observers.RegisterObserver(ctx, observer)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"observerId": observer.GetObserverId(ctx),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8011, err, details)
		}()
	}
	return err
}

/*
The SetObserverOrigin method sets the "origin" value in future Observer messages.

Input
  - ctx: A context to control lifecycle.
  - origin: The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szconfig) SetObserverOrigin(ctx context.Context, origin string) {
	client.observerOrigin = origin
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szconfig) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(29, observer.GetObserverId(ctx))
		defer func() { client.traceExit(30, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerId": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8013, err, details)
	}
	err = client.observers.UnregisterObserver(ctx, observer)
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	return err
}
