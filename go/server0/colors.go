package main

import (
	"fmt"
)

type tRGB struct {
	r, g, b float64
}

var (
	setcolor [][]int
	color    []tRGB

	colornames = map[string]tRGB{
		"white":      {1, 1, 1},
		"grey":       {.5, .5, .5},
		"red":        {.8, 0, 0},
		"green":      {0, .6, 0},
		"blue":       {0, 0, 1},
		"brown":      {.5, .5, 0},
		"lightgrey":  {.7, .7, .7},
		"lightred":   {1, .5, .5},
		"lightgreen": {.4, 1, .4},
		"lightblue":  {.6, .6, 1},
		"yellow":     {1, 1, 0},
	}
)

func initColors() {
	setcolor = make([][]int, len(cubes))
	color = make([]tRGB, len(cubes))
	for i, cube := range cubes {
		setcolor[i] = make([]int, 0, len(cubes)-1)
		color[i] = cube.color
		for j := range cubes {
			if i != j {
				setcolor[i] = append(setcolor[i], j)
			}
		}
	}
}

func resetColors(user int) {
	setcolor[user] = setcolor[user][0:0]
	for i := range cubes {
		if i != user {
			setcolor[user] = append(setcolor[user], i)
		}
	}
}

func showColors(ch chan string, idx int) {
	for _, i := range setcolor[idx] {
		users[idx].n[cntrColor]++
		ch <- fmt.Sprintf("color %s %d %g %g %g\n", cubes[i].uid, users[idx].n[cntrColor], color[i].r, color[i].g, color[i].b)
	}
	setcolor[idx] = setcolor[idx][0:0]
}

func setColor(idx int, r, g, b float64) {
	color[idx] = tRGB{r, g, b}
	for i := range cubes {
		if i != idx {
			setcolor[i] = append(setcolor[i], idx)
		}
	}
}
