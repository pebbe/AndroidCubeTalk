package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func controller() {

	for {
		select {
		case <-chQuit:
			return
		case cmd := <-chCmd:
			chLog <- "C " + cmd
			handleCmd(cmd)
		case req := <-chIn:
			chLog <- "I " + req.uid + " " + req.req
			handleIn(req)
		}
	}
}

func handleIn(req tRequest) {
	defer close(req.chClose)
	cmd := req.req
	uid := req.uid
	ch := chOut[uid]
	user := users[uid]
	words := strings.Fields(cmd)
	switch words[0] {
	case "reset":
		user.n0 = 0
		user.n1 = 0
		user.n2 = 0
		user.n3 = 0
		user.n4 = 0
		user.n5 = 0
		user.init = false
	case "lookat":
		if len(words) != 5 {
			w(fmt.Errorf("Invalid number of arguments from %q: %s", uid, cmd))
			return
		}
		X, err := strconv.ParseFloat(words[1], 64)
		if w(err) != nil {
			return
		}
		Y, err := strconv.ParseFloat(words[2], 64)
		if w(err) != nil {
			return
		}
		Z, err := strconv.ParseFloat(words[3], 64)
		if w(err) != nil {
			return
		}
		roll, err := strconv.ParseFloat(words[4], 64)
		if w(err) != nil {
			return
		}

		user.lookat.x = X
		user.lookat.y = Y
		user.lookat.z = Z
		user.roll = roll

		if !user.init {
			// this must be in one batch to make sure that the order is preserved
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "self %g\n", user.selfZ)
			for _, cube := range user.cubes {
				user.n1++
				fmt.Fprintf(&buf, "enter %s %d\n", cube.uid, user.n1)
				user.n2++
				fmt.Fprintf(&buf, "moveto %s %d %g %g %g\n", cube.uid, user.n2, cube.pos.x, cube.pos.y, cube.pos.z)
				user.n4++
				fmt.Fprintf(&buf, "color %s %d %g %g %g\n", cube.uid, user.n2, cube.red, cube.green, cube.blue)
			}
			ch <- buf.String()
			user.init = true
		}

		user.n3++
		for _, cube := range user.cubes {
			if cube.uid != uid {

				l := users[cube.uid].lookat
				f := cube.forward

				// assumption: forward is horizontal

				rotH := math.Atan2(l.x, l.z) - math.Atan2(f.x, -f.z)
				rotV := between(math.Atan2(l.y, math.Sqrt(l.x*l.x+l.z*l.z))*cube.nod, -math.Pi/2+.001, math.Pi/2-.001)

				ch <- fmt.Sprintf("lookat %s %d %g %g %g %g\n",
					cube.uid,
					user.n3,
					math.Sin(rotH)*math.Cos(rotV),
					math.Sin(rotV),
					math.Cos(rotH)*math.Cos(rotV),
					users[cube.uid].roll)
			}
		}

	default:
		w(fmt.Errorf("Invalid command from %q: %s", uid, cmd))
	}
}

func handleCmd(cmd string) {
	fmt.Println("Command:", cmd)
	words := strings.Fields(cmd)
	switch words[0] {
	case "recenter":
		if len(words) != 2 {
			w(fmt.Errorf("Invalid number of arguments from GUI: %s", cmd))
			return
		}
		if ch, ok := chOut[words[1]]; ok {
			select {
			case ch <- "recenter\n":
			default:
				// channel is full and nobody is reading from the channel
			}
		} else {
			w(fmt.Errorf("Invalid user in command from GUI: %s", cmd))
		}
	case "globalnod":
		if len(words) != 2 {
			w(fmt.Errorf("Invalid number of arguments from GUI: %s", cmd))
			return
		}
		f, err := strconv.ParseFloat(words[1], 64)
		if w(err) != nil {
			return
		}
		for _, user := range users {
			for i := range user.cubes {
				user.cubes[i].nod = f
			}
		}
	default:
		w(fmt.Errorf("Invalid command from GUI: %s", cmd))
	}
}

func between(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
