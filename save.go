package main

import (
	"io/ioutil"
	"log"
	"os"
)

func save(responseChan <-chan string, conf config) {
	var f *os.File
	var err error
	if f, err = ioutil.TempFile("", "hostsmerge.hosts."); err != nil {
		log.Panic(err)
		return
	}

	fqdnsCount := 0
	for response := range responseChan {
		if _, err = f.WriteString(conf.IP + " " + response + "\n"); err != nil {
			f.Close()
			log.Panic(err)
			return
		}
		fqdnsCount++
	}
	fn := f.Name()
	f.Close()

	if err = os.Rename(fn, conf.TargetFile); err != nil {
		log.Printf("Unable to rename tmp file to target file: %v", err)
	}

	log.Printf("Saved %v domains\n", fqdnsCount)
}
