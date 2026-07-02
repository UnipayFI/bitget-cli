# Futures Module

> Every command in this module accepts a global `--json` flag that prints the
> raw API response as indented JSON instead of a table.

The product line is selected with the persistent `--category` / `-C` flag
(default `usdt-futures`); it applies to every subcommand:
`usdt-futures` (`usdt`), `coin-futures` (`coin`), `usdc-futures` (`usdc`).

Exec: `./bitget-cli UTA futures [--category=usdt-futures] [Subcommand] [Arguments]`

## Quick Navigation
- [order create](#order---create)
- [order cancel](#order---cancel)
- [order modify](#order---modify)
- [order get](#order---get)
- [order open](#order---open)
- [order history](#order---history)
- [order cancel-all](#order---cancel-all)
- [position list](#position---list)
- [position history](#position---history)
- [position adl-rank](#position---adl-rank)
- [position close](#position---close)
- [fills](#fills)

## order - create
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Place-Order>

Place a new futures order. The created order's identifiers are printed back.

**Limit, hedge mode:**
```shell
./bitget-cli UTA futures order create --symbol=BTCUSDT --side=buy --type=limit --qty=0.001 --price=50000 --posSide=long
```
**Market, USDC futures:**
```shell
./bitget-cli UTA futures --category=usdc-futures order create --symbol=BTCPERP --side=sell --type=market --qty=0.001
```
**Supported parameters:**
- `--symbol, -s`: symbol (required)
- `--side, -S`: buy or sell (required)
- `--type, -t`: limit or market (required)
- `--qty, -q`: order quantity, decimal (required)
- `--price, -p`: order price, decimal (required for limit)
- `--tif, -T`: gtc, post_only, fok, ioc
- `--posSide, -P`: long or short (required in hedge mode)
- `--reduceOnly, -r`: yes or no
- `--marginMode, -m`: crossed or isolated
- `--clientOid`: client order id

## order - cancel
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Cancel-Order>
```shell
./bitget-cli UTA futures order cancel --orderId=xxx
```

## order - modify
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Modify-Order>
```shell
./bitget-cli UTA futures order modify --orderId=xxx --price=51000
```

## order - get
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Get-Order-Details>
```shell
./bitget-cli UTA futures order get --orderId=xxx
```

## order - open
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Get-Order-Pending>

Exec: `./bitget-cli UTA futures order open [--symbol=BTCUSDT]`

Columns: Order ID, Symbol, Side, Type, Status, Price, Qty, Filled, Avg Price,
TIF, Pos Side, Reduce, Created.

## order - history
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Get-Order-History>

Exec: `./bitget-cli UTA futures order history [--symbol=BTCUSDT] [--limit=20]`

## order - cancel-all
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Cancel-All-Order>

Cancel all open orders in the category, optionally limited to one `--symbol`.

Exec: `./bitget-cli UTA futures order cancel-all [--symbol=BTCUSDT]`

## position - list
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Get-Position>

List the account's open positions in the category.

Exec: `./bitget-cli UTA futures position list [--symbol=BTCUSDT] [--posSide=long]`
```shell
┌────────┬──────────┬─────────────┬──────────┬───────┬───────────┬───────────┬────────────┬───────────┬────────────────┬─────────────┬─────────────┐
│ SYMBOL │ POS SIDE │ MARGIN MODE │ LEVERAGE │ TOTAL │ AVAILABLE │ AVG PRICE │ MARK PRICE │ LIQ PRICE │ UNREALISED PNL │ PROFIT RATE │ MARGIN COIN │
└────────┴──────────┴─────────────┴──────────┴───────┴───────────┴───────────┴────────────┴───────────┴────────────────┴─────────────┴─────────────┘
```

## position - history
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Get-Position-History>

List closed/historical positions (90-day window).

Exec: `./bitget-cli UTA futures position history [--symbol=BTCUSDT] [--limit=20]`

Columns: Symbol, Pos Side, Open Avg, Close Avg, Open Qty, Close Qty,
Realised PNL, Net Profit, Funding, Created, Updated.

## position - adl-rank
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Get-Position-ADL-Rank>

Show the auto-deleveraging (ADL) rank for each open position.

Exec: `./bitget-cli UTA futures position adl-rank`

## position - close
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Close-All-Positions>

Market-close positions. Without `--symbol` closes all in the category; without
`--posSide` closes both sides.

Exec: `./bitget-cli UTA futures position close [--symbol=BTCUSDT] [--posSide=long]`

## fills
Docs Link: <https://www.bitget.com/api-doc/uta/trade/Get-Order-Fills>

List futures trade fills (90-day access window).

Exec: `./bitget-cli UTA futures fills [--orderId=xxx] [--limit=3]`
```shell
┌─────────────────────┬─────────────────────┬─────────┬──────┬────────────┬──────────┬────────────┬───────┬─────────────────┬────────┬─────────────────────┐
│       EXEC ID       │      ORDER ID       │ SYMBOL  │ SIDE │ EXEC PRICE │ EXEC QTY │ EXEC VALUE │ SCOPE │       FEE       │  PNL   │        TIME         │
├─────────────────────┼─────────────────────┼─────────┼──────┼────────────┼──────────┼────────────┼───────┼─────────────────┼────────┼─────────────────────┤
│ 1452849335278223360 │ 1452849335267192832 │ ETHUSDT │ sell │ 1740.42    │ 0.01     │ 17.4042    │ taker │ 0.01044252 USDT │ 0.0019 │ 2026-06-22 07:09:54 │
└─────────────────────┴─────────────────────┴─────────┴──────┴────────────┴──────────┴────────────┴───────┴─────────────────┴────────┴─────────────────────┘
```

**Supported parameters:**
- `--orderId, -i`: order id filter
- `--startTime, -a` / `--endTime, -e`: unix ms or "YYYY-MM-DD HH:MM:SS"
- `--limit, -l`: max records
