package crawler

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/romain-h/freedom-countries/internal/score"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

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

func CollectGlobalScores(table *Table) map[string]score.Country {
	records := table.Rows[1:]
	scores := make(map[string]score.Country)

	for _, record := range records {
		// Split Score & Status from same column
		re := regexp.MustCompile(`(\d+)(\D+)`)
		res := re.FindAllStringSubmatch(record[1], -1)
		s, _ := strconv.Atoi(res[0][1])
		pr, _ := strconv.Atoi(record[2])
		cl, _ := strconv.Atoi(record[3])
		name := strings.ReplaceAll(record[0], "*", "")

		scores[name] = score.Country{
			Name:            name,
			IsTerritory:     strings.Contains(record[0], "*"),
			Score:           &s,
			Status:          &res[0][2],
			PoliticalRights: &pr,
			CivilLiberties:  &cl,
		}
	}

	return scores
}

func CollectNetScores(table *Table, collection map[string]score.Country) {
	records := table.Rows[1:]
	for _, record := range records {
		// Split Score & Status from same column
		re := regexp.MustCompile(`(\d+)(\D+)`)
		res := re.FindAllStringSubmatch(record[1], -1)
		s, _ := strconv.Atoi(res[0][1])
		ota, _ := strconv.Atoi(record[2])
		lim, _ := strconv.Atoi(record[3])
		vour, _ := strconv.Atoi(record[4])
		name := strings.ReplaceAll(record[0], "*", "")

		sc, found := collection[name]
		if !found {
			sc = score.Country{
				Name:        name,
				IsTerritory: strings.Contains(record[0], "*"),
			}
		}

		sc.NetScore = &s
		sc.NetStatus = &res[0][2]
		sc.ObstaclesToAccess = &ota
		sc.LimitsOnContent = &lim
		sc.ViolationsOfUR = &vour
		collection[name] = sc
	}
}

func getTable(category string) (*Table, error) {
	url := fmt.Sprintf("https://freedomhouse.org/countries/freedom-%s/scores", category)
	resp, err := http.Get(url)
	if err != nil {
		panic("Cannot fetch Freedom house page")
	}
	defer resp.Body.Close()
	doc, _ := html.Parse(resp.Body)
	return CollectTable(doc)
}

func ScrapData() map[string]score.Country {
	global, _ := getTable("world")
	net, _ := getTable("net")
	collection := CollectGlobalScores(global)
	CollectNetScores(net, collection)

	return collection
}
