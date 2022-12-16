# Using promqlinter in GitHub Actions

## Inputs

### `root_dir`

**Required** The directory that the linter recursively searches the k8s manifests.

### `ansi_colors`

Determine whether the promqlinter's reports are colored with ANSI codes.

### `denied_labels`

the not-allowed label-matchers `<label> %PAIR% <value-pattern-regexp>` separated by comma.

example: `job %PAIR% node_exporter, instance %PAIR% .*`.
this example matches `<vector>{job="node_exporter", instance=".*"}`.

## Outputs

## Example usage

```yaml
uses: drumato/promqlinter@v0.1.3
with:
  root_dir: .
```