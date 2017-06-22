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

	if withReplay {
		for {
			select {
			case <-chLog:
			case <-chQuit:
				close(chLogDone)
				return
			}
		}
	}

	stamp := func() string {
		t := time.Now()
		return fmt.Sprintf("%02d:%02d.%03d", t.Minute(), t.Second(), t.Nanosecond()/1000000)
	}

	start := time.Now()
	s := stamp()

	filename := fmt.Sprintf("%s-log-%s.txt.gz",
		filepath.Base(os.Args[0]),
		start.Format("2006.01.02-15.04.05"))

	fp, err := os.Create(filename)
	x(err)

	zw := gzip.NewWriter(fp)

	fmt.Fprintln(zw, s, "I Start:", start.Format(time.RFC1123Z))
	fmt.Fprintf(zw, s+" I Command line: %#v\n", os.Args)

	defer func() {
		stop := time.Now()
		s := stamp()
		fmt.Fprintln(zw, s, "I Stop:", stop.Format(time.RFC1123Z))
		fmt.Fprintln(zw, s, "I Uptime:", time.Since(start))
		zw.Close()
		fp.Close()
		close(chLogDone)
	}()

	for {
		select {
		case line := <-chLog:
			fmt.Fprintln(zw, stamp(), line)
		case <-chQuit:
			return
		}
	}
}
