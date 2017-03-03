package main

import (
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func logger() {

	filename := fmt.Sprintf("%s-log-%s.txt.gz",
		filepath.Base(os.Args[0]),
		time.Now().Format("2006.01.02-15.04.05"))

	fp, err := os.Create(filename)
	x(err)

	zw := gzip.NewWriter(fp)

	defer func() {
		zw.Flush()
		zw.Close()
		fp.Close()
		chLogDone <- true
	}()

	fmt.Fprintln(zw, time.Now().Format("T 15:04:05"))

	ticker := time.Tick(1 * time.Second)

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
