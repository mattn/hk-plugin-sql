package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

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

	b, err := exec.Command("hk", "env").Output()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal("DATABASE_URL not found")
	}

	br := bufio.NewReader(os.Stdin)
	var line string
	for {
		fmt.Print("SQL> ")
		b, _, err = br.ReadLine()
		line = string(b)
		if err != nil || line == "exit" {
			break
		}
		if len(b) == 0 {
			continue
		}
		param := make(url.Values)
		param.Set("database_url", database_url)
		param.Set("sql", line)
		res, err := http.PostForm("https://sql-console.heroku.com/query", param)
		if err != nil {
			log.Println(err)
			continue
		}
		defer res.Body.Close()
		io.Copy(os.Stdout, res.Body)
	}
}
