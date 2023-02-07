package g2configclient

import (
	pb "github.com/senzing/g2-sdk-proto/go/g2config"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/subject"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2configClient struct {
	GrpcClient pb.G2ConfigClient
	isTrace    bool
	logger     messagelogger.MessageLoggerInterface
	observers  subject.Subject
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2configclient package found messages having the format "senzing-6021xxxx".
const ProductId = 6021

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2configclient package.
var IdMessages = map[int]string{
	1:    "Enter AddDataSource(%v, %s).",
	2:    "Exit  AddDataSource(%v, %s) returned (%s, %v).",
	3:    "Enter ClearLastException().",
	4:    "Exit  ClearLastException() returned (%v).",
	5:    "Enter Close(%v).",
	6:    "Exit  Close(%v) returned (%v).",
	7:    "Enter Create().",
	8:    "Exit  Create() returned (%v, %v).",
	9:    "Enter DeleteDataSource(%v, %s).",
	10:   "Exit  DeleteDataSource(%v, %s) returned (%v).",
	11:   "Enter Destroy().",
	12:   "Exit  Destroy() returned (%v).",
	13:   "Enter GetLastException().",
	14:   "Exit  GetLastException() returned (%s, %v).",
	15:   "Enter GetLastExceptionCode().",
	16:   "Exit  GetLastExceptionCode() returned (%d, %v).",
	17:   "Enter Init(%s, %s, %d).",
	18:   "Exit  Init(%s, %s, %d) returned (%v).",
	19:   "Enter ListDataSources(%v).",
	20:   "Exit  ListDataSources(%v) returned (%s, %v).",
	21:   "Enter Load(%v, %s).",
	22:   "Exit  Load(%v, %s) returned (%v).",
	23:   "Enter Save(%v).",
	24:   "Exit  Save(%v) returned (%s, %v).",
	25:   "Enter SetLogLevel(%v).",
	26:   "Exit  SetLogLevel(%v) returned (%v).",
	4001: "Call to G2Config_addDataSource(%v, %s) failed. Return code: %d",
	4002: "Call to G2Config_close(%v) failed. Return code: %d",
	4003: "Call to G2Config_create() failed. Return code: %d",
	4004: "Call to G2Config_deleteDataSource(%v, %s) failed. Return code: %d",
	4005: "Call to G2Config_getLastException() failed. Return code: %d",
	4006: "Call to G2Config_destroy() failed. Return code: %d",
	4007: "Call to G2Config_init(%s, %s, %d) failed. Return code: %d",
	4008: "Call to G2Config_listDataSources() failed. Return code: %d",
	4009: "Call to G2Config_load(%v, %s) failed. Return code: %d",
	4010: "Call to G2Config_save(%v) failed. Return code: %d",
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
	8001: "AddDataSource",
	8002: "Close",
	8003: "Create",
	8004: "DeleteDataSource",
	8005: "Destroy",
	8006: "Init",
	8007: "ListDataSources",
	8008: "Load",
	8009: "Save",
}

// Status strings for specific g2configclient messages.
var IdStatuses = map[int]string{}
