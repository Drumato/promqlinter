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
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Drumato/promqlinter/pkg/linter"
	"github.com/Drumato/promqlinter/pkg/linter/plugin"
	"github.com/spf13/cobra"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var (
	GlobalK8sManifestsRO          []string
	GlobalRecursiveRO             bool
	GlobalDiagnosticLevelFilterRO string
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

// NewCLI initializez the root application of promqlinter.
func NewCLI() *cobra.Command {
	c := &cobra.Command{
		Use:     "promqlinter",
		Short:   "A PromQL linter with CLI/GitHub Actions",
		Example: cliExample,
		RunE:    run,
	}

	c.Flags().StringVarP(
		&GlobalDiagnosticLevelFilterRO,
		"level-filter",
		"f",
		"error",
		"the diagnostic level filter(info/warning/error)",
	)

	c.Flags().StringSliceVarP(
		&GlobalK8sManifestsRO,
		"input-k8s-manifest",
		"i",
		[]string{},
		"the target PrometheusRule resource",
	)

	c.Flags().BoolVarP(
		&GlobalRecursiveRO,
		"recursive",
		"r",
		false,
		"determine whether the manifest search process should be recursive",
	)

	return c
}

func run(cmd *cobra.Command, args []string) error {
	filter, err := determineLevelFilter(GlobalDiagnosticLevelFilterRO)
	if err != nil {
		return err
	}

	if len(GlobalK8sManifestsRO) == 0 {
		return runExprFromStdinMode(cmd, args, filter)
	}

	return runK8sManifestsMode(cmd, args, filter)
}

// runExprFromStdinMode runs the linter process with the given input from stdin.
func runExprFromStdinMode(
	cmd *cobra.Command,
	args []string,
	filter linter.DiagnosticLevel,
) error {
	scanner := bufio.NewScanner(os.Stdin)
	l := linter.New(linter.WithPlugins(plugin.Defaults()...))

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	expr := strings.Join(lines, "\n")

	ok, err := l.Execute(expr, filter)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("some of linter plugins detects the filtered rules")
	}

	fmt.Println("ok")
	return nil
}

// runK8sManifestsMode runs the lint process with the k8s manifests.
func runK8sManifestsMode(
	cmd *cobra.Command,
	args []string,
	filter linter.DiagnosticLevel,
) error {
	l := linter.New(
		linter.WithOutStream(os.Stdout),
		linter.WithPlugins(plugin.Defaults()...),
	)

	var manifests []string
	if GlobalRecursiveRO {
		manifests = GlobalK8sManifestsRO
	} else {
		var err error
		if manifests, err = searchAllTargetManifests(GlobalK8sManifestsRO); err != nil {
			return err
		}
	}

	for _, manifestPath := range manifests {
		f, err := os.Open(manifestPath)
		if err != nil {
			return err
		}
		out, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}

		ruleManifest := monitoringv1.PrometheusRule{}
		if err := yaml.Unmarshal(out, &ruleManifest); err != nil {
			return err
		}

		for _, rg := range ruleManifest.Spec.Groups {
			for _, rule := range rg.Rules {
				ok, err := l.Execute(rule.Expr.StrVal, filter)
				if err != nil {
					return err
				}
				if !ok {
					return fmt.Errorf("some of linter plugins detects the filtered rules")
				}
			}
		}
	}

	fmt.Println("ok")
	return nil
}

// searchAlTargetManifests searches the k8s manifests recursively.
func searchAllTargetManifests(
	inputPathsFlagValue []string,
) ([]string, error) {
	manifests := make([]string, 0)
	queue := inputPathsFlagValue

	// Breadth-First-Search
	for len(queue) != 0 {
		dir := queue[0]
		queue = queue[1:]

		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			entryPath := path.Join(dir, entry.Name())
			if entry.IsDir() {
				queue = append(queue, entryPath)
			} else {
				if path.Ext(entry.Name()) != ".yaml" && path.Ext(entry.Name()) != ".yml" {
					continue
				}

				manifests = append(manifests, entryPath)
			}
		}
	}

	return manifests, nil
}

func determineLevelFilter(filter string) (linter.DiagnosticLevel, error) {
	const (
		lInfo    = "info"
		lWarning = "warning"
		lError   = "error"
	)
	if filter != lInfo && filter != lWarning && filter != lError {
		return linter.DiagnosticLevelInfo, fmt.Errorf("--level-filter must be one of info/warning/error")
	}

	switch filter {
	case lInfo:
		return linter.DiagnosticLevelInfo, nil
	case lWarning:
		return linter.DiagnosticLevelWarning, nil
	case lError:
		return linter.DiagnosticLevelError, nil
	default:
		panic("unreachable")
	}
}
