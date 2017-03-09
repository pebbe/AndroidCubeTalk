package main

import (
	"github.com/kr/pretty"

	"fmt"
	"math"
	"strings"
)

const (
	DISTANCE       = 4
	NR_OF_COUNTERS = 6 // see counters in API
)

type tXYZ struct {
	x, y, z float64
}

type tRGB struct {
	r, g, b float64
}

// This has data on how a user sees another cube, except for actual head movement
type tCube struct {
	uid         string
	pos         tXYZ // position
	forward     tXYZ // neutral forward direction, unit vector, with y=0
	towards     tXYZ // unit vector from user to this cube
	color       tRGB
	lookingatme bool
	nod         float64
}

type tUser struct {
	uid    string
	init   bool     // init is done?
	selfZ  float64  // position on z-axis
	lookat tXYZ     // direction the user is looking at, unit vector
	roll   float64  // rotation around the direction of lookat, between -180 and 180
	cubes  []*tCube // other cubes, where and how as seen by this user
	n      [NR_OF_COUNTERS]uint64
}

var (
	// layout is built from this list
	cubes = []tCube{
		tCube{
			uid:   "A",
			pos:   tXYZ{0, 0, DISTANCE},
			color: tRGB{1, 1, 1}, // white
		},
		tCube{
			uid:   "B",
			pos:   tXYZ{DISTANCE, 0, 0},
			color: tRGB{1, 1, 0}, // yellow
		},
		tCube{
			uid:   "C",
			pos:   tXYZ{0, 0, -DISTANCE},
			color: tRGB{0, .6, 0}, // green
		},
		tCube{
			uid:   "D",
			pos:   tXYZ{-DISTANCE, 0, 0},
			color: tRGB{.4, .7, 1}, // blue
		},
		tCube{
			uid:   "E",
			pos:   tXYZ{2, 3, -4},
			color: tRGB{1, .6, .6}, // red
		},
		tCube{
			uid:   "F",
			pos:   tXYZ{-2, -3, 4},
			color: tRGB{.5, .5, .5}, // grey
		},
	}

	users  = make([]*tUser, len(cubes))
	labels = make(map[string]int)
)

func init() {

	labelstrings := make([]string, 0)

	// create layout for each user from list of cubes
	for i, cube := range cubes {

		labels[cube.uid] = i
		labelstrings = append(labelstrings, fmt.Sprint(cube.uid, ":", i))

		user := tUser{
			uid:    cube.uid,
			init:   true,                                                     // done at first, but undone when user sends 'reset' command
			selfZ:  math.Sqrt(cube.pos.x*cube.pos.x + cube.pos.z*cube.pos.z), // horizontal distance from y-axis
			lookat: tXYZ{0, 0, -1},                                           // initially looking at y-axis
			roll:   0,                                                        // initially no roll
			cubes:  make([]*tCube, len(cubes)),
		}

		rotH0 := math.Atan2(cube.pos.x, cube.pos.z)
		Y0 := cube.pos.y

		for j, cube := range cubes {
			if i != j {
				rotH := math.Atan2(cube.pos.x, cube.pos.z) - rotH0
				l := math.Sqrt(cube.pos.x*cube.pos.x + cube.pos.z*cube.pos.z)
				c := tCube{
					uid:   cube.uid,
					color: cube.color,

					pos: tXYZ{
						l * math.Sin(rotH),
						cube.pos.y - Y0,
						l * math.Cos(rotH),
					},

					// assumption: each cube is looking horizontally towards its own y-axis
					forward: tXYZ{
						-math.Sin(rotH),
						0,
						-math.Cos(rotH),
					},

					nod: 1, // standard value for nodding: no amplification
				}
				dx := c.pos.x
				dy := c.pos.y
				dz := c.pos.z - user.selfZ
				ln := math.Sqrt(dx*dx + dy*dy + dz*dz)
				c.towards = tXYZ{dx / ln, dy / ln, dz / ln}
				user.cubes[j] = &c
			}
		}

		users[i] = &user

	}

	chLog <- fmt.Sprintf("I UIDs: map[%v]", strings.Join(labelstrings, " "))

	chLog <- fmt.Sprintf("I Global layout: %# v", pretty.Formatter(cubes))

	// Send layout for user to logger
	chLog <- fmt.Sprintf("I User layout: %# v", pretty.Formatter(users))

}