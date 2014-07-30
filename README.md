# Reserved usernames list

A list of reserved usernames to prevent url collision with resource paths.
This repository hosts the list in multiple formats like JSON, CSV, SQL and plain text.
You can use its just download its by `wget`.

## Downloading

You can download the list by `wget`.

CSV:

```
wget https://github.com/shouldbee/reserved-usernames/raw/master/reserved-usernames.csv
```

JSON:

```
wget https://github.com/shouldbee/reserved-usernames/raw/master/reserved-usernames.json
```

SQL:

```
wget https://github.com/shouldbee/reserved-usernames/raw/master/reserved-usernames.sql
```

Plain text:

```
wget https://github.com/shouldbee/reserved-usernames/raw/master/reserved-usernames.txt
```

## Contribution

### How to add new usernames

You need `go` to compile multiple formats.

1. Edit reserved-usernames.txt.
2. Run `make build`.
3. Commit it.
4. Then pull request it!

