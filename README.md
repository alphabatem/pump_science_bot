# pump_science_bot

Simple Pump Science bot to monitor changes on data via their API. 

## Environment Setup
* Create `.env` file
* Set `RPC_URL` to a solana rpc provider
* Configure other settings

```bash
RPC_URL=""
VALIDATOR_TIP="0.003"
SLIPPAGE=5
IN_AMOUNT=0.1
WALLET="./wallet_keypair.json"

PUMP_SCIENCE_KEY="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im5wY25uaHBxdHFqcWxnd3Fxb2ljIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MzE3ODUzNDgsImV4cCI6MjA0NzM2MTM0OH0.0MbSmni2Avr5IIOBZM9JMvm41E_71qqZC5OiCE_JUQY"
```


## Compile
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w'  ./runtime/main.go
```

## Run
```bash
./main
```

OR run from source:
```bash
go run ./runtime/main.go
```