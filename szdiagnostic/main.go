package szdiagnostic

import "errors"

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

/*
ComponentID is the identifier of the szdiagnostic package.
Package szdiagnostic messages will have the format "SZSDK6023eeee" where "eeee" is the error identifier.
*/
const ComponentID = 6023

var errForPackage = errors.New("szdiagnostic")
