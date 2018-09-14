package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNoCSVFileFails(t *testing.T) {
	_, err := checkFile("sdfsdfsd")

	assert.Error(t, err)
}

func TestReadCSVs(t *testing.T) {
	tt := []struct {
		TestName       string
		Contents       string
		ShouldError    bool
		ExpectedOutput []QuestAns
	}{
		{"Single field", "5", true, nil},
		{"Too many fields", "5,5,6", true, nil},
		{"Correct number of fields", "5,7", false, []QuestAns{QuestAns{Question: "5", Answer: "7"}}},
	}

	for _, x := range tt {
		t.Run(x.TestName, func(t *testing.T) {
			testFile, err := ioutil.TempFile("", "exOneTestSuite")
			assert.NoError(t, err)

			defer os.Remove(testFile.Name()) // clean up

			_, err = testFile.Write([]byte(x.Contents))
			assert.NoError(t, err)

			err = testFile.Close()
			assert.NoError(t, err)

			result, err := checkFile(testFile.Name())

			if x.ShouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, x.ExpectedOutput, result)
		})
	}

}

func TestQuizQuestion(t *testing.T) {
	tt := []struct {
		TestName       string
		QandAs         QuestAns
		Index          int
		UserAnswer     string
		ShouldError    bool
		ExpectedOutput int
	}{
		{"User gets answer correct", QuestAns{Question: "5+2", Answer: "7"}, 0, "7", false, 1},
		{"User gets answer wrong", QuestAns{Question: "5+2", Answer: "7"}, 1, "8", false, 0},
		{"User gives no answer", QuestAns{Question: "5+2", Answer: "7"}, 2, "", false, 0},
		{"User gives correct answer with lots of whitespace", QuestAns{Question: "5+2", Answer: "7"}, 2, "  \t  7    ", false, 1},
		{"User gives correct answer with all caps", QuestAns{Question: "Do you like dogs", Answer: "Yes"}, 2, "YES", false, 1},
		{"User gives correct answer with no caps", QuestAns{Question: "Do you like dogs", Answer: "YES"}, 2, "yes", false, 1},
	}

	for _, x := range tt {
		t.Run(x.TestName, func(t *testing.T) {

			reader := bufio.NewReader(strings.NewReader(x.UserAnswer + "\n"))

			result, err := quizQuestion(x.QandAs, x.Index, reader)

			if x.ShouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, x.ExpectedOutput, result)
		})
	}

}

func TestRunQuiz(t *testing.T) {
	tt := []struct {
		TestName    string
		QandAs      []QuestAns
		TimeOut     time.Duration
		UserInput   string
		Score       int
		ShouldError bool
	}{
		{"User gets answer correct", []QuestAns{QuestAns{Question: "5+2", Answer: "7"}, {Question: "1+1", Answer: "2"}}, time.Millisecond * 100, "\n7\n2\n", 2, false},
		{"User gets answer correct with user whitespace", []QuestAns{QuestAns{Question: "15+2", Answer: "17"}, {Question: "1+1", Answer: "2"}}, time.Millisecond * 100, "\n17\n2 \n", 2, false},
		{"User gets answer correct with csv whitespace", []QuestAns{QuestAns{Question: "15+2", Answer: "17"}, {Question: "1+1", Answer: "2 "}}, time.Millisecond * 100, "\n17\n2\n", 2, false},
	}

	for _, x := range tt {
		t.Run(x.TestName, func(t *testing.T) {

			reader := bufio.NewReader(strings.NewReader(x.UserInput))

			score, err := runQuiz(x.QandAs, x.TimeOut, reader)

			if x.ShouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, x.Score, score)
		})
	}

}
