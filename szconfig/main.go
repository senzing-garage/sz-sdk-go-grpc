package szconfig

import "errors"

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

/*
ComponentID is the identifier of the szconfig package.
Package szconfig messages will have the format "SZSDK6021eeee" where "eeee" is the error identifier.
*/
const ComponentID = 6021

var errForPackage = errors.New("szconfig")
