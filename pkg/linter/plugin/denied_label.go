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

package plugin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Drumato/promqlinter/pkg/linter"
	"github.com/prometheus/prometheus/promql/parser"
)

const (
	baseMetricNamePlaceholder = "__name__"
)

type LabelName string
type LabelValuePattern string

type deniedLabel struct {
	labels map[LabelName]LabelValuePattern
}

// Execute implements linter.PromQLinterPlugin
func (d *deniedLabel) Execute(expr parser.Expr) (linter.Diagnostics, error) {
	ds := linter.NewDiagnostics()
	parser.Inspect(expr, func(n parser.Node, path []parser.Node) error {
		switch node := n.(type) {
		case *parser.VectorSelector:
			for _, lm := range node.LabelMatchers {
				if lm.Name == baseMetricNamePlaceholder {
					continue
				}

				pattern, ok := d.labels[LabelName(lm.Name)]
				if !ok {
					continue
				}

				exp, err := regexp.Compile(string(pattern))
				if err != nil {
					return err
				}

				if exp.MatchString(lm.Value) {
					msg := fmt.Sprintf("matched to the denied label rule `%s`", pattern)
					ds.Add(linter.NewDiagnostic(
						linter.DiagnosticLevelError,
						node.PosRange,
						msg,
					))
				}
			}

			return nil
		default:
			// traverse all the non-nil children.
			return nil
		}
	})

	return ds, nil
}

// Name implements linter.PromQLinterPlugin
func (*deniedLabel) Name() string {
	return "denied-labels"
}

func NewDeniedLabelPlugin(
	deniedLabels string,
) linter.PromQLinterPlugin {
	labels := splitDeniedLabelsFlag(deniedLabels)
	return &deniedLabel{labels}
}

func splitDeniedLabelsFlag(value string) map[LabelName]LabelValuePattern {
	matchers := map[LabelName]LabelValuePattern{}

	if value == "" {
		return matchers
	}

	labels := strings.Split(value, ",")
	for _, label := range labels {
		pair := strings.Split(label, " %PAIR% ")
		name := LabelName(pair[0])
		matchers[name] = LabelValuePattern(pair[1])
	}

	return matchers
}