# Languages

Languages are configured in [config/languages.yaml](../config/languages.yaml). Adding a standard language should be a YAML change plus tests.

Each language has an `id`, a display `name`, a source filename rule, and a `run` command. Compiled languages also define `artifact` and `build`.

Supported filename strategies:

- `source_filename`: fixed source filename written by the server.
- `source_filename_strategy: from_request`: request must provide `source_filename`.
- `artifact`: fixed build artifact name.
- `artifact_filename_strategy: from_request`: request must provide `artifact_filename`.

Command arguments may use these placeholders:

- `{{source}}`: source filename inside the sandbox.
- `{{artifact}}`: compiled artifact filename inside the sandbox.
- `{{flags}}`: request flags after allow-list filtering.

Flags are denied by default. Add safe flags to `flag_allowlist`. A trailing `*` permits a prefix, for example `-std=*`.
