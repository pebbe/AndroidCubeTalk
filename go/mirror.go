package main

import (
	"github.com/pebbe/util"

	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type lookat struct {
	id   string
	x    float64
	y    float64
	z    float64
	resp chan string
}

var (
	x = util.CheckErr
	w = util.WarnErr

	lookats = make(chan lookat, 100)
	users   = make(map[string]uint64)
)

func main() {
	ln, err := net.Listen("tcp", ":8448")
	x(err)
	defer ln.Close()

	go handleLookats()

	for {
		conn, err := ln.Accept()
		if w(err) != nil {
			break
		}
		go handleConnection(conn)
	}
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
	fmt.Fprintln(conn, "ok")

	fmt.Println("     ", name, "=", id)

LOOP:
	for scanner.Scan() {
		line := scanner.Text()
		a := strings.Fields(line)
		switch a[0] {
		case "quit":
			break LOOP
		case "lookat":
			if len(a) == 4 {
				resp := make(chan string)
				x, err := strconv.ParseFloat(a[1], 64)
				if err != nil {
					continue
				}
				y, err := strconv.ParseFloat(a[2], 64)
				if err != nil {
					continue
				}
				z, err := strconv.ParseFloat(a[3], 64)
				if err != nil {
					continue
				}
				r := lookat{
					id:   id,
					x:    x,
					y:    y,
					z:    z,
					resp: resp,
				}
				lookats <- r

				for line := range resp {
					fmt.Fprintln(conn, line)
				}
			}
			fmt.Fprintln(conn, ".")
		}
	}
	w(scanner.Err())

}

func handleLookats() {
	for {
		req := <-lookats

		u, ok := users[req.id]
		if !ok {
			u = 0
			req.resp <- "self 0 -4"
			req.resp <- "enter B"
			req.resp <- "moveto B 0 0 0 4"
		}

		req.resp <- fmt.Sprintf("lookat B %d %g %g %g", u, req.x, req.y, -req.z)
		u++
		users[req.id] = u

		close(req.resp)
	}
}
