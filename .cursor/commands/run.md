# Run Local

Start the actx0 platform locally: Apple Container system, Qdrant, Postgres, then the API server.

## Steps

1. From the workspace root, run these in order via the Shell tool with `all` permissions:

```bash
container system start
```

```bash
./bin/qdrant --api-key 871c1580-9f7a-4f89-bd24-536a4b570c2b
```

```bash
./bin/pgsql
```

2. Start the API server in the background (`block_until_ms: 0`) with `all` permissions:

```bash
make run
```

3. Confirm each prior step succeeded, then report that the server is starting. If any step fails, stop and show the error output.

## Notes

- Do not ask for confirmation.
- Run from the repo root.
- `make run` is long-lived; leave it running in the background.
- Qdrant and Postgres scripts default to `start` when no subcommand is given.
