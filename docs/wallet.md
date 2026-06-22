# Wallet Module

> Every command in this module accepts a global `--json` flag that prints the
> raw API response as indented JSON instead of a table.

Exec: `./bitget-cli wallet [Subcommand] [Arguments]`

Account types used by transfers: `spot`, `p2p`, `coin_futures`, `usdt_futures`,
`usdc_futures`, `crossed_margin`, `isolated_margin`, `uta`.

## Quick Navigation
- [transfer](#transfer)
- [transferable-coins](#transferable-coins)
- [deposit address](#deposit---address)
- [deposit records](#deposit---records)
- [withdraw create](#withdraw---create)
- [withdraw records](#withdraw---records)
- [withdraw address](#withdraw---address)
- [withdraw cancel](#withdraw---cancel)

## transfer
Docs Link: <https://www.bitget.com/api-doc/uta/account/Transfer>

Transfer a coin between two account types within the same account.

Exec: `./bitget-cli wallet transfer --fromType=spot --toType=usdt_futures --coin=USDT --amount=10`
```shell
┌─────────────────────┬─────────────────────┐
│     TRANSFER ID     │     CLIENT OID      │
├─────────────────────┼─────────────────────┤
│ 1452850000000000000 │                     │
└─────────────────────┴─────────────────────┘
```
**Supported parameters:**
- `--fromType, -f`: source account type (required)
- `--toType, -t`: target account type (required)
- `--coin, -c`: coin (required)
- `--amount, -m`: amount, decimal (required)
- `--symbol, -s`: isolated spot-margin symbol
- `--allowBorrow`: yes or no (auto-borrow when insufficient)
- `--clientOid`: client transaction id

## transferable-coins
Docs Link: <https://www.bitget.com/api-doc/uta/account/Transfer>

List the coins transferable between two account types.

Exec: `./bitget-cli wallet transferable-coins --fromType=spot --toType=usdt_futures`
```shell
┌──────┐
│ COIN │
├──────┤
│ USDT │
│ BGB  │
└──────┘
```

## deposit - address
Docs Link: <https://www.bitget.com/api-doc/uta/account/Deposit>

Get the on-chain deposit address for a coin, optionally on a chain.

Exec: `./bitget-cli wallet deposit address --coin=USDT [--chain=trc20]`
```shell
┌──────┬───────┬────────────────────────────────────────────┬─────┐
│ COIN │ CHAIN │                  ADDRESS                   │ TAG │
├──────┼───────┼────────────────────────────────────────────┼─────┤
│ USDT │ ERC20 │ 0x88c379e744c0c297b9110d70d706bd4545f0542f │     │
└──────┴───────┴────────────────────────────────────────────┴─────┘
```

## deposit - records
Docs Link: <https://www.bitget.com/api-doc/uta/account/Deposit>

List deposit records. The time window defaults to the last 30 days.

Exec: `./bitget-cli wallet deposit records [--coin=USDT] [--limit=20]`

**Supported parameters:**
- `--coin, -c`: coin filter
- `--startTime, -a` / `--endTime, -e`: unix ms or "YYYY-MM-DD HH:MM:SS"
- `--limit, -l`: max records

## withdraw - create
Docs Link: <https://www.bitget.com/api-doc/uta/account/Withdrawal>

Submit a withdrawal (on-chain or internal). Use with care.

**On-chain:**
```shell
./bitget-cli wallet withdraw create --coin=USDT --transferType=on_chain --chain=trc20 --address=Txxxx --size=10
```
**Internal (by UID):**
```shell
./bitget-cli wallet withdraw create --coin=USDT --transferType=internal_transfer --address=123456 --size=10
```
**Supported parameters:**
- `--coin, -c`: coin (required)
- `--transferType, -T`: on_chain or internal_transfer (required)
- `--address, -d`: destination address / UID / email / mobile (required)
- `--size, -m`: amount, decimal (required)
- `--chain`: chain, e.g. trc20 (required for on-chain)
- `--tag`: address tag/memo
- `--innerToType`: internal address type — uid, email, mobile
- `--remark`, `--clientOid`

## withdraw - records
Docs Link: <https://www.bitget.com/api-doc/uta/account/Withdrawal>

List withdrawal records. The time window defaults to the last 30 days.

Exec: `./bitget-cli wallet withdraw records [--coin=USDT] [--limit=20]`

## withdraw - address
Docs Link: <https://www.bitget.com/api-doc/uta/account/Withdrawal>

List the saved withdrawal address book entries.

Exec: `./bitget-cli wallet withdraw address [--coin=USDT] [--type=EVM]`

## withdraw - cancel
Docs Link: <https://www.bitget.com/api-doc/uta/account/Withdrawal>

Cancel a withdrawal still within its cooling-off period.

Exec: `./bitget-cli wallet withdraw cancel --orderId=xxx`

Identify by `--orderId` or `--clientOid`.
