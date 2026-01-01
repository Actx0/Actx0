// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package main

import (
	"embed"
	"io/fs"

	"github.com/clivern/actx0/cli"
	"github.com/clivern/actx0/locale"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

//go:embed web/dist/*
var static embed.FS

//go:embed locale/locales
var localeFS embed.FS

func main() {
	cli.Version = version
	cli.Commit = commit
	cli.Date = date
	cli.BuiltBy = builtBy
	cli.Static = static

	if sub, err := fs.Sub(localeFS, "locale/locales"); err == nil {
		if err := locale.Load(sub); err != nil {
			panic(err)
		}
	}

	cli.Execute()
}
