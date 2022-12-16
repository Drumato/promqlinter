# Build Your Own PromQL Linter

This document describes the way to implement your custom PromQL linter.
Maybe the [dummy example](examples/dummy/main.go) is useful to understand the framework briefly.

## The PromQL Linter Framework

this project provides an unified framework called **PromQLinter Framework**.
All of the built-in lint rule follows the `PromQLinterPlugin` interface.
If you want to customize the linter rules, you only should write a plugin that implements the interface.

```go
// pkg/linter/plugin.go

// PromQLinterPlugin is an interface that all linter plugin must implement.
type PromQLinterPlugin interface {
	// Name represents the name of the plugin.
	// the name is used in the reporting message from the linter.
	Name() string
	// Execute lints the PromQL expression.
	Execute(expr parser.Expr) (Diagnostics, error)
}
```

The `PromQLinter` struct has a set of the plugins and use them to lint a PromQL expression.
so you should instantiate the struct and inject your own plugin to the linter.

```go
// yourproject/main.go
package main

import (
	"github.com/Drumato/promqlinter/pkg/linter"
	// "github.com/Drumato/promqlinter/pkg/linter/plugin"
	"github.com/prometheus/prometheus/promql/parser"
)

func main() {
	const sampleExpr = `http_requests_total{job="prometheus"}[5m]`

	l := linter.New(
		// you can pass the default linter plugins into New() 
		// linter.WithPlugins(plugin.Defaults()),
		linter.WithPlugin(&yourPlugin{}),
		linter.WithOutStream(os.Stdout),
	)

	result, err := l.Execute(sampleExpr, linter.DiagnosticLevelWarning)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	if result.Failed() {
		os.Exit(1)
	}
}
```
