package main

import (
	"bufio"
	"fmt"
	"net"
	"runtime"
	"strings"
	"sync"
)

type tRequest struct {
	req     string
	chClose chan bool
}

var (
	port = ":8448"

	chIn   = make(map[string]chan tRequest)
	chOut  = make(map[string]chan string)
	chGet  = make(map[string]chan chan [4]float64)
	chSet  = make(map[string]chan [4]float64)
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
		chIn[user] = make(chan tRequest, 100)
		chOut[user] = make(chan string, 100)
		chGet[user] = make(chan chan [4]float64, 100)
		chSet[user] = make(chan [4]float64, 100)
		go handleUser(user)
		go controller(user)
	}

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
	in := chIn[id]

	fmt.Fprintln(conn, ".")

	fmt.Println("     ", name, "=", id)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "quit" {
			break
		}

		ch := make(chan bool)
		req := tRequest{
			req:     line, // no newline
			chClose: ch,
		}
		in <- req

		for busy := true; busy; {
			select {
			case txt := <-out: // including newline
				fmt.Fprint(conn, txt) // no newline
			case <-ch:
				busy = false
			}
		}
		fmt.Fprintln(conn, ".")
		runtime.Gosched()
	}
	w(scanner.Err())
}

func handleUser(uid string) {
	getCh := chGet[uid]
	setCh := chSet[uid]
	user := users[uid]
	for {
		select {
		case data := <-setCh:
			user.lookat = tVector{data[0], data[1], data[2]}
			user.roll = data[3]
		case setter := <-getCh:
			setter <- [4]float64{user.lookat.x, user.lookat.y, user.lookat.z, user.roll}
		}
	}
}
