package main

import (
	"github.com/kr/pretty"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

type tXYZ struct {
	x, y, z float64
}

type jsLayout struct {
	AudioHandler  string   `json:"audio_handler"`
	ClickHandler  string   `json:"click_handler"`
	ChoiceHandler string   `json:"choice_handler"`
	Looking       bool     `json:"looking"`
	Looked        bool     `json:"looked"`
	Cubes         []jsCube `json:"cubes"`
}

type jsCube struct {
	Uid   string    `json:"uid"`
	Pos   []float64 `json:"pos"`
	Color string    `json:"color"`
	Head  int       `json:"head"`
	Face  int       `json:"face"`
	Gui   bool      `json:"gui"`
}

// This has data on how a user sees another cube, except for actual head movement
type tCube struct {
	uid     string
	pos     tXYZ // position
	forward tXYZ // neutral forward direction, unit vector, with y=0
	towards tXYZ // unit vector from user to this cube
	color   tRGB
	head    int // texture number
	face    int // texture number
	sees    []string
	gui     bool
}

type tUser struct {
	uid       string
	needSetup bool
	selfZ     float64 // position on z-axis
	lookat    tXYZ    // direction the user is looking at, unit vector
	roll      float64 // rotation around the direction of lookat, between -180 and 180
	audio     float64
	cubes     []*tCube // other cubes, where and how as seen by this user
	n         [numberOfCtrs]uint64
}

var (
	// layout is built from this list
	cubes = []tCube{
		tCube{
			uid:   "A",
			pos:   tXYZ{0, 0, 1},
			color: tRGB{1, .6, .6}, // red
			head:  0,
			face:  0,
			gui:   true,
		},
		tCube{
			uid:   "B",
			pos:   tXYZ{1, 0, 0},
			color: tRGB{0, .6, 0}, // green
			head:  1,
			face:  1,
			gui:   true,
		},
		tCube{
			uid:   "C",
			pos:   tXYZ{0, 0, -1},
			color: tRGB{.4, .7, 1}, // blue
			head:  2,
			face:  2,
			gui:   true,
		},
		tCube{
			uid:   "D",
			pos:   tXYZ{-1, 0, 0},
			color: tRGB{.7, .7, .7}, // grey
			head:  3,
			face:  3,
			gui:   true,
		},
	}

	users  = make([]*tUser, len(cubes))
	labels = make(map[string]int)

	firstMakeUsers = true
)

func makeUsers() {

	oldCounters := make(map[string][numberOfCtrs]uint64)

	if firstMakeUsers {
		firstMakeUsers = false

		loadUsers()

		for i := range cubes {
			cubes[i].pos.x *= *opt_d
			cubes[i].pos.y *= *opt_d
			cubes[i].pos.z *= *opt_d
		}

		// this will be redone by robotUserSetup() if masking is used
		for i := range cubes {
			sees := make([]string, 0, len(cubes)-1)
			for j := range cubes {
				if i != j {
					sees = append(sees, cubes[j].uid)
				}
			}
			cubes[i].sees = sees
		}
	} else {
		for _, user := range users {
			oldCounters[user.uid] = user.n
		}
	}

	// If with robot: layout shuffling
	// If with masking: redo 'sees'
	robotUserSetup()

	labelstrings := make([]string, 0)

	for i, cube := range cubes {
		labels[cube.uid] = i
		labelstrings = append(labelstrings, fmt.Sprint(cube.uid, ":", i))
	}

	// create layout for each user from list of cubes
	for i, cube := range cubes {

		user := tUser{
			uid:    cube.uid,
			selfZ:  math.Sqrt(cube.pos.x*cube.pos.x + cube.pos.z*cube.pos.z), // horizontal distance from y-axis
			lookat: tXYZ{0, 0, -1},                                           // initially looking at y-axis
			roll:   0,                                                        // initially no roll
			cubes:  make([]*tCube, len(cubes)),
			n:      oldCounters[cube.uid],
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

func loadUsers() {
	if *opt_l == "" {
		return
	}

	data, err := ioutil.ReadFile(*opt_l)
	x(err)
	var layout jsLayout
	x(json.Unmarshal(data, &layout))

	markLookingAtMe = layout.Looked
	markLookingAtThem = layout.Looking

	if markLookingAtMe && markLookingAtThem {
		x(fmt.Errorf("You can't use both options 'looking' and 'looked'"))
	}

	var ok bool
	if audioHandle, ok = audioHandlers[layout.AudioHandler]; !ok {
		x(fmt.Errorf("Unknown audio handler"))
	}
	if clickHandle, ok = clickHandlers[layout.ClickHandler]; !ok {
		x(fmt.Errorf("Unknown click handler"))
	}
	if choiceHandle, ok = choiceHandlers[layout.ChoiceHandler]; !ok {
		x(fmt.Errorf("Unknown choice handler"))
	}

	cubes = cubes[0:0]

	for i, c := range layout.Cubes {

		c.Uid = strings.TrimSpace(c.Uid)
		c.Color = strings.TrimSpace(c.Color)

		if c.Color == "" {
			c.Color = "white"
		}

		if c.Uid == "" {
			x(fmt.Errorf("Missing uid in file %q, item number %d", *opt_l, i))
		}

		if len(c.Pos) != 3 {
			x(fmt.Errorf("Wrong number of position values in file %q for item %q", *opt_l, c.Uid))
		}

		color, ok := colornames[c.Color]
		if !ok {
			x(fmt.Errorf("Unknown color in file %q for item %q", *opt_l, c.Uid))
		}

		if c.Head < 0 || c.Head > 9 {
			x(fmt.Errorf("Invalid head number in file %q for item %q (must be 0 - 9)", *opt_l, c.Uid))
		}

		if c.Face < 0 || c.Face > 9 {
			x(fmt.Errorf("Invalid face number in file %q for item %q (must be 0 - 9)", *opt_l, c.Uid))
		}

		cube := tCube{
			uid:   c.Uid,
			pos:   tXYZ{c.Pos[0], c.Pos[1], c.Pos[2]},
			color: color,
			head:  c.Head,
			face:  c.Face,
			gui:   c.Gui,
		}

		cubes = append(cubes, cube)

	}

	users = make([]*tUser, len(cubes))
}
