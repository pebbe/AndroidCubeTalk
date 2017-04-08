package main

import (
	"fmt"
)

var (
	setface [][]int
	face    []int
)

func initFaces() {
	setface = make([][]int, len(cubes))
	face = make([]int, len(cubes))
	for i, cube := range cubes {
		setface[i] = make([]int, 0, len(cubes)-1)
		face[i] = cube.face
		for j := range cubes {
			if i != j {
				setface[i] = append(setface[i], j)
			}
		}
	}
}

func resetFaces(user int) {
	setface[user] = setface[user][0:0]
	for i := range cubes {
		if i != user {
			setface[user] = append(setface[user], i)
		}
	}
}

func showFaces(ch chan string, idx int) {
	for _, i := range setface[idx] {
		users[idx].n[cntrFace]++
		ch <- fmt.Sprintf("face %s %d %d\n", cubes[i].uid, users[idx].n[cntrFace], face[i])
	}
	setface[idx] = setface[idx][0:0]
}

func setFace(idx int, f int) {
	face[idx] = f
	for i := range cubes {
		if i != idx {
			setface[i] = append(setface[i], idx)
		}
	}
}
