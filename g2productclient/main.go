/*
 *
 */

// Package g2productclient implements a client for the service.
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

// Identfier of the g2productclient component found messages having the format "senzing-6026xxxx".
const ProductId = 6026

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2productclient package.
var IdMessages = map[int]string{}

// Status strings for specific g2productclient messages.
var IdStatuses = map[int]string{}
