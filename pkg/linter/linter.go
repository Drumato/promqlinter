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

type PromQLinter struct {
	out     io.Writer
	plugins []PromQLinterPlugin
}

type PromQLinterOption func(*PromQLinter)

func New(options ...PromQLinterOption) *PromQLinter {
	pq := &PromQLinter{}
	for _, opt := range options {
		opt(pq)
	}

	return pq
}

func (pq *PromQLinter) Execute(
	rawExpr string,
	filter DiagnosticLevel,
) (bool, error) {
	ok := true

	expr, err := parser.ParseExpr(rawExpr)
	if err != nil {
		return false, err
	}

	for _, p := range pq.plugins {
		ds, err := p.Execute(expr)
		if err != nil {
			return false, err
		}

		for _, d := range ds.items {
			if d.level >= filter {
				if err := d.Report(pq.out); err != nil {
					return false, err
				}
				ok = false
			}
		}
	}

	return ok, nil
}

func WithPlugins(plugins ...PromQLinterPlugin) PromQLinterOption {
	return func(pq *PromQLinter) {
		pq.plugins = plugins
	}
}

func WithPlugin(plugin PromQLinterPlugin) PromQLinterOption {
	return func(pq *PromQLinter) {
		pq.plugins = append(pq.plugins, plugin)
	}
}

func WithOutStream(out io.Writer) PromQLinterOption {
	return func(pq *PromQLinter) {
		pq.out = out
	}
}
