# promqlinter

A PromQL parser/type-checker/linter in GitHub Actions/CLI

## Inputs

## `source_directory`

**Required** The directory that the linter recursively searches the k8s manifests.

## Outputs

## Example usage

```yaml
uses: drumato/promqlinter@v0.1.1
with:
  source_directory: manifest/
```
