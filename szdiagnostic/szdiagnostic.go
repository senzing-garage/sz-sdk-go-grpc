/*
Package szdiagnostic implements a client for the service.
*/
package szdiagnostic

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
	"github.com/senzing-garage/sz-sdk-go/szdiagnostic"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szdiagnostic"
)

type Szdiagnostic struct {
	GrpcClient     szpb.SzDiagnosticClient
	isTrace        bool
	logger         logging.Logging
	observerOrigin string
	observers      subject.Subject
}

const (
	baseCallerSkip       = 4
	baseTen              = 10
	initialByteArraySize = 65535
	noError              = 0
)

// ----------------------------------------------------------------------------
// sz-sdk-go.SzDiagnostic interface methods
// ----------------------------------------------------------------------------

/*
The CheckDatastorePerformance method runs performance tests on the Senzing datastore.

Input
  - ctx: A context to control lifecycle.
  - secondsToRun: Duration of the test in seconds.

Output

  - A JSON document containing performance results.
    Example: `{"numRecordsInserted":0,"insertTime":0}`
*/
func (client *Szdiagnostic) CheckDatastorePerformance(ctx context.Context, secondsToRun int) (string, error) {
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(1, secondsToRun)
		defer func() { client.traceExit(2, secondsToRun, result, err, time.Since(entryTime)) }()
	}
	result, err = client.checkDatastorePerformance(ctx, secondsToRun)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8001, err, details)
		}()
	}
	return result, err
}

/*
The Destroy method is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szdiagnostic) Destroy(ctx context.Context) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(5)
		defer func() { client.traceExit(6, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8002, err, details)
		}()
	}
	return err
}

/*
The GetDatastoreInfo method returns information about the Senzing datastore.

Input
  - ctx: A context to control lifecycle.

Output

  - A JSON document containing Senzing datastore metadata.
*/
func (client *Szdiagnostic) GetDatastoreInfo(ctx context.Context) (string, error) {
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(7)
		defer func() { client.traceExit(8, result, err, time.Since(entryTime)) }()
	}
	result, err = client.getDatastoreInfo(ctx)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8003, err, details)
		}()
	}
	return result, err
}

/*
The GetFeature method is an experimental method that returns diagnostic information of a feature.

Input
  - ctx: A context to control lifecycle.
  - featureID: The identifier of the feature to describe.

Output

  - A JSON document containing feature metadata.
*/
func (client *Szdiagnostic) GetFeature(ctx context.Context, featureID int64) (string, error) {
	var err error
	var result string
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(9, featureID)
		defer func() { client.traceExit(10, featureID, result, err, time.Since(entryTime)) }()
	}
	result, err = client.getFeature(ctx, featureID)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"featureID": strconv.FormatInt(featureID, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8004, err, details)
		}()
	}
	return result, err
}

/*
WARNING: The PurgeRepository method removes every record in the Senzing datastore.
This is a destructive method that cannot be undone.
Before calling purgeRepository(), all programs using Senzing MUST be terminated.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szdiagnostic) PurgeRepository(ctx context.Context) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(17)
		defer func() { client.traceExit(18, err, time.Since(entryTime)) }()
	}
	err = client.purgeRepository(ctx)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8007, err, details)
		}()
	}
	return err
}

/*
The Reinitialize method re-initializes the Senzing SzDiagnostic object.

Input
  - ctx: A context to control lifecycle.
  - configID: The Senzing configuration JSON document identifier used for the initialization.
*/
func (client *Szdiagnostic) Reinitialize(ctx context.Context, configID int64) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(19, configID)
		defer func() { client.traceExit(20, configID, err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configID": strconv.FormatInt(configID, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8008, err, details)
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
func (client *Szdiagnostic) GetObserverOrigin(ctx context.Context) string {
	_ = ctx
	return client.observerOrigin
}

/*
The Initialize method is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - configID: The configuration ID used for the initialization.  0 for current default configuration.
  - verboseLogging: A flag to enable deeper logging of the Sz processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szdiagnostic) Initialize(ctx context.Context, instanceName string, settings string, configID int64, verboseLogging int64) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(15, instanceName, settings, configID, verboseLogging)
		defer func() {
			client.traceExit(16, instanceName, settings, configID, verboseLogging, err, time.Since(entryTime))
		}()
	}
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configID":       strconv.FormatInt(configID, baseTen),
				"instanceName":   instanceName,
				"settings":       settings,
				"verboseLogging": strconv.FormatInt(verboseLogging, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8005, err, details)
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
func (client *Szdiagnostic) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(703, observer.GetObserverID(ctx))
		defer func() { client.traceExit(704, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}
	if client.observers == nil {
		client.observers = &subject.SimpleSubject{}
	}
	err = client.observers.RegisterObserver(ctx, observer)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"observerID": observer.GetObserverID(ctx),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8702, err, details)
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
func (client *Szdiagnostic) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(705, logLevelName)
		defer func() { client.traceExit(706, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	err = client.getLogger().SetLogLevel(logLevelName)
	client.isTrace = (logLevelName == logging.LevelTraceName)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"logLevelName": logLevelName,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8703, err, details)
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
func (client *Szdiagnostic) SetObserverOrigin(ctx context.Context, origin string) {
	_ = ctx
	client.observerOrigin = origin
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szdiagnostic) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(707, observer.GetObserverID(ctx))
		defer func() { client.traceExit(708, observer.GetObserverID(ctx), err, time.Since(entryTime)) }()
	}
	if client.observers != nil {
		// Tricky code:
		// client.notify is called synchronously before client.observers is set to nil.
		// In client.notify, each observer will get notified in a goroutine.
		// Then client.observers may be set to nil, but observer goroutines will be OK.
		details := map[string]string{
			"observerID": observer.GetObserverID(ctx),
		}
		notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8704, err, details)
		err = client.observers.UnregisterObserver(ctx, observer)
		if !client.observers.HasObservers(ctx) {
			client.observers = nil
		}
	}
	return err
}

// ----------------------------------------------------------------------------
// Private methods for gRPC request/response
// ----------------------------------------------------------------------------

func (client *Szdiagnostic) checkDatastorePerformance(ctx context.Context, secondsToRun int) (string, error) {
	request := szpb.CheckDatastorePerformanceRequest{
		SecondsToRun: int32(secondsToRun), //nolint:gosec
	}
	response, err := client.GrpcClient.CheckDatastorePerformance(ctx, &request)
	result := response.GetResult()
	err = helper.ConvertGrpcError(err)
	return result, err
}

func (client *Szdiagnostic) getDatastoreInfo(ctx context.Context) (string, error) {
	request := szpb.GetDatastoreInfoRequest{}
	response, err := client.GrpcClient.GetDatastoreInfo(ctx, &request)
	result := response.GetResult()
	err = helper.ConvertGrpcError(err)
	return result, err
}

func (client *Szdiagnostic) getFeature(ctx context.Context, featureID int64) (string, error) {
	request := szpb.GetFeatureRequest{
		FeatureId: featureID,
	}
	response, err := client.GrpcClient.GetFeature(ctx, &request)
	result := response.GetResult()
	err = helper.ConvertGrpcError(err)
	return result, err
}

func (client *Szdiagnostic) purgeRepository(ctx context.Context) error {
	request := szpb.PurgeRepositoryRequest{}
	_, err := client.GrpcClient.PurgeRepository(ctx, &request)
	err = helper.ConvertGrpcError(err)
	return err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szdiagnostic) getLogger() logging.Logging {
	if client.logger == nil {
		client.logger = helper.GetLogger(ComponentID, szdiagnostic.IDMessages, baseCallerSkip)
	}
	return client.logger
}

// Trace method entry.
func (client *Szdiagnostic) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *Szdiagnostic) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}
