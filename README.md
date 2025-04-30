# batchlog

batchlog is a tool to retrieve Ethereum event logs for specific contracts, particularly designed for Swarm's Postage Stamp contract on the Gnosis Chain. It fetches logs within a specified block range using the `export` command and saves them to a file.

## Features

- Retrieve event logs for a specified contract address and block range.
- Handles large block ranges by querying in smaller chunks.
- Supports rate limiting for RPC requests.
- Saves retrieved logs to a specified output file (default: `export.ndjson`) in NDJSON format.
- Graceful shutdown on interrupt signals (Ctrl+C).

## Requirements

- Go 1.24 or later

## Installation

```sh
git clone https://github.com/gacevicljubisa/batchlog.git
cd batchlog
make binary
# The binary will be located in the dist/ folder
```

## Usage

The primary command is export.

```sh
./dist/batchlog export --help
```

```sh
./dist/batchlog export \
    --start 31306381 \
    --endpoint <YOUR_GNOSIS_RPC_ENDPOINT> \
    --output my_logs.ndjson
```

## Flags

```sh
  -b, --block-range-limit uint32   Max blocks per log query (default 5)
      --end uint                   End block (optional, uses latest block if 0) (default 39810670)
  -e, --endpoint string            Ethereum RPC endpoint URL.
  -h, --help                       help for export
  -m, --max-request int            Max RPC requests/sec (default 15)
  -o, --output string              Output file path (NDJSON) (default "export.ndjson")
      --start uint                 Start block (optional, uses contract start block if 0) (default 31306381)
```
