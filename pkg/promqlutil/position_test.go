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

package promqlutil_test

import (
	"testing"

	"github.com/Drumato/promqlinter/pkg/promqlutil"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/stretchr/testify/assert"
)

func TestConvertPosTo2d_Single(t *testing.T) {
	s := "    abcde"
	pos := parser.PositionRange{
		Start: 4,
		End:   8,
	}
	pos2d := promqlutil.ConvertPosTo2d(&s, pos)
	assert.Equal(t, 1, pos2d.Line)
	assert.Equal(t, 5, pos2d.Column)
}

func TestConvertPosTo2d_Multi(t *testing.T) {
	s := "\n    abcde"
	pos := parser.PositionRange{
		Start: 5,
		End:   9,
	}
	pos2d := promqlutil.ConvertPosTo2d(&s, pos)
	assert.Equal(t, 2, pos2d.Line)
	assert.Equal(t, 5, pos2d.Column)
}
