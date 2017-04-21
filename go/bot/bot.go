package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

func main() {

	running := false

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("join BOT")
	if !scanner.Scan() {
		return
	}

	fmt.Println("reset")
	if !scanner.Scan() {
		return
	}

	r := 0.0
	i := 0
	for {

		r += .05
		if r > math.Pi {
			r -= 2 * math.Pi
		}
		fmt.Fprintf(os.Stdout, "lookat %g 0 %g 0 0\n", math.Sin(r), math.Cos(r))
		os.Stdout.Sync()

		for scanner.Scan() {
			line := scanner.Text()
			if !running && strings.HasPrefix(line, "enter") {
				running = true
				fmt.Fprintln(os.Stderr, "Experiment started")
			} else if running && strings.HasPrefix(line, "exit") {
				running = false
				fmt.Fprintln(os.Stderr, "Experiment stopped")
			}
			if line == "." {
				break
			}
		}

		i++
		if i == 500 {
			fmt.Fprintln(os.Stderr, "I AM ROBOT")
			os.Stderr.Sync()
			i = 0
		}

		if i%10 == 0 {
			fmt.Fprintf(os.Stdout, "command_quiet face BOT %d\n", i%3)
			fmt.Fprintf(os.Stdout, "command_quiet head BOT %d\n", i%3)
			os.Stdout.Sync()
			for scanner.Scan() {
				line := scanner.Text()
				if line == "." {
					break
				}
			}
		}

		time.Sleep(40 * time.Millisecond)
	}

}
