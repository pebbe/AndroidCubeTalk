package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func controller() {
	chInA := chIn["A"]
	chInB := chIn["B"]
	chInC := chIn["C"]
	chInD := chIn["D"]
	chInE := chIn["E"]
	chInF := chIn["F"]
	chOutA := chOut["A"]
	chOutB := chOut["B"]
	chOutC := chOut["C"]
	chOutD := chOut["D"]
	chOutE := chOut["E"]
	chOutF := chOut["F"]

	for {
		select {
		case <-chQuit:
			return
		case cmd := <-chCmd:
			chLog <- "C " + cmd
			handleCmd(cmd)
		case cmd := <-chInA:
			chLog <- "I A " + cmd
			handleIn(cmd, "A", chOutA)
		case cmd := <-chInB:
			chLog <- "I B " + cmd
			handleIn(cmd, "B", chOutB)
		case cmd := <-chInC:
			chLog <- "I C " + cmd
			handleIn(cmd, "C", chOutC)
		case cmd := <-chInD:
			chLog <- "I D " + cmd
			handleIn(cmd, "D", chOutD)
		case cmd := <-chInE:
			chLog <- "I E " + cmd
			handleIn(cmd, "E", chOutE)
		case cmd := <-chInF:
			chLog <- "I F " + cmd
			handleIn(cmd, "F", chOutF)
		}
	}
}

func handleIn(cmd string, uid string, chOut chan string) {
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
			select {
			case chOut <- buf.String():
				user.init = true
			default:
			}
		}

		user.n3++
		for _, cube := range user.cubes {
			if cube.uid != uid {

				l := users[cube.uid].lookat
				f := cube.forward

				rotH0 := math.Atan2(l.x, -l.z)
				rotH := math.Atan2(f.x, f.z) - rotH0

				rotV0 := math.Atan2(-l.y, math.Sqrt(l.x*l.x+l.z*l.z))
				rotV := between(
					math.Atan2(f.y, math.Sqrt(f.x*f.x+f.z*f.z))*cube.nod-rotV0,
					-math.Pi/2+.001,
					math.Pi/2-.001)
				select {
				case chOut <- fmt.Sprintf("lookat %s %d %g %g %g %g\n",
					cube.uid,
					user.n3,
					math.Sin(rotH)*math.Cos(rotV),
					math.Sin(rotV),
					math.Cos(rotH)*math.Cos(rotV),
					users[cube.uid].roll):
				default:
				}
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
