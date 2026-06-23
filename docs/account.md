# Account Module

> Every command in this module accepts a global `--json` flag that prints the
> raw API response as indented JSON instead of a table.

Exec: `./bitget-cli UTA account [Subcommand] [Arguments]`

## Quick Navigation
- [assets](#assets---per-coin-balances)
- [equity](#equity---aggregate-equity--margin)
- [health](#health---account-health--risk)
- [info](#info---identity--permissions)
- [settings](#settings---account-mode)
- [leverage-config](#leverage-config---per-symbol-leverage)
- [fee-rate](#fee-rate---makertaker-fee-rate)
- [funding-assets](#funding-assets---funding-p2p-balances)
- [bills](#bills---financial-ledger-records)
- [max-transferable](#max-transferable)
- [max-withdrawal](#max-withdrawal)
- [set-leverage](#set-leverage)
- [set-hold-mode](#set-hold-mode)
- [set-margin](#set-margin)

## assets - per-coin balances
Docs Link: <https://www.bitget.com/api-doc/uta/account/Get-Account>

Shows the unified account's per-coin balances (non-zero only).

Exec: `./bitget-cli UTA account assets`
```shell
┌──────┬─────────────┬─────────────┬─────────────┬────────┬──────┬─────────────┐
│ COIN │   EQUITY    │   BALANCE   │  AVAILABLE  │ LOCKED │ DEBT │  USD VALUE  │
├──────┼─────────────┼─────────────┼─────────────┼────────┼──────┼─────────────┤
│ USDT │ 98.90701886 │ 98.90701886 │ 98.90701886 │ 0      │ 0    │ 98.79203489 │
│ ETH  │ 0.0000991   │ 0.0000991   │ 0.0000991   │ 0      │ 0    │ 0.1722732   │
└──────┴─────────────┴─────────────┴─────────────┴────────┴──────┴─────────────┘
```

## equity - aggregate equity & margin
Docs Link: <https://www.bitget.com/api-doc/uta/account/Get-Account>

Same endpoint as `assets`, presenting the account-level aggregates.

Exec: `./bitget-cli UTA account equity`
```shell
┌────────────────┬─────────────┬────────────┬────────────────┬─────────────┬─────┬─────┬───────────┐
│ ACCOUNT EQUITY │ USDT EQUITY │ BTC EQUITY │ UNREALISED PNL │ EFF EQUITY  │ IMR │ MMR │ MGN RATIO │
├────────────────┼─────────────┼────────────┼────────────────┼─────────────┼─────┼─────┼───────────┤
│ 99.0243065     │ 99.13763545 │ 0.0015464  │ 0              │ 98.79395354 │ 0   │ 0   │ 0         │
└────────────────┴─────────────┴────────────┴────────────────┴─────────────┴─────┴─────┴───────────┘
```

## health - account health & risk
Docs Link: <https://www.bitget.com/api-doc/uta/account/Get-Account>

Same endpoint as `assets`/`equity`, presenting only the risk/health metrics:
account and effective equity, unrealised PnL, initial/maintenance margin
requirements (IMR/MMR) and the margin ratio. A margin ratio approaching 1
(100%) signals liquidation risk.

Exec: `./bitget-cli UTA account health`
```shell
┌────────────────┬─────────────┬────────────────┬─────┬─────┬───────────┐
│ ACCOUNT EQUITY │ EFF EQUITY  │ UNREALISED PNL │ IMR │ MMR │ MGN RATIO │
├────────────────┼─────────────┼────────────────┼─────┼─────┼───────────┤
│ 99.0243065     │ 98.79395354 │ 0              │ 0   │ 0   │ 0         │
└────────────────┴─────────────┴────────────────┴─────┴─────┴───────────┘
```

## info - identity & permissions
Docs Link: <https://www.bitget.com/api-doc/uta/account/Get-Account-Info>

Exec: `./bitget-cli UTA account info`
```shell
┌────────────┬───────────┬────────────────┬───────────────────┬─────────┬─────────────────────┐
│  USER ID   │ PARENT ID │   PERM TYPE    │    PERMISSIONS    │ IP LIST │    REGISTER TIME    │
├────────────┼───────────┼────────────────┼───────────────────┼─────────┼─────────────────────┤
│ 1234567890 │           │ read_and_write │ uta_mgt,uta_trade │         │ 2024-11-26 08:12:07 │
└────────────┴───────────┴────────────────┴───────────────────┴─────────┴─────────────────────┘
```

## settings - account mode
Docs Link: <https://www.bitget.com/api-doc/uta/account/Get-Account-Setting>

Exec: `./bitget-cli UTA account settings`
```shell
┌────────────┬──────────────┬──────────────┬───────────────┬──────────────┬──────────┐
│    UID     │ ACCOUNT MODE │  ASSET MODE  │ ACCOUNT LEVEL │  HOLD MODE   │ STP MODE │
├────────────┼──────────────┼──────────────┼───────────────┼──────────────┼──────────┤
│ 1234567890 │ unified      │ multi_assets │ basic         │ one_way_mode │ none     │
└────────────┴──────────────┴──────────────┴───────────────┴──────────────┴──────────┘
```

## leverage-config - per-symbol leverage
Docs Link: <https://www.bitget.com/api-doc/uta/account/Get-Account-Setting>

Same endpoint as `settings`, listing the per-symbol leverage / margin-mode
configuration (`symbolConfigList`).

Exec: `./bitget-cli UTA account leverage-config`

## fee-rate - maker/taker fee rate
Docs Link: <https://www.bitget.com/api-doc/uta/account/Get-Account-Fee-Rate>

Exec: `./bitget-cli UTA account fee-rate --category=spot --symbol=BTCUSDT`
```shell
┌─────────┬────────────────┬────────────────┐
│ SYMBOL  │ MAKER FEE RATE │ TAKER FEE RATE │
├─────────┼────────────────┼────────────────┤
│ BTCUSDT │ 0.001          │ 0.001          │
└─────────┴────────────────┴────────────────┘
```
**Supported parameters:**
- `--category, -C`: spot, usdt-futures, coin-futures, usdc-futures (required)
- `--symbol, -s`: trading pair symbol (required)

## funding-assets - funding (P2P) balances
Docs Link: <https://www.bitget.com/api-doc/uta/account/Get-Account-Funding-Assets>

Exec: `./bitget-cli UTA account funding-assets [--coin=USDT]`
```shell
┌──────┬─────────┬───────────┬────────┐
│ COIN │ BALANCE │ AVAILABLE │ FROZEN │
├──────┼─────────┼───────────┼────────┤
│ USDT │ 1       │ 1         │ 0      │
└──────┴─────────┴───────────┴────────┘
```

## bills - financial (ledger) records
Docs Link: <https://www.bitget.com/api-doc/uta/account/Get-Financial-Records>

Exec: `./bitget-cli UTA account bills --category=usdt-futures --limit=3`
```shell
┌─────────────────────┬──────────────┬──────┬───────────┬─────────────┬─────────────┬───────────────┬─────────┐
│        TIME         │   CATEGORY   │ COIN │   TYPE    │   AMOUNT    │     FEE     │    BALANCE    │ SYMBOL  │
├─────────────────────┼──────────────┼──────┼───────────┼─────────────┼─────────────┼───────────────┼─────────┤
│ 2026-06-22 07:09:48 │ USDT-FUTURES │ USDT │ BUY_DEAL  │ -0.01044138 │ -0.01044138 │ 98.8965774852 │ ETHUSDT │
└─────────────────────┴──────────────┴──────┴───────────┴─────────────┴─────────────┴───────────────┴─────────┘
```
**Supported parameters:**
- `--category, -C`: spot, usdt-futures, coin-futures, usdc-futures (required)
- `--coin, -c`: coin filter
- `--type, -t`: record type filter
- `--startTime, -a` / `--endTime, -e`: unix ms or "YYYY-MM-DD HH:MM:SS" (90-day window)
- `--limit, -l`: max records

## max-transferable
Docs Link: <https://www.bitget.com/api-doc/uta/account/Get-Max-Transferable>

Exec: `./bitget-cli UTA account max-transferable --coin=USDT`
```shell
┌──────┬─────────────────────┬─────────────────────┐
│ COIN │    MAX TRANSFER     │ BORROW MAX TRANSFER │
├──────┼─────────────────────┼─────────────────────┤
│ USDT │ 80.6670165561479888 │ 0                   │
└──────┴─────────────────────┴─────────────────────┘
```

## max-withdrawal
Docs Link: <https://www.bitget.com/api-doc/uta/account/Get-Max-Withdrawal>

Exec: `./bitget-cli UTA account max-withdrawal --coin=USDT`
```shell
┌──────┬─────────────────────┬──────────┬─────────┬─────────────────────┐
│ COIN │       UTA MAX       │ SPOT MAX │ OTC MAX │      TOTAL MAX      │
├──────┼─────────────────────┼──────────┼─────────┼─────────────────────┤
│ USDT │ 80.6670165561479888 │ 1        │ 0       │ 81.6670165561479888 │
└──────┴─────────────────────┴──────────┴─────────┴─────────────────────┘
```

## set-leverage
Docs Link: <https://www.bitget.com/api-doc/uta/account/Change-Leverage>

Exec: `./bitget-cli UTA account set-leverage --category=usdt-futures --symbol=BTCUSDT --leverage=10`

**Supported parameters:**
- `--category, -C`: usdt-futures, coin-futures, usdc-futures, margin (required)
- `--leverage, -L`: leverage multiple, e.g. 10 (required)
- `--symbol, -s`: futures symbol (for futures)
- `--coin, -c`: margin coin (for margin)
- `--marginMode, -m`: crossed, isolated
- `--posSide, -p`: long, short (isolated)
- `--longLeverage` / `--shortLeverage`: isolated two-way leverage

## set-hold-mode
Docs Link: <https://www.bitget.com/api-doc/uta/account/Change-Position-Mode>

Exec: `./bitget-cli UTA account set-hold-mode --holdMode=hedge_mode`

`--holdMode, -H`: `one_way_mode` or `hedge_mode` (required)

## set-margin
Docs Link: <https://www.bitget.com/api-doc/uta/account/Set-Margin>

Add or reduce isolated-position margin for a futures symbol.

Exec: `./bitget-cli UTA account set-margin --category=usdt-futures --symbol=BTCUSDT --posSide=long --operation=add --amount=5`

**Supported parameters:**
- `--category, -C`: usdt-futures, coin-futures, usdc-futures (required)
- `--symbol, -s`: futures symbol (required)
- `--posSide, -p`: long, short (default long)
- `--operation, -o`: add or reduce (required)
- `--amount, -m`: margin amount, decimal (required)
