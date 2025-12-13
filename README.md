# Gator

## [boot.dev](https://www.boot.dev) guided project.

### Installation

This is instruction only for macOS(sorry).
If you have something already installed then just skip the step.

1. Install **Go**. Use this [link](https://go.dev/doc/install) or
```bash
brew install go
```

2.Install **PostgreSQL 15**
```bash
brew install postgresql@15
```
This will install [PostgreSQL](https://www.postgresql.org)
After installation try to run `psql --version`, if it says something like "command not found", then you have to put it into your path with this:

If you are using **zsh**
```bash
'export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"' >> 
~/.zshrc
```
Or this, if you are using **bash**
```bash
'export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"' >> 
~/.bash_profile
```
When it's done, start the postgres server with:
```bash
brew services start postgresql@15
```

3. Create database **gator**.
```bash
psql postgres
CREATE DATABASE gator;
\q
```
Cool, you have a db now.
To confirm run `psql -l` and you should see `gator` among others.

4.Install [goose](https://github.com/pressly/goose).
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```
or
```bash
brew install goose
```
This will be used to do all the migrations.

5. Put your **DB** url into `.env`

Before calling migrations you have to put your DB url inside `.env` file.
Currently it holds placeholder: `GOOSE_DBSTRING=postgres://[username]@localhost:5432/gator`
You have to replace `[username]` with your actual user name. If for some reason you don't know it run:
```bash
users
```
The output will be user name that you have to use in this step.

6. Run migrations.
```bash
goose up
```
You can confirm that everything worked as intented with `goose status`.

7. Finally you can **compile/install** this app for usage:
```bash
go build .
```
or
```bash
go install .
```
`build` will create a binary file inside current directory, you can call it with `./gator`.

`install` will create a binary file inside go default directory. It should be `~/go/bin/`. You can call with `gator` from anywhere.

### Usage commands

#### `login username`

#### `register username`

#### `reset`
This is basically for testing, but you can reset all users and everything that relates to them from `gator` database.

#### `users`
List all users in the db showing who are you logged in as.

#### `agg time_interval`
Run aggregations to fetch all posts for saved feeds in the database. It cares only about saved feeds and url that it provides for this operation.

#### `addfeed name url`
Add new feed providing its name and url, so it can later be aggregated. Currently logged user will be subscribed to this feed.

#### `feeds`
Display all saved feeds in the database.

#### `follow url`
Follow (subscribe) to the given feed. Worth noting: the feed should exist and is not automatically created.

#### `following`
Show which feeds you are currently following.

#### `unfollow url`
Unfollow the feed by given url.

#### `browse [limit]`
Show all posts that are saved after `agg` command has happened and only for feeds that you follow.

### Migrations

If for some reason you want to migrate down or up, you can use these commands:
```bash
goose up
goose down
```
This will ask goose to migrate up or down.

### Uninstallation

```bash
./uninstall.sh
```
This will remove the binary from your machine.

