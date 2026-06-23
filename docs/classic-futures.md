# Classic Futures Module

> Classic-account futures commands (v2 `/api/v2/mix/*`). For the unified account
> use `./bitget-cli UTA futures` instead.
>
> Every command accepts a global `--json` flag that prints the raw API response
> as indented JSON instead of a table.

Exec: `./bitget-cli futures [--product=usdt-futures] [Subcommand] [Arguments]`

## Product
A persistent `--product` / `-P` flag selects the product line (default
`usdt-futures`):

| alias | product       |
|-------|---------------|
| `usdt` / `usdt-futures` | USDT-FUTURES |
| `coin` / `coin-futures` | COIN-FUTURES |
| `usdc` / `usdc-futures` | USDC-FUTURES |

## Quick Navigation
- [account](#account---balances--equity)
- [health](#health---account-health--risk)
- [position list](#position-list)
- [order create](#order-create)
- [order cancel](#order-cancel)
- [order get](#order-get)
- [order open](#order-open)

## account - balances & equity
Docs Link: <https://www.bitget.com/api-doc/contract/account/Get-Account-List>

Shows the classic futures account balances and equity per margin coin.

Exec: `./bitget-cli futures account [--product=usdt-futures]`

## health - account health & risk
Docs Link: <https://www.bitget.com/api-doc/contract/account/Get-Account-List>

Shows the account-health / risk picture per margin coin: equity, crossed risk
rate, maintenance margin and unrealised PnL. A higher crossed risk rate means a
higher risk of liquidation.

Exec: `./bitget-cli futures health [--product=usdt-futures]`

## position list
Docs Link: <https://www.bitget.com/api-doc/contract/position/get-all-position>

List the account's open futures positions in the product line, optionally
filtered by `--marginCoin` and/or `--symbol`.

```shell
./bitget-cli futures position list [--marginCoin=USDT] [--symbol=BTCUSDT]
```

## order create
Docs Link: <https://www.bitget.com/api-doc/contract/trade/Place-Order>

Place a new futures order. `--price` is required for limit orders. `--marginCoin`
defaults to USDT; set it for coin/usdc lines. `--tradeSide` (open/close) and
`--reduceOnly` apply in hedge mode.

```shell
./bitget-cli futures order create --symbol=BTCUSDT --side=buy --type=limit --size=0.001 --price=50000
./bitget-cli futures --product=usdc-futures order create --symbol=BTCPERP --side=sell --type=market --size=0.001 --marginCoin=USDC
```
Flags: `--symbol`, `--side` (buy/sell), `--type` (limit/market), `--size`,
`--price`, `--marginCoin` (default USDT), `--marginMode` (crossed/isolated),
`--force`, `--tradeSide`, `--reduceOnly`, `--clientOid`.

## order cancel
Docs Link: <https://www.bitget.com/api-doc/contract/trade/Cancel-Order>

Cancel a single futures order by `--orderId` or `--clientOid` (`--symbol` required).

```shell
./bitget-cli futures order cancel --symbol=BTCUSDT --orderId=xxx
```

## order get
Docs Link: <https://www.bitget.com/api-doc/contract/trade/Get-Order-Details>

Query a single futures order by `--orderId` or `--clientOid` (`--symbol` required).

```shell
./bitget-cli futures order get --symbol=BTCUSDT --orderId=xxx
```

## order open
Docs Link: <https://www.bitget.com/api-doc/contract/trade/Get-Orders-Pending>

List currently open (unfilled / partially filled) futures orders.

```shell
./bitget-cli futures order open [--symbol=BTCUSDT] [--limit=20]
```
