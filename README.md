# Gator - Blog Aggregator

## What is Gator?
**Gator** is the Blog Aggregator App that works on the Command Line Interface(CLI) and is written in _**Go**_. It collects blog posts from RSS Feeds that are specified by the user which is stored inside a PostgreSQL Database.

## Pre-Requisites
You must have the Go Toolchain and PostgreSQL installed for building the binary executable as well as well as the database to maintain our data.

The installation instructions can be found here:
* [Go Installation Instructions](https://go.dev/doc/install)

For **PostgreSQL**:

**macOS** with brew: 
```
brew install postgresql@15
```

**Linux / WSL (Debian)**: 
   ```
   sudo apt update
   sudo apt install postgresql postgresql-contrib
   ```

Ensure the installation worked. Make sure you're on version 15+ of Postgres:
```
psql --version
```
(Linux / WSL only) Update postgres password:
```
sudo passwd postgres
```
**Start the Postgres Server in the background:**
```
brew services start postgresql@15 (for MacOS)
sudo service postgresql start     (for Linux)
```

## How To Install Gator
* Clone the Github Repo using the following command:
    ```
    git clone https://github.com/QueenQ14/gator.git
    ```
* Navigate inside the directory and run the following command:
    ```
    go install
    ```
* **Alternatively** if you do not wish to download the source code:
    ```
    go install https://github.com/QueenQ14/gator
    ```

## Configuration

We are almost done but we need to setup just a few more things before we can start using `gator`.
Firstly, we will need a DB Connection String which is in the following format:
```
protocol://username:password@host:port/database
```
For a local setup you can use the following example:
```
postgres://username:@localhost:5432/gator
```
* We will need to create the actual Database, to do that we will require `goose` to apply the migration schemas in the `/sql/schemas/` . Install `goose` using the following command:
    ```
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```
    We will apply the migrations using the following commands:
    ```
    cd sql/schema
    goose postgres '<connection_string>' up
    ```
* In your HOME Directory, create a `~/.gatorconfig.json` with the following value, this will allow `gator` to interact with this database(SSL Mode is disabled in local setup):
    ```
    {"db_url":"postgres://<username>:@localhost:5432/gator?sslmode=disable"}
    ```

## Functionality

Here are the comprehensive list of commands that are available in Gator:
* `gator register <user>`: Registers the user and sets as current user.
* `gator login <user>`: Logs in as the user(if they exist).
* `gator users` : Lists all the users available.
* `gator addfeed <name> <url>`: Adds the feed as the current user.
* `gator delfeed <url>`: Deletes the feed entry for all users.
* `gator follow <feed_url>`: Current user will now follow the feed.
* `gator following`: Lists all the feeds that the current user is following.
* `gator unfollow <feed_url>`: Unfollow from specified feed URL.
* `gator feeds`: Lists all the Feeds that are available.
* `gator browse [limit]`: Browse all the posts from the feeds that the user is following. Limit is an optional argument(default is 3) that limits how many posts are shown at a time in a page. You can use the `n`, `p` and `q` commands to move next, previous across the pages and exit respectively.
* `gator reset`: Resets the user database; use this only if you want to clear ALL data.
* `gator agg <time> [num of workers]`: Aggregates posts from the following feeds at given frequency. Time can be specified in the formats of `s/m/h`. Number of workers are defined for concurrency on aggregation. Default value is 3.

## Credits
This project was completed as part of the Backend Development Path on [Boot.Dev](https://www.boot.dev).
If you want to build this project by yourself while getting amazing guidance which improves your coding check the website out!