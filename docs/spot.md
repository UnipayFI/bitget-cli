# Spot Module

> Every command in this module accepts a global `--json` flag that prints the
> raw API response as indented JSON instead of a table.

All spot commands target the `SPOT` category.

Exec: `./bitget-cli UTA spot [Subcommand] [Arguments]`

## Quick Navigation
- [order create](#order---create)
- [order cancel](#order---cancel)
- [order modify](#order---modify)
- [order get](#order---get)
- [order open](#order---open)
- [order history](#order---history)
- [order cancel-all](#order---cancel-all)
- [fills](#fills)

## order - create
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Place-Order>

Place a new spot order. The created order's identifiers are printed back.

**Limit order:**
```shell
./bitget-cli UTA spot order create --symbol=BTCUSDT --side=buy --type=limit --qty=0.0001 --price=30000
```
**Market order:**
```shell
./bitget-cli UTA spot order create --symbol=ETHUSDT --side=sell --type=market --qty=0.01
```
```shell
┌─────────────────────┬─────────────────────┐
│      ORDER ID       │     CLIENT OID      │
├─────────────────────┼─────────────────────┤
│ 1452849501609095168 │ 1452849501609095169 │
└─────────────────────┴─────────────────────┘
```
**Supported parameters:**
- `--symbol, -s`: trading pair symbol (required)
- `--side, -S`: buy or sell (required)
- `--type, -t`: limit or market (required)
- `--qty, -q`: order quantity, decimal (required)
- `--price, -p`: order price, decimal (required for limit)
- `--tif, -T`: time in force — gtc, post_only, fok, ioc
- `--clientOid`: client order id

## order - cancel
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Cancel-Order>

Cancel a single spot order by `--orderId` or `--clientOid`.
```shell
./bitget-cli UTA spot order cancel --orderId=1452849501609095168
```

## order - modify
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Modify-Order>

Amend a spot order's quantity and/or price.
```shell
./bitget-cli UTA spot order modify --orderId=xxx --price=31000 --qty=0.0002
```
Identify the order by `--orderId` or `--clientOid`; supply at least one of
`--qty` / `--price`.

## order - get
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Get-Order-Details>

```shell
./bitget-cli UTA spot order get --orderId=1452849501609095168
```

## order - open
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Get-Order-Pending>

List currently open (unfilled / partially filled) spot orders.

Exec: `./bitget-cli UTA spot order open [--symbol=BTCUSDT]`
```shell
┌─────────────────────┬─────────┬──────┬───────┬────────┬───────┬────────┬────────┬───────────┬─────┬──────────┬────────┬─────────────────────┐
│      ORDER ID       │ SYMBOL  │ SIDE │ TYPE  │ STATUS │ PRICE │  QTY   │ FILLED │ AVG PRICE │ TIF │ POS SIDE │ REDUCE │       CREATED       │
├─────────────────────┼─────────┼──────┼───────┼────────┼───────┼────────┼────────┼───────────┼─────┼──────────┼────────┼─────────────────────┤
│ 1452849501609095168 │ BTCUSDT │ buy  │ limit │ new    │ 30000 │ 0.0001 │ 0      │ 0         │ gtc │          │ NO     │ 2026-06-22 07:10:34 │
└─────────────────────┴─────────┴──────┴───────┴────────┴───────┴────────┴────────┴───────────┴─────┴──────────┴────────┴─────────────────────┘
```

## order - history
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Get-Order-History>

List historical spot orders (90-day access window).

Exec: `./bitget-cli UTA spot order history [--symbol=BTCUSDT] [--limit=20]`

**Supported parameters:**
- `--symbol, -s`: symbol filter
- `--startTime, -a` / `--endTime, -e`: unix ms or "YYYY-MM-DD HH:MM:SS"
- `--limit, -l`: max records

## order - cancel-all
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Cancel-All-Order>

Cancel all open spot orders, optionally limited to one `--symbol`. The response
lists each attempted cancellation with its per-order result code.

Exec: `./bitget-cli UTA spot order cancel-all [--symbol=BTCUSDT]`

## fills
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Get-Order-Fills>

List spot trade fills (90-day access window).

Exec: `./bitget-cli UTA spot fills [--orderId=xxx] [--limit=20]`

**Supported parameters:**
- `--orderId, -i`: order id filter
- `--startTime, -a` / `--endTime, -e`: unix ms or "YYYY-MM-DD HH:MM:SS"
- `--limit, -l`: max records
