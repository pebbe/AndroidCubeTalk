package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"
)

func main() {

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

		r += .1
		if r > math.Pi {
			r -= 2 * math.Pi
		}
		fmt.Fprintf(os.Stdout, "lookat %g 0 %g 0 0\n", math.Sin(r), math.Cos(r))
		os.Stdout.Sync()

		for scanner.Scan() {
			line := scanner.Text()
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

		if i%20 == 0 {
			r := float64(i%100) / 100
			g := 1.0 - r
			b := .7
			fmt.Fprintf(os.Stdout, "command_quiet color BOT %f %f %f\n", r, g, b)
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
