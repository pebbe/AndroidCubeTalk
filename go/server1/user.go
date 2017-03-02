package main

import (
	"math"
)

type tVector struct {
	x, y, z float64
}

type tCube struct {
	uid     string
	pos     tVector
	forward tVector // unit vector
	red     float64
	green   float64
	blue    float64
	nod     float64
}

type tUser struct {
	uid    string
	selfZ  float64
	lookat tVector // unit vector
	roll   float64 // between -90 and 90
	cubes  []tCube
	init   bool
	n0     uint64
	n1     uint64
	n2     uint64
	n3     uint64
	n4     uint64
	n5     uint64
}

var (
	users = make(map[string]*tUser)
	cubes = []tCube{
		tCube{
			uid:   "A",
			pos:   tVector{0, 0, 4},
			red:   1,
			green: 1,
			blue:  1,
		},
		tCube{
			uid:   "B",
			pos:   tVector{4, 0, 0},
			red:   1,
			green: 1,
			blue:  0,
		},
		tCube{
			uid:   "C",
			pos:   tVector{0, 0, -4},
			red:   0,
			green: .6,
			blue:  0,
		},
		tCube{
			uid:   "D",
			pos:   tVector{-4, 0, 0},
			red:   .4,
			green: .7,
			blue:  1,
		},
		tCube{
			uid:   "E",
			pos:   tVector{2, 3, -4},
			red:   1,
			green: .6,
			blue:  .6,
		},
	}
)

func init() {
	for i, cube := range cubes {
		user := tUser{
			uid:    cube.uid,
			selfZ:  length(cube.pos),
			lookat: tVector{0, 0, -1},
			roll:   0,
			cubes:  make([]tCube, 0, len(cubes)-1),
		}
		rotH0 := math.Atan2(cube.pos.x, cube.pos.z)
		rotH1 := math.Atan2(cube.pos.x, -cube.pos.z)
		rotV0 := math.Atan2(cube.pos.y, math.Sqrt(cube.pos.x*cube.pos.x+cube.pos.z*cube.pos.z))
		for j, cube := range cubes {
			if j != i {
				rotH := math.Atan2(cube.pos.x, cube.pos.z) - rotH0
				rotH2 := math.Atan2(cube.pos.x, -cube.pos.z) - rotH1
				rotV := math.Atan2(cube.pos.y, math.Sqrt(cube.pos.x*cube.pos.x+cube.pos.z*cube.pos.z)) - rotV0
				l := length(cube.pos)
				c := tCube{
					uid:   cube.uid,
					red:   cube.red,
					green: cube.green,
					blue:  cube.blue,
					pos: tVector{
						l * math.Sin(rotH) * math.Cos(rotV),
						l * math.Sin(rotV),
						l * math.Cos(rotH) * math.Cos(rotV),
					},
					forward: tVector{
						math.Sin(rotH2) * math.Cos(rotV),
						math.Sin(rotV),
						-math.Cos(rotH2) * math.Cos(rotV),
					},
					nod: 1,
				}
				user.cubes = append(user.cubes, c)
			}
		}
		users[cube.uid] = &user
	}
}

func length(v tVector) float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}
