package main

var (
	shaking1 [][]float64
)

func initShaking1() {
	shaking1 = make([][]float64, len(users))
	for i := range users {
		shaking1[i] = make([]float64, len(users))
		for j := range users {
			shaking1[i][j] = 1
		}
	}
}

func setShake1All(shake float64) {
	for i := range users {
		for j := range users {
			setShake1(i, j, shake)
		}
	}
}

func setShake1(me, them int, shake float64) {
	shaking1[me][them] = shake
}

func doShake1(me, them int, angle float64) float64 {
	return angle * shaking1[me][them]
}
