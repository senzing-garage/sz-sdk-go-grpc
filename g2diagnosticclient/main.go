/*
 *
 */

// Package main implements a client for the service.
package g2diagnosticclient

import (
	"github.com/senzing/g2-sdk-go/g2diagnostic"
	pb "github.com/senzing/g2-sdk-proto/go/g2diagnostic"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type G2diagnosticClient struct {
	G2DiagnosticGrpcClient pb.G2DiagnosticClient
	G2DiagnosticInterface  g2diagnostic.G2diagnostic
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2diagnosticclient component found messages having the format "senzing-6023xxxx".
const ProductId = 6023

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2diagnostic package.
var IdMessages = map[int]string{}

// Status strings for specific g2diagnostic messages.
var IdStatuses = map[int]string{}
