# Translation API 2

Translation API written in Go.

## Routes

You can check out the different routes of the API [here](ROUTES.md).

## Running

This project is configured to run inside docker. The database connection can be configured using the following environment variables:
* DATABASE_TYPE (sqlite3, mysql or postgres)
* DATABASE_HOST
* DATABASE_PORT
* DATABASE_NAME (if database type is sqlite3 this is the path to the database file)
* DATABASE_USER
* DATABASE_PASSWORD
* DATABASE_SSL (for mysql and postgres connection, default is disabled/false)

Easiest way of running is using a docker image.