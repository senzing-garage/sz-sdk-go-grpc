/*
 *
 */

// Package g2diagnosticclient implements a client for the service.
package g2diagnosticclient

import (
	pb "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2diagnosticClient struct {
	GrpcClient pb.G2DiagnosticClient
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2diagnosticclient component found messages having the format "senzing-6023xxxx".
const ProductId = 6023

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2diagnosticclient package.
var IdMessages = map[int]string{}

// Status strings for specific g2diagnosticclient messages.
var IdStatuses = map[int]string{}
