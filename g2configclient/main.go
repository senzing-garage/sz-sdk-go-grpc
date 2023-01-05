/*
 *
 */

// Package main implements a client for the service.
package g2configclient

import (
	pb "github.com/senzing/g2-sdk-proto/go/g2config"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2configClient struct {
	G2ConfigGrpcClient pb.G2ConfigClient
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2configclient component found messages having the format "senzing-6021xxxx".
const ProductId = 6021

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2configclient package.
var IdMessages = map[int]string{}

// Status strings for specific g2configclient messages.
var IdStatuses = map[int]string{}
