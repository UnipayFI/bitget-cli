# Bitget CLI

A command-line tool for the Bitget API developed in Go, supporting both of
Bitget's account systems:

- **Unified Trading Account (UTA)** — the v3 REST API (`/api/v3/*`), reached
  under the `UTA` subcommand (`./bitget-cli UTA ...`). The unified account
  serves spot and every futures line from one client; the product is chosen per
  command via the *category*.
- **Classic account** — the v2 REST API (`/api/v2/*`), reached through the
  top-level `spot` and `futures` commands.

Both use HMAC-SHA256 signing (`ACCESS-KEY` / `ACCESS-SIGN` /
`ACCESS-TIMESTAMP` / `ACCESS-PASSPHRASE`) with the same credentials. Only
authenticated (private) endpoints are covered — public market data is out of
scope.

## Installation and Configuration

### Install (prebuilt binary)
```shell
curl -sSL https://raw.githubusercontent.com/UnipayFI/bitget-cli/refs/heads/main/download.sh | bash
```
Downloads the latest release for your platform/arch from GitHub Releases.

### Build from source
```shell
go build -o bitget-cli .
```
Releases are produced by the `Release` GitHub Action (`.github/workflows/release.yml`),
which cross-compiles for Linux/macOS/Windows (amd64 + arm64) on every `v*` tag and
injects the version via ldflags.

### Environment variables
Before using, set your Bitget API credentials (from the Bitget API-management
page):
```shell
export BITGET_API_KEY="bg_..."        # API key
export BITGET_API_SECRET="..."        # API secret
export BITGET_PASSPHRASE="..."        # passphrase set when the key was created

# Optional
export HTTPS_PROXY="socks5://127.0.0.1:1080"    # route REST traffic through a proxy
                                                # (HTTPS_PROXY > ALL_PROXY > HTTP_PROXY;
                                                #  supports http/https/socks5/socks5h)
export BITGET_LOCALE="en-US"                    # error-message language (default en-US)
export BITGET_BASE_URL="https://api.bitget.com" # override REST base URL
export BITGET_DEMO="true"                       # use demo (paper) trading
```

### Output format
Every command supports a global `--json` flag. Without it, results render as a
table; with it, the raw API response is printed as indented JSON, e.g.:
```shell
./bitget-cli UTA account assets --json
```

## Usage
Top-level commands split by account system:
```
./bitget-cli [Command] [Subcommand] [Arguments]

Available Commands:
  UTA         Unified Trading Account (UTA) commands (v3 /api/v3/*)
  spot        Classic spot trading (account, balances & orders) (v2)
  futures     Classic futures trading (account, positions & orders) (v2)
  version     Print version information
```

Each leaf subcommand's `-h` output includes a `Docs Link:` pointing to the
official Bitget API documentation page for that endpoint.

## UTA (Unified Trading Account)
All unified-account features live under the `UTA` subcommand:
```
./bitget-cli UTA [Module] [Subcommand] [Arguments]

Available Commands:
  account     Account info, balances, settings and leverage
  spot        Spot trading (orders & fills)
  futures     Futures trading (orders, positions & fills)
  wallet      Funds: transfer, deposit and withdrawal
```

### Category
UTA spot commands always target the `SPOT` category. UTA futures commands accept
a persistent `--category` / `-C` flag (default `usdt-futures`):

| alias | category      |
|-------|---------------|
| `usdt` / `usdt-futures` | USDT-FUTURES |
| `coin` / `coin-futures` | COIN-FUTURES |
| `usdc` / `usdc-futures` | USDC-FUTURES |
| `spot` | SPOT |
| `margin` | MARGIN |

### Account Module
Exec: `./bitget-cli UTA account [Subcommand] [Arguments]`
```shell
Available Commands:
  assets            Show unified account per-coin balances (non-zero)
  equity            Show unified account aggregate equity and margin
  health            Show unified account health (margin ratio & risk)
  info              Show account identity and permissions
  settings          Show account mode settings
  leverage-config   Show per-symbol leverage / margin configuration
  fee-rate          Show maker/taker fee rate for a symbol
  funding-assets    Show funding (P2P) account balances
  bills             Show account financial (ledger) records
  max-transferable  Show max transferable amount for a coin
  max-withdrawal    Show max withdrawable amount for a coin
  set-leverage      Set leverage for a coin or futures symbol
  set-hold-mode     Set futures position hold mode (one-way / hedge)
  set-margin        Add/reduce isolated position margin
```
**[View detailed documentation](docs/account.md)**

### Spot Module
Exec: `./bitget-cli UTA spot [Subcommand] [Arguments]`
```shell
Available Commands:
  order       Create, modify, cancel and query spot orders
  fills       List spot trade fills
```
**[View detailed documentation](docs/spot.md)**

### Futures Module
Exec: `./bitget-cli UTA futures [Subcommand] [Arguments]`
```shell
Available Commands:
  order       Create, modify, cancel and query futures orders
  position    Query and close futures positions
  fills       List futures trade fills
```
**[View detailed documentation](docs/futures.md)**

### Wallet Module
Exec: `./bitget-cli UTA wallet [Subcommand] [Arguments]`
```shell
Available Commands:
  transfer            Transfer a coin between account types
  transferable-coins  List coins transferable between two account types
  deposit             Deposit address and records
  withdraw            Withdraw, list records and address book
```
**[View detailed documentation](docs/wallet.md)**

## Classic account
The top-level `spot` and `futures` commands target the classic account
(v2 `/api/v2/*`) and cover the basic operations: account info, balances,
positions, account health and order create/cancel/query.

### Product (classic futures)
Classic futures commands accept a persistent `--product` / `-P` flag
(default `usdt-futures`): `usdt` / `coin` / `usdc` → `USDT-FUTURES` /
`COIN-FUTURES` / `USDC-FUTURES`.

### Classic Spot
Exec: `./bitget-cli spot [Subcommand] [Arguments]`
```shell
Available Commands:
  account     Spot account info and balances
    info        Show spot account identity and permissions
    assets      Show spot per-coin balances (non-zero)
  order       Create, cancel and query spot orders
    create      Create a spot order
    cancel      Cancel a spot order
    get         Query a single spot order
    open        List open spot orders
```
**[View detailed documentation](docs/classic-spot.md)**

### Classic Futures
Exec: `./bitget-cli futures [Subcommand] [Arguments]`
```shell
Available Commands:
  account     Show futures account balances and equity
  health      Show futures account health (risk rate & maintenance margin)
  position    Query futures positions
    list        List open positions
  order       Create, cancel and query futures orders
    create      Create a futures order
    cancel      Cancel a futures order
    get         Query a single futures order
    open        List open futures orders
```
**[View detailed documentation](docs/classic-futures.md)**

## Official API documentation
- UTA: <https://www.bitget.com/api-doc/uta/intro>
- Classic spot: <https://www.bitget.com/api-doc/spot/intro>
- Classic futures: <https://www.bitget.com/api-doc/contract/intro>
