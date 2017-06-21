package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type jsConfig struct {
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

	Script       []string `json:"script"`
	ScriptRepeat bool     `json:"script_repeat"`

	Users          []string `json:"users"`
	UnitDistance   float64  `json:"unit_distance"`
	DefaultColor   string   `json:"default_color"`
	DefaultFace    int      `json:"default_face"`
	DefaultHead    int      `json:"default_head"`
	DefaultSkipGui bool     `json:"default_skip_gui"`
	Cubes          []jsCube `json:"cubes"`
}

type jsCube struct {
	Uid   string    `json:"uid"`
	Pos   []float64 `json:"pos"`
	Color string    `json:"color"`
	Face  int       `json:"face"`
	Head  int       `json:"head"`
	Gui   bool      `json:"gui"`
}

var (
	config jsConfig
)

func readConfig(filename string) {
	var ok bool

	data, err := ioutil.ReadFile(filename)
	x(err)

	if data[0] == 0x1F && data[1] == 0x8B {
		configReplay(filename)
		return
	}

	x(json.Unmarshal(data, &config))

	if config.Port < 1 {
		config.Port = 8448
	}

	markLookingAtMe = config.Looked
	markLookingAtThem = config.Looking
	if markLookingAtMe && markLookingAtThem {
		x(fmt.Errorf("You can't use both options 'looking' and 'looked'"))
	}
	if config.Tolerance <= 0 {
		config.Tolerance = .99
	}

	if config.Audio {
		withAudio = true
		config.AudioHandler = strings.TrimSpace(config.AudioHandler)
		if config.AudioHandler == "" {
			config.AudioHandler = "none"
		}
		if audioHandle, ok = audioHandlers[config.AudioHandler]; !ok {
			x(fmt.Errorf("Unknown audio handler"))
		}
	} else {
		config.AudioHandler = "none"
	}

	config.ClickHandler = strings.TrimSpace(config.ClickHandler)
	if config.ClickHandler == "" {
		config.ClickHandler = "none"
	}
	if clickHandle, ok = clickHandlers[config.ClickHandler]; !ok {
		x(fmt.Errorf("Unknown click handler"))
	}

	config.ChoiceHandler = strings.TrimSpace(config.ChoiceHandler)
	if config.ChoiceHandler == "" {
		config.ChoiceHandler = "none"
	}
	if choiceHandle, ok = choiceHandlers[config.ChoiceHandler]; !ok {
		x(fmt.Errorf("Unknown choice handler"))
	}

	config.Robot = strings.TrimSpace(config.Robot)
	if config.Robot == "" {
		config.RobotMasking = false
	} else {
		withRobot = true
	}
	withMasking = config.RobotMasking

	hasUsers := (config.Users != nil && len(config.Users) > 0)
	hasCubes := (config.Cubes != nil && len(config.Cubes) > 0)
	if hasCubes && hasUsers {
		x(fmt.Errorf("You can't define both users and cubes"))
	}
	if !(hasUsers || hasCubes) {
		x(fmt.Errorf("You need to define users or cubes in oyur config"))
	}

	if config.UnitDistance <= 0 {
		config.UnitDistance = 4
	}

	config.DefaultColor = strings.TrimSpace(config.DefaultColor)
	if config.DefaultColor == "" {
		config.DefaultColor = "lightgrey"
	}

	if hasCubes {
		config.Users = make([]string, len(config.Cubes))
		for i, c := range config.Cubes {
			config.Users[i] = c.Uid
		}
	} else {
		config.Cubes = make([]jsCube, len(config.Users))
		n := len(config.Users)
		if config.RobotMasking {
			n--
		}
		for i, u := range config.Users {
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
			config.Cubes[i] = jsCube{
				Uid:  u,
				Pos:  []float64{x, y, z},
				Face: config.DefaultFace,
				Head: config.DefaultHead,
				Gui:  !config.DefaultSkipGui,
			}
		}
	}

	for _, c := range config.Cubes {
		if c.Color == "" {
			c.Color = config.DefaultColor
		}
	}

	cubes = make([]tCube, 0, len(config.Cubes))

	for i, c := range config.Cubes {

		c.Uid = strings.TrimSpace(c.Uid)
		c.Color = strings.TrimSpace(c.Color)

		if c.Color == "" {
			c.Color = config.DefaultColor
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

	b, err := json.MarshalIndent(config, "    ", "    ")
	x(err)
	chLog <- fmt.Sprint("I Config:\n    ", string(b))
}
