# goboxd

goboxd is a Go HTTP service for compiling and running submitted code inside nsjail.

HTTP framework choice: goboxd uses Go's standard `net/http` package. The required API is small (`GET /healthz` and `POST /run`), so the standard library keeps the server easy to audit and avoids an extra dependency.

## Getting started

You need Docker with Compose v2. The host does not need Go, nsjail, or language toolchains.

Build and run:

```sh
make build
make run
```

Check the liveness endpoint:

```sh
curl http://localhost:8080/healthz
```

Common commands:

```sh
make test
make integration
make load
make lint
```

## Docs

- [API](docs/api.md)
- [Languages](docs/languages.md)
- [Security](docs/security.md)
- [Benchmarks](docs/benchmarks.md)
- [Architecture](docs/architecture.md)

## License

GPL-3.0. See [LICENSE](LICENSE).
