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
	"fmt"
	"io"

	"github.com/prometheus/prometheus/promql/parser"
)

type Diagnostics struct {
	items []Diagnostic
}

func NewDiagnostics() *Diagnostics {
	return &Diagnostics{
		items: make([]Diagnostic, 0),
	}
}

func (ds *Diagnostics) Add(d Diagnostic) {
	ds.items = append(ds.items, d)
}

type Diagnostic struct {
	level    DiagnosticLevel
	position parser.PositionRange
	message  string
}

func (d *Diagnostic) Report(out io.Writer) error {
	_, err := fmt.Fprintf(out, "[%s] (%d..%d) %s\n", d.level.String(), d.position.Start, d.position.End, d.message)
	return err
}

func NewDiagnostic(
	level DiagnosticLevel,
	position parser.PositionRange,
	message string,
) Diagnostic {
	return Diagnostic{level, position, message}
}

type DiagnosticLevel uint

const (
	DiagnosticLevelInfo DiagnosticLevel = iota
	DiagnosticLevelWarning
	DiagnosticLevelError
)

// String implements fmt.Stringer
func (d DiagnosticLevel) String() string {
	switch d {
	case DiagnosticLevelInfo:
		return "INFO"
	case DiagnosticLevelWarning:
		return "WARN"
	case DiagnosticLevelError:
		return "ERROR"
	default:
		// unreachable
		return ""
	}
}