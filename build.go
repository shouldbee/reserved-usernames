package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"encoding/json"
	"github.com/codegangsta/cli"
)

const WORDS_FILE = "reserved-usernames.txt"

func main() {
	app := cli.NewApp()
	app.Name = "reserved_words"
	app.Usage = "reserved words convert to json, sql, csv"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "format", Value: "json", Usage: "format for the reserved words"},
	}

	app.Action = func(c *cli.Context) {
		format := c.String("format")

		file, err := os.Open(WORDS_FILE)

		if err != nil {
			panic(err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}()

		lines := scanForInput(file)

		fin := make(chan bool)

		var formatter Formatter

		if format == "json" {
			formatter = NewJsonFormatter()
		} else if format == "csv" {
			formatter = NewCsvFormatter()
		} else if format == "sql" {
			formatter = NewSqlFormatter()
		} else {
			panic("no such formatter")
		}

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

type Formatter interface {
	start() string
	end() string
	format(line string) string
	join(lines []string) string
}

func NewJsonFormatter() *JsonFormatter {
	return &JsonFormatter{}
}

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
func NewCsvFormatter() *CsvFormatter {
	return &CsvFormatter{}
}

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
func NewSqlFormatter() *SqlFormatter {
	return &SqlFormatter{}
}

type SqlFormatter struct {
}

func (this *SqlFormatter) start() string {
	return "INSERT INTO ${reserved_words_table} (${reserved_words_column}) VALUES \n"
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
