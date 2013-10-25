package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mattn/go-runewidth"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

const describe = "select column_name as column,data_type as type,character_maximum_length as length,is_nullable,column_default as default from information_schema.columns where table_name=$1 order by ordinal_position"

func runQuery(conn *sql.DB, query string, args []interface{}) error {
	stmt, err := conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	lines := make([][]string, 0)
	cellWidths := make([][]int, 0)
	maxWidths := make([]int, len(columns))

	cells := make([]string, len(columns))
	widths := make([]int, len(columns))
	for i, column := range columns {
		cells[i] = column
		widths[i] = runewidth.StringWidth(column)
		if widths[i] > maxWidths[i] {
			maxWidths[i] = widths[i]
		}
	}
	lines = append(lines, cells)
	cellWidths = append(cellWidths, widths)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanValues := make([]interface{}, len(columns))
		for i := range values {
			scanValues[i] = &values[i]
		}
		err = rows.Scan(scanValues...)
		if err != nil {
			return err
		}
		cells = make([]string, len(columns))
		widths = make([]int, len(columns))
		for i, value := range values {
			var cell string
			leftalign := true
			switch value.(type) {
			case int, uint, int32, uint32, int64, uint64, float32, float64:
				leftalign = false
				cell = fmt.Sprint(value)
			case nil:
				cell = ""
			case []byte:
				cell = string(value.([]byte))
			default:
				cell = fmt.Sprint(value)
			}
			cells[i] = cell
			widths[i] = runewidth.StringWidth(cell)
			if widths[i] > maxWidths[i] {
				maxWidths[i] = widths[i]
			}
			if leftalign {
				widths[i] = -widths[i]
			}
		}
		lines = append(lines, cells)
		cellWidths = append(cellWidths, widths)
	}
	for i, cell := range lines[0] {
		spaces := strings.Repeat(" ", maxWidths[i]-cellWidths[0][i])
		fmt.Print("|" + cell + spaces)
	}
	fmt.Println("|")
	for i := range lines[0] {
		dashes := strings.Repeat("-", maxWidths[i])
		if i == 0 {
			fmt.Print("|" + dashes)
		} else {
			fmt.Print("+" + dashes)
		}
	}
	fmt.Println("|")
	for n, cells := range lines[1:] {
		for i, cell := range cells {
			width := cellWidths[n+1][i]
			if width < 0 {
				spaces := strings.Repeat(" ", maxWidths[i]+cellWidths[n+1][i])
				fmt.Print("|" + cell + spaces)
			} else {
				spaces := strings.Repeat(" ", maxWidths[i]-cellWidths[n+1][i])
				fmt.Print("|" + spaces + cell)
			}
		}
		fmt.Println("|")
	}
	return nil
}

func main() {
	app := os.Getenv("HKAPP")
	host := os.Getenv("HKHOST")
	pluginMode := os.Getenv("HKPLUGINMODE")
	description := "Plugin that provides a sql console for your heroku app."
	switch pluginMode {
	case "info":
		fmt.Printf("sql 0.1: SQL console for %s.%s\n\n%s", app, host, description)
		return
	}

	b, err := exec.Command("hk", "env").CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(b))
		os.Exit(1)
	}
	var database_url string
	for _, line := range strings.Split(string(b), "\n") {
		token := strings.SplitN(line, "=", 2)
		if len(token) != 2 {
			continue
		}
		if token[0] == "DATABASE_URL" {
			database_url = token[1]
			break
		}
	}
	if database_url == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL not found")
		os.Exit(1)
	}

	conn, err := sql.Open("postgres", database_url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		<-sc
		fmt.Println()
		conn.Close()
		os.Exit(0)
	}()

	br := bufio.NewReader(os.Stdin)
	var query string
	for {
		fmt.Print("SQL> ")
		b, _, err = br.ReadLine()
		query = string(b)
		if err != nil || query == "exit" {
			break
		}
		token := strings.SplitN(query, " ", 2)
		args := []interface{}{}
		if len(token) == 2 && token[0] == "\\d" || token[0] == "desc" {
			query = describe
			println(strings.TrimRight(token[1], "; "))
			args = append(args, strings.TrimRight(token[1], "; "))
		}
		if len(b) == 0 {
			continue
		}
		err = runQuery(conn, query, args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
