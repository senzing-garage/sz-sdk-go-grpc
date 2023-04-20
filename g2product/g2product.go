/*
 *
 */

// Package g2product implements a client for the service.
package g2product

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/senzing/g2-sdk-go-grpc/helper"
	g2productapi "github.com/senzing/g2-sdk-go/g2product"
	g2pb "github.com/senzing/g2-sdk-proto/go/g2product"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-observing/notifier"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2product struct {
	GrpcClient g2pb.G2ProductClient
	isTrace    bool // Performance optimization
	logger     logging.LoggingInterface
	observers  subject.Subject
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (client *G2product) getLogger() logging.LoggingInterface {
	var err error = nil
	if client.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		client.logger, err = logging.NewSenzingSdkLogger(ProductId, g2productapi.IdMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return client.logger
}

// Trace method entry.
func (client *G2product) traceEntry(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// Trace method exit.
func (client *G2product) traceExit(errorNumber int, details ...interface{}) {
	client.getLogger().Log(errorNumber, details...)
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Destroy method will destroy and perform cleanup for the Senzing G2Product object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2product) Destroy(ctx context.Context) error {
	if client.isTrace {
		client.traceEntry(3)
	}
	entryTime := time.Now()
	request := g2pb.DestroyRequest{}
	_, err := client.GrpcClient.Destroy(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8001, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(4, err, time.Since(entryTime))
	}
	return err
}

/*
The GetSdkId method returns the identifier of this particular Software Development Kit (SDK).
It is handy when working with multiple implementations of the same G2productInterface.
For this implementation, "grpc" is returned.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2product) GetSdkId(ctx context.Context) string {
	if client.isTrace {
		client.traceEntry(25)
	}
	entryTime := time.Now()
	var err error = nil
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8007, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(26, err, time.Since(entryTime))
	}
	return "grpc"
}

/*
The Init method initializes the Senzing G2Product object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - moduleName: A name for the auditing node, to help identify it within system logs.
  - iniParams: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *G2product) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	if client.isTrace {
		client.traceEntry(9, moduleName, iniParams, verboseLogging)
	}
	entryTime := time.Now()
	request := g2pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	_, err := client.GrpcClient.Init(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"iniParams":      iniParams,
				"moduleName":     moduleName,
				"verboseLogging": strconv.Itoa(verboseLogging),
			}
			notifier.Notify(ctx, client.observers, ProductId, 8002, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(10, moduleName, iniParams, verboseLogging, err, time.Since(entryTime))
	}
	return err
}

/*
The License method retrieves information about the currently used license by the Senzing API.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing Senzing license metadata.
    See the example output.
*/
func (client *G2product) License(ctx context.Context) (string, error) {
	if client.isTrace {
		client.traceEntry(11)
	}
	entryTime := time.Now()
	request := g2pb.LicenseRequest{}
	response, err := client.GrpcClient.License(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8003, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(12, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The RegisterObserver method adds the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2product) RegisterObserver(ctx context.Context, observer observer.Observer) error {
	if client.isTrace {
		client.traceEntry(21, observer.GetObserverId(ctx))
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
			notifier.Notify(ctx, client.observers, ProductId, 8008, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(22, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *G2product) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil
	if client.isTrace {
		entryTime := time.Now()
		client.traceEntry(13, logLevelName)
		defer func() { client.traceExit(14, logLevelName, err, time.Since(entryTime)) }()
	}
	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}
	client.getLogger().SetLogLevel(logLevelName)
	client.isTrace = (logLevelName == logging.LevelTraceName)
	if client.observers != nil {
		go func() {
			details := map[string]string{
				"logLevel": logLevelName,
			}
			notifier.Notify(ctx, client.observers, ProductId, 8009, err, details)
		}()
	}
	return err
}

/*
The UnregisterObserver method removes the observer to the list of observers notified.

Input
  - ctx: A context to control lifecycle.
  - observer: The observer to be added.
*/
func (client *G2product) UnregisterObserver(ctx context.Context, observer observer.Observer) error {
	if client.isTrace {
		client.traceEntry(23, observer.GetObserverId(ctx))
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
		notifier.Notify(ctx, client.observers, ProductId, 8010, err, details)
	}
	err = client.observers.UnregisterObserver(ctx, observer)
	if !client.observers.HasObservers(ctx) {
		client.observers = nil
	}
	if client.isTrace {
		defer client.traceExit(24, observer.GetObserverId(ctx), err, time.Since(entryTime))
	}
	return err
}

/*
The ValidateLicenseFile method validates the licence file has not expired.

Input
  - ctx: A context to control lifecycle.
  - licenseFilePath: A fully qualified path to the Senzing license file.

Output
  - if error is nil, license is valid.
  - If error not nil, license is not valid.
  - The returned string has additional information.
*/
func (client *G2product) ValidateLicenseFile(ctx context.Context, licenseFilePath string) (string, error) {
	if client.isTrace {
		client.traceEntry(15, licenseFilePath)
	}
	entryTime := time.Now()
	request := g2pb.ValidateLicenseFileRequest{
		LicenseFilePath: licenseFilePath,
	}
	response, err := client.GrpcClient.ValidateLicenseFile(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8004, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(16, licenseFilePath, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The ValidateLicenseStringBase64 method validates the licence, represented by a Base-64 string, has not expired.

Input
  - ctx: A context to control lifecycle.
  - licenseString: A Senzing license represented by a Base-64 encoded string.

Output
  - if error is nil, license is valid.
  - If error not nil, license is not valid.
  - The returned string has additional information.
    See the example output.
*/
func (client *G2product) ValidateLicenseStringBase64(ctx context.Context, licenseString string) (string, error) {
	if client.isTrace {
		client.traceEntry(17, licenseString)
	}
	entryTime := time.Now()
	request := g2pb.ValidateLicenseStringBase64Request{
		LicenseString: licenseString,
	}
	response, err := client.GrpcClient.ValidateLicenseStringBase64(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8005, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(18, licenseString, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}

/*
The Version method returns the version of the Senzing API.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing metadata about the Senzing Engine version being used.
    See the example output.
*/
func (client *G2product) Version(ctx context.Context) (string, error) {
	if client.isTrace {
		client.traceEntry(19)
	}
	entryTime := time.Now()
	request := g2pb.VersionRequest{}
	response, err := client.GrpcClient.Version(ctx, &request)
	err = helper.ConvertGrpcError(err)
	if client.observers != nil {
		go func() {
			details := map[string]string{}
			notifier.Notify(ctx, client.observers, ProductId, 8006, err, details)
		}()
	}
	if client.isTrace {
		defer client.traceExit(20, response.GetResult(), err, time.Since(entryTime))
	}
	return response.GetResult(), err
}
