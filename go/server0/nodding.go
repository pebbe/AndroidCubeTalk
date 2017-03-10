package main

import (
	"math"
)

var (
	nodding [][]float64
)

func initNodding() {
	nodding = make([][]float64, len(users))
	for i := range users {
		nodding[i] = make([]float64, len(users))
		for j := range users {
			nodding[i][j] = 1
		}
	}
}

func setNodAll(nod float64) {
	for i := range users {
		for j := range users {
			setNod(i, j, nod)
		}
	}
}

func setNod(me, them int, nod float64) {
	nodding[me][them] = nod
}

func doNod(me, them int, angle float64) float64 {

	nod := nodding[me][them]

	if nod >= -1.0 && nod <= 1.0 {
		return angle * nod
	}

	sign := 1.0
	if math.Signbit(nod) {
		sign = -1.0
		nod = -nod
	}

	var v float64
	if angle < 0 {
		v = -0.5 * math.Pi * (1.0 - math.Pow(1.0+angle*2.0/math.Pi, nod))
	} else {
		v = 0.5 * math.Pi * (1.0 - math.Pow(1.0-angle*2.0/math.Pi, nod))
	}
	return sign * v
}
