/*
 *
 */

// Package g2configmgrclient implements a client for the service.
package g2configmgrclient

import (
	pb "github.com/senzing/g2-sdk-proto/go/g2configmgr"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2configmgrClient struct {
	GrpcClient pb.G2ConfigMgrClient
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2configmgrclient component found messages having the format "senzing-6022xxxx".
const ProductId = 6022

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2configmgrclient package.
var IdMessages = map[int]string{}

// Status strings for specific g2configmgrclient messages.
var IdStatuses = map[int]string{}
