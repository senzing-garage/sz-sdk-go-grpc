/*
 *
 */

// Package g2productclient implements a client for the service.
package g2productclient

import (
	"context"

	pb "github.com/senzing/g2-sdk-proto/go/g2product"
	"github.com/senzing/go-logging/logger"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Destroy method will destroy and perform cleanup for the Senzing G2Product object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2productClient) Destroy(ctx context.Context) error {
	request := pb.DestroyRequest{}
	_, err := client.GrpcClient.Destroy(ctx, &request)
	return err
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
func (client *G2productClient) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	request := pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	_, err := client.GrpcClient.Init(ctx, &request)
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
func (client *G2productClient) License(ctx context.Context) (string, error) {
	request := pb.LicenseRequest{}
	response, err := client.GrpcClient.License(ctx, &request)
	return response.GetResult(), err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *G2productClient) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	var err error = nil
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
func (client *G2productClient) ValidateLicenseFile(ctx context.Context, licenseFilePath string) (string, error) {
	request := pb.ValidateLicenseFileRequest{
		LicenseFilePath: licenseFilePath,
	}
	response, err := client.GrpcClient.ValidateLicenseFile(ctx, &request)
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
func (client *G2productClient) ValidateLicenseStringBase64(ctx context.Context, licenseString string) (string, error) {
	request := pb.ValidateLicenseStringBase64Request{
		LicenseString: licenseString,
	}
	response, err := client.GrpcClient.ValidateLicenseStringBase64(ctx, &request)
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
func (client *G2productClient) Version(ctx context.Context) (string, error) {
	request := pb.VersionRequest{}
	response, err := client.GrpcClient.Version(ctx, &request)
	return response.GetResult(), err
}
