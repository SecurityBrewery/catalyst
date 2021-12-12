package catalyst

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var VERSION string

func GetVersion() string {
	return strings.TrimSpace(VERSION)
}
