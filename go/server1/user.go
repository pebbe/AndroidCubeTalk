package main

type tVector struct {
	x, y, z float64
}

type tCube struct {
	uid     string
	pos     tVector
	forward tVector // unit vector
}

type tUser struct {
	uid    string
	selfZ  float64
	lookat tVector // unit vector
	roll   float64 // between -90 and 90
	cubes  []tCube
}

var (
	users = map[string]tUser{
		"A": {
			uid:    "A",
			selfZ:  4,
			lookat: tVector{0, 0, -1}, // unit vector
			roll:   0,
			cubes: []tCube{
				tCube{
					uid:     "B",
					pos:     tVector{0, 0, -4},
					forward: tVector{0, 0, 1}, // unit vector
				},
			},
		},
		"B": {
			uid:    "B",
			selfZ:  4,
			lookat: tVector{0, 0, -1}, // unit vector
			roll:   0,
			cubes: []tCube{
				tCube{
					uid:     "A",
					pos:     tVector{0, 0, -4},
					forward: tVector{0, 0, 1}, // unit vector
				},
			},
		},
	}
)
