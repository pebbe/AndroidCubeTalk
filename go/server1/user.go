package main

import (
	"github.com/kr/pretty"

	"fmt"
	"math"
)

const (
	DISTANCE = 4
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
	roll   float64 // between -180 and 180
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
			pos:   tVector{0, 0, DISTANCE},
			red:   1, // white
			green: 1,
			blue:  1,
		},
		tCube{
			uid:   "B",
			pos:   tVector{DISTANCE, 0, 0},
			red:   1, // yellow
			green: 1,
			blue:  0,
		},
		tCube{
			uid:   "C",
			pos:   tVector{0, 0, -DISTANCE},
			red:   0, // green
			green: .6,
			blue:  0,
		},
		tCube{
			uid:   "D",
			pos:   tVector{-DISTANCE, 0, 0},
			red:   .4, // blue
			green: .7,
			blue:  1,
		},
		tCube{
			uid:   "E",
			pos:   tVector{2, 3, -4},
			red:   1, // red
			green: .6,
			blue:  .6,
		},
		tCube{
			uid:   "F",
			pos:   tVector{-2, -3, 4},
			red:   .5, // grey
			green: .5,
			blue:  .5,
		},
	}
)

func init() {
	for i, cube := range cubes {
		user := tUser{
			uid:    cube.uid,
			init:   true,
			selfZ:  math.Sqrt(cube.pos.x*cube.pos.x + cube.pos.z*cube.pos.z),
			lookat: tVector{0, 0, -1},
			cubes:  make([]tCube, 0, len(cubes)-1),
			roll:   0,
		}
		rotH0 := math.Atan2(cube.pos.x, cube.pos.z)
		Y0 := cube.pos.y
		for j, cube := range cubes {
			if j != i {
				rotH := math.Atan2(cube.pos.x, cube.pos.z) - rotH0
				l := math.Sqrt(cube.pos.x*cube.pos.x + cube.pos.z*cube.pos.z)
				c := tCube{
					uid:   cube.uid,
					red:   cube.red,
					green: cube.green,
					blue:  cube.blue,
					pos: tVector{
						l * math.Sin(rotH),
						cube.pos.y - Y0,
						l * math.Cos(rotH),
					},
					// assumption: each cube is looking horizontally towards its own y-axis
					forward: tVector{
						-math.Sin(rotH),
						0,
						-math.Cos(rotH),
					},
					nod: 1,
				}
				user.cubes = append(user.cubes, c)
			}
		}
		users[cube.uid] = &user
		chLog <- fmt.Sprintf("User %s: %# v\n", cube.uid, pretty.Formatter(user))
	}
}
