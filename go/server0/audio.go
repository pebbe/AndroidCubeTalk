package main

import (
	"fmt"
)

var (
	setaudio []bool

	audioHandlers = map[string]func(int){
		"":     audioNone,
		"none": audioNone,
	}
	audioHandle = audioHandlers[""]
)

func initAudio() {
	if !withAudio {
		return
	}
	setaudio = make([]bool, len(cubes))
}

func resetAudio(idx int) {
	if !withAudio {
		return
	}
	setaudio[idx] = withAudio
}

func showAudio(ch chan string, idx int) {
	if !withAudio {
		return
	}
	if setaudio[idx] {
		setaudio[idx] = false
		user := users[idx]
		user.n[cntrAudio]++
		ch <- fmt.Sprintf("audio %d on\n", user.n[cntrAudio])
	}
}

func doAudio(idx int) {
	if !withAudio {
		return
	}
	audioHandle(idx)
}

func audioNone(idx int) {
}
