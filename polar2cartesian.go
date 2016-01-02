package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
)

type polar struct {
	radius float64
	O			 float64
}

type cartesian struct {
	x float64
	y float64
}

var prompt = "Enter a radius and an angle (in degrees), e.g., 12.5 90, " + "or %s to quit."
func init() {
	if runtime.GOOS == "windows" {
		prompt = fmt.Sprintf(prompt, "Ctrl+Z, Enter")
	} else {
		prompt = fmt.Sprintf(prompt, "Ctrl+D")
	}
}

func main() {
	questions := make(chan polar)
	defer close(questions)

	answers := createSlover(questions)
	defer close(answers)

	interact(questions, answers)
}

func createSlover(questions chan polar) chan cartesian {
	answers := make(chan cartesian)

	go func() {
		for {
			polarCoord := <-questions
			O := polarCoord.O * math.Pi / 180

			x := polarCoord.radius * math.Cos(O)
			y := polarCoord.radius * math.Sin(O)

			answers <- cartesian{x, y}
		}
	}()

	return answers
}

const result = "Polar radius=%.02f angle=%.02f° ® Cartesian x=%.02f y=%.02f\n"
func interact(questions chan polar, answers chan cartesian) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)

	for {
		fmt.Printf("Radius and angle: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		var radius, O float64
		if _, err := fmt.Sscanf(line, "%f %f", &radius, &O); err != nil {
			fmt.Fprintf(os.Stderr, "invalid input \n")
			continue
		}
		questions <-polar{radius, O}
		coord := <-answers
		fmt.Printf(result, radius, O, coord.x, coord.y )
	}

	fmt.Println()
}