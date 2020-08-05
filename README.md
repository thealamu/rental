# Rental
A car rental app using [Flutterwave](https://flutterwave.com/)'s payment technologies

![Go](https://github.com/thealamu/rental/workflows/Go/badge.svg)

**NOTE: This repo is for the backend API and its implementation alone**

## API
The [OpenAPI](http://spec.openapis.org/oas/v3.0.3) definition can be found in the [oas](https://github.com/thealamu/rental/tree/master/oas) directory. Please note that the definition is a WIP, hence, the document is susceptible to change.

## Proposed Features
- [ ] Customer and Merchant Accounts
- [ ] Merchant Resource(Car) management
- [ ] Resource(Car) provisioning for customers
- [ ] Customer payment methods:
    - [ ] Card
    - [ ] Momo
    - [ ] Alternative / PreAuth
- [ ] Merchant payout methods:
    - [ ] Transfers(Payout Bulk)
    - [ ] Refunds
- [ ] Merchant payment-aniticipation webhook notification

## Deployment
### Auth
Some endpoints are protected, to allow users log in, create an Auth0 account, register [an application](https://manage.auth0.com/#/applications) and obtain your domain, client ID and client Secret (It's free). Export them in your environment as RTL_DOMAIN, RTL_CLIENT_ID and RTL_CLIENT_SECRET respectively.
```bash
export RTL_DOMAIN={YOUR_DOMAIN}
export RTL_CLIENT_ID={YOUR_CLIENT_ID}
export RTL_CLIENT_SECRET={YOUR_CLIENT_SECRET}
```
For more info, please see https://auth0.com/docs/quickstart/webapp/golang#configure-auth0
### Sessions
Export an environment variable RTL_STOREKEY for the session store, this helps keep your users logged in.
```bash
export RTL_STOREKEY={YOUR_STOREKEY}
```
### Database
This app uses a MySQL store, to setup your local instance, see https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/
Export a mysql username, password and database using the RTL_USER, RTL_PASS and RTL_DB environment variables.
```bash
export RTL_USER={MYSQL_USER}
export RTL_PASS={MYSQL_PASS}
export RTL_DB={MYSQL_DB}
```
### Run :rocket:
To start the server:
```bash
go build -o rental .
./rental
```
To start the server on a different port, say 1024:
```bash
./rental --port 1024
```

## License
[MIT](https://github.com/thealamu/rental/blob/master/LICENSE)