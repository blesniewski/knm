# Readme

A simple go app for fetching OpenExchangeRates rates for a set of given currencies
Some mocked 'crypto exchange functionality' as well

## Usage

Two endpoints:

- `GET /rates?currencies=USD,GBP,EUR`: Requires at least two currencies available in the openexchangerates '/latest' api
- `GET /exchange?from=WBTC&to=USDT&amount=1.0`, Requires all 3 parameters, handles only hardcoded cryptos

## Running it

### Running locally

- Requires a OPENEXCHANGERATES_APP_ID env variable with APP ID for the OXR API
- Listens on port 8080, not configurable right now (should be if this was going to be something else than a simple project)

### Running in a docker container:

`docker build . -t kryptonim-app`

`docker run --rm --name kryptonim-app --publish 8080:8080 --env OPENEXCHANGERATES_APP_ID=<your_app_id> kryptonim-app`
