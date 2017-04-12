package main

import (
	"fmt"
)

var (
	// don't set both to true
	markLookingAtMe   = false
	markLookingAtThem = true

	lookColor = "1.0 0.7 0.0" // orange

	lookMarked [][]bool
	useLookAt  bool
)

func isLookingAt(from, to int) bool {
	if users[from].cubes == nil {
		return false
	}
	v := users[from].lookat
	w := users[from].cubes[to].towards
	return v.x*w.x+v.y*w.y+v.z*w.z > settings.Tolerance
}

func initLooking() {
	lookMarked = make([][]bool, len(users))
	for i := range users {
		lookMarked[i] = make([]bool, len(users))
	}
	useLookAt = true
}

func resetLooking(user int) {
	for i := range users {
		lookMarked[user][i] = false
	}
}

func showLooking(ch chan string, me int) {

	if !useLookAt {
		return
	}

	user := users[me]

	for i, cube := range user.cubes {
		if cube == nil {
			continue
		}

		if markLookingAtMe {

			if isLookingAt(i, me) {
				if !lookMarked[me][i] {
					lookMarked[me][i] = true
					user.n[cntrColor]++
					ch <- fmt.Sprintf("color %s %d %s\n", cube.uid, user.n[cntrColor], lookColor)
					chLog <- fmt.Sprintf("I Begin %s looked at by %s", user.uid, cube.uid)
				}
			} else {
				if lookMarked[me][i] {
					lookMarked[me][i] = false
					user.n[cntrColor]++
					ch <- fmt.Sprintf("color %s %d %g %g %g\n", cube.uid, user.n[cntrColor], cube.color.r, cube.color.g, cube.color.b)
					chLog <- fmt.Sprintf("I End %s looked at by %s", user.uid, cube.uid)
				}
			}

		}

		if markLookingAtThem {

			if isLookingAt(me, i) {
				if !lookMarked[me][i] {
					lookMarked[me][i] = true
					user.n[cntrColor]++
					ch <- fmt.Sprintf("color %s %d %s\n", cube.uid, user.n[cntrColor], lookColor)
					chLog <- fmt.Sprintf("I Begin %s looking at %s", user.uid, cube.uid)
				}
			} else {
				if lookMarked[me][i] {
					lookMarked[me][i] = false
					user.n[cntrColor]++
					ch <- fmt.Sprintf("color %s %d %g %g %g\n", cube.uid, user.n[cntrColor], cube.color.r, cube.color.g, cube.color.b)
					chLog <- fmt.Sprintf("I End %s looking at %s", user.uid, cube.uid)
				}
			}

		}
	}
}
