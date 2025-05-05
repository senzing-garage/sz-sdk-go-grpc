package helper

import "errors"

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

/*
MessageIDPrefix is the message prefix for `SZSDKcccceeee` message identifers
where "cccc" is the component ID and "eeee" is the error identifier.
*/
const (
	MessageIDPrefix = "SZSDK"
)

var errPackage = errors.New("helper")
