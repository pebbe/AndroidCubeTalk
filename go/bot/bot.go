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

		time.Sleep(40 * time.Millisecond)
	}

}
