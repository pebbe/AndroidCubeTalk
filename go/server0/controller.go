package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	cntrSelfZ = iota
	cntrEnterExit
	cntrMoveto
	cntrLookat
	cntrColor
	cntrInfo
	cntrCubesize
	cntrHead
	cntrFace
	cntrAudio
	numberOfCtrs // this one MUST be last
)

var (
	started = false
)

func restart() {
	initRobot()
	initFaces()
	initSize()
	initNodding()
	initShaking()
	initTilting()
	initLooking()
}

func controller() {

	restart()

	for {
		select {
		case <-chQuit:
			return
		case req := <-chIn:
			chLog <- "R " + req.uid + " " + req.req
			handleReq(req)
		case cmd := <-chCmd:
			chLog <- "C " + cmd
			handleCmd(cmd)
		}
	}
}

// handleReq is not run concurrently, so it must be fast
func handleReq(req tRequest) {

	cmd := req.req
	idx := req.idx
	ch := chOut[idx]
	user := users[idx]

	words := strings.Fields(cmd)
	switch words[0] {

	case "log":

		fmt.Println("Log", req.uid+":", strings.Join(words[1:], " "))

	case "reset":

		// This commands tells the server that the client needs the setup.
		// The setup is not sent now, but with the following 'lookat' command.

		user.needSetup = true

		resetSize(idx)

		resetLooking(idx)

		resetFaces(idx)

	case "lookat":

		if !started {
			return
		}

		if len(words) != 6 && len(words) != 7 {
			w(fmt.Errorf("Invalid number of arguments from %q: %s", req.uid, cmd))
			return
		}

		X, err := strconv.ParseFloat(words[1], 64)
		if w(err) != nil {
			return
		}
		Y, err := strconv.ParseFloat(words[2], 64)
		if w(err) != nil {
			return
		}
		Z, err := strconv.ParseFloat(words[3], 64)
		if w(err) != nil {
			return
		}
		roll, err := strconv.ParseFloat(words[4], 64)
		if w(err) != nil {
			return
		}

		if useAudio {
			/*
				audio, err := strconv.ParseFloat(words[5], 64)
				if w(err) != nil {
					return
				}
			*/
		}

		// only one goroutine modifying these variables, so no sync needed
		user.lookat.x = X
		user.lookat.y = Y
		user.lookat.z = Z
		user.roll = roll

		marked := len(words) == 7
		if marked {
			fmt.Printf("Mark %s -> %g %g %g    %.0f° right   %.0f° up\n",
				req.uid,
				X, Y, Z,
				math.Atan2(X, -Z)/math.Pi*180,
				math.Atan2(Y, math.Sqrt(X*X+Z*Z))/math.Pi*180)
		}

		if user.needSetup {
			// this must be in one batch to make sure that the order is preserved
			var buf bytes.Buffer
			user.n[cntrSelfZ]++
			fmt.Fprintf(&buf, "self %d %g\n", user.n[cntrSelfZ], user.selfZ)
			if useAudio {
				user.n[cntrAudio]++
				fmt.Fprintf(&buf, "audio %d on\n", user.n[cntrAudio])
			}
			for _, cube := range user.cubes {
				if cube != nil {
					user.n[cntrEnterExit]++
					fmt.Fprintf(&buf, "enter %s %d\n", cube.uid, user.n[cntrEnterExit])
					user.n[cntrMoveto]++
					fmt.Fprintf(&buf, "moveto %s %d %g %g %g\n", cube.uid, user.n[cntrMoveto], cube.pos.x, cube.pos.y, cube.pos.z)
					user.n[cntrColor]++
					fmt.Fprintf(&buf, "color %s %d %g %g %g\n", cube.uid, user.n[cntrColor], cube.color.r, cube.color.g, cube.color.b)
					user.n[cntrHead]++
					fmt.Fprintf(&buf, "head %s %d %d\n", cube.uid, user.n[cntrHead], cube.head)
				}
			}
			ch <- buf.String()
			user.needSetup = false
		}

		user.n[cntrLookat]++
		for i, cube := range user.cubes {

			if cube != nil {

				if marked && isLookingAt(idx, i) {
					chLog <- fmt.Sprintf("I Mark %s -> %s", req.uid, cube.uid)
					fmt.Printf("Mark %s -> %s\n", req.uid, cube.uid)
					marked = false
					clickHandle(idx, i)
				}

				l := users[i].lookat
				f := cube.forward

				// assumption: forward is horizontal

				rotH := math.Atan2(l.x, l.z) - math.Atan2(f.x, -f.z)
				rotV := math.Atan2(l.y, math.Sqrt(l.x*l.x+l.z*l.z))
				tilt := users[i].roll

				rotH = doShake(idx, i, rotH)
				rotV = doNod(idx, i, rotV)
				tilt = doTilt(idx, i, tilt)

				ch <- fmt.Sprintf("lookat %s %d %g %g %g %g\n",
					cube.uid,
					user.n[cntrLookat],
					math.Sin(rotH)*math.Cos(rotV),
					math.Sin(rotV),
					math.Cos(rotH)*math.Cos(rotV),
					tilt)
			}
		}

		doRobot(idx)

		showLooking(ch, idx)

		showSize(ch, idx)

		showFaces(ch, idx)

	case "info":

		if len(words) != 3 {
			w(fmt.Errorf("Invalid number of arguments from %q: %s", req.uid, cmd))
			return
		}

		chLog <- fmt.Sprintf("Choice by %s for %s: %s", req.uid, words[1], words[2])
		fmt.Printf("Choice by %s for %s: %s\n", req.uid, words[1], words[2])

		infoHandleChoice(idx, words[1], words[2])

	default:

		w(fmt.Errorf("Invalid command from %q: %s", req.uid, cmd))

	}
}

// handleCmd is not run concurrently, so it must be fast
func handleCmd(cmd string) {

	number_args := "Invalid number of arguments from GUI: %s"

	fmt.Println("Command:", cmd)

	words := strings.Fields(cmd)
	switch words[0] {

	case "start":

		started = true

	case "restart":

		makeUsers()
		restart()
		for idx, user := range users {
			user.needSetup = true
			resetSize(idx)
			resetLooking(idx)
			resetFaces(idx)
		}

		started = true

	case "stop":

		started = false

	case "hideall":

		for i, user := range users {
			ch := chOut[i]
			for _, cube := range user.cubes {
				if cube == nil {
					continue
				}
				user.n[cntrEnterExit]++
				select {
				case ch <- fmt.Sprintf("exit %s %d\n", cube.uid, user.n[cntrEnterExit]):
				default:
					// drop if channel is full
				}
			}
		}

	case "face":

		if len(words) != 3 {
			w(fmt.Errorf(number_args, cmd))
			return
		}

		idx, ok := labels[words[1]]
		if !ok {
			w(fmt.Errorf("Illegal label in command: %s", cmd))
			return
		}

		f, err := strconv.ParseInt(words[2], 10, 16)
		if w(err) != nil {
			return
		}

		setFace(idx, int(f))

	case "cubesize":

		if len(words) != 4 {
			w(fmt.Errorf(number_args, cmd))
			return
		}
		var f [3]float64
		for i := 0; i < 3; i++ {
			var err error
			f[i], err = strconv.ParseFloat(words[i+1], 64)
			if w(err) != nil {
				f[i] = 1
			}
		}
		setSize(f[0], f[1], f[2])

	case "recenter":

		// Orders the headset of a particular user to recenter
		// the direction the head is currently pointing.

		if len(words) != 2 {
			w(fmt.Errorf(number_args, cmd))
			return
		}

		if idx, ok := labels[words[1]]; ok {
			select {
			case chOut[idx] <- "recenter\n":
			default:
				// drop if channel is full
			}
		} else {
			w(fmt.Errorf("Invalid user in command from GUI: %s", cmd))
		}

	case "nod":

		if len(words) != 4 {
			w(fmt.Errorf(number_args, cmd))
			return
		}

		i, oki := labels[words[1]]
		j, okj := labels[words[2]]
		if !(oki && okj) {
			w(fmt.Errorf("Invalid users in command from GUI: %s", cmd))
			return
		}
		f, err := strconv.ParseFloat(words[3], 64)
		if w(err) != nil {
			return
		}
		setNod(i, j, f)

	case "globalnod":

		// Set amplification of nodding for all users, for all cubes they see

		if len(words) != 2 {
			w(fmt.Errorf(number_args, cmd))
			return
		}
		f, err := strconv.ParseFloat(words[1], 64)
		if w(err) != nil {
			return
		}
		setNodAll(f)

	case "shake":

		if len(words) != 4 {
			w(fmt.Errorf(number_args, cmd))
			return
		}

		i, oki := labels[words[1]]
		j, okj := labels[words[2]]
		if !(oki && okj) {
			w(fmt.Errorf("Invalid users in command from GUI: %s", cmd))
			return
		}
		f, err := strconv.ParseFloat(words[3], 64)
		if w(err) != nil {
			return
		}
		setShake(i, j, f)

	case "globalshake":

		// Set amplification of shaking for all users, for all cubes they see

		if len(words) != 2 {
			w(fmt.Errorf(number_args, cmd))
			return
		}
		f, err := strconv.ParseFloat(words[1], 64)
		if w(err) != nil {
			return
		}
		setShakeAll(f)

	case "tilt":

		if len(words) != 4 {
			w(fmt.Errorf(number_args, cmd))
			return
		}

		i, oki := labels[words[1]]
		j, okj := labels[words[2]]
		if !(oki && okj) {
			w(fmt.Errorf("Invalid users in command from GUI: %s", cmd))
			return
		}
		f, err := strconv.ParseFloat(words[3], 64)
		if w(err) != nil {
			return
		}
		setTilt(i, j, f)

	case "globaltilt":

		// Set amplification of tilting for all users, for all cubes they see

		if len(words) != 2 {
			w(fmt.Errorf(number_args, cmd))
			return
		}
		f, err := strconv.ParseFloat(words[1], 64)
		if w(err) != nil {
			return
		}
		setTiltAll(f)

	case "turn":

		if len(words) != 5 {
			w(fmt.Errorf(number_args, cmd))
			return
		}

		i, oki := labels[words[1]]
		j, okj := labels[words[2]]
		k, okk := labels[words[3]]
		if !(oki && okj && okk) {
			w(fmt.Errorf("Invalid users in command from GUI: %s", cmd))
			return
		}
		setTurn(i, j, k, words[4] == "on")

	default:

		w(fmt.Errorf("Invalid command from GUI: %s", cmd))

	}
}
