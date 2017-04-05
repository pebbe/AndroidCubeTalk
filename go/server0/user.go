package main

import (
	"github.com/kr/pretty"

	"fmt"
	"math"
	"strings"
)

type tXYZ struct {
	x, y, z float64
}

type tRGB struct {
	r, g, b float64
}

// This has data on how a user sees another cube, except for actual head movement
type tCube struct {
	uid     string
	pos     tXYZ // position
	forward tXYZ // neutral forward direction, unit vector, with y=0
	towards tXYZ // unit vector from user to this cube
	color   tRGB
	head    int      // texture number
	face    int      // texture number
	sees    []string // names of other cubes seen by this one
	isRobot bool
}

type tUser struct {
	uid       string
	needSetup bool
	selfZ     float64 // position on z-axis
	lookat    tXYZ    // direction the user is looking at, unit vector
	roll      float64 // rotation around the direction of lookat, between -180 and 180
	isRobot   bool
	cubes     []*tCube // other cubes, where and how as seen by this user
	n         [numberOfCtrs]uint64
}

var (
	// layout is built from this list
	cubes = []tCube{
		tCube{
			uid:   "A",
			pos:   tXYZ{0, 0, 1},
			color: tRGB{.6, .6, .6},
			head:  0,
			face:  0,
			sees:  []string{"B", "C"},
		},
		tCube{
			uid:   "B",
			pos:   tXYZ{.866, 0, -.5},
			color: tRGB{.6, .6, 6},
			head:  0,
			face:  0,
			sees:  []string{"A", "C"},
		},
		tCube{
			uid:     "C",
			pos:     tXYZ{-.866, 0, -.5},
			color:   tRGB{.6, .6, .6},
			head:    0,
			face:    0,
			sees:    []string{"A", "B"},
			isRobot: true,
		},
	}

	users  = make([]*tUser, len(cubes))
	labels = make(map[string]int)
)

func makeUsers() {

	for i := range cubes {
		cubes[i].pos.x *= *opt_d
		cubes[i].pos.y *= *opt_d
		cubes[i].pos.z *= *opt_d
	}

	labelstrings := make([]string, 0)

	for i, cube := range cubes {
		labels[cube.uid] = i
		labelstrings = append(labelstrings, fmt.Sprint(cube.uid, ":", i))
	}

	// create layout for each user from list of cubes
	for i, cube := range cubes {

		user := tUser{
			uid:     cube.uid,
			selfZ:   math.Sqrt(cube.pos.x*cube.pos.x + cube.pos.z*cube.pos.z), // horizontal distance from y-axis
			lookat:  tXYZ{0, 0, -1},                                           // initially looking at y-axis
			roll:    0,                                                        // initially no roll
			cubes:   make([]*tCube, len(cubes)),
			isRobot: cube.isRobot,
		}

		rotH0 := math.Atan2(cube.pos.x, cube.pos.z)
		Y0 := cube.pos.y

		for _, see := range cube.sees {
			j := labels[see]
			cube2 := cubes[j]
			rotH := math.Atan2(cube2.pos.x, cube2.pos.z) - rotH0
			l := math.Sqrt(cube2.pos.x*cube2.pos.x + cube2.pos.z*cube2.pos.z)
			c := tCube{
				uid:   cube2.uid,
				color: cube2.color,
				head:  cube2.head,
				face:  cube2.face,

				pos: tXYZ{
					l * math.Sin(rotH),
					cube2.pos.y - Y0,
					l * math.Cos(rotH),
				},

				// assumption: each cube is looking horizontally towards its own y-axis
				forward: tXYZ{
					-math.Sin(rotH),
					0,
					-math.Cos(rotH),
				},
			}
			dx := c.pos.x
			dy := c.pos.y
			dz := c.pos.z - user.selfZ
			ln := math.Sqrt(dx*dx + dy*dy + dz*dz)
			c.towards = tXYZ{dx / ln, dy / ln, dz / ln}
			user.cubes[j] = &c
		}

		users[i] = &user

	}

	chLog <- fmt.Sprintf("I UIDs: map[%v]", strings.Join(labelstrings, " "))

	chLog <- fmt.Sprintf("I Global layout: %# v", pretty.Formatter(cubes))

	// Send layout for user to logger
	chLog <- fmt.Sprintf("I User layout: %# v", pretty.Formatter(users))

}
