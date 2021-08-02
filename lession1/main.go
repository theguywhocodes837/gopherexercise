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
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,abswer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		Exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	r.Read()
	lines, err := r.ReadAll()
	if err != nil {
		Exit("Failed to parse the provided CSV files")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		answerChan := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			answerChan <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d\n", correct, len(problems))
			return
		case answer := <-answerChan:
			if strings.TrimSpace(answer) == p.a {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{q: line[0], a: line[1]}
	}
	return ret
}

type problem struct {
	q, a string
}

func Exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
