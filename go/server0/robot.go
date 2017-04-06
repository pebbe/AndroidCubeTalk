package main

import (
	"fmt"
	"time"
)

var (
	robotDelay  = 3 * time.Second // time to confirm
	robotResult = 3 * time.Second // time to display result
	robotBlank  = 1 * time.Second // time to blank after result before reset

	useRobot     bool
	robotThen    time.Time
	robotRunning bool
	rLookAt      = make([]int, len(cubes))
	rCurrent     int
)

func hasRobots() bool {
	for _, cube := range cubes {
		if cube.isRobot {
			return true
		}
	}
	return false
}

func initRobot() {
	useRobot = hasRobots()
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
	if users[rCurrent].isRobot {
		a = "Correct"
	}
	chLog <- fmt.Sprintf("I Users selected %s: %s", users[rCurrent].uid, a)
	fmt.Printf("Users selected %s: %s\n", users[rCurrent].uid, a)

	for i, user := range users {
		ch := chOut[i]
		if user.isRobot {
			setFace(i, 9)
		}
		for j, cube := range user.cubes {
			if cube == nil {
				continue
			}
			if cube.isRobot {
				user.n[cntrHead]++
				select {
				case ch <- fmt.Sprintf("head %s %d 9\n", cube.uid, user.n[cntrHead]):
				default:
					// drop if channel is full
				}
			}
			color := lightgrey
			if j == rCurrent && !cube.isRobot {
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
