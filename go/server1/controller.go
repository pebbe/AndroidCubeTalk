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
		case req := <-chIn:
			chLog <- "R " + req.uid + " " + req.req
			handleReq(req)
		case cmd := <-chCmd:
			chLog <- "C " + cmd
			handleCmd(cmd)
		}
	}
}

// handleReq is not run concurrently, so it must be fast
func handleReq(req tRequest) {

	cmd := req.req
	uid := req.uid
	ch := chOut[uid]
	user := users[uid]

	words := strings.Fields(cmd)
	switch words[0] {

	case "reset":

		for i := 0; i < NR_OF_COUNTERS; i++ {
			user.n[i] = 0
		}

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

		// only one goroutine modifying these variables, so no sync needed
		user.lookat.x = X
		user.lookat.y = Y
		user.lookat.z = Z
		user.roll = roll

		if !user.init {
			// this must be in one batch to make sure that the order is preserved
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "self %g\n", user.selfZ)
			for _, cube := range user.cubes {
				user.n[1]++
				fmt.Fprintf(&buf, "enter %s %d\n", cube.uid, user.n[1])
				user.n[2]++
				fmt.Fprintf(&buf, "moveto %s %d %g %g %g\n", cube.uid, user.n[2], cube.pos.x, cube.pos.y, cube.pos.z)
				user.n[4]++
				fmt.Fprintf(&buf, "color %s %d %g %g %g\n", cube.uid, user.n[4], cube.color.r, cube.color.g, cube.color.b)
			}
			ch <- buf.String()
			user.init = true
		}

		user.n[3]++
		for _, cube := range user.cubes {

			l := users[cube.uid].lookat
			f := cube.forward

			// assumption: forward is horizontal

			// vertical movement (nodding) is amplified by cube.nod
			// cube.nod can by modified for each user and each cube s/he sees individually
			// currently, the GUI only allows setting all cubes for all users to the same value

			rotH := math.Atan2(l.x, l.z) - math.Atan2(f.x, -f.z)
			rotV := between(
				math.Atan2(l.y, math.Sqrt(l.x*l.x+l.z*l.z))*cube.nod,
				-math.Pi/2+.001,
				math.Pi/2-.001)

			ch <- fmt.Sprintf("lookat %s %d %g %g %g %g\n",
				cube.uid,
				user.n[3],
				math.Sin(rotH)*math.Cos(rotV),
				math.Sin(rotV),
				math.Cos(rotH)*math.Cos(rotV),
				users[cube.uid].roll)

		}

	default:

		w(fmt.Errorf("Invalid command from %q: %s", uid, cmd))

	}
}

// handleCmd is not run concurrently, so it must be fast
func handleCmd(cmd string) {

	fmt.Println("Command:", cmd)

	words := strings.Fields(cmd)
	switch words[0] {

	case "recenter":

		// Orders the headset of a particular user to recenter
		// the direction the head is currently pointing.

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

		// Set amplification of nodding for all users, for all cubes they see

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
