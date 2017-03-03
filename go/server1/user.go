package main

import (
	"github.com/kr/pretty"

	"fmt"
	"math"
)

const (
	DISTANCE       = 4
	NR_OF_COUNTERS = 6 // see counters in API
)

type tVector struct {
	x, y, z float64
}

// This has data on how a user sees another cube, except for actual head movement
type tCube struct {
	uid     string
	pos     tVector // position
	forward tVector // neutral forward direction, unit vector, with y=0
	red     float64
	green   float64
	blue    float64
	nod     float64
}

type tUser struct {
	init   bool    // init is done?
	selfZ  float64 // position on z-axis
	lookat tVector // direction the user is looking at, unit vector
	roll   float64 // rotation around the direction of lookat, between -180 and 180
	cubes  []tCube // other cubes, where and how as seen by this user
	n      [NR_OF_COUNTERS]uint64
}

var (
	users = make(map[string]*tUser)

	// lay-out is built from this list
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

	chLog <- fmt.Sprintf("I Global lay-out: %# v", pretty.Formatter(cubes))

	// create lay-out for each user from list of cubes
	for i, cube := range cubes {

		user := tUser{
			init:   true,                                                     // done at first, but undone when user sends 'reset' command
			selfZ:  math.Sqrt(cube.pos.x*cube.pos.x + cube.pos.z*cube.pos.z), // horizontal distance from y-axis
			lookat: tVector{0, 0, -1},                                        // initially looking at y-axis
			roll:   0,                                                        // initially no roll
			cubes:  make([]tCube, 0, len(cubes)-1),
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

					nod: 1, // standard value for nodding: no amplification
				}
				user.cubes = append(user.cubes, c)
			}
		}

		users[cube.uid] = &user

		// Send lay-out for user to logger
		chLog <- fmt.Sprintf("I User %s: %# v", cube.uid, pretty.Formatter(user))
	}
}
