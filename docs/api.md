# API

## GET /healthz

Returns `200 OK` with body `OK` when the server process is alive.

## POST /run

The execution endpoint accepts one source submission and one or more tests. Validation failures return `400`. User-code failures return `200` with a status in the response body.

Request fields:

- `language`: required. Must match an id in the language registry.
- `source`: required UTF-8 source text. The default limit is 256 KiB.
- `source_filename`: optional unless the language uses `source_filename_strategy: from_request`.
- `artifact_filename`: optional unless the language uses `artifact_filename_strategy: from_request`.
- `build`: optional build limit and flag overrides.
- `run`: optional run limit and flag overrides.
- `tests`: required. Must contain at least one entry with `stdin` and `expected_stdout`.

Example:

```json
{
  "language": "py3",
  "source": "print(input().upper())",
  "tests": [
    {
      "stdin": "hi\n",
      "expected_stdout": "HI\n"
    }
  ]
}
```

Bad requests use this shape:

```json
{
  "error": {
    "code": "unknown_language",
    "message": "language is not configured"
  }
}
```

Sample request and reply pairs live under [tests/samples](../tests/samples).
