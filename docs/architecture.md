# Architecture

The binary entry point is [cmd/goboxd/main.go](../cmd/goboxd/main.go). It builds an HTTP server on port `8080` and uses `server.NewMux()` as the router.

The `server` package owns route registration. Handlers live under `server/handler`.

Language configuration is loaded by `internal/languages`. The default registry is embedded from `config/languages.yaml` so the binary has a known config at startup.

Run request validation lives in `internal/runs`. It checks request shape against the language registry before execution code touches the filesystem or starts nsjail.

The intended execution flow is:

```text
HTTP request
  -> route handler
  -> JSON decode
  -> request validation
  -> per-job sandbox directory
  -> optional build command through nsjail
  -> test run commands through nsjail
  -> response status mapping
```
