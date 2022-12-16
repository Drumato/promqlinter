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
	"strings"

	"github.com/Drumato/promqlinter/pkg/promqlutil"
	"github.com/fatih/color"
	"github.com/prometheus/prometheus/promql/parser"
)

// Diagnostics holds the set of the linter diagnostic.
type Diagnostics interface {
	// Slice returns the list of the diagnostic.
	Slice() []Diagnostic
}

// Diagnostic is the detailed message from linter plugin's rules.
type Diagnostic interface {
	// Level returns the diagnostic level.
	Level() DiagnosticLevel
	// Report outputs the lint result to the out stream.
	Report(pluginName string, rawExpr *string, out io.Writer) error
}

// diagnostics is the default implementation of Diagnostics.
type diagnostics struct {
	items []Diagnostic
}

// NewDiagnostics creates a new default Diagnostics.
func NewDiagnostics() *diagnostics {
	return &diagnostics{
		items: make([]Diagnostic, 0),
	}
}

// Add appends the given diagnostic to the set.
func (ds *diagnostics) Add(d Diagnostic) {
	ds.items = append(ds.items, d)
}

// Slice implements Diagnostics
func (d *diagnostics) Slice() []Diagnostic {
	return []Diagnostic(d.items)
}

// diagnostic is the default implementation of Diagnostic.
type diagnostic struct {
	level    DiagnosticLevel
	position parser.PositionRange
	message  string
	color    PromQLinterColorMode
}

// Level implements Diagnostic.
func (d *diagnostic) Level() DiagnosticLevel {
	return d.level
}

// Report implements Diagnostic.
func (d *diagnostic) Report(
	pluginName string,
	rawExpr *string,
	out io.Writer,
) error {
	pos2d := promqlutil.ConvertPosTo2d(rawExpr, d.position)

	if d.color == PromQLinterColorModeEnable {
		return d.coloredReport(pluginName, rawExpr, out, pos2d)
	}

	topMsg := fmt.Sprintf("%s<[%s] %s %s", pluginName, d.level.String(), pos2d, d.message)
	if _, err := fmt.Fprintln(out, topMsg); err != nil {
		return err
	}

	// prefix <- "L1| "
	prefix := fmt.Sprintf("L%d| ", pos2d.Line)

	// line <- "L1| <the contents at the line>"
	line := strings.Split(*rawExpr, "\n")[pos2d.Line-1]
	if _, err := fmt.Fprintf(out, "%s%s\n", prefix, line); err != nil {
		return err
	}

	arrowSpaces := strings.Repeat(" ", len(prefix)+pos2d.Column-1)
	arrow := strings.Repeat("^", int(d.position.End)-int(d.position.Start))
	// arrow <- "<prefix-len><^ * <content-length>>"
	arrow = fmt.Sprintf("%s%s", arrowSpaces, arrow)

	if _, err := fmt.Fprintf(out, "%s %s\n", arrow, d.message); err != nil {
		return err
	}

	return nil
}

func (d *diagnostic) coloredReport(
	pluginName string,
	rawExpr *string,
	out io.Writer,
	pos2d *promqlutil.Source2dPosition,
) error {
	topMsg := fmt.Sprintf("%s<[%s] %s", pluginName, d.level.coloredString(), pos2d)
	if _, err := fmt.Fprintln(out, topMsg); err != nil {
		return err
	}

	// prefix <- "L1| "
	prefix := fmt.Sprintf("L%d| ", pos2d.Line)

	errSubExpr := getSpecifiedSubExpr(rawExpr, &d.position)

	// line <- "L1| <the contents at the line>"
	line := strings.Split(*rawExpr, "\n")[pos2d.Line-1]
	line = strings.ReplaceAll(line, errSubExpr, coloredString(d.level, errSubExpr))
	if _, err := fmt.Fprintf(out, "%s%s\n", prefix, line); err != nil {
		return err
	}

	arrowSpaces := strings.Repeat(" ", len(prefix)+pos2d.Column-1)
	arrow := strings.Repeat("^", int(d.position.End)-int(d.position.Start))
	// arrow <- "<prefix-len><^ * <content-length>>"
	arrow = fmt.Sprintf("%s%s", arrowSpaces, coloredString(d.level, arrow))

	if _, err := fmt.Fprintf(out, "%s %s\n", arrow, coloredString(d.level, d.message)); err != nil {
		return err
	}

	return nil
}

func getSpecifiedSubExpr(
	rawExpr *string,
	source *parser.PositionRange,
) string {
	return (*rawExpr)[int(source.Start):int(source.End)]
}

func InfoDiagnostic(
	position parser.PositionRange,
	message string,
	color PromQLinterColorMode,
) *diagnostic {
	return &diagnostic{
		level:    DiagnosticLevelInfo,
		position: position,
		message:  message,
		color:    color,
	}
}

func WarningDiagnostic(
	position parser.PositionRange,
	message string,
	color PromQLinterColorMode,
) *diagnostic {
	return &diagnostic{
		level:    DiagnosticLevelWarning,
		position: position,
		message:  message,
		color:    color,
	}
}

func ErrorDiagnostic(
	position parser.PositionRange,
	message string,
	color PromQLinterColorMode,
) *diagnostic {
	return &diagnostic{
		level:    DiagnosticLevelError,
		position: position,
		message:  message,
		color:    color,
	}
}

func ColoredInfoDiagnostic(
	position parser.PositionRange,
	message string,
) *diagnostic {
	return InfoDiagnostic(position, message, PromQLinterColorModeEnable)
}

func ColoredWarningDiagnostic(
	position parser.PositionRange,
	message string,
) *diagnostic {
	return WarningDiagnostic(position, message, PromQLinterColorModeEnable)
}

func ColoredErrorDiagnostic(
	position parser.PositionRange,
	message string,
) *diagnostic {
	return ErrorDiagnostic(position, message, PromQLinterColorModeEnable)
}

func NoncoloredInfoDiagnostic(
	position parser.PositionRange,
	message string,
) *diagnostic {
	return InfoDiagnostic(position, message, PromQLinterColorModeDisable)
}

func NoncoloredWarningDiagnostic(
	position parser.PositionRange,
	message string,
) *diagnostic {
	return WarningDiagnostic(position, message, PromQLinterColorModeDisable)
}

func NoncoloredErrorDiagnostic(
	position parser.PositionRange,
	message string,
) *diagnostic {
	return ErrorDiagnostic(position, message, PromQLinterColorModeDisable)
}

// DiagnosticLevel represents the level of a diagnostic.
type DiagnosticLevel uint

const (
	// DiagnosticLevelInfo represents the "information" level.
	DiagnosticLevelInfo DiagnosticLevel = iota
	// DiagnosticLevelWarning represents the "warning" level.
	DiagnosticLevelWarning
	// DiagnosticLevelError represents the "error" level.
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

func (d DiagnosticLevel) coloredString() string {
	switch d {
	case DiagnosticLevelInfo:
		return color.HiBlueString("INFO")
	case DiagnosticLevelWarning:
		return color.HiYellowString("WARN")
	case DiagnosticLevelError:
		return color.HiRedString("ERROR")
	default:
		// unreachable
		return ""
	}
}

func convertParseErrorToDiagnostics(
	err error,
	color PromQLinterColorMode,
) Diagnostics {
	if err == nil {
		return nil
	}

	errs, ok := err.(parser.ParseErrors)
	if !ok {
		return nil
	}

	ds := NewDiagnostics()
	for _, e := range errs {
		d := ErrorDiagnostic(
			e.PositionRange,
			e.Error(),
			color,
		)

		ds.Add(d)
	}

	return ds
}

func coloredString(level DiagnosticLevel, s string) string {
	switch level {
	case DiagnosticLevelInfo:
		return color.HiBlueString(s)
	case DiagnosticLevelWarning:
		return color.HiYellowString(s)
	case DiagnosticLevelError:
		return color.HiRedString(s)
	default:
		// unreachable
		return ""
	}

}
