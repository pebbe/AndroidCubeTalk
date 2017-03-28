package main

import (
	"math"
)

type shakePar struct {
	shake    float64
	prevTrue float64
	prevFake float64
	useTurn  bool
	turn     float64
}

var (
	shaking [][]shakePar
)

func initShaking() {
	shaking = make([][]shakePar, len(users))
	for i := range users {
		shaking[i] = make([]shakePar, len(users))
		for j := range users {
			shaking[i][j] = shakePar{1, 0, 0, false, 0}
		}
	}
}

func setShakeAll(shake float64) {
	for i := range users {
		for j := range users {
			setShake(i, j, shake)
		}
	}
}

func setShake(me, them int, shake float64) {
	shaking[me][them].shake = shake
}

func setTurn(me, them, to int, useTurn bool) {
	shaking[me][them].useTurn = useTurn
	if useTurn {
		u := users[me]
		shaking[me][them].turn = math.Atan2(
			u.cubes[to].pos.x-u.cubes[them].pos.x,
			u.cubes[to].pos.z-u.cubes[them].pos.z)
	}
}

func doShake(me, them int, currentTrue float64) float64 {

	// do immediate shake
	dr := (currentTrue - shaking[me][them].prevTrue)
	for dr > math.Pi {
		dr -= 2 * math.Pi
	}
	for dr < -math.Pi {
		dr += 2 * math.Pi
	}
	currentFake := shaking[me][them].prevFake + dr*shaking[me][them].shake

	// delay to actual angle
	if shaking[me][them].useTurn {
		dr = shaking[me][them].turn - currentFake
	} else {
		dr = currentTrue - currentFake
	}
	for dr > math.Pi {
		dr -= 2 * math.Pi
	}
	for dr < -math.Pi {
		dr += 2 * math.Pi
	}
	currentFake += dr * .1

	shaking[me][them].prevFake = currentFake
	shaking[me][them].prevTrue = currentTrue

	return currentFake
}
