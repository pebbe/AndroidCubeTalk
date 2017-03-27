package main

import (
	"math"
)

type shakePar struct {
	shake   float64
	prevIn  float64
	prevOut float64
}

var (
	shaking [][]shakePar
)

func initShaking() {
	shaking = make([][]shakePar, len(users))
	for i := range users {
		shaking[i] = make([]shakePar, len(users))
		for j := range users {
			shaking[i][j] = shakePar{1, 0, 0}
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

func doShake(me, them int, currentIn float64) float64 {

	// do immediate shake
	dr := (currentIn - shaking[me][them].prevIn)
	for dr > math.Pi {
		dr -= 2 * math.Pi
	}
	for dr < -math.Pi {
		dr += 2 * math.Pi
	}
	currentOut := shaking[me][them].prevOut + dr*shaking[me][them].shake

	// delay to actual angle
	dr = currentIn - currentOut
	for dr > math.Pi {
		dr -= 2 * math.Pi
	}
	for dr < -math.Pi {
		dr += 2 * math.Pi
	}
	currentOut += dr * .1

	shaking[me][them].prevOut = currentOut
	shaking[me][them].prevIn = currentIn

	return currentOut
}
