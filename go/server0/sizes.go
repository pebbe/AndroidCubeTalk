package main

import (
	"fmt"
)

var (
	setsize  []bool
	cubesize [3]float64
)

func initSize() {
	setsize = make([]bool, len(users))
	cubesize[0] = 1
	cubesize[1] = 1
	cubesize[2] = 1
}

func resetSize(user int) {
	setsize[user] = true
}

func setSize(w, h, d float64) {
	cubesize[0] = w
	cubesize[1] = h
	cubesize[2] = d
	for i := range users {
		setsize[i] = true
		users[i].n[6]++
	}
}

func showSize(ch chan string, user int) {
	if setsize[user] {
		setsize[user] = false
		ch <- fmt.Sprintf("cubesize %d %g %g %g\n", users[user].n[6], cubesize[0], cubesize[1], cubesize[2])
	}
}
