/*
 *
 */

// Package szproduct implements a client for the service.
package szproduct

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
	szproductapi "github.com/senzing-garage/sz-sdk-go/szproduct"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szproduct"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type SzProduct struct {
	GrpcClient     szpb.SzProductClient
	isTrace        bool // Performance optimization
	logger         logging.LoggingInterface
	observerOrigin string
	observers      subject.Subject
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Destroy method is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
*/
func (client *SzProduct) Destroy(ctx context.Context) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(3)
		defer func() { client.traceExit(4, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8001, err, details)
		}()
	}
	return err
}

/*
The Initialize method is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *SzProduct) Initialize(ctx context.Context, instanceName string, settings string, verboseLogging int64) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(9, instanceName, settings, verboseLogging)
		defer func() { client.traceExit(10, instanceName, settings, verboseLogging, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"iniParams":      settings,
				"moduleName":     instanceName,
				"verboseLogging": strconv.FormatInt(verboseLogging, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8002, err, details)
		}()
	}
	return err
}

/*
The GetLicense method retrieves information about the currently used license by the Senzing API.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing Senzing license metadata.
    See the example output.
*/
func (client *SzProduct) GetLicense(ctx context.Context) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(11)
		defer func() { client.traceExit(12, result, err, time.Since(entryTime)) }()
	}
	request := szpb.GetLicenseRequest{}
	response, err := client.GrpcClient.GetLicense(ctx, &request)
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
The GetVersion method returns the version of the Senzing API.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing metadata about the Senzing Engine version being used.
    See the example output.
*/
func (client *SzProduct) GetVersion(ctx context.Context) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(19)
		defer func() { client.traceExit(20, result, err, time.Since(entryTime)) }()
	}
	request := szpb.GetVersionRequest{}
	response, err := client.GrpcClient.GetVersion(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8006, err, details)
		}()
	}
	return result, err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *SzProduct) getLogger() logging.LoggingInterface {
	var err error = nil
	if client.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		client.logger, err = logging.NewSenzingSdkLogger(ComponentId, szproductapi.IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return client.logger
}

// Trace method entry.
func (client *SzProduct) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *SzProduct) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *SzProduct) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(13, logLevelName)
		defer func() { client.traceExit(14, logLevelName, err, time.Since(entryTime)) }()
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
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8009, err, details)
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
func (client *SzProduct) GetObserverOrigin(ctx context.Context) string {
	return client.observerOrigin
}

/*
The GetSdkId method returns the identifier of this particular Software Development Kit (SDK).
It is handy when working with multiple implementations of the same SzProduct Interface.
For this implementation, "grpc" is returned.

Input
  - ctx: A context to control lifecycle.
*/
func (client *SzProduct) GetSdkId(ctx context.Context) string {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(25)
		defer func() { client.traceExit(26, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8007, err, details)
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
func (client *SzProduct) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(21, observer.GetObserverId(ctx))
		defer func() { client.traceExit(22, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
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
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8008, err, details)
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
func (client *SzProduct) SetObserverOrigin(ctx context.Context, origin string) {
	client.observerOrigin = origin
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *SzProduct) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(23, observer.GetObserverId(ctx))
		defer func() { client.traceExit(24, observer.GetObserverId(ctx), err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerId": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8010, err, details)
	}
	err = client.observers.UnregisterObserver(ctx, observer)
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	return err
}
