package main

import (
	"fmt"
)

var (
	markLookingAtMe   = true
	markLookingAtThem = false

	lookColor = "1.0 0.7 0.0" // orange

	lookMarked [][]bool
)

func initLooking() {
	lookMarked = make([][]bool, len(users))
	for i := range users {
		lookMarked[i] = make([]bool, len(users))
	}
}

func resetLooking(user int) {
	for i := range users {
		lookMarked[user][i] = false
	}
}

func showLooking(ch chan string, me int) {

	user := users[me]

	for i, cube := range user.cubes {
		if i == me {
			continue
		}

		if markLookingAtMe {

			v := users[i].lookat
			w := users[i].cubes[me].towards
			if v.x*w.x+v.y*w.y+v.z*w.z > *opt_t {
				if !lookMarked[me][i] {
					lookMarked[me][i] = true
					user.n[4]++
					ch <- fmt.Sprintf("color %s %d %s\n", cube.uid, user.n[4], lookColor)
					chLog <- fmt.Sprintf("I Begin %s looked at by %s", user.uid, cube.uid)
				}
			} else {
				if lookMarked[me][i] {
					lookMarked[me][i] = false
					user.n[4]++
					ch <- fmt.Sprintf("color %s %d %g %g %g\n", cube.uid, user.n[4], cube.color.r, cube.color.g, cube.color.b)
					chLog <- fmt.Sprintf("I End %s looked at by %s", user.uid, cube.uid)
				}
			}

		}

		if markLookingAtThem {

			v := user.lookat
			w := user.cubes[i].towards
			if v.x*w.x+v.y*w.y+v.z*w.z > *opt_t {
				if !lookMarked[me][i] {
					lookMarked[me][i] = true
					user.n[4]++
					ch <- fmt.Sprintf("color %s %d %s\n", cube.uid, user.n[4], lookColor)
					chLog <- fmt.Sprintf("I Begin %s looking at %s", user.uid, cube.uid)
				}
			} else {
				if lookMarked[me][i] {
					lookMarked[me][i] = false
					user.n[4]++
					ch <- fmt.Sprintf("color %s %d %g %g %g\n", cube.uid, user.n[4], cube.color.r, cube.color.g, cube.color.b)
					chLog <- fmt.Sprintf("I End %s looking at %s", user.uid, cube.uid)
				}
			}

		}
	}
}
