package csv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/romain-h/freedom-countries/internal/score"
)

func GenerateBTListWithScore(countries map[string]score.Country, filename string) bytes.Buffer {
	file, _ := os.Open(filename)
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))

	var fullLines [][]string
	i := -1
	for {
		i++
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if i == 0 {
			fullLines = append(fullLines, line)
			continue
		}
		country, ok := countries[line[7]]
		if ok {
			line = append(line, *country.BtStatus)
		}
		fullLines = append(fullLines, line)
	}

	var b bytes.Buffer
	writer := csv.NewWriter(bufio.NewWriter(&b))
	defer writer.Flush()

	err := writer.WriteAll(fullLines)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

func ReadScores(filename string) []score.Country {
	file, _ := os.Open(filename)
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))

	var allScores []score.Country

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("ReadScores: %e", err)
		}

		s := score.Country{Name: line[0]}

		if line[1] != "" {
			res, _ := strconv.Atoi(line[1])
			s.PoliticalRights = &res
		}
		if line[2] != "" {
			res, _ := strconv.Atoi(line[2])
			s.CivilLiberties = &res
		}
		if line[3] != "" {
			res, _ := strconv.Atoi(line[3])
			s.Score = &res
		}
		if line[4] != "" {
			s.Status = &line[4]
		}
		if line[5] != "" {
			res, _ := strconv.Atoi(line[5])
			s.ObstaclesToAccess = &res
		}
		if line[6] != "" {
			res, _ := strconv.Atoi(line[6])
			s.LimitsOnContent = &res
		}
		if line[7] != "" {
			res, _ := strconv.Atoi(line[7])
			s.ViolationsOfUR = &res
		}
		if line[8] != "" {
			res, _ := strconv.Atoi(line[8])
			s.NetScore = &res
		}
		if line[9] != "" {
			s.NetStatus = &line[9]
		}
		if line[10] != "" {
			s.BtStatus = &line[10]
		}
		allScores = append(allScores, s)
	}

	return allScores
}

// func WriteScores(filename string, scores []score.Country) {
// file, _ := os.Create(filename)
// defer file.Close()

// writer := csv.NewWriter(file)
// defer writer.Flush()

// var collection [][]string
// for i, s := range scores {

// politicalRights := strconv.Itoa(*s.PoliticalRights)
// civilLiberties := strconv.Itoa(*s.CivilLiberties)
// score := strconv.Itoa(*s.Score)
// obstacle := strconv.Itoa(*s.ObstaclesToAccess)
// line := []string{s.Name, politicalRights}
// collection = append(collection, line)
// }
// err := writer.WriteAll(collection)
// if err != nil {
// log.Fatal(err)
// }
// }
