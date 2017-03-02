package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func logger() {

	filename := fmt.Sprintf("%s-log-%s.txt",
		filepath.Base(os.Args[0]),
		time.Now().Format("2006.01.02-15.04.05"))

	fp, err := os.Create(filename)
	x(err)
	defer fp.Close()

	fmt.Fprintln(fp, time.Now().Format("T 15:04:05"))

	ticker := time.Tick(1 * time.Second)

	for {
		select {
		case t := <-ticker:
			fmt.Fprintln(fp, t.Format("T 15:04:05"))
		case line := <-chLog:
			fmt.Fprintln(fp, line)
		case <-chQuit:
			return
		}
	}
}
