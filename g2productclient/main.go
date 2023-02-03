package g2productclient

import (
	pb "github.com/senzing/g2-sdk-proto/go/g2product"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2productClient struct {
	GrpcClient pb.G2ProductClient
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2productclient package found messages having the format "senzing-6026xxxx".
const ProductId = 6026

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2productclient package.
var IdMessages = map[int]string{
	1:    "Enter ClearLastException().",
	2:    "Exit  ClearLastException() returned (%v).",
	3:    "Enter Destroy().",
	4:    "Exit  Destroy() returned (%v).",
	5:    "Enter GetLastException().",
	6:    "Exit  GetLastException() returned (%s, %v).",
	7:    "Enter GetLastExceptionCode().",
	8:    "Exit  GetLastExceptionCode() returned (%d, %v).",
	9:    "Enter Init(%s, %s, %s).",
	10:   "Exit  Init(%s, %s, %s) returned (%v).",
	11:   "Enter License().",
	12:   "Exit  License() returned (%s, %v).",
	13:   "Enter SetLogLevel(%v).",
	14:   "Exit  SetLogLevel(%v) returned (%v).",
	15:   "Enter ValidateLicenseFile(%s).",
	16:   "Exit  ValidateLicenseFile(%s) returned (%s, %v).",
	17:   "Enter ValidateLicenseStringBase64(%s).",
	18:   "Exit  ValidateLicenseStringBase64(%s) returned (%s, %v).",
	19:   "Enter Version().",
	20:   "Exit  Version() returned (%s, %v).",
	4001: "Call to G2Product_destroy() failed. Return code: %d",
	4002: "Call to G2Product_getLastException() failed. Return code: %d",
	4003: "Call to G2Product_init(%s, %s, %s) failed. Return code: %d",
	4004: "Call to G2Product_validateLicenseFile(%s) failed. Return code: %d",
	4005: "Call to G2Product_validateLicenseStringBase64(%s) failed. Return code: %d",
	5901: "During setup, call to messagelogger.NewSenzingApiLogger() failed.",
	5902: "During setup, call to g2eg2engineconfigurationjson.BuildSimpleSystemConfigurationJson() failed.",
	5903: "During setup, call to g2engine.Init() failed.",
	5904: "During setup, call to g2engine.PurgeRepository() failed.",
	5905: "During setup, call to g2engine.Destroy() failed.",
	5906: "During setup, call to g2config.Init() failed.",
	5907: "During setup, call to g2config.Create() failed.",
	5908: "During setup, call to g2config.AddDataSource() failed.",
	5909: "During setup, call to g2config.Save() failed.",
	5910: "During setup, call to g2config.Close() failed.",
	5911: "During setup, call to g2config.Destroy() failed.",
	5912: "During setup, call to g2configmgr.Init() failed.",
	5913: "During setup, call to g2configmgr.AddConfig() failed.",
	5914: "During setup, call to g2configmgr.SetDefaultConfigID() failed.",
	5915: "During setup, call to g2configmgr.Destroy() failed.",
	5916: "During setup, call to g2engine.Init() failed.",
	5917: "During setup, call to g2engine.AddRecord() failed.",
	5918: "During setup, call to g2engine.Destroy() failed.",
	5931: "During setup, call to g2engine.Init() failed.",
	5932: "During setup, call to g2engine.PurgeRepository() failed.",
	5933: "During setup, call to g2engine.Destroy() failed.",
	8001: "Destroy",
	8002: "Init",
	8003: "License",
	8004: "ValidateLicenseFile",
	8005: "ValidateLicenseStringBase64",
	8006: "Version",
}

// Status strings for specific g2productclient messages.
var IdStatuses = map[int]string{}
