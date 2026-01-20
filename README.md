# GoBlogAggregator (gator)

## Requirements
- Postgres
- Go

## How to install:
1. Download the files
2. run `go install GoBlogAggregator` to instal gator
3. create a `.gatorconfig.json` file in your home directory
```
{
"db_url": "postgres://example"
}
```


## How to Run:
- Run `./GoBlogAggregator <command> <args>`

### Commands:
- `register <username>`: registers sets the current user to given username
- `login <username>`: switches current user to given registered user
- `addfeed <title> <url>`: adds a feed to the database under the given title, and adds the feed to the current user's follow list
- `follow <url>`: adds the given url to the current user's follow list
- `following`: lists all the feeds the user is following
- `agg <time_between requests>`: parses all feeds in the database
- `browse <limit>`: lists the latest x amount of posts
