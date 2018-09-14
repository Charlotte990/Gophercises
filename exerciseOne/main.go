package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// QuestAns contains the Question and Answer
type QuestAns struct {
	Question string
	Answer   string
}

func main() {
	var file string
	var shuffle bool
	var timeOut time.Duration

	flag.StringVar(&file, "f", "defaultQuiz.csv", "Enter a CSV File name with .csv")
	flag.BoolVar(&shuffle, "s", false, "Randomise the order of questions")
	flag.DurationVar(&timeOut, "t", time.Second*30, "Max Quiz Duration e.g. '-t 30s'")
	flag.Parse()

	fmt.Printf("Playing quiz file: %s\n", file)

	qAndA, err := checkFile(file)
	if err != nil {
		log.Fatalf("Invalid CSV, %s", err)
	}
	if len(qAndA) == 0 {
		fmt.Println("File is empty, submit another file")
		return
	}
	if shuffle {
		rand.Shuffle(len(qAndA), func(i, j int) {
			qAndA[i], qAndA[j] = qAndA[j], qAndA[i]
		})
	}

	reader := bufio.NewReader(os.Stdin)

	score, err := runQuiz(qAndA, timeOut, reader)
	if err != nil {
		log.Fatal(err)
	}
	if score > 6 {
		fmt.Println("Well done, your score is:", score)
	} else {
		fmt.Println("Bad luck. Your score is:", score)
	}
}

func runQuiz(qAndA []QuestAns, timeOut time.Duration, reader *bufio.Reader) (int, error) {
	var score int

	fmt.Printf("Press Enter to start the quiz")
	_, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	time.AfterFunc(timeOut, func() {
		fmt.Println("Ran out of time")
		fmt.Printf("Your score was: %d\n", score)
		os.Exit(1)
	})

	for index, question := range qAndA {
		result, err := quizQuestion(question, index, reader)
		if err != nil {
			return 0, err
		}
		score += result
	}
	return score, nil
}

func checkFile(csvFileName string) ([]QuestAns, error) {
	var qsAndAs []QuestAns

	csvFile, err := os.Open(csvFileName)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.FieldsPerRecord = 2
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		qsAndAs = append(qsAndAs, QuestAns{
			Question: line[0],
			Answer:   line[1],
		})
	}
	return qsAndAs, nil
}

func quizQuestion(qAndA QuestAns, index int, reader *bufio.Reader) (int, error) {

	fmt.Printf("Question %d: %s\nAnswer: ", index+1, qAndA.Question)
	userAnswer, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	answer := strings.ToLower(strings.TrimSpace(userAnswer))
	if answer == strings.ToLower(strings.TrimSpace(qAndA.Answer)) {
		return 1, nil
	}
	return 0, nil
}
