package utils

import "regexp"

var REGEX_EMAIL, _ = regexp.Compile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
