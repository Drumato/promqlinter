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

package linter

import (
	"io"

	"github.com/prometheus/prometheus/promql/parser"
)

// PromQLinter is the actual PromQL linter.
// users can configure this struct for controlling the behaviors of the linter.
type PromQLinter struct {
	outStream io.Writer

	plugins []PromQLinterPlugin
	colored bool
}

// PromQLinterOption enables the initialization of the PromQLinter by FOP(Functional-Options-Pattern)
type PromQLinterOption func(*PromQLinter)

// New creates a new PromQLinter.
func New(options ...PromQLinterOption) *PromQLinter {
	pq := &PromQLinter{
		plugins: make([]PromQLinterPlugin, 0),
	}
	for _, opt := range options {
		opt(pq)
	}

	return pq
}

// Execute starts the lint process.
// the rawExpr parameter is a PromQL expression.
// filter determines whether the reported diagnostics from plugin(s) are ignored.
func (pq *PromQLinter) Execute(
	rawExpr string,
	filter DiagnosticLevel,
) (bool, error) {
	ok := true
	expr, err := parser.ParseExpr(rawExpr)
	parserDs := convertParseErrorToDiagnostics(err, pq.colored)
	if parserDs != nil {
		for _, d := range parserDs.Slice() {
			if d.Level() >= filter {
				if err := d.Report("promql/parser", &rawExpr, pq.outStream); err != nil {
					return false, err
				}
				ok = false
			}
		}
	}
	// if any parse errors are found, we quickly quit the lint process.
	if !ok {
		return false, nil
	}

	for _, p := range pq.plugins {
		ds, err := p.Execute(expr)
		if err != nil {
			return false, err
		}

		for _, d := range ds.Slice() {
			if d.Level() >= filter {
				if err := d.Report(p.Name(), &rawExpr, pq.outStream); err != nil {
					return false, err
				}
				ok = false
			}
		}
	}

	return ok, nil
}

// WithPlugins sets the set of the linter plugin to the linter.
// Note that this function should be called before WithPlugin().
// Because this function updates the plugin set entirely.
func WithPlugins(plugins ...PromQLinterPlugin) PromQLinterOption {
	return func(pq *PromQLinter) {
		pq.plugins = plugins
	}
}

// WithPlugins appends the given plugin to the linter plugins.
func WithPlugin(plugin PromQLinterPlugin) PromQLinterOption {
	return func(pq *PromQLinter) {
		pq.plugins = append(pq.plugins, plugin)
	}
}

// WithOutStream sets the output stream to the linter.
func WithOutStream(out io.Writer) PromQLinterOption {
	return func(pq *PromQLinter) {
		pq.outStream = out
	}
}

// WithANSIColored sets the colored flag to the linter.
func WithANSIColored(colored bool) PromQLinterOption {
	return func(pq *PromQLinter) {
		pq.colored = colored
	}
}
