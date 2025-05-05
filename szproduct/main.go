package szproduct

import "errors"

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

/*
ComponentID is the identifier of the szproduct package.
Package szproduct messages will have the format "SZSDK6026eeee" where "eeee" is the error identifier.
*/
const ComponentID = 6026

var errForPackage = errors.New("szproduct")
