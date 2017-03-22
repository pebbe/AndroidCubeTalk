package main

var (
	tilting [][]float64
)

func initTilting() {
	tilting = make([][]float64, len(users))
	for i := range users {
		tilting[i] = make([]float64, len(users))
		for j := range users {
			tilting[i][j] = 1
		}
	}
}

func setTiltAll(shake float64) {
	for i := range users {
		for j := range users {
			setTilt(i, j, shake)
		}
	}
}

func setTilt(me, them int, shake float64) {
	tilting[me][them] = shake
}

func doTilt(me, them int, angle float64) float64 {
	return angle * tilting[me][them]
}
