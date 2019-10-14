package pingdom

import "errors"

// ErrMissingID is an error for when a required Id field is missing.
var ErrMissingID = errors.New("required field 'Id' missing")

// ErrBadResolution is an error for when an invalid resolution is specified.
var ErrBadResolution = errors.New("resolution must be either 'hour', 'day' or 'week'")
