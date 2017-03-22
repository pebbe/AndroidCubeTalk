package main

var (
	shaking [][]float64
)

func initShaking() {
	shaking = make([][]float64, len(users))
	for i := range users {
		shaking[i] = make([]float64, len(users))
		for j := range users {
			shaking[i][j] = 1
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
	shaking[me][them] = shake
}

func doShake(me, them int, angle float64) float64 {
	return angle * shaking[me][them]
}
