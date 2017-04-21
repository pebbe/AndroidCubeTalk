package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

var (
	robotDelay  = 3 * time.Second // time to confirm
	robotResult = 3 * time.Second // time to display result
	robotBlank  = 1 * time.Second // time to blank after result before reset

	masked        = -1
	robotUID      string
	robotThen     time.Time
	robotSelected bool
	rLookAt       []int
	rCurrent      int
)

func initRobot() {
	if !withRobot {
		return
	}

	rLookAt = make([]int, len(cubes))
	for i := 0; i < len(cubes); i++ {
		rLookAt[i] = -1
	}
	robotSelected = false
	rCurrent = -1
}

func robotUserSetup() {

	if !withRobot {
		return
	}

	initRobot()

	rand.Seed(time.Now().UnixNano())

	if !withMasking {
		// shuffle positions
		for i := len(cubes) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			cubes[i].pos, cubes[j].pos = cubes[j].pos, cubes[i].pos
		}
		return
	}

	// shuffle positions without the last one, which is supposed to be the robot
	for i := len(cubes) - 2; i > 0; i-- {
		j := rand.Intn(i + 1)
		cubes[i].pos, cubes[j].pos = cubes[j].pos, cubes[i].pos
	}

	bot := len(cubes) - 1
	masked = rand.Intn(bot) // range: [0, bot>

	chLog <- "B Robot is masking " + cubes[masked].uid
	fmt.Println("Robot is masking", cubes[masked].uid)

	cubes[bot].pos = cubes[masked].pos

	// masked cube sees all other cubes, except for bot
	sees := make([]string, 0)
	for i := 0; i < len(cubes)-1; i++ {
		if i != masked {
			sees = append(sees, cubes[i].uid)
		}
	}
	cubes[masked].sees = sees

	// other cubes (including bot) see all other cubes, except masked cube
	for i := range cubes {
		if i == masked {
			continue
		}
		sees := make([]string, 0)
		for j := range cubes {
			if i != j && j != masked {
				sees = append(sees, cubes[j].uid)
			}
		}
		cubes[i].sees = sees
	}

}

func doRobot(me int) {
	if !withRobot || robotSelected || me == masked {
		return
	}

	// who is 'me' looking at?
	rLookAt[me] = -1
	found := -1
	for them, cube := range users[me].cubes {
		if cube != nil && isLookingAt(me, them) {
			rLookAt[me] = them
			found = them
			break
		}
	}

	required := len(cubes) - 1
	if withMasking {
		required--
	}

	// test
	counts := make([]int, len(users))
	found = -1
	for _, i := range rLookAt {
		if i >= 0 {
			counts[i]++
			if counts[i] == required {
				found = i
				break
			}
		}
	}
	if found < 0 {
		rCurrent = -1
		return
	}

	if found != rCurrent {
		rCurrent = found
		robotThen = time.Now().Add(robotDelay)
	}

	if time.Now().Before(robotThen) {
		return
	}

	robotSelected = true

	a := "Wrong"
	if users[rCurrent].uid == robotUID {
		a = "Correct"
	}
	chLog <- fmt.Sprintf("B Users selected %s: %s", users[rCurrent].uid, a)
	fmt.Printf("Users selected %s: %s\n", users[rCurrent].uid, a)

	rcube := cubes[len(cubes)-1]
	chCmdQuiet <- "face " + rcube.uid + " 9"
	chCmdQuiet <- "head " + rcube.uid + " 9"

	useLookAt = false
	for i, cube := range cubes {
		color := cube.color
		if i == rCurrent && cube.uid != robotUID {
			color = colornames["red"]
		} else if i == len(cubes)-1 {
			color = colornames["lightred"]
		}
		chCmdQuiet <- fmt.Sprintf("color %s %g %g %g", cube.uid, color.r, color.g, color.b)
	}

	go func(current int) {

		time.Sleep(robotResult)

		chCmdQuiet <- "stop"
		chCmdQuiet <- "hideall"

		time.Sleep(robotBlank)

		chCmdQuiet <- fmt.Sprintf("face %s %d", rcube.uid, rcube.face)
		chCmdQuiet <- fmt.Sprintf("head %s %d", rcube.uid, rcube.head)

		cube := cubes[current]
		color := cube.color
		chCmdQuiet <- fmt.Sprintf("color %s %g %g %g", cube.uid, color.r, color.g, color.b)

		initRobot()

		useLookAt = true

		chCmdQuiet <- "restart"

	}(rCurrent)
}

func runRobot() {
	if !withRobot {
		return
	}

	chLog <- "B Starting robot: " + config.Robot

	words := strings.Fields(config.Robot)
	cmd := exec.Command(words[0], words[1:]...)
	stdin, err := cmd.StdinPipe()
	x(err)
	stdout, err := cmd.StdoutPipe()
	x(err)
	stderr, err := cmd.StderrPipe()
	x(err)

	w(cmd.Start())

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			chLog <- "B Error: " + line
			fmt.Println("ROBOT error:", line)
		}
		if err := scanner.Err(); err != nil {
			chLog <- "B Error reading from stderr: " + err.Error()
			fmt.Println("ROBOT Error reading from stderr:", err)
		}
	}()

	scanner := bufio.NewScanner(stdout)
	scanner.Scan()
	line := scanner.Text()
	a := strings.Fields(line)
	if len(a) != 2 || a[0] != "join" {
		return
	}
	robotUID = a[1]

	uid := a[1]

	idx, ok := labels[uid]
	if !ok {
		return
	}

	out := chOut[idx]

	fmt.Fprintln(stdin, ".")

	fmt.Println("      ROBOT =", uid)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "quit" {
			break
		}

		chIn <- tRequest{
			uid: uid,
			idx: idx,
			req: line, // no newline
		}

		for busy := true; busy; {
			select {
			case txt := <-out: // including newline
				fmt.Fprint(stdin, txt) // no newline
			default:
				busy = false
			}
		}
		fmt.Fprintln(stdin, ".")
	}
	w(scanner.Err())

	w(cmd.Wait())

	chLog <- "B Stopped robot: " + config.Robot
	fmt.Println("Stopped robot:", config.Robot)
}
