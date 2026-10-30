# Golang Rest Client for Meridian API

[Go Reference](https://pkg.go.dev/github.com/roundinternetmoney/ethereal-rest)
[Go Report Card](https://goreportcard.com/report/github.com/roundinternetmoney/ethereal-rest)
[Release](https://github.com/roundinternetmoney/ethereal-rest/actions/workflows/release.yaml)

Lightweight golang client for interacting with the [Meridian](https://meridian.xyz) API (formerly Ethereal).

## Migrating from Ethereal

- `rest.Testnet` / `rest.Mainnet` now point at `api.meridiantest.net` / `api.meridian.xyz`. The old hosts remain as deprecated `rest.EtherealTestnet` / `rest.EtherealMainnet`.
- Meridian is a separate deployment: Ethereal subaccounts, balances and IDs do not carry over. Register and fund a new subaccount.
- The EIP-712 domain is still loaded from `/v1/rpc/config`, so signing needs no changes.
- See Meridian's [Migrate from Ethereal](https://docs.meridian.xyz/developer-guides/trading-api/migrate-from-ethereal) guide for the full list of API changes.

## Features

- Batch execution support (concurrent, ordered, type-safe)
- EIP-712 data signing
- Automatic nonce, timestamp, and subaccount handling for requests
- Minimal dependencies

## Getting started

- Requires Go 1.25+.
- Install from GitHub: `go get github.com/roundinternetmoney/ethereal-rest`
- Import path: `github.com/roundinternetmoney/ethereal-rest`

## Example Usage

From the client directory:

- `make examples`
- Run a binary with your key set, for example: `MERIDIAN_PK=<hex> ./bin/example_account_balance`

All examples read **`MERIDIAN_PK`** (hex private key, with or without `0x` prefix). The legacy `ETHEREAL_PK` is still accepted as a fallback. They target **testnet** by default (`rest.Testnet`).

| Example | Binary (`make examples`) | What it shows |
|--------|---------------------------|---------------|
| [examples/account_balance](./examples/account_balance/) | `bin/example_account_balance` | Read-only balances after client init |
| [examples/limit_single](./examples/limit_single/) | `bin/example_limit_single` | Single limit order, cancel, `Send` helper |
| [examples/limit_multiple](./examples/limit_multiple/) | `bin/example_limit_multiple` | `CreateOrders` batch + multi-ID cancel |
| [examples/positions](./examples/positions/) | `bin/example_positions` | Open positions (`GetPosition`) |
| [examples/cancel_replace](./examples/cancel_replace/) | `bin/example_cancel_replace` | Cancel then submit a new order (replace) |
| [examples/twap](./examples/twap/) | `bin/example_twap` | **Composition:** time-sliced orders (not a venue TWAP type) |
| [examples/chase](./examples/chase/) | `bin/example_chase` | **Composition:** cancel/replace loop with a stub reference price |

For more detail, see the [examples/](./examples/) folder and the file comments in `twap` and `chase`.

## Configuration Notes

- If no private key is passed to the rest client, an error will be returned.
- All signable request messages implement the `Signable` interface.
- Only one subaccount is currently supported; by default the first one discovered is used.

## Contributing

Contributions are welcome! Please open issues or pull requests as needed.

## Todo

- add missing rest methods
- add Meridian-only fields (`marginMode`, quote-token balances, position fees) and endpoints (`order/cancel-all`, `token/convert`)
- rename the module path and package from `ethereal` to `meridian`