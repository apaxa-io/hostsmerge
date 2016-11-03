package main

import (
	"bufio"
	"log"
	"net/http"
	"regexp"
)

var uselessLine = regexp.MustCompile(`^[[:space:]]*(#.*)?$`) // Empty lines and|or comments

var usefullLine = regexp.MustCompile(`^[[:space:]]*[0-9.:]+[[:space:]]+([0-9A-Za-z_.-]+)(?:[[:space:]]+(?:#.*)?)?$`)

func parseLine(l []byte) (fqdn string, skip, bad bool) {
	if uselessLine.Match(l) {
		skip = true
		return
	}
	parts := usefullLine.FindSubmatch(l)
	bad = parts == nil || len(parts) != 2
	if bad {
		skip = true
	} else {
		fqdn = string(parts[1])
	}
	return
}

func parse(urls <-chan string, fqdns chan<- string) {
	for url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Unable to fetch %v : %v\n", url, err)
			continue
		}

		reader := bufio.NewReader(resp.Body)
		lineNum := -1
		fqdnsCount := 0
		for line, prefix, err := reader.ReadLine(); err == nil; line, prefix, err = reader.ReadLine() {
			lineNum++
			if prefix {
				// Skip long line
				for _, prefix, err = reader.ReadLine(); err == nil && prefix; _, prefix, err = reader.ReadLine() {
				}
				log.Printf("Skip too long line #%v in %v : %v\n", lineNum, url, string(line))
				continue
			}

			fqdn, skip, bad := parseLine(line)
			if bad {
				log.Printf("Skip bad line #%v in %v : %v\n", lineNum, url, string(line))
				continue
			}
			if skip {
				continue
			}
			fqdns <- fqdn
			fqdnsCount++
		}
		resp.Body.Close()

		log.Printf("In %v %v domains\n", url, fqdnsCount)
	}
}
