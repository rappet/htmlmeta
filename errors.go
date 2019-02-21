package htmlmeta

import "errors"

// ErrorBaseURLNotSet is returned, if BaseURL is nil and ConvertURLs is true
var ErrorBaseURLNotSet = errors.New("BaseURL not set")
