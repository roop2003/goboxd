# AI Prompt Log

This file follows the format requested in the goboxd hackathon discussion: log non-trivial AI use, summarize the response, and say what was accepted or rejected.

## 2026-05-31 · Fixing initial Go compile errors

**Prompt:**
`./main.go:17:2: undefined: registerRoutes` and `./main.go:29:4: undefined: log` facing these errors.

**Response summary:**
The AI inspected the Go files and identified two problems: `log` was missing from imports, and `registerRoutes` was not visible when running only one file instead of the package. It suggested moving toward the repo's intended `cmd/goboxd` entrypoint and a `server` package.

**What we used / didn't use:**
Used the package/entrypoint split and `log` import fix. Did not keep the duplicate server startup path that was present in the early `server/main.go`.

## 2026-05-31 · Understanding `go run` package behavior

**Prompt:**
`./main.go:13:2: undefined: registerRoutes`

**Response summary:**
The AI explained that `go run main.go` compiles only that file, so Go does not include `register.go`. It recommended running a package path and aligning the code with `go run ./cmd/goboxd`.

**What we used / didn't use:**
Used `cmd/goboxd/main.go` as the binary entrypoint and changed routing code into an importable `server` package. Did not continue with a single-file `go run main.go` workflow.

## 2026-05-31 · Running `/healthz`

**Prompt:**
How will it run `/healthz`?

**Response summary:**
The AI explained the flow: start the server, register `GET /healthz`, and test with `curl http://localhost:8080/healthz`.

**What we used / didn't use:**
Used the explanation to confirm the health route and kept `/healthz` as a simple liveness endpoint returning `OK`.

## 2026-05-31 · Fixing missing local package import

**Prompt:**
`go run ./cmd/goboxd` fails with `cmd/goboxd/main.go:11:12: undefined: server`.

**Response summary:**
The AI found that `cmd/goboxd/main.go` referenced `server.NewMux()` without importing the local module package.

**What we used / didn't use:**
Used the import `github.com/thesouldev/goboxd/server`. This fixed the missing symbol while keeping the server construction in the `server` package.

## 2026-05-31 · Freeing port 8080

**Prompt:**
Server failed with `bind: address already in use`; kill 8080.

**Response summary:**
The AI used `lsof` to identify the process listening on port 8080 and killed it. A normal kill did not stop it, so a force kill was used after checking the same PID was still listening.

**What we used / didn't use:**
Used the port cleanup commands. This was an operational fix only; no code was changed for this prompt.

## 2026-05-31 · Understanding the local module import

**Prompt:**
What is `import "github.com/thesouldev/goboxd/server"` doing?

**Response summary:**
The AI explained that the module path in `go.mod` makes `github.com/thesouldev/goboxd/server` refer to the local `server` directory.

**What we used / didn't use:**
Used this understanding to keep package boundaries clear: `cmd/goboxd` starts the binary, while `server` owns the HTTP router.

## 2026-05-31 · Understanding execution flow

**Prompt:**
How is the code execution flow working?

**Response summary:**
The AI described the current request flow: `main()` starts the HTTP server, `server.NewMux()` registers routes, `/healthz` maps to `HealthzHandler`, and the handler writes `OK`.

**What we used / didn't use:**
Used the mental model to guide later route additions. At that point, code execution for submissions was not implemented yet.

## 2026-05-31 · Asking whether Go can execute C code

**Prompt:**
Can Go execute C lang code?

**Response summary:**
The AI explained two approaches: calling C through cgo, or compiling/running C as a separate process with tools like `gcc`. It noted that a code execution service should compile and run source in a sandbox rather than interpret C directly.

**What we used / didn't use:**
Used the compile-and-run model as the direction for `/run`. Did not use cgo because user-submitted code should run as a separate process.

## 2026-05-31 · Defining request and response validation rules

**Prompt:**
Provided field rules for `language`, `source`, filenames, build/run flags and limits, required tests, `200` responses for user-code outcomes, and `400` responses for bad input.

**Response summary:**
The AI summarized the API contract: validation errors should return `400`; user-code failures should return `200` with a status in JSON; server/sandbox failures should be `5xx`.

**What we used / didn't use:**
Used this contract to shape `internal/runs` request types, validation errors, status values, and sample response files.

## 2026-05-31 · Writing a validator

**Prompt:**
How do I write validator for this?

**Response summary:**
The AI proposed request structs, an API error type, filename validation, source-size checks, test validation, and flag allow-list checks.

**What we used / didn't use:**
Used the validator structure in `internal/runs/validator.go`. Revised the language lookup later so it comes from YAML config instead of a hardcoded map.

## 2026-05-31 · Designing plug-and-play language config

**Prompt:**
Provided a YAML shape for plug-and-play languages with `py3`, `cpp`, and `java`, including source filename strategy, build/run commands, limits, and flag allow-lists.

**Response summary:**
The AI implemented a YAML-driven language registry and validator integration using `gopkg.in/yaml.v3`.

**What we used / didn't use:**
Used the YAML-based registry and flag allow-list pattern. Later changed from duplicate YAML files to one canonical `config/languages.yaml`.

## 2026-05-31 · Checking against the hackathon submission requirements

**Prompt:**
Asked whether the repo follows requirements for Dockerfile, Docker Compose, Makefile targets, docs layout, README style, tests, linting, and Docker-based execution.

**Response summary:**
The AI audited the repo and patched several gaps: added required Make targets, shortened README, added docs files, added `.dockerignore`, installed language toolchains in Docker, and added unit/integration/load test scaffolding.

**What we used / didn't use:**
Used the Makefile, README, docs, and test structure changes. The AI could not verify Docker or Go tests locally because the environment did not have `docker` or `go` installed.

## 2026-05-31 · Removing duplicate language config files

**Prompt:**
Why are there two `languages.yaml` files? Then confirmed with `ok`.

**Response summary:**
The AI explained the drift risk and moved to a single canonical config file at `config/languages.yaml`, embedded through a small `config` package.

**What we used / didn't use:**
Used one canonical YAML file and removed the duplicate runtime/docs copies.

## 2026-05-31 · Understanding the `internal/runs` package

**Prompt:**
What is `runs` folder doing?

**Response summary:**
The AI explained that `internal/runs` contains the request types, validation logic, status mapping, and output truncation helpers for the future run endpoint.

**What we used / didn't use:**
Used this package boundary. Later extended it with an executor when `/run` needed to execute submitted code.

## 2026-05-31 · Adding sample request and reply pairs

**Prompt:**
Sample request and reply pairs should be inside the reference repo under `tests/`.

**Response summary:**
The AI added JSON sample request/reply pairs under `tests/samples` for Python, C++, Java, wrong output, and unknown language.

**What we used / didn't use:**
Used the `tests/samples/<language>/<case>/request.json` and `reply.json` layout. These are contract samples, not full execution fixtures yet.

## 2026-05-31 · Stopping and rerunning the server

**Prompt:**
How do I stop the existing server, free the port, and rerun the entire code?

**Response summary:**
The AI gave `lsof` and `kill` commands for freeing port 8080 and recommended rerunning through `make build` and `make run`.

**What we used / didn't use:**
Used the Docker/Makefile workflow as the documented path. Did not add new code for this prompt.

## 2026-05-31 · Making `/run` execute the submitted body

**Prompt:**
`/run` is not executing the code passed in the body.

**Response summary:**
The AI found that `server/handler/run.go` was only a stub returning `"Run handler"`. It wired the run handler to decode JSON, validate it, execute configured commands, and return structured JSON.

**What we used / didn't use:**
Used the basic executor path for local command execution and added tests around validation/routing and language execution. Discarded the temporary plural `/runs` alias after checking the official spec. This still needs nsjail wrapping before it satisfies the sandboxing requirement.

## 2026-05-31 · Reading the AI usage discussion and creating this log

**Prompt:**
Understand `https://github.com/intern-iitm/goboxd-hackathon/discussions/5` and document the prompts used so far.

**Response summary:**
The AI fetched the discussion, identified that `docs/ai/prompts.md` is required, and extracted the expected format: prompt, response summary, and what was used or rejected.

**What we used / didn't use:**
Used the required `docs/ai/prompts.md` format and recorded the non-trivial prompts from this thread. Did not create all optional AI docs yet.

## 2026-05-31 · Correcting the endpoint name to match the spec

**Prompt:**
The project details say it is never `/runs`, it is `/run` only, with a link to the official spec.

**Response summary:**
The AI fetched the spec and confirmed that the API contract and judging notes use `POST /run`. It found plural endpoint references in the router, API docs, and integration tests.

**What we used / didn't use:**
Used the correction and removed the `/runs` route, docs text, and test calls. Kept the internal package name `internal/runs` because it is an implementation package, not a public endpoint.

## 2026-05-31 · Fixing Docker build failure from missing `go.sum`

**Prompt:**
Docker build failed during `go build` because `internal/languages/registry.go` imports `gopkg.in/yaml.v3`, but there was no `go.sum` entry.

**Response summary:**
The AI identified that `go.mod` declared the YAML dependency without a committed checksum file. It fetched the official Go checksum entry and updated the Dockerfile dependency layer to copy both `go.mod` and `go.sum` before `go mod download`.

**What we used / didn't use:**
Used the added `go.sum` file and the Dockerfile change `COPY go.mod go.sum ./`. Did not follow the error's suggested `go get github.com/thesouldev/goboxd/internal/languages` command because the issue was dependency metadata, not a new module dependency.

## 2026-05-31 · Keeping the AI prompt log current

**Prompt:**
Don't forget to update the prompts used in the file in the entire project.

**Response summary:**
The AI treated the prompt log as a required deliverable and added entries for the latest repo-changing prompts.

**What we used / didn't use:**
Used this reminder to update `docs/ai/prompts.md` immediately after the Docker build fix. No code changes were needed for this prompt.

## 2026-05-31 · Auditing Stage 1 readiness

**Prompt:**
Asked whether the repo meets the Stage 1 prototype requirements: Docker build/run, `/healthz` returns 200, `POST /run` runs one interpreted and one compiled language end-to-end, unit tests, readable README, and clean commit history.

**Response summary:**
The AI checked the Dockerfile, Compose file, Makefile, routes, executor, tests, git status, and recent commit history. It found that the implementation is mostly in place but still needs Docker/Go verification and a final commit for the latest changes.

**What we used / didn't use:**
Used the audit to identify remaining verification steps. Did not mark Stage 1 complete because this environment cannot run Docker/Go and the working tree still has uncommitted changes.
