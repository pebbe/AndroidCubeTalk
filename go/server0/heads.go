package main

import (
	"fmt"
)

var (
	sethead = make([][]int, len(cubes))
	head    = make([]int, len(cubes))
)

func initHeads() {
	for i, cube := range cubes {
		sethead[i] = make([]int, 0, len(cubes)-1)
		head[i] = cube.head
		for j := range cubes {
			if i != j {
				sethead[i] = append(sethead[i], j)
			}
		}
	}
}

func resetHeads(user int) {
	sethead[user] = sethead[user][0:0]
	for i := range cubes {
		if i != user {
			sethead[user] = append(sethead[user], i)
		}
	}
}

func showHeads(ch chan string, idx int) {
	for _, i := range sethead[idx] {
		users[idx].n[cntrHead]++
		ch <- fmt.Sprintf("head %s %d %d\n", cubes[i].uid, users[idx].n[cntrHead], head[i])
	}
	sethead[idx] = sethead[idx][0:0]
}

func setHead(idx int, f int) {
	head[idx] = f
	for i := range cubes {
		if i != idx {
			sethead[i] = append(sethead[i], idx)
		}
	}
}
