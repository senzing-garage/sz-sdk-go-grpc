/*
Package szconfig implements a client for the service.
*/
package szconfig

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go/szconfig"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
)

type Szconfig struct {
	configDefinition string
	GrpcClient       szpb.SzConfigClient
	isTrace          bool
	logger           logging.Logging
	observerOrigin   string
	observers        subject.Subject
}

const (
	baseCallerSkip       = 4
	baseTen              = 10
	initialByteArraySize = 65535
	noError              = 0
)

// ----------------------------------------------------------------------------
// sz-sdk-go.SzConfig interface methods
// ----------------------------------------------------------------------------

/*
Method AddDataSource adds a new data source to the Senzing configuration.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Unique identifier of the data source (e.g. "TEST_DATASOURCE").

Output
  - A JSON document listing the newly created data source.
*/
func (client *Szconfig) AddDataSource(ctx context.Context, dataSourceCode string) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(1, dataSourceCode)

		entryTime := time.Now()
		defer func() {
			client.traceExit(2, dataSourceCode, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.addDataSource(ctx, dataSourceCode)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
				"return":         result,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8001, err, details)
		}()
	}

	return result, wraperror.Errorf(err, "szconfig.AddDataSource error: %w", err)
}

/*
Method DeleteDataSource removes a data source from the Senzing configuration.

Input
  - ctx: A context to control lifecycle.
  - dataSourceCode: Unique identifier of the data source (e.g. "TEST_DATASOURCE").

Output
  - A JSON document listing the newly created data source. Currently an empty string.
*/
func (client *Szconfig) DeleteDataSource(ctx context.Context, dataSourceCode string) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(9, dataSourceCode)

		entryTime := time.Now()
		defer func() { client.traceExit(10, dataSourceCode, err, time.Since(entryTime)) }()
	}

	result, err = client.deleteDataSource(ctx, dataSourceCode)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"dataSourceCode": dataSourceCode,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8004, err, details)
		}()
	}

	return result, wraperror.Errorf(err, "szconfig.DeleteDataSource error: %w", err)
}

/*
Method Export retrieves the Senzing configuration JSON document.

Input
  - ctx: A context to control lifecycle.

Output
  - configDefinition: A Senzing configuration JSON document representation of the in-memory configuration.
*/
func (client *Szconfig) Export(ctx context.Context) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(13)

		entryTime := time.Now()
		defer func() { client.traceExit(14, result, err, time.Since(entryTime)) }()
	}

	result, err = client.export(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8006, err, details)
		}()
	}

	return result, wraperror.Errorf(err, "szconfig.Export error: %w", err)
}

/*
Method GetDataSources returns a JSON document containing data sources defined in the Senzing configuration.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document listing data sources in the in-memory configuration.
*/
func (client *Szconfig) GetDataSources(ctx context.Context) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(15)

		entryTime := time.Now()
		defer func() { client.traceExit(16, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getDataSources(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8008, err, details)
		}()
	}

	return result, wraperror.Errorf(err, "szconfig.GetDataSources error: %w", err)
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

/*
Method Destroy is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szconfig) Destroy(ctx context.Context) error {
	var err error

	if client.isTrace {
		client.traceEntry(11)

		entryTime := time.Now()
		defer func() { client.traceExit(12, err, time.Since(entryTime)) }()
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8005, err, details)
		}()
	}

	return wraperror.Errorf(err, "szconfig.Destroy error: %w", err)
}

/*
Method GetObserverOrigin returns the "origin" value of past Observer messages.

Input
  - ctx: A context to control lifecycle.

Output
  - The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szconfig) GetObserverOrigin(ctx context.Context) string {
	_ = ctx

	return client.observerOrigin
}

/*
Method Import sets the value of the Senzing configuration to be operated upon.

Input
  - ctx: A context to control lifecycle.
  - configDefinition: A Senzing configuration JSON document.
*/
func (client *Szconfig) Import(ctx context.Context, configDefinition string) error {
	var err error

	if client.isTrace {
		client.traceEntry(21, configDefinition)

		entryTime := time.Now()
		defer func() { client.traceExit(22, configDefinition, err, time.Since(entryTime)) }()
	}

	err = client.importConfigDefinition(ctx, configDefinition)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8009, err, details)
		}()
	}

	return wraperror.Errorf(err, "szconfig.Import error: %w", err)
}

/*
Method ImportTemplate retrieves a Senzing configuration from the default template.
The default template is the Senzing configuration JSON document file,
g2config.json, located in the PIPELINE.RESOURCEPATH path.

Input
  - ctx: A context to control lifecycle.

Output
  - configDefinition: A Senzing configuration JSON document.
*/
func (client *Szconfig) ImportTemplate(ctx context.Context) error {
	var (
		err              error
		configDefinition string
	)

	if client.isTrace {
		client.traceEntry(7)

		entryTime := time.Now()
		defer func() { client.traceExit(8, configDefinition, err, time.Since(entryTime)) }()
	}

	// TODO: Implement.

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8003, err, details)
		}()
	}

	return wraperror.Errorf(err, "szconfig.ImportTemplate error: %w", err)
}

/*
Method Initialize is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the Sz processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szconfig) Initialize(
	ctx context.Context,
	instanceName string,
	settings string,
	verboseLogging int64) error {
	var err error

	if client.isTrace {
		client.traceEntry(23, instanceName, settings, verboseLogging)

		entryTime := time.Now()
		defer func() { client.traceExit(24, instanceName, settings, verboseLogging, err, time.Since(entryTime)) }()
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"instanceName":   instanceName,
				"settings":       settings,
				"verboseLogging": strconv.FormatInt(verboseLogging, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8007, err, details)
		}()
	}

	return wraperror.Errorf(err, "szconfig.Initialize error: %w", err)
}

/*
Method RegisterObserver adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szconfig) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if client.isTrace {
		client.traceEntry(703, observer.GetObserverID(ctx))

		entryTime := time.Now()
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

	return wraperror.Errorf(err, "szconfig.RegisterObserver error: %w", err)
}

/*
Method SetLogLevel sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevelName: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *Szconfig) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error

	if client.isTrace {
		client.traceEntry(705, logLevelName)

		entryTime := time.Now()
		defer func() { client.traceExit(706, logLevelName, err, time.Since(entryTime)) }()
	}

	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s; %w", logLevelName, szerror.ErrSzSdk)
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

	return wraperror.Errorf(err, "szconfig.SetLogLevel error: %w", err)
}

/*
Method SetObserverOrigin sets the "origin" value in future Observer messages.

Input
  - ctx: A context to control lifecycle.
  - origin: The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szconfig) SetObserverOrigin(ctx context.Context, origin string) {
	_ = ctx
	client.observerOrigin = origin
}

/*
Method UnregisterObserver removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szconfig) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	var err error

	if client.isTrace {
		client.traceEntry(707, observer.GetObserverID(ctx))

		entryTime := time.Now()
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

	return wraperror.Errorf(err, "szconfig.UnregisterObserver error: %w", err)
}

/*
Method VerifyConfigDefinition determines if the Senzing configuration JSON document is syntactically correct.
If no error is returned, the JSON document is valid.

Input
  - ctx: A context to control lifecycle.
  - configDefinition: A Senzing configuration JSON document.
*/
func (client *Szconfig) VerifyConfigDefinition(ctx context.Context, configDefinition string) error {
	var err error

	if client.isTrace {
		client.traceEntry(25, configDefinition)

		entryTime := time.Now()
		defer func() { client.traceExit(26, configDefinition, err, time.Since(entryTime)) }()
	}

	// TODO: Implement.

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8010, err, details)
		}()
	}

	return wraperror.Errorf(err, "szconfig.VerifyConfigDefinition error: %w", err)
}

// ----------------------------------------------------------------------------
// Private methods for gRPC request/response
// ----------------------------------------------------------------------------

func (client *Szconfig) addDataSource(
	ctx context.Context,
	dataSourceCode string,
) (string, error) {
	var (
		result string
	)

	request := &szpb.AddDataSourceRequest{
		ConfigDefinition: client.configDefinition,
		DataSourceCode:   dataSourceCode,
	}

	response, err := client.GrpcClient.AddDataSource(ctx, request)
	if err != nil {
		return result, helper.ConvertGrpcError(err)
	}

	result = response.GetResult()
	err = client.importConfigDefinition(ctx, response.GetConfigDefinition())

	return result, helper.ConvertGrpcError(err)
}

func (client *Szconfig) deleteDataSource(ctx context.Context, dataSourceCode string) (string, error) {
	var (
		result string
	)

	request := &szpb.DeleteDataSourceRequest{
		ConfigDefinition: client.configDefinition,
		DataSourceCode:   dataSourceCode,
	}

	response, err := client.GrpcClient.DeleteDataSource(ctx, request)
	if err != nil {
		return result, helper.ConvertGrpcError(err)
	}

	result = response.GetResult()
	err = client.importConfigDefinition(ctx, response.GetConfigDefinition())

	return result, helper.ConvertGrpcError(err)
}

func (client *Szconfig) export(ctx context.Context) (string, error) {
	_ = ctx
	return client.configDefinition, nil
}

func (client *Szconfig) getDataSources(ctx context.Context) (string, error) {
	var (
		result string
	)

	request := &szpb.GetDataSourcesRequest{
		ConfigDefinition: client.configDefinition,
	}

	response, err := client.GrpcClient.GetDataSources(ctx, request)
	if err != nil {
		return result, helper.ConvertGrpcError(err)
	}

	result = response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szconfig) importConfigDefinition(ctx context.Context, configDefinition string) error {
	_ = ctx
	client.configDefinition = configDefinition
	return nil
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szconfig) getLogger() logging.Logging {
	if client.logger == nil {
		client.logger = helper.GetLogger(ComponentID, szconfig.IDMessages, baseCallerSkip)
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
