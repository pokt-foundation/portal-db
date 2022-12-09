<div align="center">
    <img src=".github/banner.png" alt="Pocket Network logo" width="600"/>
    <h1>Portal API Postgres Driver</h1>
    <big>Database driver and repository struct definitions for use with the Portal API</big>
    <div>
    <br/>
        <a href="https://github.com/pokt-foundation/node-nanny/pulse"><img src="https://img.shields.io/github/last-commit/pokt-foundation/node-nanny.svg"/></a>
        <a href="https://github.com/pokt-foundation/node-nanny/pulls"><img src="https://img.shields.io/github/issues-pr/pokt-foundation/node-nanny.svg"/></a>
        <a href="https://github.com/pokt-foundation/node-nanny/issues"><img src="https://img.shields.io/github/issues-closed/pokt-foundation/node-nanny.svg"/></a>
    </div>
</div>
<br/>

# Modules

## Postgres Driver

Used to interact with the Postgres database. Reads and writes are done using the repository structs.

## Repository

Contains the structs and methods used across the Portal backend Go repos.

# Development

## Packages in Use

- [SQLC](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html) - Generates idiomatic type-safe Go code from SQL schema & queries.
- [Mockery](https://github.com/vektra/mockery) - Generates mock code for all interfaces in the code for testing purposes.

## Generating code

**Before committing any code to the repo, run the default Make target (`make`)**

This will generate SQLC code. This is a useful way to check the database `schema.sql` and `query.sql` files for SQL errors.

It will also generate as well a mock of the `IPostgresDriver` interface for testing purposes. This mock will automatically reflect changes made to the SQL schema files.

## Pre-Commit Installation

Before starting development work on this repo, `pre-commit` must be installed. In order to do so, run the command **`make init-pre-commit`** from the repository root.

Once this is done, the following commands will be performed on every commit to the repo and must pass before the commit is allowed:

Basic checks

- **check-yaml** - Checks YAML files for errors
- **check-merge-conflict** - Ensures there are no merge conflict markers
- **end-of-file-fixer** - Adds a newline to end of files
- **trailing-whitespace** - Trims trailing whitespace
- **no-commit-to-branch** - Ensures commits are not made directly to `main`

Go-specific checks

- **go-fmt** - Runs `gofmt`
- **go-imports** - Runs `goimports`
- **golangci-lint** - run `golangci-lint run ./...`
- **go-critic** - run `gocritic check ./...`
- **go-build** - run `go build`
- **go-mod-tidy** - run `go mod tidy -v`
