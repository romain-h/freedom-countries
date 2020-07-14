package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/romain-h/freedom-countries/internal/csv"
	"github.com/romain-h/freedom-countries/internal/score"
)

func main() {

	var previousScoresArr []score.Country
	prev, _ := ioutil.ReadFile("./last.json")
	// Store as previous
	previousFile, _ := os.Create("./previous.json")
	defer previousFile.Close()
	previousFile.Write(prev)

	err := json.Unmarshal(prev, &previousScoresArr)
	if err != nil {
		log.Fatal(err)
	}
	previousScores := score.ToCollection(previousScoresArr)
	buf := csv.GenerateBTListWithScore(previousScores, "./bt-countries.csv")
	fmt.Println(buf.String())
	// scores := crawler.ScrapData()
	// score.Preprocess(scores)

	// diff := score.GetDiff(scores, previousScores)
	// if len(diff) > 0 {
	// fmt.Println(diff)
	// // Store current score to file
	// t, _ := json.Marshal(scores)
	// f, _ := os.Create("./last.json")
	// defer f.Close()
	// f.Write(t)
	// }
}
