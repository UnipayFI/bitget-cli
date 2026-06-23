# Classic Spot Module

> Classic-account spot commands (v2 `/api/v2/spot/*`). For the unified account
> use `./bitget-cli UTA spot` instead.
>
> Every command accepts a global `--json` flag that prints the raw API response
> as indented JSON instead of a table.

Exec: `./bitget-cli spot [Subcommand] [Arguments]`

## Quick Navigation
- [account info](#account-info---identity--permissions)
- [account assets](#account-assets---per-coin-balances)
- [order create](#order-create)
- [order cancel](#order-cancel)
- [order get](#order-get)
- [order open](#order-open)

## account info - identity & permissions
Docs Link: <https://www.bitget.com/api-doc/spot/account/Get-Account-Info>

Shows the classic spot account's identity and API permission metadata.

Exec: `./bitget-cli spot account info`

## account assets - per-coin balances
Docs Link: <https://www.bitget.com/api-doc/spot/account/Get-Account-Assets>

Shows the classic spot account's per-coin balances (non-zero only). Optionally
filter to one coin.

Exec: `./bitget-cli spot account assets [--coin=USDT]`

## order create
Docs Link: <https://www.bitget.com/api-doc/spot/trade/Place-Order>

Place a new spot order. `--price` is required for limit orders. For market buy
orders, `--size` is denominated in the quote currency (e.g. USDT).

```shell
./bitget-cli spot order create --symbol=BTCUSDT --side=buy --type=limit --size=0.0001 --price=30000
./bitget-cli spot order create --symbol=ETHUSDT --side=sell --type=market --size=0.01
```
Flags: `--symbol`, `--side` (buy/sell), `--type` (limit/market), `--size`,
`--price`, `--force` (gtc/post_only/fok/ioc, default gtc), `--clientOid`.

## order cancel
Docs Link: <https://www.bitget.com/api-doc/spot/trade/Cancel-Order>

Cancel a single spot order by `--orderId` or `--clientOid` (`--symbol` required).

```shell
./bitget-cli spot order cancel --symbol=BTCUSDT --orderId=1452849501609095168
```

## order get
Docs Link: <https://www.bitget.com/api-doc/spot/trade/Get-Order-Info>

Query a single spot order by `--orderId` or `--clientOid`.

```shell
./bitget-cli spot order get --orderId=1452849501609095168
```

## order open
Docs Link: <https://www.bitget.com/api-doc/spot/trade/Get-Unfilled-Orders>

List currently open (unfilled / partially filled) spot orders.

```shell
./bitget-cli spot order open [--symbol=BTCUSDT] [--limit=20]
```
