package build

import (
	"strings"
)

const (
	// VersionNumber 版本号
	VersionNumber = "1.0.0"
)

// Copyright 版权信息
const Copyright = "Copyright (c) 2021-2040, XXXX. All rights reserved."

// The value of variables come form `gb build -ldflags '-X "build.Date=xxxxx" -X "build.CommitID=xxxx"' `
var (
	// Date build time
	Date string
	// Branch current git branch
	Branch string
	// Commit git commit id
	Commit string
	// Plugins business plugin list
	Plugins string
)

// Version 生成版本信息
func Version() string {
	var buf strings.Builder
	buf.WriteString(VersionNumber)

	if Date != "" {
		buf.WriteByte('\n')
		buf.WriteString("date: ")
		buf.WriteString(Date)
	}
	if Branch != "" {
		buf.WriteByte('\n')
		buf.WriteString("branch: ")
		buf.WriteString(Branch)
	}
	if Commit != "" {
		buf.WriteByte('\n')
		buf.WriteString("commit: ")
		buf.WriteString(Commit)
	}
	if Plugins != "" {
		buf.WriteByte('\n')
		buf.WriteString("plugins: ")
		buf.WriteString(Plugins)
	}
	return buf.String()
}
