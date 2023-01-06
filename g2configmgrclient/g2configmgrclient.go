/*
 *
 */

// Package g2configmgrclient implements a client for the service.
package g2configmgrclient

import (
	"context"

	pb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing/go-logging/logger"
)

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
func (client *G2configmgrClient) AddConfig(ctx context.Context, configStr string, configComments string) (int64, error) {
	request := pb.AddConfigRequest{
		ConfigStr:      configStr,
		ConfigComments: configComments,
	}
	response, err := client.GrpcClient.AddConfig(ctx, &request)
	return response.GetResult(), err
}

/*
The Destroy method will destroy and perform cleanup for the Senzing G2ConfigMgr object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2configmgrClient) Destroy(ctx context.Context) error {
	request := pb.DestroyRequest{}
	_, err := client.GrpcClient.Destroy(ctx, &request)
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
func (client *G2configmgrClient) GetConfig(ctx context.Context, configID int64) (string, error) {
	request := pb.GetConfigRequest{
		ConfigID: configID,
	}
	response, err := client.GrpcClient.GetConfig(ctx, &request)
	return response.GetResult(), err
}

/*
The GetConfigList method retrieves a list of Senzing configurations from the Senzing database.

Input
  - ctx: A context to control lifecycle.

Output
  - A JSON document containing Senzing configurations.
    See the example output.
*/
func (client *G2configmgrClient) GetConfigList(ctx context.Context) (string, error) {
	request := pb.GetConfigListRequest{}
	response, err := client.GrpcClient.GetConfigList(ctx, &request)
	return response.GetResult(), err
}

/*
The GetDefaultConfigID method retrieves from the Senzing database the configuration identifier of the default Senzing configuration.

Input
  - ctx: A context to control lifecycle.

Output
  - A configuration identifier which identifies the current configuration in use.
*/
func (client *G2configmgrClient) GetDefaultConfigID(ctx context.Context) (int64, error) {
	request := pb.GetDefaultConfigIDRequest{}
	response, err := client.GrpcClient.GetDefaultConfigID(ctx, &request)
	return response.GetConfigID(), err
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
func (client *G2configmgrClient) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	request := pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	_, err := client.GrpcClient.Init(ctx, &request)
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
func (client *G2configmgrClient) ReplaceDefaultConfigID(ctx context.Context, oldConfigID int64, newConfigID int64) error {
	request := pb.ReplaceDefaultConfigIDRequest{
		OldConfigID: oldConfigID,
		NewConfigID: newConfigID,
	}
	_, err := client.GrpcClient.ReplaceDefaultConfigID(ctx, &request)
	return err
}

/*
The SetDefaultConfigID method replaces the sets a new configuration identifier in the Senzing database.
To serialize modifying of the configuration identifier, see ReplaceDefaultConfigID().

Input
  - ctx: A context to control lifecycle.
  - configID: The configuration identifier of the Senzing Engine configuration to use as the default.
*/
func (client *G2configmgrClient) SetDefaultConfigID(ctx context.Context, configID int64) error {
	request := pb.SetDefaultConfigIDRequest{
		ConfigID: configID,
	}
	_, err := client.GrpcClient.SetDefaultConfigID(ctx, &request)
	return err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *G2configmgrClient) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	var err error = nil
	return err
}
