package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type tRequest struct {
	uid string
	idx int
	req string
}

var (
	port = ":8448"

	chIn      = make(chan tRequest, 100)
	chOut     = make([]chan string, 0)
	chCmd     = make(chan string, 100)
	chLog     = make(chan string, 100)
	chLogDone = make(chan bool)
	chQuit    = make(chan bool)
)

func main() {

	for range users {
		chOut = append(chOut, make(chan string, 100))
	}

	go func() {
		chSignal := make(chan os.Signal, 1)
		signal.Notify(chSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		sig := <-chSignal
		chLog <- "I Signal: " + sig.String()
		finish()
		os.Exit(0)
	}()

	go logger()

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

	finish()
}

func finish() {
	close(chQuit)
	<-chLogDone
}

func handleConnection(conn net.Conn) {

	r := conn.RemoteAddr()
	name := r.Network() + "/" + r.String()

	// log to console, not to logfile
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
	uid := a[1]

	idx, ok := labels[uid]
	if !ok {
		return
	}

	out := chOut[idx]

	fmt.Fprintln(conn, ".")

	fmt.Println("     ", name, "=", uid)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "quit" {
			break
		}

		chIn <- tRequest{
			uid: uid,
			idx: idx,
			req: line, // no newline
		}

		for busy := true; busy; {
			select {
			case txt := <-out: // including newline
				fmt.Fprint(conn, txt) // no newline
			default:
				busy = false
			}
		}
		fmt.Fprintln(conn, ".")
	}
	w(scanner.Err())
}
