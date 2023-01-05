/*
 *
 */

// Package main implements a client for the service.
package g2configclient

import (
	"context"

	pb "github.com/senzing/g2-sdk-proto/go/g2config"
	"github.com/senzing/go-logging/logger"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The AddDataSource method adds a data source to an existing in-memory configuration.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
  - inputJson: A JSON document in the format `{"DSRC_CODE": "NAME_OF_DATASOURCE"}`.

Output
  - A string containing a JSON document listing the newly created data source.
    See the example output.
*/
func (client *G2configClient) AddDataSource(ctx context.Context, configHandle uintptr, inputJson string) (string, error) {
	request := pb.AddDataSourceRequest{
		ConfigHandle: int64(configHandle),
		InputJson:    inputJson,
	}
	response, err := client.G2ConfigGrpcClient.AddDataSource(ctx, &request)
	result := response.GetResult()
	return result, err
}

/*
The Close method cleans up the Senzing G2Config object pointed to by the handle.
The handle was created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
*/
func (client *G2configClient) Close(ctx context.Context, configHandle uintptr) error {
	request := pb.CloseRequest{
		ConfigHandle: int64(configHandle),
	}
	_, err := client.G2ConfigGrpcClient.Close(ctx, &request)
	return err
}

/*
The Create method creates an in-memory Senzing configuration from the g2config.json
template configuration file located in the PIPELINE.RESOURCEPATH path.
A handle is returned to identify the in-memory configuration.
The handle is used by the AddDataSource(), ListDataSources(), DeleteDataSource(), Load(), and Save() methods.
The handle is terminated by the Close() method.

Input
  - ctx: A context to control lifecycle.

Output
  - A Pointer to an in-memory Senzing configuration.
*/
func (client *G2configClient) Create(ctx context.Context) (uintptr, error) {
	request := pb.CreateRequest{}
	response, err := client.G2ConfigGrpcClient.Create(ctx, &request)
	result := response.GetResult()
	return uintptr(result), err
}

/*
The DeleteDataSource method removes a data source from an existing configuration.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
  - inputJson: A JSON document in the format `{"DSRC_CODE": "NAME_OF_DATASOURCE"}`.
*/
func (client *G2configClient) DeleteDataSource(ctx context.Context, configHandle uintptr, inputJson string) error {
	request := pb.DeleteDataSourceRequest{
		ConfigHandle: int64(configHandle),
		InputJson:    inputJson,
	}
	_, err := client.G2ConfigGrpcClient.DeleteDataSource(ctx, &request)
	return err
}

/*
The Destroy method will destroy and perform cleanup for the Senzing G2Config object.
It should be called after all other calls are complete.

Input
  - ctx: A context to control lifecycle.
*/
func (client *G2configClient) Destroy(ctx context.Context) error {
	request := pb.DestroyRequest{}
	_, err := client.G2ConfigGrpcClient.Destroy(ctx, &request)
	return err
}

/*
The Init method initializes the Senzing G2Config object.
It must be called prior to any other calls.

Input
  - ctx: A context to control lifecycle.
  - moduleName: A name for the auditing node, to help identify it within system logs.
  - iniParams: A JSON string containing configuration parameters.
  - verboseLogging: A flag to enable deeper logging of the G2 processing. 0 for no Senzing logging; 1 for logging.
*/
func (client *G2configClient) Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error {
	request := pb.InitRequest{
		ModuleName:     moduleName,
		IniParams:      iniParams,
		VerboseLogging: int32(verboseLogging),
	}
	_, err := client.G2ConfigGrpcClient.Init(ctx, &request)
	return err
}

/*
The ListDataSources method returns a JSON document of data sources.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.

Output
  - A string containing a JSON document listing all of the data sources.
    See the example output.
*/
func (client *G2configClient) ListDataSources(ctx context.Context, configHandle uintptr) (string, error) {
	request := pb.ListDataSourcesRequest{
		ConfigHandle: int64(configHandle),
	}
	response, err := client.G2ConfigGrpcClient.ListDataSources(ctx, &request)
	result := response.GetResult()
	return result, err
}

/*
The Load method initializes the Senzing G2Config object from a JSON string.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.
  - jsonConfig: A JSON document containing the Senzing configuration.
*/
func (client *G2configClient) Load(ctx context.Context, configHandle uintptr, jsonConfig string) error {
	request := pb.LoadRequest{
		ConfigHandle: int64(configHandle),
		JsonConfig:   jsonConfig,
	}
	_, err := client.G2ConfigGrpcClient.Load(ctx, &request)
	return err
}

/*
The Save method creates a JSON string representation of the Senzing G2Config object.
The configHandle is created by the Create() method.

Input
  - ctx: A context to control lifecycle.
  - configHandle: An identifier of an in-memory configuration.

Output
  - A string containing a JSON Document representation of the Senzing G2Config object.
    See the example output.
*/
func (client *G2configClient) Save(ctx context.Context, configHandle uintptr) (string, error) {
	request := pb.SaveRequest{
		ConfigHandle: int64(configHandle),
	}
	response, err := client.G2ConfigGrpcClient.Save(ctx, &request)
	result := response.GetResult()
	return result, err
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (client *G2configClient) SetLogLevel(ctx context.Context, logLevel logger.Level) error {
	var err error = nil
	return err
}
