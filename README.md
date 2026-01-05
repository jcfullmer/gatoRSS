# GatoRSS - RSS Feed Aggregator CLI - my take on the gator cli in the [Boot.dev Blog Aggregator in Go Course](https://www.boot.dev/courses/build-blog-aggregator-golang)

A command-line RSS feed aggregator built with Go and PostgreSQL.

## Prerequisites

Before you can run gatoRSS, you'll need to have the following installed:

### Go (1.19 or higher)
- **Official Installation Guide**: https://go.dev/doc/install
- **Quick Install**:
  - macOS: `brew install go`
  - Linux: Download from https://go.dev/dl/
  - Windows: Download installer from https://go.dev/dl/

### PostgreSQL
- **Official Installation Guide**: https://www.postgresql.org/download/
- **Quick Install**:
  - macOS: `brew install postgresql`
  - Linux (Ubuntu/Debian): `sudo apt-get install postgresql postgresql-contrib`
  - Windows: Download installer from https://www.postgresql.org/download/windows/

After installing PostgreSQL, make sure the service is running:
```bash
# macOS/Linux
sudo service postgresql start

# Or check status
sudo service postgresql status

## Installation

Install the GatoRSS CLI using Go:

```bash
go install github.com/jcfullmer/gatoRSS@latest
```

Make sure your `$GOPATH/bin` is in your system's PATH. You can verify the installation by running:

```bash
gatoRSS
```

## Configuration

Create a configuration file at `.gatorconfig.json`
Windows: "C:Users:<current_user>/.gatorconfig.json"
Linux: "/home/<user>/.gaatorconfig.json"
MacOS: "/Users/<current_user>"
with the following structure:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": "",
  "current_user_id": "" 
}
```

### Configuration Fields

- **db_url**: Your PostgreSQL connection string
  - Format: `postgres://username:password@host:port/database?sslmode=disable`
  - Replace `username` and `password` with your PostgreSQL credentials
  - The database name (`gator` in the example) should match your database
- **current_user_name & current_user_id**: Will be set automatically when you register/login

### Setting up the Database

Create the database in PostgreSQL:

```bash
# Connect to PostgreSQL
psql postgres

# Create the database
CREATE DATABASE gator;

# Exit psql
\q
```

## Usage

### User Management

#### Register a new user
```bash
gatoRSS register <username>
```

#### Login as a user
```bash
gatoRSS login <username>
```

#### List all users
```bash
gatoRSS users
```

### Feed Management

#### Add a new RSS feed
```bash
gatoRSS addfeed <feed_name> <feed_url>
```

Example:
```bash
gatoRSS addfeed "TechCrunch" https://techcrunch.com/feed/
```

#### List all feeds
```bash
gatoRSS feeds
```

#### Follow a feed
```bash
gatoRSS follow <feed_url>
```

#### Unfollow a feed
```bash
gatoRSS unfollow <feed_url>
```

#### List feeds you're following
```bash
gatoRSS following
```

### Reading Posts

#### Browse recent posts
```bash
gatoRSS browse [limit]
```

Example:
```bash
gatoRSS browse 10
```

#### Aggregate feeds (fetch new posts)
```bash
gatoRSS agg <time_between_requests>
```

Example (fetch every 60 seconds):
```bash
gatoRSS agg 60s
```

## Commands Reference

| Command | Description | Usage |
|---------|-------------|-------|
| `register` | Create a new user account | `gatoRSS register <username>` |
| `login` | Switch to a different user | `gatoRSS login <username>` |
| `users` | List all registered users | `gatoRSS users` |
| `addfeed` | Add a new RSS feed | `gatoRSS addfeed <name> <url>` |
| `feeds` | List all available feeds | `gator feeds` |
| `follow` | Follow an RSS feed | `gatoRSS follow <url>` |
| `unfollow` | Unfollow an RSS feed | `gatoRSS unfollow <url>` |
| `following` | List feeds you're following | `gatoRSS following` |
| `browse` | View recent posts from your feeds | `gatoRSS browse [limit]` |
| `agg` | Start aggregating feeds | `gatoRSS agg <interval>` |

## Example Workflow

```bash
# 1. Register a user
gatoRSS register john

# 2. Add some feeds
gatoRSS addfeed "Hacker News" https://news.ycombinator.com/rss
gatoRSS addfeed "Go Blog" https://go.dev/blog/feed.atom

# 3. Follow the feeds
gatoRSS follow https://news.ycombinator.com/rss
gatoRSS follow https://go.dev/blog/feed.atom

# 4. Start aggregating (in a separate terminal)
gatoRSS agg 1m

# 5. Browse posts
gatoRSS browse 20
```

## Troubleshooting

### Database Connection Issues
- Verify PostgreSQL is running: `sudo service postgresql status`
- Check your credentials in `~/.gatorconfig.json`
- Ensure the database exists: `psql -l`

### Command Not Found
- Ensure `$GOPATH/bin` is in your PATH
- Try running: `export PATH=$PATH:$(go env GOPATH)/bin`
