package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type tRequest struct {
	uid     string
	req     string
	chClose chan bool
}

var (
	port = ":8448"

	chIn   = make(chan tRequest, 100)
	chOut  = make(map[string]chan string)
	chCmd  = make(chan string, 100)
	chLog  = make(chan string, 100)
	chQuit = make(chan bool)
)

func main() {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		logger()
		wg.Done()
	}()

	for user := range users {
		chOut[user] = make(chan string, 100)
	}

	go controller()

	go func() {
		ln, err := net.Listen("tcp", port)
		x(err)
		for {
			conn, err := ln.Accept()
			x(err)
			go handleConnection(conn)
		}
	}()

	gui()
	close(chQuit)

	wg.Wait() // wait for logger to flush and close file
}

func handleConnection(conn net.Conn) {

	r := conn.RemoteAddr()
	name := r.Network() + "/" + r.String()

	fmt.Println("Open ", name)
	defer func() {
		fmt.Println("Close", name)
		conn.Close()
	}()

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	a := strings.Fields(scanner.Text())
	if len(a) != 2 || a[0] != "join" {
		return
	}
	id := a[1]

	if _, ok := users[id]; !ok {
		return
	}

	out := chOut[id]

	fmt.Fprintln(conn, ".")

	fmt.Println("     ", name, "=", id)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "quit" {
			break
		}

		ch := make(chan bool)
		req := tRequest{
			uid:     id,
			req:     line, // no newline
			chClose: ch,
		}
		chIn <- req

		for busy := true; busy; {
			select {
			case txt := <-out: // including newline
				fmt.Fprint(conn, txt) // no newline
			case <-ch:
				busy = false
			}
		}
		fmt.Fprintln(conn, ".")
	}
	w(scanner.Err())
}
