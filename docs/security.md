# Security

The container is the unit of execution. The Docker image builds nsjail from the upstream `3.4` tag during image build and does not assume nsjail exists on the host.

Request validation happens before the server writes source files or starts a sandbox. The validator rejects unknown languages, oversized source, malformed filenames, unsupported build or run options, and flags outside the configured allow-list.

Filenames must be one path component. Separators, leading dots, and names over the length cap are rejected.

The API treats bad requests differently from user-code outcomes:

- Invalid request data returns `400`.
- Server and sandbox setup failures return `5xx`.
- Compile errors, runtime errors, timeouts, memory exits, and wrong output return `200` with a status in the JSON body.

Docker Compose runs the service with `privileged: true` because nsjail uses Linux namespace features that are not available in a default container profile.
