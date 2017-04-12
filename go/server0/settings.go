package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type jsSettings struct {
	Port int `json:"port"`

	Looking   bool    `json:"looking"`
	Looked    bool    `json:"looked"`
	Tolerance float64 `json:"tolerance"`

	Audio        bool   `json:"audio"`
	AudioHandler string `json:"audio_handler"`

	ClickHandler  string `json:"click_handler"`
	ChoiceHandler string `json:"choice_handler"`

	Robot        string `json:"robot"`
	RobotMasking bool   `json:"robot_masking"`

	Users          []string `json:"users"`
	UnitDistance   float64  `json:"unit_distance"`
	DefaultColor   string   `json:"default_color"`
	DefaultSkipGui bool     `json:"default_skip_gui"`
	Cubes          []jsCube `json:"cubes"`
}

type jsCube struct {
	Uid   string    `json:"uid"`
	Pos   []float64 `json:"pos"`
	Color string    `json:"color"`
	Head  int       `json:"head"`
	Face  int       `json:"face"`
	Gui   bool      `json:"gui"`
}

var (
	settings jsSettings
)

func readSettings(filename string) {
	var ok bool

	data, err := ioutil.ReadFile(filename)
	x(err)
	x(json.Unmarshal(data, &settings))

	if settings.Port < 1 {
		settings.Port = 8448
	}

	markLookingAtMe = settings.Looked
	markLookingAtThem = settings.Looking
	if markLookingAtMe && markLookingAtThem {
		x(fmt.Errorf("You can't use both options 'looking' and 'looked'"))
	}
	if settings.Tolerance <= 0 {
		settings.Tolerance = .99
	}

	if settings.Audio {
		withAudio = true
		settings.AudioHandler = strings.TrimSpace(settings.AudioHandler)
		if settings.AudioHandler == "" {
			settings.AudioHandler = "none"
		}
		if audioHandle, ok = audioHandlers[settings.AudioHandler]; !ok {
			x(fmt.Errorf("Unknown audio handler"))
		}
	} else {
		settings.AudioHandler = "none"
	}

	settings.ClickHandler = strings.TrimSpace(settings.ClickHandler)
	if settings.ClickHandler == "" {
		settings.ClickHandler = "none"
	}
	if clickHandle, ok = clickHandlers[settings.ClickHandler]; !ok {
		x(fmt.Errorf("Unknown click handler"))
	}

	settings.ChoiceHandler = strings.TrimSpace(settings.ChoiceHandler)
	if settings.ChoiceHandler == "" {
		settings.ChoiceHandler = "none"
	}
	if choiceHandle, ok = choiceHandlers[settings.ChoiceHandler]; !ok {
		x(fmt.Errorf("Unknown choice handler"))
	}

	settings.Robot = strings.TrimSpace(settings.Robot)
	if settings.Robot == "" {
		settings.RobotMasking = false
	} else {
		withRobot = true
	}
	withMasking = settings.RobotMasking

	hasUsers := (settings.Users != nil && len(settings.Users) > 0)
	hasCubes := (settings.Cubes != nil && len(settings.Cubes) > 0)
	if hasCubes && hasUsers {
		x(fmt.Errorf("You can't define both users and cubes"))
	}
	if !(hasUsers || hasCubes) {
		x(fmt.Errorf("You need to define users or cubes in oyur settings"))
	}

	if settings.UnitDistance <= 0 {
		settings.UnitDistance = 4
	}

	settings.DefaultColor = strings.TrimSpace(settings.DefaultColor)
	if settings.DefaultColor == "" {
		settings.DefaultColor = "lightgrey"
	}

	if hasCubes {
		settings.Users = make([]string, len(settings.Cubes))
		for i, c := range settings.Cubes {
			settings.Users[i] = c.Uid
		}
	} else {
		settings.Cubes = make([]jsCube, len(settings.Users))
		n := len(settings.Users)
		if settings.RobotMasking {
			n--
		}
		for i, u := range settings.Users {
			var x, y, z float64
			if i == n {
				x = 0
				y = 1
				z = 2
			} else {
				r := math.Pi * 2.0 / float64(n) * float64(i)
				x = math.Sin(r)
				y = 0
				z = math.Cos(r)
			}
			settings.Cubes[i] = jsCube{
				Uid: u,
				Pos: []float64{x, y, z},
				Gui: !settings.DefaultSkipGui,
			}
		}
	}

	for _, c := range settings.Cubes {
		if c.Color == "" {
			c.Color = settings.DefaultColor
		}
	}

	cubes = make([]tCube, 0, len(settings.Cubes))

	for i, c := range settings.Cubes {

		c.Uid = strings.TrimSpace(c.Uid)
		c.Color = strings.TrimSpace(c.Color)

		if c.Color == "" {
			c.Color = settings.DefaultColor
		}

		if c.Uid == "" {
			x(fmt.Errorf("Missing uid in file %q, item number %d", os.Args[1], i))
		}

		if len(c.Pos) != 3 {
			x(fmt.Errorf("Wrong number of position values in file %q for item %q", os.Args[1], c.Uid))
		}

		color, ok := colornames[c.Color]
		if !ok {
			x(fmt.Errorf("Unknown color in file %q for item %q", os.Args[1], c.Uid))
		}

		if c.Head < 0 || c.Head > 9 {
			x(fmt.Errorf("Invalid head number in file %q for item %q (must be 0 - 9)", os.Args[1], c.Uid))
		}

		if c.Face < 0 || c.Face > 9 {
			x(fmt.Errorf("Invalid face number in file %q for item %q (must be 0 - 9)", os.Args[1], c.Uid))
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

	b, err := json.MarshalIndent(settings, "    ", "    ")
	x(err)
	chLog <- fmt.Sprint("I Settings:\n    ", string(b))
}
