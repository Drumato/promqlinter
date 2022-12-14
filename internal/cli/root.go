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

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Drumato/promqlinter/pkg/linter"
	"github.com/Drumato/promqlinter/pkg/linter/plugin"
	"github.com/spf13/cobra"
)

var (
	GlobalK8sManifestRO string
	GlobalRecursiveRO   bool
)

const (
	cliExample = `
	# lint a raw PromQL expression that is given from stdin
	echo -n 'http_requests_total{job="prometheus"}' | promqlinter

	# lint a raw PromQL expression in the PrometheusRule manifest
	promqlinter -i manifest/sample.yaml

	# lint each raw PromQL expression in the PrometheusRule manifests in ./manifest
	promqlinter -r -i ./manifest/
	`
)

func NewCLI() *cobra.Command {
	c := &cobra.Command{
		Use:     "promqlinter",
		Short:   "A PromQL linter with CLI/GitHub Actions",
		Example: cliExample,
		RunE:    run,
	}

	c.Flags().StringVarP(&GlobalK8sManifestRO, "--input-k8s-manifest", "i", "", "the target PrometheusRule resource")
	c.Flags().BoolVarP(&GlobalRecursiveRO, "--recursive", "r", false, "determine whether the manifest search process should be recursive")

	return c
}

func run(cmd *cobra.Command, args []string) error {
	if GlobalK8sManifestRO == "" {
		return runExprFromStdinMode(cmd, args)
	}

	return nil
}

func runExprFromStdinMode(cmd *cobra.Command, args []string) error {
	scanner := bufio.NewScanner(os.Stdin)
	l := linter.New(linter.WithPlugins(plugin.Defaults()...))

	for scanner.Scan() {
		ok, err := l.Execute(scanner.Text(), linter.DiagnosticLevelInfo)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("some of linter plugins detects the filtered rules")
		}

	}

	return nil
}
