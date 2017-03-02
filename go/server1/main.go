package main

import (
	"sync"
	"time"
)

var (
	chLog  = make(chan string, 100)
	chQuit = make(chan bool)
	wg     sync.WaitGroup
)

func main() {

	wg.Add(1)
	go func() {
		logger()
		wg.Done()
	}()

	chLog <- "I dit is een test"
	time.Sleep(30 * time.Second)

	close(chQuit)
	wg.Wait()
}
