/*
MIT License

# Copyright (c) 2022 Drumato

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"fmt"
	"os"

	"github.com/Drumato/promqlinter/pkg/linter"
	"github.com/prometheus/prometheus/promql/parser"
)

type samplePlugin struct{}

// Execute implements linter.PromQLinterPlugin
func (*samplePlugin) Execute(expr parser.Expr) (linter.Diagnostics, error) {
	ds := linter.NewDiagnostics()
	ds.Add(linter.ColoredInfoDiagnostic(
		parser.PositionRange{},
		"foo",
	))
	ds.Add(linter.ColoredInfoDiagnostic(
		parser.PositionRange{},
		"bar",
	))
	ds.Add(linter.ColoredInfoDiagnostic(
		parser.PositionRange{},
		"baz",
	))

	return ds, nil
}

// Name implements linter.PromQLinterPlugin
func (*samplePlugin) Name() string {
	return "sample-plugin"
}

var _ linter.PromQLinterPlugin = &samplePlugin{}

func main() {
	l := linter.New(
		linter.WithPlugin(&samplePlugin{}),
		linter.WithOutStream(os.Stdout),
	)
	result, err := l.Execute("http_requests_total", linter.DiagnosticLevelWarning)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	if result.Failed() {
		os.Exit(1)
	}
}
