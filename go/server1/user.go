package main

import (
	"github.com/kr/pretty"

	"fmt"
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
			pos:   tVector{0, 0, 4},
			red:   1, // white
			green: 1,
			blue:  1,
		},
		tCube{
			uid:   "B",
			pos:   tVector{4, 0, 0},
			red:   1, // yellow
			green: 1,
			blue:  0,
		},
		tCube{
			uid:   "C",
			pos:   tVector{0, 0, -4},
			red:   0, // green
			green: .6,
			blue:  0,
		},
		tCube{
			uid:   "D",
			pos:   tVector{-4, 0, 0},
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
			selfZ:  length(cube.pos),
			lookat: tVector{0, 0, -1},
			roll:   0,
			cubes:  make([]tCube, 0, len(cubes)-1),
		}
		rotH0 := math.Atan2(cube.pos.x, cube.pos.z)
		rotV0 := math.Atan2(-cube.pos.y, math.Sqrt(cube.pos.x*cube.pos.x+cube.pos.z*cube.pos.z))
		for j, cube := range cubes {
			if j != i {
				rotH := math.Atan2(cube.pos.x, cube.pos.z) - rotH0
				rotV := math.Atan2(cube.pos.y, math.Sqrt(cube.pos.x*cube.pos.x+cube.pos.z*cube.pos.z)) - rotV0
				fmt.Println(cubes[i].uid, cubes[j].uid, rotH/math.Pi*180, rotV/math.Pi*180)
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
					// assumption: each cube is looking towards its own origin
					forward: tVector{
						-math.Sin(rotH) * math.Cos(rotV),
						-math.Sin(rotV),
						-math.Cos(rotH) * math.Cos(rotV),
					},
					nod: 1,
				}
				user.cubes = append(user.cubes, c)
			}
		}
		users[cube.uid] = &user
	}
	pretty.Println(users)
}

func length(v tVector) float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}
