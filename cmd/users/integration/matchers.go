package integration

import "regexp"

var IS_UUID = regexp.MustCompile(`[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}`)
