/*
Package szconfigmanager implements a client for the service.
*/
package szconfigmanager

import (
	"context"
	"strconv"
	"time"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
	"github.com/senzing-garage/sz-sdk-go-grpc/helper"
	"github.com/senzing-garage/sz-sdk-go-grpc/szconfig"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"github.com/senzing-garage/sz-sdk-go/szconfigmanager"
	"github.com/senzing-garage/sz-sdk-go/szerror"
	szconfigpb "github.com/senzing-garage/sz-sdk-proto/go/szconfig"
	szpb "github.com/senzing-garage/sz-sdk-proto/go/szconfigmanager"
)

type Szconfigmanager struct {
	GrpcClient         szpb.SzConfigManagerClient
	GrpcClientSzConfig szconfigpb.SzConfigClient

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
// sz-sdk-go.SzConfigManager interface methods
// ----------------------------------------------------------------------------

/*
Method CreateConfigFromConfigID creates a new SzConfig instance for a configuration ID.

Input
  - ctx: A context to control lifecycle.
  - configID: The identifier of the desired Senzing configuration JSON document to retrieve.

Output
  - senzing.SzConfig:
*/
func (client *Szconfigmanager) CreateConfigFromConfigID(ctx context.Context, configID int64) (senzing.SzConfig, error) {
	var (
		err    error
		result senzing.SzConfig
	)

	if client.isTrace {
		client.traceEntry(7, configID)

		entryTime := time.Now()

		defer func() { client.traceExit(8, configID, result, err, time.Since(entryTime)) }()
	}

	result, err = client.createConfigFromConfigIDChoreography(ctx, configID)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8003, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method CreateConfigFromString creates a new SzConfig instance from a configuration definition.

Input
  - ctx: A context to control lifecycle.
  - configDefinition: The Senzing configuration JSON document.

Output
  - senzing.SzConfig:
*/
func (client *Szconfigmanager) CreateConfigFromString(
	ctx context.Context,
	configDefinition string,
) (senzing.SzConfig, error) {
	var (
		err    error
		result senzing.SzConfig
	)

	if client.isTrace {
		client.traceEntry(23, configDefinition)

		entryTime := time.Now()

		defer func() { client.traceExit(24, configDefinition, result, err, time.Since(entryTime)) }()
	}

	result, err = client.createConfigFromString(ctx, configDefinition)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8009, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method CreateConfigFromTemplate creates a new SzConfig instance from the template configuration definition.

Input
  - ctx: A context to control lifecycle.

Output
  - senzing.SzConfig:
*/
func (client *Szconfigmanager) CreateConfigFromTemplate(ctx context.Context) (senzing.SzConfig, error) {
	var (
		err    error
		result senzing.SzConfig
	)

	if client.isTrace {
		client.traceEntry(25)

		entryTime := time.Now()

		defer func() { client.traceExit(26, result, err, time.Since(entryTime)) }()
	}

	result, err = client.createConfigFromTemplateChoreography(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8010, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method Destroy will destroy and perform cleanup for the Senzing SzConfigMgr object.

It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *Szconfigmanager) Destroy(ctx context.Context) error {
	var err error

	if client.isTrace {
		client.traceEntry(5)

		entryTime := time.Now()

		defer func() { client.traceExit(6, err, time.Since(entryTime)) }()
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8002, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method GetConfigRegistry gets the configuration registry.

The registry contains the original timestamp, original comment, and configuration ID of all configurations ever
registered with the repository.

Registered configurations cannot be unregistered.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document listing Senzing configuration JSON document metadata.
*/
func (client *Szconfigmanager) GetConfigRegistry(ctx context.Context) (string, error) {
	var (
		err    error
		result string
	)

	if client.isTrace {
		client.traceEntry(9)

		entryTime := time.Now()

		defer func() { client.traceExit(10, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getConfigRegistry(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8004, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method GetDefaultConfigID gets the default configuration ID for the repository.

Unless an explicit configuration ID is specified at initialization, the default configuration ID is used.

This may not be the same as the active configuration ID.

Input
  - ctx: A context to control lifecycle.

Output
  - configID: The default Senzing configuration JSON document identifier. If none exists, zero (0) is returned.
*/
func (client *Szconfigmanager) GetDefaultConfigID(ctx context.Context) (int64, error) {
	var (
		err    error
		result int64
	)

	if client.isTrace {
		client.traceEntry(11)

		entryTime := time.Now()

		defer func() { client.traceExit(12, result, err, time.Since(entryTime)) }()
	}

	result, err = client.getDefaultConfigID(ctx)

	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8005, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method RegisterConfig registers a configuration definition in the repository.

Registered configurations do not become immediately active nor do they become the default.

Registered configurations cannot be unregistered.

Input
  - ctx: A context to control lifecycle.
  - configDefinition: The Senzing configuration JSON document.
  - configComment: A free-form string describing the Senzing configuration JSON document.

Output
  - configID: A Senzing configuration JSON document identifier.
*/
func (client *Szconfigmanager) RegisterConfig(
	ctx context.Context,
	configDefinition string,
	configComment string,
) (int64, error) {
	var (
		err    error
		result int64
	)

	if client.isTrace {
		client.traceEntry(1, configDefinition, configComment)

		entryTime := time.Now()

		defer func() {
			client.traceExit(2, configDefinition, configComment, result, err, time.Since(entryTime))
		}()
	}

	result, err = client.registerConfig(ctx, configDefinition, configComment)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configComment": configComment,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8001, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method ReplaceDefaultConfigID replaces the existing default configuration ID with a new configuration ID.

The change is prevented if the current default configuration ID value is not as expected.

Use this in place of setDefaultConfigID() to handle race conditions.

Input
  - ctx: A context to control lifecycle.
  - currentDefaultConfigID: The Senzing configuration JSON document identifier to replace.
  - newDefaultConfigID: The Senzing configuration JSON document identifier to use as the default.
*/
func (client *Szconfigmanager) ReplaceDefaultConfigID(
	ctx context.Context,
	currentDefaultConfigID int64,
	newDefaultConfigID int64,
) error {
	var err error

	if client.isTrace {
		client.traceEntry(19, currentDefaultConfigID, newDefaultConfigID)

		entryTime := time.Now()

		defer func() { client.traceExit(20, currentDefaultConfigID, newDefaultConfigID, err, time.Since(entryTime)) }()
	}

	err = client.replaceDefaultConfigID(ctx, currentDefaultConfigID, newDefaultConfigID)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"newDefaultConfigID": strconv.FormatInt(newDefaultConfigID, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8007, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method SetDefaultConfig registers a configuration in the repository and sets its ID as the default for the repository.

Convenience method for registerConfig() followed by setDefaultConfigId().

Input
  - ctx: A context to control lifecycle.
  - configDefinition: The Senzing configuration JSON document.
  - configComment: A free-form string describing the Senzing configuration JSON document.
*/
func (client *Szconfigmanager) SetDefaultConfig(
	ctx context.Context,
	configDefinition string,
	configComment string,
) (int64, error) {
	var (
		err    error
		result int64
	)

	if client.isTrace {
		client.traceEntry(27, configDefinition, configComment)

		entryTime := time.Now()

		defer func() { client.traceExit(28, configDefinition, configComment, err, time.Since(entryTime)) }()
	}

	result, err = client.setDefaultConfigChoreography(ctx, configDefinition, configComment)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configDefinition":   configDefinition,
				"configComment":      configComment,
				"newDefaultConfigID": strconv.FormatInt(result, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8011, err, details)
		}()
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method SetDefaultConfigID sets the default configuration ID.

Usually this method is sufficient for setting the default configuration ID.
However in concurrent environments that could encounter race conditions,
consider using replaceDefaultConfigId() instead.

Input
  - ctx: A context to control lifecycle.
  - configID: The Senzing configuration JSON document identifier to use as the default.
*/
func (client *Szconfigmanager) SetDefaultConfigID(ctx context.Context, configID int64) error {
	var err error

	if client.isTrace {
		client.traceEntry(21, configID)

		entryTime := time.Now()

		defer func() { client.traceExit(22, configID, err, time.Since(entryTime)) }()
	}

	err = client.setDefaultConfigID(ctx, configID)

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configID": strconv.FormatInt(configID, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8008, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Public non-interface methods
// ----------------------------------------------------------------------------

/*
Method GetObserverOrigin returns the "origin" value of past Observer messages.

Input
  - ctx: A context to control lifecycle.

Output
  - The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szconfigmanager) GetObserverOrigin(ctx context.Context) string {
	_ = ctx

	return client.observerOrigin
}

/*
Method Initialize is a Null function for sz-sdk-go-grpc.

Input
  - ctx: A context to control lifecycle.
  - instanceName: A name for the auditing node, to help identify it within system logs.
  - settings: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the Sz processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *Szconfigmanager) Initialize(
	ctx context.Context,
	instanceName string,
	settings string,
	verboseLogging int64,
) error {
	var err error

	if client.isTrace {
		client.traceEntry(17, instanceName, settings, verboseLogging)

		entryTime := time.Now()

		defer func() { client.traceExit(18, instanceName, settings, verboseLogging, err, time.Since(entryTime)) }()
	}

	if client.observers != nil {
		go func() {
			details := map[string]string{
				"instanceName":   instanceName,
				"settings":       settings,
				"verboseLogging": strconv.FormatInt(verboseLogging, baseTen),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentID, 8006, err, details)
		}()
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method RegisterObserver adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szconfigmanager) RegisterObserver(ctx context.Context, observer observer.Observer) error {
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

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method SetLogLevel sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevelName: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *Szconfigmanager) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error

	if client.isTrace {
		client.traceEntry(705, logLevelName)

		entryTime := time.Now()

		defer func() { client.traceExit(706, logLevelName, err, time.Since(entryTime)) }()
	}

	if !logging.IsValidLogLevelName(logLevelName) {
		return wraperror.Errorf(szerror.ErrSzSdk, "invalid error level: %s", logLevelName)
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

	return wraperror.Errorf(err, wraperror.NoMessage)
}

/*
Method SetObserverOrigin sets the "origin" value in future Observer messages.

Input
  - ctx: A context to control lifecycle.
  - origin: The value sent in the Observer's "origin" key/value pair.
*/
func (client *Szconfigmanager) SetObserverOrigin(ctx context.Context, origin string) {
	_ = ctx
	client.observerOrigin = origin
}

/*
Method UnregisterObserver removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *Szconfigmanager) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
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

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func (client *Szconfigmanager) createConfigFromConfigIDChoreography(
	ctx context.Context,
	configID int64,
) (senzing.SzConfig, error) {
	var err error

	configDefinition, err := client.getConfig(ctx, configID)
	if err != nil {
		return nil, wraperror.Errorf(err, "getConfig")
	}

	return client.createConfigFromString(ctx, configDefinition)
}

func (client *Szconfigmanager) createConfigFromString(
	ctx context.Context,
	configDefinition string,
) (senzing.SzConfig, error) {
	var err error

	result := &szconfig.Szconfig{
		GrpcClient: client.GrpcClientSzConfig,
	}

	err = result.VerifyConfigDefinition(ctx, configDefinition)
	if err != nil {
		return result, wraperror.Errorf(err, "VerifyConfigDefinition")
	}

	err = result.Import(ctx, configDefinition)

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

func (client *Szconfigmanager) createConfigFromTemplateChoreography(ctx context.Context) (senzing.SzConfig, error) {
	var err error

	request := szpb.GetTemplateConfigRequest{}
	response, err := client.GrpcClient.GetTemplateConfig(ctx, &request)

	err = helper.ConvertGrpcError(err)
	if err != nil {
		return nil, wraperror.Errorf(err, "ConvertGrpcError")
	}

	return client.createConfigFromString(ctx, response.GetResult())
}

func (client *Szconfigmanager) setDefaultConfigChoreography(
	ctx context.Context,
	configDefinition string,
	configComment string,
) (int64, error) {
	var (
		err    error
		result int64
	)

	result, err = client.registerConfig(ctx, configDefinition, configComment)
	if err != nil {
		return 0, wraperror.Errorf(err, "registerConfig")
	}

	err = client.setDefaultConfigID(ctx, result)

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private methods for gRPC request/response
// ----------------------------------------------------------------------------

func (client *Szconfigmanager) registerConfig(
	ctx context.Context,
	configDefinition string,
	configComment string,
) (int64, error) {
	request := szpb.RegisterConfigRequest{
		ConfigDefinition: configDefinition,
		ConfigComment:    configComment,
	}
	response, err := client.GrpcClient.RegisterConfig(ctx, &request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szconfigmanager) getConfig(ctx context.Context, configID int64) (string, error) {
	request := szpb.GetConfigRequest{
		ConfigId: configID,
	}
	response, err := client.GrpcClient.GetConfig(ctx, &request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szconfigmanager) getConfigRegistry(ctx context.Context) (string, error) {
	request := szpb.GetConfigRegistryRequest{}
	response, err := client.GrpcClient.GetConfigRegistry(ctx, &request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szconfigmanager) getDefaultConfigID(ctx context.Context) (int64, error) {
	request := szpb.GetDefaultConfigIdRequest{}
	response, err := client.GrpcClient.GetDefaultConfigId(ctx, &request)
	result := response.GetResult()

	return result, helper.ConvertGrpcError(err)
}

func (client *Szconfigmanager) replaceDefaultConfigID(
	ctx context.Context,
	currentDefaultConfigID int64,
	newDefaultConfigID int64,
) error {
	request := szpb.ReplaceDefaultConfigIdRequest{
		CurrentDefaultConfigId: currentDefaultConfigID,
		NewDefaultConfigId:     newDefaultConfigID,
	}
	_, err := client.GrpcClient.ReplaceDefaultConfigId(ctx, &request)

	return helper.ConvertGrpcError(err)
}

func (client *Szconfigmanager) setDefaultConfigID(ctx context.Context, configID int64) error {
	request := szpb.SetDefaultConfigIdRequest{
		ConfigId: configID,
	}
	_, err := client.GrpcClient.SetDefaultConfigId(ctx, &request)

	return helper.ConvertGrpcError(err)
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *Szconfigmanager) getLogger() logging.Logging {
	if client.logger == nil {
		client.logger = helper.GetLogger(ComponentID, szconfigmanager.IDMessages, baseCallerSkip)
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
