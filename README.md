# Gator 🐊

Gator is a custom command-line interface (CLI) tool built in Go for managing and aggregating RSS feeds directly from your terminal.

## Prerequisites

Before you can install and run Gator, you will need to have the following installed on your system:

* **[Go](https://go.dev/doc/install)**: The programming language used to build this project.
* **[PostgreSQL](https://www.postgresql.org/download/)**: The relational database used to store users, feeds, and posts.

## Installation

Because Go programs compile into static binaries, you don't need to run `go run .` every time you want to use the app in production. You can install it globally on your machine using the `go install` command!

Run the following command in your terminal:

    go install github.com/matavak/gator@latest


This will compile the program and place the executable `gator` binary into your Go bin directory (usually `~/go/bin`). Make sure this directory is added to your system's `$PATH` so you can run the `gator` command from anywhere!

## Configuration

Before running Gator, you need to set up your configuration file so the CLI knows how to connect to your PostgreSQL database.

1. Create a file named `.gatorconfig.json` in your home directory (e.g., `~/.gatorconfig.json`).
2. Add your database connection string to the file like this:

    {
      "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
    }

*(Update the username, password, and port if your local Postgres setup is different).*

## Usage

Once installed and configured, you can use the `gator` command. Here are a few essential commands to get you started:

**User Management:**
* `gator register <username>`: Registers a new user in the database.
* `gator login <username>`: Logs in as an existing user.
* `gator users`: Lists all registered users.

**Feed Management:**
* `gator add <name> <url>`: Adds a new RSS feed to the database for the logged-in user.
* `gator feeds`: Lists all available RSS feeds.
* `gator follow <url>`: Follows an existing RSS feed.
* `gator following`: Lists all feeds the current user is following.

**Aggregation:**
* `gator agg <time_between_requests>`: Starts the continuous aggregation server to fetch new posts (e.g., `gator agg 1m`).

