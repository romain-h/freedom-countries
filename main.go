package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var file = "/Users/rom/Google Drive/bt-countries.xlsx"

type Table struct {
	Headers []string
	Rows    [][]string
}

func isHidden(n *html.Node) bool {
	for _, att := range n.Attr {
		if att.Key == "class" && att.Val == "visually-hidden" {
			return true
		}
	}
	return false
}

func innerText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	result := ""

	for x := n.FirstChild; x != nil; x = x.NextSibling {
		// fmt.Println(x.Attr)
		if !isHidden(x) {
			result += innerText(x)
		}
	}

	return result
}

func CollectTable(doc *html.Node) (*Table, error) {
	var table *Table
	var crawler func(*html.Node)

	crawler = func(n *html.Node) {
		switch n.DataAtom {
		case atom.Table:
			table = &Table{}
		case atom.Th:
			table.Headers = append(table.Headers, innerText(n))
		case atom.Tr:
			table.Rows = append(table.Rows, []string{})
		case atom.Td:
			l := len(table.Rows) - 1
			table.Rows[l] = append(table.Rows[l], innerText(n))
			return
		}

		// if node.Type == html.ElementNode && node.Data == "table" {
		// body = node
		// return
		// }
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if table != nil {
		return table, nil
	}
	return nil, errors.New("Missing <table> in the node tree")
}

// func renderNode(n *html.Node) string {
// var buf bytes.Buffer
// w := io.Writer(&buf)
// html.Render(w, n)
// return buf.String()
// }

type Score struct {
	Name            string
	MatchedName     string
	IsTerritory     bool
	Score           int
	Status          string
	PoliticalRights int
	CivilLiberties  int
}

func tableToScores(table *Table) map[string]Score {
	records := table.Rows[1:]
	scores := make(map[string]Score)
	for _, record := range records {
		// Split Score & Status from same column
		re := regexp.MustCompile(`(\d+)(\D+)`)
		res := re.FindAllStringSubmatch(record[1], -1)
		score, _ := strconv.Atoi(res[0][1])
		pr, _ := strconv.Atoi(record[2])
		cl, _ := strconv.Atoi(record[3])
		name := strings.ReplaceAll(record[0], "*", "")

		scores[name] = Score{
			Name:            name,
			IsTerritory:     strings.Contains(record[0], "*"),
			MatchedName:     "",
			Score:           score,
			Status:          res[0][2],
			PoliticalRights: pr,
			CivilLiberties:  cl,
		}
	}

	return scores
}

func getSheetByName(workbook *xlsx.File, name string) *xlsx.Sheet {
	for _, sh := range workbook.Sheet {
		if sh.Name == name {
			return sh
		}
	}

	return nil
}

type Country struct {
	Name   string
	FHName string
	Code   string
	Row    int
}

func Contains(countries []Country, x string) *int {
	for i, c := range countries {
		if x == c.Name {
			return &i
		}
	}
	return nil
}

func main() {
	resp, err := http.Get("https://freedomhouse.org/countries/freedom-world/scores")
	if err != nil {
		panic("Cannot fetch Freedom house page")
	}
	defer resp.Body.Close()
	doc, _ := html.Parse(resp.Body)
	table, _ := CollectTable(doc)
	scores := tableToScores(table)
	fmt.Println(scores)

	workbook, _ := xlsx.OpenFile(file)
	sheet := getSheetByName(workbook, "countries")
	fmt.Println(sheet.Name)

	// Grab country list from official list
	countries := make(map[string]Country)
	fhCountries := make(map[string]Country)

	for i := 1; i < 252; i++ {
		name := sheet.Cell(i, 2).String()
		fhName := sheet.Cell(i, 7).String()
		country := Country{Name: name, FHName: fhName, Code: sheet.Cell(i, 1).String(), Row: i}
		countries[strings.ToLower(name)] = country
		fhCountries[strings.ToLower(fhName)] = country
	}

	for _, score := range scores {

		_, found := countries[strings.ToLower(score.Name)]
		_, f := fhCountries[strings.ToLower(score.Name)]
		if !found && !f {
			fmt.Println(score.Name, score.IsTerritory)
			// } else {
			// sheet.Cell(c.Row, 7).SetString(score.Name)
		}
	}
	workbook.Save(file)

	// w := csv.NewWriter(os.Stdout)

	// records := append([][]string{table.Headers}, table.Rows[1:]...)
	// for i, record := range records {
	// if i == 0 {
	// record = append([]string{record[0], "Score", "Status"}, record[2:]...)
	// } else {

	// re := regexp.MustCompile(`(\d+)(\D+)`)
	// res := re.FindAllStringSubmatch(record[1], -1)
	// record = append([]string{record[0], res[0][1], res[0][2]}, record[2:]...)
	// }
	// if err := w.Write(record); err != nil {
	// log.Fatalln("error writing record to CSV:", err)
	// }
	// }

	// w.Flush()
}
