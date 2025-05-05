package szconfigmanager

import "errors"

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

/*
ComponentID is the identifier of the szconfigmanager package.
Package szconfigmanager messages will have the format "SZSDK6022eeee" where "eeee" is the error identifier.
*/
const ComponentID = 6022

var errForPackage = errors.New("szconfigmanager")
