package szengine

import "errors"

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

/*
ComponentID is the identifier of the szengine package.
Package szengine messages will have the format "SZSDK6024eeee" where "eeee" is the error identifier.
*/
const ComponentID = 6024

var errForPackage = errors.New("szengine")
