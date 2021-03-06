package main

import (
	"fmt"
	"strings"
	"time"
)

var (
	choiceHandlers = map[string]func(int, string, string){
		"":     choiceNone,
		"none": choiceNone,
		"demo": choiceDemo,
	}
	choiceHandle = choiceHandlers[""]
)

func infoMakeNotice(user int, lines []string) {
	users[user].n[cntrInfo]++
	chOut[user] <- fmt.Sprintf("info %d %d\n%s\n",
		users[user].n[cntrInfo],
		len(lines),
		strings.Join(lines, "\n"))
}

func infoMakeChoice(user int, infoID string, opt1, opt2 string, lines []string) {
	users[user].n[cntrInfo]++
	chOut[user] <- fmt.Sprintf("info %d %d %s %s %s\n%s\n",
		users[user].n[cntrInfo],
		len(lines),
		infoID,
		opt1,
		opt2,
		strings.Join(lines, "\n"))
}

func choiceNone(user int, infoID string, choice string) {
}

func choiceDemo(user int, infoID string, choice string) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		infoMakeNotice(user, []string{fmt.Sprintf("You clicked %s for %s", choice, infoID)})
	}()
}
