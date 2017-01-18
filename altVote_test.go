package altVote

import (
	"bufio"
	"os"
	"testing"

	"fmt"

	"io/ioutil"
	"path"

	"strings"

	"strconv"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {

	// for i := 0; i < 1000; i++ {
	// 	testCandidates := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	// 	fmt.Println(fmt.Sprintf("%v, %v, %v, %v", testCandidates[rand.Intn(10)], testCandidates[rand.Intn(10)], testCandidates[rand.Intn(10)], testCandidates[rand.Intn(10)]))

	// }

	// Only pass t into top-level Convey calls
	Convey("altVote", t, func() {

		files, err := ioutil.ReadDir("./testElections")
		if err != nil {
			t.Log("Error reading testElections directory:")
			panic(err)
		}

		for i := range files {
			file, err := os.Open(path.Join("./testElections/", files[i].Name()))
			if err != nil {
				panic(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)

			candidates := []string{}
			ballots := [][]string{}

			//do one scan to get the expected winner
			scanner.Scan()
			expectedWinner := scanner.Text()

			//do another scan to get the expected number of rounds
			scanner.Scan()
			expectedNumberOfRounds, err := strconv.Atoi(scanner.Text())
			if err != nil {
				panic(err)
			}

			//now get the ballots
			for scanner.Scan() {
				line := scanner.Text()
				line = strings.Replace(line, " ", "", -1)
				votes := strings.Split(line, ",")

				//dynamically build up candidate list
				for _, vote := range votes {
					candidates = uniqueAppend(candidates, strings.TrimSpace(vote))
				}
				ballots = append(ballots, votes)

			}

			if err := scanner.Err(); err != nil {
				panic(err)
			}

			Convey(files[i].Name(), func() {
				fmt.Println("running file: ", files[i].Name())
				runTest(expectedWinner, expectedNumberOfRounds, ballots, candidates)
				fmt.Println("")
			})

		}

	})
}

func runTest(expectedWinner string, expectedNumberOfRounds int, ballots [][]string, candidates []string) {

	fmt.Println("expectedWinner:")
	fmt.Println(expectedWinner)
	fmt.Println("candidates:")
	fmt.Println(candidates)
	fmt.Println("ballots:")
	fmt.Println(ballots)
	results, _ := GetResults(candidates, ballots)
	So(len(results.Rounds), ShouldEqual, expectedNumberOfRounds)
	So(results.Winner, ShouldEqual, expectedWinner)
}

func uniqueAppend(ss []string, s string) []string {
	for i := range ss {
		if ss[i] == s {
			return ss
		}
	}

	return append(ss, s)
}
