/*
 *
 */

// Package g2configmgrclient implements a client for the service.
package g2configmgr

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/senzing-garage/g2-sdk-go-grpc/helper"
	g2configmgrapi "github.com/senzing-garage/g2-sdk-go/g2configmgr"
	g2pb "github.com/senzing-garage/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/notifier"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2configmgr struct {
	GrpcClient     g2pb.G2ConfigMgrClient
	isTrace        bool // Performance optimization
	logger         logging.LoggingInterface
	observerOrigin string
	observers      subject.Subject
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *G2configmgr) getLogger() logging.LoggingInterface {
	var err error = nil
	if client.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		client.logger, err = logging.NewSenzingSdkLogger(ComponentId, g2configmgrapi.IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return client.logger
}

// Trace method entry.
func (client *G2configmgr) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *G2configmgr) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The AddConfig method adds a Senzing configuration JSON document to the Senzing database.

Input
  - ctx: A context to control lifecycle.
  - configStr: The Senzing configuration JSON document.
  - configComments: A free-form string of comments describing the configuration document.

Output
  - A configuration identifier.
*/
func (client *G2configmgr) AddConfig(ctx context.Context, configStr string, configComments string) (int64, error) {
	var err error = nil
	var result int64 = 0
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(1, configStr, configComments)
		defer func() { client.traceExit(2, configStr, configComments, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.AddConfigRequest{
		ConfigStr:      configStr,
		ConfigComments: configComments,
	}
	response, err := client.GrpcClient.AddConfig(ctx, &request)
	result = response.GetResult()
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configComments": configComments,
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8001, err, details)
		}()
	}
	return result, err
}

/*
The Destroy method will destroy and perform cleanup for the Senzing G2ConfigMgr object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2configmgr) Destroy(ctx context.Context) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(5)
		defer func() { client.traceExit(6, err, time.Since(entryTime)) }()
	}
	request := g2pb.DestroyRequest{}
	_, err = client.GrpcClient.Destroy(ctx, &request)
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
The GetConfig method retrieves a specific Senzing configuration JSON document from the Senzing database.

Input
  - ctx: A context to control lifecycle.
  - configID: The configuration identifier of the desired Senzing Engine configuration JSON document to retrieve.

Output
  - A JSON document containing the Senzing configuration.
    See the example output.
*/
func (client *G2configmgr) GetConfig(ctx context.Context, configID int64) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(7, configID)
		defer func() { client.traceExit(8, configID, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetConfigRequest{
		ConfigID: configID,
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
func (client *G2configmgr) GetConfigList(ctx context.Context) (string, error) {
	var err error = nil
	var result string = ""
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(9)
		defer func() { client.traceExit(10, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetConfigListRequest{}
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
The GetDefaultConfigID method retrieves from the Senzing database the configuration identifier of the default Senzing configuration.

Input
  - ctx: A context to control lifecycle.

Output
  - A configuration identifier which identifies the current configuration in use.
*/
func (client *G2configmgr) GetDefaultConfigID(ctx context.Context) (int64, error) {
	var err error = nil
	var result int64 = 0
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(11)
		defer func() { client.traceExit(12, result, err, time.Since(entryTime)) }()
	}
	request := g2pb.GetDefaultConfigIDRequest{}
	response, err := client.GrpcClient.GetDefaultConfigID(ctx, &request)
	result = response.GetConfigID()
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
The GetObserverOrigin method returns the "origin" value of past Observer messages.

Input
  - ctx: A context to control lifecycle.

Output
  - The value sent in the Observer's "origin" key/value pair.
*/
func (client *G2configmgr) GetObserverOrigin(ctx context.Context) string {
	return client.observerOrigin
}

/*
The GetSdkId method returns the identifier of this particular Software Development Kit (SDK).
It is handy when working with multiple implementations of the same G2configmgrInterface.
For this implementation, "grpc" is returned.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2configmgr) GetSdkId(ctx context.Context) string {
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
The Init method initializes the Senzing G2ConfigMgr object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - moduleName: A name for the auditing node, to help identify it within system logs.
  - iniParams: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *G2configmgr) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int64) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(17, moduleName, iniParams, verboseLogging)
		defer func() { client.traceExit(18, moduleName, iniParams, verboseLogging, err, time.Since(entryTime)) }()
	}
	request := g2pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: verboseLogging,
	}
	_, err = client.GrpcClient.Init(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"iniParams":      iniParams,
				"moduleName":     moduleName,
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
func (client *G2configmgr) RegisterObserver(ctx context.Context, observer observer.Observer) error {
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
				"observerID": observer.GetObserverId(ctx),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8010, err, details)
		}()
	}
	return err
}

/*
The ReplaceDefaultConfigID method replaces the old configuration identifier with a new configuration identifier in the Senzing database.
It is like a "compare-and-swap" instruction to serialize concurrent editing of configuration.
If oldConfigID is no longer the "old configuration identifier", the operation will fail.
To simply set the default configuration ID, use SetDefaultConfigID().

Input
  - ctx: A context to control lifecycle.
  - oldConfigID: The configuration identifier to replace.
  - newConfigID: The configuration identifier to use as the default.
*/
func (client *G2configmgr) ReplaceDefaultConfigID(ctx context.Context, oldConfigID int64, newConfigID int64) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(19, oldConfigID, newConfigID)
		defer func() { client.traceExit(20, oldConfigID, newConfigID, err, time.Since(entryTime)) }()
	}
	request := g2pb.ReplaceDefaultConfigIDRequest{
		OldConfigID: oldConfigID,
		NewConfigID: newConfigID,
	}
	_, err = client.GrpcClient.ReplaceDefaultConfigID(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"newConfigID": strconv.FormatInt(newConfigID, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8007, err, details)
		}()
	}
	return err
}

/*
The SetDefaultConfigID method replaces the sets a new configuration identifier in the Senzing database.
To serialize modifying of the configuration identifier, see ReplaceDefaultConfigID().

Input
  - ctx: A context to control lifecycle.
  - configID: The configuration identifier of the Senzing Engine configuration to use as the default.
*/
func (client *G2configmgr) SetDefaultConfigID(ctx context.Context, configID int64) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(21, configID)
		defer func() { client.traceExit(22, configID, err, time.Since(entryTime)) }()
	}
	request := g2pb.SetDefaultConfigIDRequest{
		ConfigID: configID,
	}
	_, err = client.GrpcClient.SetDefaultConfigID(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"configID": strconv.FormatInt(configID, 10),
			}
			notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8008, err, details)
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
func (client *G2configmgr) SetLogLevel(ctx context.Context, logLevelName string) error {
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
func (client *G2configmgr) SetObserverOrigin(ctx context.Context, origin string) {
	client.observerOrigin = origin
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2configmgr) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
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
			"observerID": observer.GetObserverId(ctx),
		}
		notifier.Notify(ctx, client.observers, client.observerOrigin, ComponentId, 8012, err, details)
	}
	err = client.observers.UnregisterObserver(ctx, observer)
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	return err
}
