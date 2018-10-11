# Translation API 2

Translation API written in Go.

## Routes

You can check out the different routes of the API [here](ROUTES.md).

## Running

For running the api a database is necessary. If you don't specify a database, it will create a sqlite database in `/tmp`.
To configure the database connection, the following environment variables are used:

* DATABASE_TYPE (sqlite3, mysql or postgres)
* DATABASE_HOST
* DATABASE_PORT
* DATABASE_NAME (if database type is sqlite3 this is the path to the database file)
* DATABASE_USER
* DATABASE_PASSWORD
* DATABASE_SSL (for mysql and postgres connection, default is disable/false)
* FRONTEND_URL (necessary for cors, default http://localhost:3000 for local development)

Easiest way of deploying is running it in a container. 
For building a container image you can have a look in the [gitlab-ci.yml](gitlab-ci.yml).