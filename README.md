# Aggregator CLI

## Prerequisites

To run this application, you need to have **Go** and **PostgreSQL** installed on your local machine.

## Installation

Install the aggregator CLI using the following command:

```sh
go install github.com/MrR0b0t/aggregator@latest
```

Once installed, you’ll be able to use various commands to interact with the system.

## Running Commands

To run a command, simply open your terminal and enter:

```sh
gator <command> [arguments]
```

For example, to register a new user, you would run:

```sh
gator register <username>
```

## Available Commands

- **login <username>** – Authenticate a user
- **register <username>** – Register a new user
- **reset** – Reset user password
- **users** – List all users
- **agg <timer>** – Aggregate data (e.g., `30s`, `1m`, `45s`, `1h`)
- **feeds** – List available feeds
- **addfeed <title> <url>** – Add a new feed (requires login)
- **follow <url>** – Follow a feed (requires login)
- **following** – List followed feeds (requires login)
- **unfollow <url>** – Unfollow a feed (requires login)
- **browse [limit]** – Browse content (requires login, default limit = 2)

There are plenty of additional features to explore—happy coding! 🚀
