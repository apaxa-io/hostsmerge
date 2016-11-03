package main

import (
	"log"
	"strings"
)

func inExcluded(fqdn string, conf config) bool {
	for _, excl := range conf.Excludes {
		if fqdn == excl {
			return true
		}
	}
	return false
}

func normalize(fqdn string) string {
	return strings.ToLower(fqdn)
}

func unique(fqdnsChan <-chan string, requestChan chan<- string, conf config) {
	skipped := 0
	s := make(map[string]struct{})
	for fqdn := range fqdnsChan {
		fqdn = normalize(fqdn)
		if inExcluded(fqdn, conf) {
			skipped++
			continue
		}
		if _, exists := s[fqdn]; !exists {
			s[fqdn] = struct{}{}
			requestChan <- fqdn
		} else {
			skipped++
		}
	}
	log.Printf("%v unique domains in all urls, %v skipped\n", len(s), skipped)
	close(requestChan)
}
