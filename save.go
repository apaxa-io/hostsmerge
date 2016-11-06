package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

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

	if err = copyFileContents(fn, conf.TargetFile); err != nil {
		log.Printf("Unable to copy tmp file to target file: %v", err)
	}
	if err = os.Remove(fn); err != nil {
		log.Printf("Unable to remove tmp file: %v", err)
	}

	log.Printf("Saved %v domains\n", fqdnsCount)
}
