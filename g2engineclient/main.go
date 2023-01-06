/*
 *
 */

// Package g2engineclient implements a client for the service.
package g2engineclient

import (
	pb "github.com/senzing/g2-sdk-proto/go/g2engine"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2engineClient struct {
	GrpcClient pb.G2EngineClient
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2engineclient component found messages having the format "senzing-6024xxxx".
const ProductId = 6024

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2engineclient package.
var IdMessages = map[int]string{}

// Status strings for specific g2engineclient messages.
var IdStatuses = map[int]string{}
