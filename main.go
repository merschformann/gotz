package main

import (
	"fmt"
	"strings"

	"github.com/merschformann/gotz/cmd"
)

func main() {
	info := GetReleaseInfo()
	cmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", strings.TrimPrefix(info.Version, "v"), info.Commit, info.Date)
	cmd.Execute()
}
