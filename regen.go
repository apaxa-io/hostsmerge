package main

import (
	"log"
	"sync"
	"time"
)

const urlsChanLen = 20

func statistics(stop <-chan struct{}, urlsChan <-chan string, fqdnsChan <-chan string, requestChan <-chan string, responseChan <-chan string) {
	ticker := time.Tick(10 * time.Second)
	for {
		select {
		case <-ticker:
			log.Printf("Status: %v urls for fetch&parse, %v FQDNs for unique filter, %v FQDNs for DNS resolve (%v in cache), %v FQDNs for write\n", len(urlsChan), len(fqdnsChan), len(requestChan), cache.Len(), len(responseChan))
		case <-stop:
			return
		}
	}
}

func main() {
	// Init
	conf, err := readConfig()
	if err != nil {
		log.Panic("Unable to read config: " + err.Error())
	}

	log.Printf("Config: %#v\n", conf)

	if err = initCache(conf); err != nil {
		log.Print("Loading cache failed: " + err.Error())
	}

	// Send urls from config
	urlsChan := make(chan string, urlsChanLen)
	go func() {
		for _, url := range conf.URLs {
			urlsChan <- url
		}
		close(urlsChan)
	}()

	// Parsers
	var wgParsers sync.WaitGroup
	wgParsers.Add(conf.NumParsers)
	fqdnsChan := make(chan string, conf.ChanLength)
	for i := 0; i < conf.NumParsers; i++ {
		go func() {
			parse(urlsChan, fqdnsChan)
			wgParsers.Done()
		}()
	}

	// Unique
	requestChan := make(chan string, conf.ChanLength)
	go unique(fqdnsChan, requestChan, conf)

	// Resolvers
	var wgResolvers sync.WaitGroup
	wgResolvers.Add(conf.NumResolvers)
	responseChan := make(chan string, conf.ChanLength)
	for i := 0; i < conf.NumResolvers; i++ {
		go func() {
			resolver(requestChan, responseChan, conf)
			wgResolvers.Done()
		}()
	}

	// Save
	go save(responseChan, conf)

	// Statistics
	stop := make(chan struct{})
	go statistics(stop, urlsChan, fqdnsChan, requestChan, responseChan)

	// Close channels
	wgParsers.Wait()
	log.Println("Urls fetched & parsed")
	close(fqdnsChan)
	wgResolvers.Wait()
	log.Println("FQDNs resolved")
	close(responseChan)

	log.Print("Saving cache ... ")
	if err := closeCache(conf); err != nil {
		log.Println("Unabla to save cache: " + err.Error())
	} else {
		log.Println("Cache saved")
	}

	close(stop)
}
