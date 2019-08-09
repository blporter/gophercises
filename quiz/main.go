package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func parseLines(lines [][]string) []problem {
	toReturn := make([]problem, len(lines))

	for i, line := range lines {
		toReturn[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return toReturn
}

func main() {
	csvFileName := flag.String(
		"csv",
		"problems.csv",
		"a csv file in the format of 'question,answer'")
	timeLimit := flag.Int(
		"limit",
		30,
		"the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV file: %s", *csvFileName))
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		exit("Failed to parse CSV file")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	totalCorrect := 0

problemLoop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, problem.question)

		answerCh := make(chan string)
		go func() {
			var answer string
			_, err := fmt.Scanf("%s\n", &answer)
			if err != nil {
				exit("There was a problem with your answer")
			}
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("Time is up")
			break problemLoop
		case answer := <-answerCh:
			if answer == problem.answer {
				totalCorrect++
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", totalCorrect, len(problems))
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
