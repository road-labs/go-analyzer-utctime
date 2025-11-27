# go-analyzer-utctime

A custom linter for [golangci-lint](https://golangci-lint.run/) that ensures all `time.Now()` calls are followed by `.UTC()`.

## Description

This linter helps prevent timezone-related bugs by ensuring that all `time.Now()` calls are immediately followed by `.UTC()`. This is particularly useful in applications where consistent timezone handling is critical.

## Installation

This linter uses the golangci-lint [Module Plugin System](https://golangci-lint.run/docs/plugins/module-plugins/).

1. Create a `.custom-gcl.yml` file in your project root:

```yaml
version: v2.6.2  # Use your desired golangci-lint version

plugins:
  - module: 'github.com/road-labs/go-analyzer-utctime'
    version: latest  # or specify a version tag
```

2. Build a custom golangci-lint binary with the plugin:

```bash
$ golangci-lint custom
```

This will create a `./custom-gcl` binary that includes the plugin.

3. Configure the plugin in your `.golangci.yml` (see Usage section below).

4. Run your custom golangci-lint binary:

```bash
$ ./custom-gcl run
```

## Usage

Add the linter to your `.golangci.yml` configuration:

```yaml
version: "2"

linters:
  enable:
    - utctime
  settings:
    custom:
      utctime:
        type: module
        description: Checks that time.Now() is followed by .UTC()
```

## Examples

```go
// Bad:
t := time.Now() // Will trigger a linter error

// Good:
t := time.Now().UTC()
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
