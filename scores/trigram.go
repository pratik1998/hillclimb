package scores

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/pratik1998/hillclimb/enigma"
)

type TrigramScorer struct {
	scores [18000]float64 // Maximum trigram possible = 26*26*26 = 17576
}

func NewTrigramScore() *TrigramScorer {
	scorer := &TrigramScorer{}
	b, err := ioutil.ReadFile("trigram_data.txt")
	if err != nil {
		fmt.Println("Error reading a file")
	}
	fileData := string(b)
	lines := strings.Split(fileData, "\n")
	for _, line := range lines {
		temp := strings.Split(line, " ")
		trigram_str := temp[0]
		trigram_score, err := strconv.ParseFloat(temp[1], 64)
		if err == nil {
			// we can use bitwise left operator to decrease time complexity
			mask := enigma.GetTrigramToMask(trigram_str)
			scorer.scores[mask] = trigram_score
			// fmt.Println(trigram_str, mask, scorer.scores[mask])
		}
	}
	return scorer
}

func (scorer *TrigramScorer) GetTrigramScore(str string) float64 {
	var score float64
	for i := 0; i < len(str)-2; i++ {
		mask := enigma.GetTrigramToMask(str[i : i+3])
		score += scorer.scores[mask]
	}
	return score
}

func GenerateTrigramProbability() {
	b, err := ioutil.ReadFile("english_trigrams.txt")
	if err != nil {
		fmt.Println("Error reading a file")
	}
	fileData := string(b)
	var total int64
	lines := strings.Split(fileData, "\n")
	for _, line := range lines {
		temp := strings.Split(line, " ")
		value, err := strconv.ParseInt(temp[1], 10, 64)
		if err == nil {
			total += value
		}
	}
	f, err := os.Create("trigram_data.txt")
	if err != nil {
		fmt.Println("Error creating a file")
	}
	defer f.Close()
	var answer string
	for _, line := range lines {
		temp := strings.Split(line, " ")
		value, err := strconv.ParseInt(temp[1], 10, 64)
		if err == nil {
			answer += (temp[0] + " " + fmt.Sprintln(math.Log(float64(value)/float64(total))))
			// answer += (temp[0] + " " + strconv.FormatFloat(math.Log(float64(value)/float64(total)), 'E', -1, 64) + "\n")
			//fmt.Println(temp[0], math.Log(float64(value)/float64(total)))
		}
	}
	f.WriteString(answer)
	fmt.Println("Total Trigrams: ", total)
}
