/*

This server provides the user with three cubes
that copy the user's movements

*/

package main

import (
	"github.com/pebbe/util"

	"bufio"
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
)

type Request struct {
	id   string
	req  string
	resp chan string
}

var (
	x = util.CheckErr
	w = util.WarnErr

	requests = make(chan Request, 100)
	users    = make(map[string]uint64)

	errArgs    = errors.New("Wrong number of arguments")
	errUnknown = errors.New("Unknown command")
	errNan     = errors.New("Not a number")
	errInf     = errors.New("Infinity")
)

func main() {
	ln, err := net.Listen("tcp", ":8448")
	x(err)
	defer ln.Close()

	go handleRequests()

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
	fmt.Fprintln(conn, ".")

	fmt.Println("     ", name, "=", id)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "quit" {
			break
		}
		resp := make(chan string)
		requests <- Request{
			id:   id,
			req:  line,
			resp: resp,
		}

		for line := range resp {
			fmt.Fprintln(conn, line)
		}
		fmt.Fprintln(conn, ".")
	}
	w(scanner.Err())

}

func handleRequests() {

	invalid := func(req Request, err error) {
		req.resp <- strings.Replace(fmt.Sprintf("error - %s - %v", req.req, err), "\n", " ", -1)
		fmt.Printf("error %s - %s - %v\n", req.id, req.req, err)
	}

	for {
		req := <-requests

		a := strings.Fields(req.req)
		if len(a) > 0 {
			switch a[0] {
			case "reset":
				if len(a) == 1 {
					delete(users, req.id)
				} else {
					invalid(req, errArgs)
				}
			case "lookat":
				if len(a) == 4 {
					u, ok := users[req.id]
					if !ok {
						u = 0
						req.resp <- "self 0 -4"
						req.resp <- "enter A"
						req.resp <- "moveto A 0 4 0 0"
						req.resp <- "enter B"
						req.resp <- "moveto B 0 0 0 4"
						req.resp <- "enter C"
						req.resp <- "moveto C 0 -4 0 0"
					}
					x, err := strconv.ParseFloat(a[1], 64)
					if err == nil {
						if math.IsNaN(x) {
							invalid(req, errNan)
						} else if math.IsInf(x, 0) {
							invalid(req, errInf)
						} else {

							z, err := strconv.ParseFloat(a[3], 64)
							if err == nil {
								if math.IsNaN(z) {
									invalid(req, errNan)
								} else if math.IsInf(z, 0) {
									invalid(req, errInf)
								} else {
									req.resp <- fmt.Sprintf("lookat A %d %g %s %g", u, z, a[2], -x)
									req.resp <- fmt.Sprintf("lookat B %d %g %s %g", u, -x, a[2], -z)
									req.resp <- fmt.Sprintf("lookat C %d %g %s %g", u, -z, a[2], x)
									u++
								}
							} else {
								invalid(req, err)
							}

						}
					} else {
						invalid(req, err)
					}
					users[req.id] = u
				} else {
					invalid(req, errArgs)
				}
			default:
				invalid(req, errUnknown)
			}
		}

		close(req.resp)
	}
}
