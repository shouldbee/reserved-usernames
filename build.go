package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

const MASTER_USERNAME_LIST = "reserved-usernames.txt"

var formatters = map[string]Formatter{
	"json": &JsonFormatter{},
	"csv":  &CsvFormatter{},
	"sql":  &SqlFormatter{},
	"php":  &PhpFormatter{},
}

func main() {
	app := cli.NewApp()
	app.Name = "reserved-usernames"
	app.Usage = "reserved usernames convert to json, sql, csv"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "format", Value: "json", Usage: "format for the reserved usernames"},
	}

	app.Action = func(c *cli.Context) {
		format := c.String("format")

		file, err := os.Open(MASTER_USERNAME_LIST)
		failOnError(err)
		defer file.Close()

		lines := scanForInput(file)
		fin := make(chan bool)
		formatter := formatters[format]

		go func(formatter Formatter, lines chan string, fin chan bool) {
			fmt.Print(formatter.start())
			buf := make([]string, 0)
			for line := range lines {
				buf = append(buf, formatter.format(line))
			}
			fmt.Print(formatter.join(buf))
			fmt.Print(formatter.end())
			fin <- true
		}(formatter, lines, fin)

		<-fin
	}

	app.Run(os.Args)
}

func scanForInput(file *os.File) chan string {
	lines := make(chan string)

	r := bufio.NewReader(file)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	go func() {
		for scanner.Scan() {
			lines <- scanner.Text()
		}

		close(lines)
	}()

	return lines
}

func failOnError(err error) {
	if err != nil {
		panic(err)
	}
}

//
// fomatters
//

type Formatter interface {
	start() string
	end() string
	format(line string) string
	join(lines []string) string
}

//
// json
//

type JsonFormatter struct {
}

func (this *JsonFormatter) start() string {
	return ""
}

func (this *JsonFormatter) end() string {
	return "\n"
}

func (this *JsonFormatter) format(line string) string {
	return line
}

func (this *JsonFormatter) join(lines []string) string {
	b, err := json.MarshalIndent(lines, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)
}

//
// csv
//

type CsvFormatter struct {
}

func (this *CsvFormatter) start() string {
	return ""
}

func (this *CsvFormatter) end() string {
	return "\n"
}

func (this *CsvFormatter) format(line string) string {
	return line
}

func (this *CsvFormatter) join(lines []string) string {
	return strings.Join(lines, "\n")
}

//
// sql
//

type SqlFormatter struct {
}

func (this *SqlFormatter) start() string {
	return `CREATE TABLE reserved_usernames (
  username varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
INSERT INTO reserved_usernames VALUES
`
}

func (this *SqlFormatter) end() string {
	return ";\n"
}

func (this *SqlFormatter) format(line string) string {
	return fmt.Sprintf("('%s')", line)
}

func (this *SqlFormatter) join(lines []string) string {
	return strings.Join(lines, ",\n")
}

//
// php
//

type PhpFormatter struct {
}

func (this *PhpFormatter) start() string {
	return `<?php

/*
How to use:

$reservedUsernames = require "reserved-usernames.php";

if (in_array($username, $reservedUsernames)) {
   // This username is reserved.
}
*/

return [
`
}

func (this *PhpFormatter) end() string {
	return "\n];\n"
}

func (this *PhpFormatter) format(line string) string {
	return fmt.Sprintf("    '%s'", line)
}

func (this *PhpFormatter) join(lines []string) string {
	return strings.Join(lines, ",\n")
}
