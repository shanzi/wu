package main

import (
	"fmt"
	"strings"
)

var Version = "Run make to build with version"

func PrintVersion() {
	fmt.Printf("wu Version: %s\n", strings.TrimPrefix(Version, "v"))
}
