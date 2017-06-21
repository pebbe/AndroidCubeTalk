package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	reCommandLine = regexp.MustCompile(` I Command line: \[\]string{".*?", "(.*)"}`)
	reLogLine     = regexp.MustCompile(`^[0-9]+:[0-9]+.[0-9]+`)
)

func configReplay(filename string) {
	useReplay = true

	fp, err := os.Open(filename)
	x(err)
	rd, err := gzip.NewReader(fp)
	x(err)

	defer func() {
		rd.Close()
		fp.Close()
	}()

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		m := reCommandLine.FindStringSubmatch(scanner.Text())
		if m != nil && len(m) == 2 {
			readConfig(m[1])
			return
		}
	}

	x(fmt.Errorf("No command line found in %s", filename))
}

func replay(filename string) {
	time.Sleep(time.Millisecond * 100)
	fmt.Print("Press ENTER to start...")
	r := bufio.NewReader(os.Stdin)
	r.ReadLine()

	fp, err := os.Open(filename)
	x(err)
	rd, err := gzip.NewReader(fp)
	x(err)

	defer func() {
		rd.Close()
		fp.Close()
	}()

	scanner := bufio.NewScanner(rd)

	scanner.Scan()
	logStart, err := time.Parse("04:05.000", strings.Fields(scanner.Text())[0])
	x(err)
	logPrev := logStart

	replayStart := time.Now()

	fmt.Println(logStart)
	fmt.Println(replayStart)

	for scanner.Scan() {
		line := scanner.Text()
		if reLogLine.MatchString(line) {
			words := strings.Fields(line)
			if len(words) < 3 {
				continue
			}
			t, err := time.Parse("04:05.000", words[0])
			if err != nil {
				continue
			}
			for t.Before(logPrev) {
				t.Add(time.Hour)
			}
			logPrev = t
			time.Sleep(t.Sub(logStart) - time.Now().Sub(replayStart))
			fmt.Println(line)

			switch words[1] {
			case "C":
				chCmd <- strings.Join(words[2:], " ")
			case "R":
				uid := words[2]
				chReplay <- tRequest{
					uid: uid,
					idx: labels[uid],
					req: strings.Join(words[3:], " "), // no newline
				}
			}
		}
	}
}
