# batchlog

batchlog is a tool to retrieve Ethereum event logs for specific contracts, particularly designed for Swarm's Postage Stamp contract on the Gnosis Chain. It fetches logs within a specified block range and saves them to a file.

## Features

- Retrieve event logs for a specified contract address and block range.
- Handles large block ranges by querying in smaller chunks.
- Supports rate limiting for RPC requests.
- Saves retrieved logs to an NDJSON file (`export.ndjson`).
- Graceful shutdown on interrupt signals (Ctrl+C).

## Requirements

- Go 1.24 or later

## Installation

```sh
git clone https://github.com/gacevicljubisa/batchlog.git
cd batchlog
make binary
```
