build:
	go run build.go --format csv  > reserved-usernames.csv
	go run build.go --format json > reserved-usernames.json
	go run build.go --format sql  > reserved-usernames.sql
