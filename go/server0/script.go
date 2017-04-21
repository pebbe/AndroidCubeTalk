package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	chScript chan bool
)

func scriptStart() {
	if config.Script == nil || len(config.Script) == 0 {
		return
	}

	scriptStop()

	chScript = make(chan bool)

	go func(ch chan bool) {
		fmt.Println("Script started")
		for {

			for _, line := range config.Script {
				select {
				case <-ch:
					fmt.Println("Script stopped")
					return
				default:
				}
				if strings.HasPrefix(line, "sleep") {
					a := strings.Fields(line)
					i, err := strconv.Atoi(a[1])
					x(err)
					time.Sleep(time.Duration(i) * time.Millisecond)
				} else {
					chCmd <- line
				}
			}

			if !config.ScriptRepeat {
				break
			}
		}
		fmt.Println("Script finished")
	}(chScript)
}

func scriptStop() {
	if chScript != nil {
		close(chScript)
		chScript = nil
	}
}
