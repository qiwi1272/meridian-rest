# Golang Rest Client for Meridian API

[Go Reference](https://pkg.go.dev/github.com/qiwi1272/meridian-rest)
[Go Report Card](https://goreportcard.com/report/github.com/qiwi1272/meridian-rest)
[Release](https://github.com/qiwi1272/meridian-rest/actions/workflows/release.yaml)

Lightweight golang client for interacting with the [Meridian](https://meridian.xyz) API (formerly Ethereal).

## Migrating from Ethereal

- The module moved from `github.com/roundinternetmoney/ethereal-rest` to `github.com/qiwi1272/meridian-rest`, and the package is now `meridianRest`. Update your imports. Versions up to `v1.0.8` stay available at the old path.
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
- Install from GitHub: `go get github.com/qiwi1272/meridian-rest`
- Import path: `github.com/qiwi1272/meridian-rest`

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

## Releasing

Releases are manual and GitHub marks them immutable, so a published tag can't be moved or deleted.

1. Merge to `main`. Merging does not publish anything.
2. Go to **Actions → Release → Run workflow**, enter the version (for example `1.1.0`) and run it.

The workflow runs the checks, then tags `v<version>` and publishes the GitHub release marked as **Latest**. Without that flag, GitHub would keep a higher old version such as `v2.0.6` as Latest. Locally, `make release NEW_VERSION=1.1.0` does the same. Only `0.x.y` / `1.x.y` versions that aren't already tagged are accepted, because the module path has no `/vN` suffix.

## Contributing

Contributions are welcome! Please open issues or pull requests as needed.

## Todo

Endpoints the client doesn't cover yet, from the [OpenAPI spec](https://api.meridian.xyz/openapi.json). **New** marks endpoints added or changed in Meridian.

**Orders**
- [ ] `POST /v1/order/cancel-all`: cancel every open order for the subaccount
- [ ] `POST /v1/order/batch`: place and cancel in one request (the client currently sends batches as concurrent single requests)
- [ ] `POST /v1/order/dry-run`: simulate an order (unsigned body)
- [ ] `GET /v1/order`, `GET /v1/order/{id}`, `GET /v1/order/{id}/group`: query orders
- [ ] `GET /v1/order/fill`, `GET /v1/order/trade`: fill and trade history

**Positions and fees**
- [ ] `GET /v1/position/active`, `GET /v1/position/{id}`, `GET /v1/position/fill`, `GET /v1/position/liquidation`
- [ ] **New** `GET /v1/position-fee`, `GET /v1/position-fee/projected-rate`: mPerp position fees
- [ ] **New** `GET /v1/funding` (now takes `startTime`/`endTime` with cursor pagination instead of `range`), `GET /v1/funding/projected-rate` (replaces `/v1/funding/projected`)

**Products and market data**
- [ ] `GET /v1/product/{id}`, `GET /v1/product/market-price`, `GET /v1/product/market-liquidity`
- [ ] **New** `GET /v1/product/mark-price-gap`: scheduled mark-price gaps for mPerps

**Tokens and transfers**
- [ ] **New** `POST /v1/token/convert`: move balance between quote-token pools (signed `ConvertToken`)
- [ ] `GET /v1/token`, `GET /v1/token/{id}`, `GET /v1/token/transfer`
- [ ] `POST /v1/token/{id}/withdraw`, `GET /v1/token/withdraw` (signed `InitiateWithdraw`)

**Accounts and signers**
- [ ] `GET /v1/subaccount/{id}`, plus support for more than one subaccount per client
- [ ] Linked signers: `GET /v1/linked-signer`, `POST /v1/linked-signer/link`, `/refresh`, `/extend`, `DELETE /v1/linked-signer/revoke`

**Utilities**
- [ ] `GET /v1/time`: sync the clock used for `nonce` / `signedAt`
- [ ] `GET /v1/rate-limit/config`, `GET /v1/maintenance`

**Models and housekeeping**
- [ ] Add Meridian fields to the existing models: `marginMode`, `quoteTokenAddress` and `DELISTED` status on `Product`; `positionFeeUsd` / `positionFeeAccruedUsd` on positions and balances
