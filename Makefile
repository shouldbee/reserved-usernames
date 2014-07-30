build:
	cat reserved-usernames.txt | sort | uniq > reserved-usernames.txt.tmp
	mv reserved-usernames.txt.tmp reserved-usernames.txt
	go run build.go --format csv  > reserved-usernames.csv
	go run build.go --format json > reserved-usernames.json
	go run build.go --format sql  > reserved-usernames.sql
	go run build.go --format php  > reserved-usernames.php
