# Ethereal → Meridian rebrand checklist

Ethereal is rebranding to **Meridian** ([meridian.xyz](https://meridian.xyz)).
This branch makes the changes that are safe now. The rest is blocked on
information from Meridian or on decisions that break things for current users.

## Done on this branch

- README now says Meridian.
- `USER_AGENT` is now `meridian-go/1.0.0dev`.
- HTTP errors now read `meridian error <status>: ...`.
- Examples read `MERIDIAN_PK`, falling back to `ETHEREAL_PK`.
- The live test reads `MERIDIAN_INTEGRATION` / `MERIDIAN_PK`, falling back to the
  `ETHEREAL_*` names. The Makefile uses the new names.

## Blocked: needs info from Meridian

- [ ] **API hosts.** `Testnet` / `Mainnet` in `client.go` still point to
      `api.etherealtest.net` / `api.ethereal.trade` (marked `TODO(meridian)`).
      Swap them for the hosts in the Perpetuals API section of
      https://docs.meridian.xyz. Keep the old values as deprecated aliases while
      the old hosts still work.
- [ ] **EIP-712 domain.** Nothing to change in code: `InitDomain` fetches the
      domain from `/v1/rpc/config`. If the server's domain `name` changes from
      `"Ethereal"`, update `testdata/rpc_config.json` and the golden domain hash
      in `client_http_test.go` to match the real response. Meridian now runs on
      Robinhood Chain, so expect `chainId` and `verifyingContract` to change too.

## Breaking changes: decide and release together

- [ ] **Repo / module path.** The module is still
      `github.com/roundinternetmoney/ethereal-rest`, while this repo lives at
      `qiwi1272/ethereal-rest`. Pick the final home (e.g. `.../meridian-rest`),
      rename the GitHub repo (GitHub redirects the old URL), then update
      `go.mod`, the example imports and the README badges.
- [ ] **Package name.** `etherealRest` → e.g. `meridianRest` (or `meridian`).
      Callers that import without an alias will break, so make this change in
      the same release as the module path.
- [ ] **Tag a release** that marks the rename, with a note telling users to
      move from `ETHEREAL_PK` to `MERIDIAN_PK`.
- [ ] Drop the `ETHEREAL_*` fallbacks in a later release.
