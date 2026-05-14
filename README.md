# indexer

A Go tool that walks a maildir-style directory tree, parses each RFC 822 email file, and bulk-uploads the parsed records to [ZincSearch](https://github.com/zincsearch/zincsearch) for full-text search.

Built originally to index the [Enron email corpus](https://www.cs.cmu.edu/~enron/) (~500k messages, 1.5 GB on disk) and feed the companion search UI in [email-search-engine](https://github.com/rdaniel1105/email-search-engine).

## What it does

- Recursively walks a maildir-style directory tree
- Parses each file as an RFC 822 email (Message-ID, From, To, Subject, headers, body)
- Repairs common header continuation bugs in Enron-style dumps before parsing
- Deduplicates identical bodies via a SHA-256 hash table
- Buffers records into ndjson chunks (capped at 5,000 entries or 10 MB, whichever comes first)
- Posts chunks concurrently to ZincSearch's `/_multi` bulk endpoint with a semaphore-bounded worker pool

## Architecture

```
main.go
├── helpers/
│   ├── directory_reader.go  / directory_checker.go   Filesystem traversal
│   ├── emails.go     RFC 822 parsing + dedup
│   ├── headers.go    Repair multi-line header bugs before mail.ReadMessage()
│   └── bulk_data.go  Buffered ndjson chunker + ZincSearch HTTP client
└── models/
    └── email.go      Email struct mapping every header to a JSON field
```

The chunk-and-stream design keeps memory bounded regardless of corpus size — the whole Enron dataset runs in a few hundred MB of resident memory rather than loading every message at once.

## Requirements

- Go 1.26+
- A reachable [ZincSearch](https://github.com/zincsearch/zincsearch) instance (defaults to `http://localhost:4080`)

## Configuration

```bash
cp .env.example .env
```

| Variable | Default | Purpose |
|---|---|---|
| `ZINCSEARCH_URL` | `http://localhost:4080` | Base URL of the ZincSearch instance |
| `ZINCSEARCH_INDEX` | `emails` | Target index name |
| `ZINCSEARCH_USERNAME` | `""` | ZincSearch admin user |
| `ZINCSEARCH_PASSWORD` | `""` | ZincSearch admin password |

## Usage

```bash
go run . /path/to/maildir
```

The argument must point at the root of a maildir-style tree. For the Enron corpus that's typically `enron_mail_<date>/maildir/`.

## Tests

```bash
go test ./...
```

A tiny synthetic email lives under `helpers/testdata/maildir/sample/sent/1` so tests pass on a fresh clone without any external data.

## License

MIT — see [LICENSE](./LICENSE).
