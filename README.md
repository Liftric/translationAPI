# Translation API 2

Translation API written in Go.

## Routes

You can check out the different routes of the API [here](ROUTES.md).

## Running

This project is configured to run inside docker. The database connection can be configured using the following environment variables:
* DATABASE_TYPE (mysql or postgres)
* DATABASE_HOST
* DATABASE_PORT
* DATABASE_NAME
* DATABASE_USER
* DATABASE_PASSWORD

There should be a docker image in the gitlab registry.