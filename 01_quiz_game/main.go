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
	csvFileLoc := flag.String("file", "problems.csv", "expects the csv file with the quiz problems in the format question, answer; default: problems.csv")
	timeLimit := flag.Int("timelimit", 30, "the time limit for answering all questions in seconds")

	flag.Parse()

	file, err := os.Open(*csvFileLoc)
	if err != nil {
		exit(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	qNA, err := r.ReadAll()
	if err != nil {
		exit(err)
	}

	problems := parseLines(qNA)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
	done := false

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s: ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			done = true
			fmt.Printf("\ntimes up! you scored %d out of %d!\n", correct, len(qNA))
			os.Exit(1)
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	
	if done == false {
		fmt.Printf("you scored %d out of %d!\n", correct, len(qNA))
	}
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))

	for i, line := range lines {
		problems[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return problems
}

type problem struct {
	q string
	a string
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}
