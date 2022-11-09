package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file is expected in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz is in seconds")
	flag.Parse()

	_ = timeLimit
	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("File couldn't be opened : %s\n", *csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit("Couldnt read CSV file")
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	//<-timer.C
	problems := ParseFile(lines)

	correct := 0
	for i, v := range problems {
		fmt.Printf("Problem #%d : %s\n", i+1, v.q)
		answerChannel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("You got %d out of %d questions right", correct, len(problems))
			return
		default:

			if <-answerChannel == v.a {
				correct++
				fmt.Println("Correct!")
			} else {
				fmt.Println("Ouch, Wrong!")
			}
		}

	}
	fmt.Printf("You got %d out of %d questions right", correct, len(problems))
	//_ = csvFileName
}

func ParseFile(lines [][]string) []Problems {
	ret := make([]Problems, len(lines))

	for i, v := range lines {
		ret[i] = Problems{
			q: v[0],
			a: strings.TrimSpace(v[1]),
		}
	}

	return ret
}

type Problems struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
