package score

import (
	"bytes"
	"html/template"
	"os"
	"reflect"
	"sort"

	"github.com/romain-h/freedom-countries/internal/storage"
)

type Diff struct {
	Type string   `json:"type"`
	Base *Country `json:"base_country"`
	New  *Country `json:"new_country"`
}

type DiffMap map[string]Diff

func (diffs *DiffMap) RenderEmail(store storage.Storage, templateFileName string) (*string, error) {
	tmplt, err := store.GetFile(templateFileName)

	t, err := template.New("email.html").Funcs(template.FuncMap{
		"Deref": func(i *int) int { return *i },
	}).Parse(string(*tmplt))
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	// Prepare data
	minors, majors := diffs.SplitMinors()
	min := minors.GetSorted()
	maj := majors.GetSorted()
	data := struct {
		Name   string
		Title  string
		Minors []Diff
		Majors []Diff
	}{
		Name:   os.Getenv("FCUP_NAME"),
		Minors: min,
		Majors: maj,
	}

	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	str := buf.String()
	return &str, nil
}

func (diffs *DiffMap) SplitMinors() (DiffMap, DiffMap) {
	minors := make(DiffMap)
	majors := make(DiffMap)

	for k, v := range *diffs {
		if v.Type != "update" {
			majors[k] = v
			continue
		}

		if *v.Base.BtStatus != *v.New.BtStatus {
			majors[k] = v
			continue
		}
		minors[k] = v
	}

	return minors, majors

}

func (diffs *DiffMap) GetSorted() []Diff {
	keys := make([]string, 0, len(*diffs))
	for k := range *diffs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	col := make([]Diff, 0, len(*diffs))
	d := *diffs
	for _, c := range keys {
		col = append(col, d[c])
	}

	return col
}

func GetDiff(c1 Countries, c2 Countries) DiffMap {
	diff := make(DiffMap)
	for key, _ := range c2 {
		newScore, _ := c2[key]
		old, found := c1[key]
		if !found {
			diff[key] = Diff{Type: "addition", Base: &newScore}
			continue
		}

		if !reflect.DeepEqual(newScore, old) {
			diff[key] = Diff{Type: "update", Base: &old, New: &newScore}
			continue

		}
	}
	// Check deletion
	for key, _ := range c1 {
		_, found := c2[key]
		if !found {
			n, _ := c1[key]
			diff[key] = Diff{Type: "deletion", Base: &n}
			continue
		}
	}

	return diff
}
