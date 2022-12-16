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

package promqlutil

import (
	"fmt"

	"github.com/prometheus/prometheus/promql/parser"
)

type Source2dPosition struct {
	Line   int
	Column int
}

// String implements fmt.Stringer
func (s *Source2dPosition) String() string {
	return fmt.Sprintf("(%d:%d)", s.Line, s.Column)
}

func ConvertPosTo2d(
	rawExpr *string,
	source parser.PositionRange,
) *Source2dPosition {
	lineIdx := 0
	offsetFromLineStart := 0
	for i, c := range *rawExpr {
		if c == '\n' {
			lineIdx++
			offsetFromLineStart = 0
		} else {
			offsetFromLineStart++
		}

		if int(source.Start) == i {
			break
		}
	}

	return &Source2dPosition{
		Line:   lineIdx + 1,
		Column: offsetFromLineStart,
	}
}
