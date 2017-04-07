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

	fmt.Println("join A")
	if !scanner.Scan() {
		return
	}

	fmt.Println("reset")
	if !scanner.Scan() {
		return
	}

	go func() {
		i := 0
		for {
			fmt.Fprintln(os.Stderr, "I AM ROBOT")
			os.Stderr.Sync()

			time.Sleep(20 * time.Second)

			i++
			if i == 3 {
				fmt.Fprintln(os.Stderr, "CRASH!")
				os.Stderr.Sync()
				os.Exit(1)
			}
		}
	}()

	r := 0.0
	for {

		r += .1
		if r > math.Pi {
			r -= 2 * math.Pi
		}
		fmt.Fprintf(os.Stdout, "lookat %g 0 %g 0 0\n", math.Sin(r), math.Cos(r))
		os.Stdout.Sync()

		for scanner.Scan() {
			if scanner.Text() == "." {
				break
			}
		}

		time.Sleep(40 * time.Millisecond)
	}

}
