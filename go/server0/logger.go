package main

import (
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	INTERVAL = 1000 // time stamp interval in milliseconds
)

func logger() {

	start := time.Now()

	filename := fmt.Sprintf("%s-log-%s.txt.gz",
		filepath.Base(os.Args[0]),
		start.Format("2006.01.02-15.04.05"))

	fp, err := os.Create(filename)
	x(err)

	zw := gzip.NewWriter(fp)

	fmt.Fprintln(zw, "I Start:", start.Format(time.RFC1123Z))
	fmt.Fprintf(zw, "I Command line: %#v\n", os.Args)

	defer func() {
		stop := time.Now()
		fmt.Fprintln(zw, "I Stop:", stop.Format(time.RFC1123Z))
		fmt.Fprintln(zw, "I Uptime:", time.Since(start))
		zw.Close()
		fp.Close()
		close(chLogDone)
	}()

	ticker := time.Tick(INTERVAL * time.Millisecond)

	for {
		select {
		case t := <-ticker:
			fmt.Fprintln(zw, t.Format("T 15:04:05"))
		case line := <-chLog:
			fmt.Fprintln(zw, line)
		case <-chQuit:
			return
		}
	}
}
