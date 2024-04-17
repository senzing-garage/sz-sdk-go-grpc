/*
Package szconfigmanager implements a client for the service.
*/
package szconfigmanager

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
	szconfigmanagerapi "github.com/senzing-garage/sz-sdk-go/szconfigmanager"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
)

type Szconfigmanager struct {
	GrpcClient     szpb.SzConfigManagerClient
	isTrace        bool
	logger         logging.LoggingInterface
	observerOrigin string
	observers      subject.Subject
}

// ----------------------------------------------------------------------------
// sz-sdk-go.SzConfigManager interface methods
// ----------------------------------------------------------------------------

/*
The AddConfig method adds a Senzing configuration JSON document to the Senzing database.

Input
  - ctx: A context to control lifecycle.
  - configDefinition: The Senzing configuration JSON document.
  - configComment: A free-form string describing the configuration document.

Output
  - A configuration identifier.
*/
func (client *Szconfigmanager) AddConfig(ctx context.Context, configDefinition string, configComment string) (int64, error) {
	var err error = nil
	var result int64 = 0
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(1, configDefinition, configComment)
		defer func() { client.traceExit(2, configDefinition, configComment, result, err, time.Since(entryTime)) }()
	}
	request := szpb.AddConfigRequest{
		ConfigDefinition: configDefinition,
		ConfigComment:    configComment,
	}
	response, err := client.GrpcClient.AddConfig(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configComment": configComment,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8001, err, details)
		}()
	}
	return result, err
}

/*
The Destroy method is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szconfigmanager) Destroy(ctx context.Context) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(5)
		defer func() { client.traceExit(6, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8002, err, details)
		}()
	}
	return err
}

/*
The GetConfig method retrieves a specific Senzing configuration JSON document from the Senzing database.

Input
  - ctx: A context to control lifecycle.
  - configId: The configuration identifier of the desired Senzing Engine configuration JSON document to retrieve.

Output
  - A JSON document containing the Senzing configuration.
    See the example output.
*/
func (client *Szconfigmanager) GetConfig(ctx context.Context, configId int64) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(7, configId)
		defer func() { client.traceExit(8, configId, result, err, time.Since(entryTime)) }()
	}
	request := szpb.GetConfigRequest{
		ConfigId: configId,
	}
	response, err := client.GrpcClient.GetConfig(ctx, &request)
	result = response.GetResult()
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
The GetConfigList method retrieves a list of Senzing configurations from the Senzing database.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing Senzing configurations.
    See the example output.
*/
func (client *Szconfigmanager) GetConfigList(ctx context.Context) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(9)
		defer func() { client.traceExit(10, result, err, time.Since(entryTime)) }()
	}
	request := szpb.GetConfigListRequest{}
	response, err := client.GrpcClient.GetConfigList(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8004, err, details)
		}()
	}
	return result, err
}

/*
The GetDefaultConfigId method retrieves from the Senzing database the configuration identifier of the default Senzing configuration.

Input
  - ctx: A context to control lifecycle.

Output
  - A configuration identifier which identifies the current configuration in use.
*/
func (client *Szconfigmanager) GetDefaultConfigId(ctx context.Context) (int64, error) {
	var err error = nil
	var result int64 = 0
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(11)
		defer func() { client.traceExit(12, result, err, time.Since(entryTime)) }()
	}
	request := szpb.GetDefaultConfigIdRequest{}
	response, err := client.GrpcClient.GetDefaultConfigId(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8005, err, details)
		}()
	}
	return result, err
}

/*
The ReplaceDefaultConfigId method replaces the old configuration identifier with a new configuration identifier in the Senzing database.
It is like a "compare-and-swap" instruction to serialize concurrent editing of configuration.
If currentDefaultConfigId is no longer the "old configuration identifier", the operation will fail.
To simply set the default configuration ID, use SetDefaultConfigId().

Input
  - ctx: A context to control lifecycle.
  - currentDefaultConfigId: The configuration identifier to replace.
  - newDefaultConfigId: The configuration identifier to use as the default.
*/
func (client *Szconfigmanager) ReplaceDefaultConfigId(ctx context.Context, currentDefaultConfigId int64, newDefaultConfigId int64) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(19, currentDefaultConfigId, newDefaultConfigId)
		defer func() { client.traceExit(20, currentDefaultConfigId, newDefaultConfigId, err, time.Since(entryTime)) }()
	}
	request := szpb.ReplaceDefaultConfigIdRequest{
		CurrentDefaultConfigId: currentDefaultConfigId,
		NewDefaultConfigId:     newDefaultConfigId,
	}
	_, err = client.GrpcClient.ReplaceDefaultConfigId(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"newDefaultConfigId": strconv.FormatInt(newDefaultConfigId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8007, err, details)
		}()
	}
	return err
}

/*
The SetDefaultConfigId method replaces the sets a new configuration identifier in the Senzing database.
To serialize modifying of the configuration identifier, see ReplaceDefaultConfigId().

Input
  - ctx: A context to control lifecycle.
  - configId: The configuration identifier of the Senzing Engine configuration to use as the default.
*/
func (client *Szconfigmanager) SetDefaultConfigId(ctx context.Context, configId int64) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(21, configId)
		defer func() { client.traceExit(22, configId, err, time.Since(entryTime)) }()
	}
	request := szpb.SetDefaultConfigIdRequest{
		ConfigId: configId,
	}
	_, err = client.GrpcClient.SetDefaultConfigId(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configId": strconv.FormatInt(configId, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8008, err, details)
		}()
	}
	return err
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

/*
The GetObserverOrigin method returns the "origin" value of past Observer messages.

Input
  - ctx: A context to control lifecycle.

Output
  - The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szconfigmanager) GetObserverOrigin(ctx context.Context) string {
	return client.observerOrigin
}

/*
The GetSdkId method returns the identifier of this particular Software Development Kit (SDK).
It is handy when working with multiple implementations of the same SzConfigManager Interface.
For this implementation, "grpc" is returned.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szconfigmanager) GetSdkId(ctx context.Context) string {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(29)
		defer func() { client.traceExit(30, err, time.Since(entryTime)) }()
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
The Initialize method is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szconfigmanager) Initialize(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
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

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szconfigmanager) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(25, observer.GetObserverId(ctx))
		defer func() { client.traceExit(26, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
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
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8010, err, details)
		}()
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevelName: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *Szconfigmanager) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(23, logLevelName)
		defer func() { client.traceExit(24, logLevelName, err, time.Since(entryTime)) }()
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
func (client *Szconfigmanager) SetObserverOrigin(ctx context.Context, origin string) {
	client.observerOrigin = origin
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szconfigmanager) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(27, observer.GetObserverId(ctx))
		defer func() { client.traceExit(28, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerId": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8012, err, details)
	}
	err = client.observers.UnregisterObserver(ctx, observer)
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	return err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szconfigmanager) getLogger() logging.LoggingInterface {
	var err error = nil
	if client.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		client.logger, err = logging.NewSenzingSdkLogger(ComponentId, szconfigmanagerapi.IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return client.logger
}

// Trace method entry.
func (client *Szconfigmanager) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *Szconfigmanager) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}
