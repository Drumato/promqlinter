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

package cli

import "github.com/spf13/cobra"

var (
	GlobalK8sManifestRO           string
	GlobalRecursiveRO             bool
	GlobalDiagnosticLevelFilterRO string
	GlobalDeniedLabelsRO          string
)

func defineCLIFlags(c *cobra.Command) {
	c.Flags().StringVarP(
		&GlobalDeniedLabelsRO,
		"denied-labels",
		"d",
		"",
		"the denied labels",
	)

	c.Flags().StringVarP(
		&GlobalDiagnosticLevelFilterRO,
		"level-filter",
		"f",
		"error",
		"the diagnostic level filter(info/warning/error)",
	)

	c.Flags().StringVarP(
		&GlobalK8sManifestRO,
		"input-k8s-manifest",
		"i",
		"",
		"the target PrometheusRule resource",
	)

	c.Flags().BoolVarP(
		&GlobalRecursiveRO,
		"recursive",
		"r",
		false,
		"determine whether the manifest search process should be recursive",
	)

}
