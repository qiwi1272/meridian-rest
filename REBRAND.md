# Ethereal → Meridian rebrand checklist

Ethereal is now **Meridian** ([docs.meridian.xyz](https://docs.meridian.xyz)).
Meridian is a separate deployment on Robinhood Chain: Ethereal subaccounts,
balances, IDs and linked signers do not carry over. See the official
[Migrate from Ethereal](https://docs.meridian.xyz/developer-guides/trading-api/migrate-from-ethereal)
guide.

## Done on this branch

- `Testnet` / `Mainnet` now point to `https://api.meridiantest.net` /
  `https://api.meridian.xyz` (from
  [API Hosts](https://docs.meridian.xyz/protocol-reference/api-hosts)). The old
  hosts remain as deprecated `EtherealTestnet` / `EtherealMainnet`.
- README now says Meridian.
- `USER_AGENT` is now `meridian-go/1.0.0dev`.
- HTTP errors now read `meridian error <status>: ...`.
- Examples read `MERIDIAN_PK`, falling back to `ETHEREAL_PK`.
- The live test reads `MERIDIAN_INTEGRATION` / `MERIDIAN_PK`, falling back to the
  `ETHEREAL_*` names. The Makefile uses the new names.

## Checked against the migration guide, no change needed

- **Signing.** The EIP-712 domain and type strings are already loaded from
  `GET /v1/rpc/config` at startup. Nonce (nanosecond string), `signedAt`
  (seconds), D9 decimals, `bytes32` subaccount and `onchainId` as `productId`
  all match [Message Signing](https://docs.meridian.xyz/developer-guides/trading-api/message-signing).
- **Endpoints used** (`rpc/config`, `subaccount`, `subaccount/balance`,
  `position`, `product`, `order`, `order/cancel`) are not on the breaking list.
  The client does not use `/v1/funding`, `/v1/whitelist` or `X-Ethereal-*`
  headers.
- **IDs** are treated as opaque strings, so UUIDv7 works.

## Still to do

- [ ] **Live testnet check.** Run `make test` and an example against
      `api.meridiantest.net`. This sandbox could reach the docs but not the API
      hosts, so nothing here has been verified against the live API.
- [ ] **Test fixture.** `testdata/rpc_config.json` still holds an Ethereal-style
      domain (`"Ethereal"`, chain `421614`). Meridian uses chain `4663`
      (mainnet) / `46630` (testnet). Replace it with a real
      `GET /v1/rpc/config` response from testnet, then update the golden domain
      hash in `client_http_test.go`.
- [ ] **New schema fields.** Meridian adds quote-token pools and mPerps:
      `marginMode`, `quoteTokenAddress`, `positionFeeUsd`,
      `positionFeeAccruedUsd`, `DELISTED` status. Compare the `Product`,
      `Position` and balance structs in `messages.go` against
      `https://api.meridiantest.net/openapi.json` and add what callers need.
- [ ] **Optional new endpoints:** `POST /v1/order/cancel-all`,
      `POST /v1/order/dry-run`, `POST /v1/token/convert` (with the new
      `ConvertToken` signed type).

## Breaking changes: decide and release together

- [ ] **Repo / module path.** The module is still
      `github.com/roundinternetmoney/ethereal-rest`, while this repo lives at
      `qiwi1272/ethereal-rest`. Meridian's
      [SDK page](https://docs.meridian.xyz/developer-guides/trading-api/sdk)
      lists the `roundinternetmoney` path as a community Go client. Pick the
      final home (e.g. `.../meridian-rest`), rename the GitHub repo (GitHub
      redirects the old URL), then update `go.mod`, the example imports and the
      README badges.
- [ ] **Package name.** `etherealRest` → e.g. `meridianRest` (or `meridian`).
      Callers that import without an alias will break, so make this change in
      the same release as the module path.
- [ ] **Tag a release** that marks the rename. Its notes should say that
      `Testnet` / `Mainnet` now point to Meridian, and that users should move
      from `ETHEREAL_PK` to `MERIDIAN_PK`.
- [ ] Drop the `ETHEREAL_*` fallbacks and the deprecated `Ethereal*` hosts in a
      later release.
