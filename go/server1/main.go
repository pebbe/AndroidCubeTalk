package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	port = ":8448"

	chIn   = make(map[string]chan string)
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
		chIn[user] = make(chan string, 100)
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

	in := chIn[id]
	out := chOut[id]

	fmt.Fprintln(conn, ".")

	fmt.Println("     ", name, "=", id)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "quit" {
			break
		}

		in <- line // no newline

	LOOP:
		for {
			select {
			case txt := <-out: // including newline
				fmt.Fprint(conn, txt) // no newline
			default:
				break LOOP
			}
		}
		fmt.Fprintln(conn, ".")
	}
	w(scanner.Err())
}
