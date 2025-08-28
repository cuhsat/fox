// The Swiss Army Knife for examining text files.
//
// Copyright 2025 Christian Uhsat. All rights reserved.
// Use of this source code is governed by the GPL-3.0
// license that can be found in the LICENSE.md file.
//
// For more information, please consult:
//
//	https://forensic-examiner.eu
package main

import (
	"runtime/debug"

	"github.com/cuhsat/fox/internal/cmd"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

// Main start and catch.
func main() {
	defer func() {
		if err := recover(); err != nil {
			sys.Trace(err, debug.Stack())
		}
	}()

	sys.Setup()

	_ = cmd.Execute()
}
