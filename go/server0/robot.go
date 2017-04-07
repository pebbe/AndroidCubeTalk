package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var (
	robotDelay  = 3 * time.Second // time to confirm
	robotResult = 3 * time.Second // time to display result
	robotBlank  = 1 * time.Second // time to blank after result before reset

	isRobot      string
	useRobot     bool
	robotThen    time.Time
	robotRunning bool
	rLookAt      = make([]int, len(cubes))
	rCurrent     int
)

func hasRobot() bool {
	return *opt_b != ""
}

func initRobot() {
	useRobot = hasRobot()
	if useRobot {
		for i := 0; i < len(cubes); i++ {
			rLookAt[i] = -1
		}
		robotRunning = false
		rCurrent = -1
	}
}

func doRobot(me int) {
	if !useRobot || robotRunning {
		return
	}

	rLookAt[me] = -1
	found := -1
	for them, cube := range users[me].cubes {
		if cube != nil && isLookingAt(me, them) {
			rLookAt[me] = them
			found = them
			break
		}
	}
	if found < -1 {
		rCurrent = -1
		return
	}

	if found != rCurrent {
		count := 0
		for _, i := range rLookAt {
			if i == found {
				count++
			}
		}
		if count == len(cubes)-1 {
			rCurrent = found
			robotThen = time.Now().Add(robotDelay)
		}
	}
	if rCurrent < 0 {
		return
	}
	if time.Now().Before(robotThen) {
		return
	}

	robotRunning = true
	useLookAt = false

	a := "Wrong"
	if users[rCurrent].uid == isRobot {
		a = "Correct"
	}
	chLog <- fmt.Sprintf("I Users selected %s: %s", users[rCurrent].uid, a)
	fmt.Printf("Users selected %s: %s\n", users[rCurrent].uid, a)

	for i, user := range users {
		ch := chOut[i]
		if user.uid == isRobot {
			setFace(i, 9)
		}
		for j, cube := range user.cubes {
			if cube == nil {
				continue
			}
			if cube.uid == isRobot {
				user.n[cntrHead]++
				select {
				case ch <- fmt.Sprintf("head %s %d 9\n", cube.uid, user.n[cntrHead]):
				default:
					// drop if channel is full
				}
			}
			color := lightgrey
			if j == rCurrent && cube.uid != isRobot {
				color = red
			}
			user.n[cntrColor]++
			select {
			case ch <- fmt.Sprintf("color %s %d %g %g %g\n", cube.uid, user.n[cntrColor], color.r, color.g, color.b):
			default:
				// drop if channel is full
			}
		}
	}
	go func() {

		time.Sleep(robotResult)

		chCmd <- "stop"
		chCmd <- "hideall"

		time.Sleep(robotBlank)

		chCmd <- "restart"

	}()
}

func runRobot() {
	if !hasRobot() {
		return
	}

	chLog <- "B Starting robot: " + *opt_b

	cmd := exec.Command(*opt_b)
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
	isRobot = a[1]

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

	chLog <- "B Robot " + *opt_b + " has stopped"
	fmt.Println("Robot", *opt_b, "has stopped")
}
