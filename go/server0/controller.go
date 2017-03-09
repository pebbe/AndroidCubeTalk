package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	setsize  = make([]bool, len(cubes))
	cubesize = [3]float64{1, 1, 1}
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
	idx := req.idx
	ch := chOut[idx]
	user := users[idx]

	words := strings.Fields(cmd)
	switch words[0] {

	case "reset":

		for i := 0; i < NR_OF_COUNTERS; i++ {
			user.n[i] = 0
		}

		setsize[idx] = true

		user.init = false

	case "lookat":

		if len(words) != 5 && len(words) != 6 {
			w(fmt.Errorf("Invalid number of arguments from %q: %s", req.uid, cmd))
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

		marked := len(words) == 6

		if !user.init {
			// this must be in one batch to make sure that the order is preserved
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "self %g\n", user.selfZ)
			for i, cube := range user.cubes {
				if i != req.idx {
					user.n[1]++
					fmt.Fprintf(&buf, "enter %s %d\n", cube.uid, user.n[1])
					user.n[2]++
					fmt.Fprintf(&buf, "moveto %s %d %g %g %g\n", cube.uid, user.n[2], cube.pos.x, cube.pos.y, cube.pos.z)
					user.n[4]++
					fmt.Fprintf(&buf, "color %s %d %g %g %g\n", cube.uid, user.n[4], cube.color.r, cube.color.g, cube.color.b)
				}
			}
			ch <- buf.String()
			user.init = true
		}

		user.n[3]++
		for i, cube := range user.cubes {

			if i != req.idx {

				if marked {
					if X*cube.towards.x+Y*cube.towards.y+Z*cube.towards.z > *opt_t {
						chLog <- fmt.Sprintf("I Mark %s -> %s", req.uid, cube.uid)
						fmt.Printf("Mark %s -> %s\n", req.uid, cube.uid)
						marked = false
					}
				}

				l := users[i].lookat
				f := cube.forward

				// assumption: forward is horizontal

				// vertical movement (nodding) is amplified by cube.nod
				// cube.nod can by modified for each user and each cube s/he sees individually
				// currently, the GUI only allows setting all cubes for all users to the same value

				rotH := math.Atan2(l.x, l.z) - math.Atan2(f.x, -f.z)
				rotV := nodEnhance(math.Atan2(l.y, math.Sqrt(l.x*l.x+l.z*l.z)), cube.nod)

				ch <- fmt.Sprintf("lookat %s %d %g %g %g %g\n",
					cube.uid,
					user.n[3],
					math.Sin(rotH)*math.Cos(rotV),
					math.Sin(rotV),
					math.Cos(rotH)*math.Cos(rotV),
					users[i].roll)

				// change color of cube to orange if it is looking at me
				v := users[i].lookat
				w := users[i].cubes[idx].towards
				if v.x*w.x+v.y*w.y+v.z*w.z > *opt_t {
					if !cube.lookingatme {
						cube.lookingatme = true
						user.n[4]++
						ch <- fmt.Sprintf("color %s %d 1 .7 0\n", cube.uid, user.n[4])
						chLog <- fmt.Sprintf("I Begin %s looking at %s", cube.uid, req.uid)
					}
				} else {
					if cube.lookingatme {
						cube.lookingatme = false
						user.n[4]++
						ch <- fmt.Sprintf("color %s %d %g %g %g\n", cube.uid, user.n[4], cube.color.r, cube.color.g, cube.color.b)
						chLog <- fmt.Sprintf("I End %s looking at %s", cube.uid, req.uid)
					}
				}
			}

		}

		if setsize[idx] {
			setsize[idx] = false
			ch <- fmt.Sprintf("cubesize %d %g %g %g\n", user.n[6], cubesize[0], cubesize[1], cubesize[2])
		}

		if marked {
			fmt.Printf("Mark %s -> %g %g %g\n", req.uid, X, Y, Z)
		}

	default:

		w(fmt.Errorf("Invalid command from %q: %s", req.uid, cmd))

	}
}

// handleCmd is not run concurrently, so it must be fast
func handleCmd(cmd string) {

	number_args := "Invalid number of arguments from GUI: %s"

	fmt.Println("Command:", cmd)

	words := strings.Fields(cmd)
	switch words[0] {

	case "cubesize":
		if len(words) != 4 {
			w(fmt.Errorf(number_args, cmd))
			return
		}
		for i := 0; i < 3; i++ {
			var err error
			cubesize[i], err = strconv.ParseFloat(words[i+1], 64)
			if w(err) != nil {
				cubesize[i] = 1
			}
		}
		for i := range cubes {
			setsize[i] = true
			users[i].n[6]++
		}

	case "recenter":

		// Orders the headset of a particular user to recenter
		// the direction the head is currently pointing.

		if len(words) != 2 {
			w(fmt.Errorf(number_args, cmd))
			return
		}

		if idx, ok := labels[words[1]]; ok {
			select {
			case chOut[idx] <- "recenter\n":
			default:
				// channel is full and nobody is reading from the channel
			}
		} else {
			w(fmt.Errorf("Invalid user in command from GUI: %s", cmd))
		}

	case "globalnod":

		// Set amplification of nodding for all users, for all cubes they see

		if len(words) != 2 {
			w(fmt.Errorf(number_args, cmd))
			return
		}
		f, err := strconv.ParseFloat(words[1], 64)
		if w(err) != nil {
			return
		}
		for _, user := range users {
			for _, cube := range user.cubes {
				if cube != nil {
					cube.nod = f
				}
			}
		}

	default:

		w(fmt.Errorf("Invalid command from GUI: %s", cmd))

	}
}

func nodEnhance(rotV, enhance float64) float64 {
	if enhance >= -1.0 && enhance <= 1.0 {
		return rotV * enhance
	}

	sign := 1.0
	if math.Signbit(enhance) {
		sign = -1.0
		enhance = -enhance
	}

	var v float64
	if rotV < 0 {
		v = -0.5 * math.Pi * (1.0 - math.Pow(1.0+rotV*2.0/math.Pi, enhance))
	} else {
		v = 0.5 * math.Pi * (1.0 - math.Pow(1.0-rotV*2.0/math.Pi, enhance))
	}
	return sign * v
}
