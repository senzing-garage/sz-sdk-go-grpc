package g2configmgrclient

import (
	pb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2configmgrClient struct {
	GrpcClient pb.G2ConfigMgrClient
	isTrace    bool
	logger     messagelogger.MessageLoggerInterface
	observers  subject.Subject
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2configmgrclient package found messages having the format "senzing-6022xxxx".
const ProductId = 6022

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2configmgrclient package.
var IdMessages = map[int]string{
	1:    "Enter AddConfig(%s, %s).",
	2:    "Exit  AddConfig(%s, %s) returned (%d, %v).",
	3:    "Enter ClearLastException().",
	4:    "Exit  ClearLastException() returned (%v).",
	5:    "Enter Destroy().",
	6:    "Exit  Destroy() returned (%v).",
	7:    "Enter GetConfig(%d).",
	8:    "Exit  GetConfig(%d) returned (%s, %v).",
	9:    "Enter GetConfigList().",
	10:   "Exit  GetConfigList() returned (%s, %v).",
	11:   "Enter GetDefaultConfigID().",
	12:   "Exit  GetDefaultConfigID() returned (%d, %v).",
	13:   "Enter GetLastException().",
	14:   "Exit  GetLastException() returned (%s, %v).",
	15:   "Enter GetLastExceptionCode().",
	16:   "Exit  GetLastExceptionCode() returned (%d, %v).",
	17:   "Enter Init(%s, %s, %d).",
	18:   "Exit  Init(%s, %s, %d) returned (%v).",
	19:   "Enter ReplaceDefaultConfigID(%d, %d).",
	20:   "Exit  ReplaceDefaultConfigID(%d, %d) returned (%v).",
	21:   "Enter SetDefaultConfigID(%d).",
	22:   "Exit  SetDefaultConfigID(%d) returned (%v).",
	23:   "Enter SetLogLevel(%v).",
	24:   "Exit  SetLogLevel(%v) returned (%v).",
	4001: "Call to G2ConfigMgr_addConfig(%s, %s) failed. Return code: %d",
	4002: "Call to G2ConfigMgr_destroy() failed. Return code: %d",
	4003: "Call to G2ConfigMgr_getConfig(%d) failed. Return code: %d",
	4004: "Call to G2ConfigMgr_getConfigList() failed. Return code: %d",
	4005: "Call to G2ConfigMgr_getDefaultConfigID() failed. Return code: %d",
	4006: "Call to G2ConfigMgr_getLastException() failed. Return code: %d",
	4007: "Call to G2ConfigMgr_init(%s, %s, %d) failed. Return code: %d",
	4008: "Call to G2ConfigMgr_replaceDefaultConfigID(%d, %d) failed. Return code: %d",
	4009: "Call to G2ConfigMgr_setDefaultConfigID(%d) failed. Return code: %d",
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
	5920: "During setup, call to setupSenzingConfig() failed.",
	5921: "During setup, call to setupPurgeRepository() failed.",
	5922: "During setup, call to setupAddRecords() failed.",
	5931: "During setup, call to g2engine.Init() failed.",
	5932: "During setup, call to g2engine.PurgeRepository() failed.",
	5933: "During setup, call to g2engine.Destroy() failed.",
	8001: "AddConfig",
	8002: "Destroy",
	8003: "GetConfig",
	8004: "GetConfigList",
	8005: "GetDefaultConfigID",
	8006: "Init",
	8007: "ReplaceDefaultConfigID",
	8008: "SetDefaultConfigID",
}

// Status strings for specific g2configmgrclient messages.
var IdStatuses = map[int]string{}
