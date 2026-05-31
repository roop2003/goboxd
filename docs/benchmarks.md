# Benchmarks

Use `make load` for the local load smoke test. It exercises the HTTP server through an in-process test server and is meant to catch obvious regressions in request handling.

Before submission, record container-based numbers here:

- Docker image build time from a clean cache.
- Startup time until `/healthz` returns `200`.
- Median and p95 latency for `/healthz`.
- Median and p95 latency for one passing submission per supported language.
- Peak memory while running the language matrix.

The benchmark run should use Docker Compose so it measures the same path as local development.
