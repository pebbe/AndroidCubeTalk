package main

import (
	//	"github.com/kr/pretty"

	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type jsonRGB struct {
	R, G, B float64
}

type jsonXYZ struct {
	X, Y, Z float64
}

// This has data on how a user sees another cube, except for actual head movement
type jsonCube struct {
	Uid     string   `json:"uid"`
	Pos     jsonXYZ  `json:"pos"`
	Forward jsonXYZ  `json:"forward"`
	Towards jsonXYZ  `json:"towards"`
	Color   jsonRGB  `json:"color"`
	Head    int      `json:"head"`
	Face    int      `json:"face"`
	Sees    []string `json:"sees"`
	Gui     bool     `json:"gui"`
}

type jsonUser struct {
	Uid       string               `json:"uid"`
	Needsetup bool                 `json:"needSetup"`
	Selfz     float64              `json:"selfZ"`
	Lookat    jsonXYZ              `json:"lookat"`
	Roll      float64              `json:"roll"`
	Audio     float64              `json:"audio"`
	Cubes     []*jsonCube          `json:"cubes"`
	N         [numberOfCtrs]uint64 `json:"n"`
}

var (
	reCommandLine = regexp.MustCompile(` I Command line: \[\]string{".*?", "(.*)"}`)
	reLogLine     = regexp.MustCompile(`^[0-9]+:[0-9]+.[0-9]+`)

	reHex        = regexp.MustCompile("0x[0-9a-fA-F]+")
	reKey        = regexp.MustCompile("[a-zA-Z]+:")
	reCubes      = regexp.MustCompile(`cubes:\s*{`)
	reN1         = regexp.MustCompile(`},\s*n:\s*{`)
	reN2         = regexp.MustCompile(`n:\s*\[.*}`)
	reComma      = regexp.MustCompile(`,\s*}`)
	reSees       = regexp.MustCompile(`sees:.*}`)
	reNil        = regexp.MustCompile(`\bnil\b`)
	reCloseBrace = regexp.MustCompile(`^\s*}\s*$`)
)

func configReplay(filename string) {
	withReplay = true

	fp, err := os.Open(filename)
	x(err)
	rd, err := gzip.NewReader(fp)
	x(err)

	defer func() {
		rd.Close()
		fp.Close()
	}()

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		m := reCommandLine.FindStringSubmatch(scanner.Text())
		if m != nil && len(m) == 2 {
			readConfig(m[1])
			config.Script = nil
			config.Robot = ""
			config.RobotMasking = false
			withRobot = false
			withMasking = false
			return
		}
	}

	x(fmt.Errorf("No command line found in %s", filename))
}

func replay(filename string) {
	time.Sleep(time.Millisecond * 100)
	fmt.Print("Press ENTER to start...")
	bufio.NewReader(os.Stdin).ReadLine()

	fp, err := os.Open(filename)
	x(err)
	rd, err := gzip.NewReader(fp)
	x(err)

	defer func() {
		rd.Close()
		fp.Close()
	}()

	scanner := bufio.NewScanner(rd)

	scanner.Scan()
	logStart, err := time.Parse("04:05.000", strings.Fields(scanner.Text())[0])
	x(err)
	logPrev := logStart

	replayStart := time.Now()

	for scanner.Scan() {
		line := scanner.Text()
		if reLogLine.MatchString(line) {
			words := strings.Fields(line)
			if len(words) < 3 {
				continue
			}
			t, err := time.Parse("04:05.000", words[0])
			if err != nil {
				continue
			}
			for t.Before(logPrev) {
				t.Add(time.Hour)
			}
			logPrev = t
			time.Sleep(t.Sub(logStart) - time.Now().Sub(replayStart))

			switch words[1] {
			case "R":
				uid := words[2]
				chReplay <- tRequest{
					uid: uid,
					idx: labels[uid],
					req: strings.Join(words[3:], " "), // no newline
				}
			case "C":
				chCmd <- strings.Join(words[2:], " ")
			case "I":
				if len(words) == 4 && words[3] == "layout:" {
					if words[2] == "Global" {
						//	replaceGlobalLayout(scanner)
					} else if words[2] == "User" {
						// this is instead of command "restart" in file "controller.go"
						replaceUserLayout(scanner)
						for _, user := range users {
							user.needSetup = true
						}
						scriptStart()
						started = true
					} else {
						fmt.Println(line)
					}
				} else {
					fmt.Println(line)
				}
			case "B":
				fmt.Println(line)
			}
		}
	}
}

func replaceUserLayout(scanner *bufio.Scanner) {
	var buf bytes.Buffer

	for scanner.Scan() {
		line := scanner.Text()
		buf.WriteString(line + "\n")
		if reCloseBrace.MatchString(line) {
			break
		}
	}

	data := buf.String()

	data = strings.Replace(data, "[]*main.tUser{", "[", -1)
	data = strings.Replace(data, "(*main.tCube)(nil)", "nil", -1)
	data = strings.Replace(data, "&main.tUser", "", -1)
	data = strings.Replace(data, "&main.tCube", "", -1)
	data = strings.Replace(data, "main.tXYZ{}", "{x:0,y:0,z:0}", -1)
	data = strings.Replace(data, "main.tRGB{}", "{r:0,g:0,b:0}", -1)
	data = strings.Replace(data, "main.tXYZ", "", -1)
	data = strings.Replace(data, "main.tRGB", "", -1)

	data = reComma.ReplaceAllStringFunc(data, func(s string) string {
		return s[1:]
	})

	data = reHex.ReplaceAllStringFunc(data, func(s string) string {
		a, _ := strconv.ParseInt(s[2:], 16, 32)
		return fmt.Sprint(a)
	})

	data = reCubes.ReplaceAllLiteralString(data, "cubes: [")

	data = reN1.ReplaceAllLiteralString(data, "],\n\t    n: [")

	data = reNil.ReplaceAllLiteralString(data, "null")

	data = reN2.ReplaceAllStringFunc(data, func(s string) string {
		return s[:len(s)-1] + `]`
	})

	data = reSees.ReplaceAllStringFunc(data, func(s string) string {
		return strings.Replace(strings.Replace(s, "{", "[", 1), "}", "]", 1)
	})

	data = reKey.ReplaceAllStringFunc(data, func(s string) string {
		return `"` + s[:len(s)-1] + `":`
	})

	data = strings.TrimSpace(data)
	data = data[:len(data)-1] + "]"

	fmt.Println(data)

	p := []jsonUser{}
	err := json.Unmarshal([]byte(data), &p)
	x(err)

	// pretty.Println(p)

	for i, user := range users {
		for j, cube := range p[i].Cubes {
			if cube == nil {
				user.cubes[j] = nil
			} else {
				user.cubes[j] = &tCube{
					uid:     cube.Uid,
					pos:     tXYZ{x: cube.Pos.X, y: cube.Pos.Y, z: cube.Pos.Z},
					forward: tXYZ{x: cube.Forward.X, y: cube.Forward.Y, z: cube.Forward.Z},
					towards: tXYZ{x: cube.Towards.X, y: cube.Towards.Y, z: cube.Towards.Z},
					color:   tRGB{r: cube.Color.R, g: cube.Color.G, b: cube.Color.B},
					head:    cube.Head,
					face:    cube.Face,
					sees:    cube.Sees,
					gui:     cube.Gui,
				}
			}
		}
	}

	// pretty.Println(users)
}

func replaceGlobalLayout(scanner *bufio.Scanner) {
	var buf bytes.Buffer

	for scanner.Scan() {
		line := scanner.Text()
		buf.WriteString(line + "\n")
		if reCloseBrace.MatchString(line) {
			break
		}
	}

	data := buf.String()

	data = strings.Replace(data, "[]main.tCube{", "[", -1)
	data = strings.Replace(data, "main.tXYZ{}", "{x:0,y:0,z:0}", -1)
	data = strings.Replace(data, "main.tRGB{}", "{r:0,g:0,b:0}", -1)
	data = strings.Replace(data, "main.tXYZ", "", -1)
	data = strings.Replace(data, "main.tRGB", "", -1)

	data = reComma.ReplaceAllStringFunc(data, func(s string) string {
		return s[1:]
	})

	data = reHex.ReplaceAllStringFunc(data, func(s string) string {
		a, _ := strconv.ParseInt(s[2:], 16, 32)
		return fmt.Sprint(a)
	})

	data = reSees.ReplaceAllStringFunc(data, func(s string) string {
		return strings.Replace(strings.Replace(s, "{", "[", 1), "}", "]", 1)
	})

	data = reKey.ReplaceAllStringFunc(data, func(s string) string {
		return `"` + s[:len(s)-1] + `":`
	})

	data = strings.TrimSpace(data)
	data = data[:len(data)-1] + "]"

	fmt.Println(data)

	p := []jsonCube{}
	err := json.Unmarshal([]byte(data), &p)
	x(err)

	// pretty.Println(p)

	for i, cube := range p {
		cubes[i] = tCube{
			uid:     cube.Uid,
			pos:     tXYZ{x: cube.Pos.X, y: cube.Pos.Y, z: cube.Pos.Z},
			forward: tXYZ{x: cube.Forward.X, y: cube.Forward.Y, z: cube.Forward.Z},
			towards: tXYZ{x: cube.Towards.X, y: cube.Towards.Y, z: cube.Towards.Z},
			color:   tRGB{r: cube.Color.R, g: cube.Color.G, b: cube.Color.B},
			head:    cube.Head,
			face:    cube.Face,
			sees:    cube.Sees,
			gui:     cube.Gui,
		}
	}

	// pretty.Println(cubes)
}
