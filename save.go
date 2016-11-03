package main

import (
	"log"
	"os"
)

const newSuffix = ".new"

func save(responseChan <-chan string, conf config) {
	var f *os.File
	var err error
	if f, err = os.Create(conf.TargetFile + newSuffix); err != nil {
		log.Panic(err)
		return
	}
	defer f.Close()

	fqdnsCount := 0
	for response := range responseChan {
		if _, err = f.WriteString(conf.IP + " " + response + "\n"); err != nil {
			log.Panic(err)
			return
		}
		fqdnsCount++
	}

	if err = os.Rename(conf.TargetFile+newSuffix, conf.TargetFile); err != nil {
		log.Printf("Unable to rename tmp file to target file: %v", err)
	}

	log.Printf("Saved %v domains\n", fqdnsCount)
}
