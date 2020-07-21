package score

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func getFromFile(filename string) (Countries, error) {
	f, _ := ioutil.ReadFile(filename)
	return ReadBuf(f)
}

func TestGetDiff(t *testing.T) {
	baseCountries, _ := getFromFile("testdata/1_countries.json")
	newCountries, _ := getFromFile("testdata/2_countries.json")

	comoros, _ := baseCountries["Comoros"]
	colombia, _ := newCountries["Colombia"]
	dominica, _ := newCountries["Dominica"]
	france, _ := baseCountries["France"]
	newFrance, _ := newCountries["France"]

	testCases := map[string]struct {
		newCountries Countries
		diff         DiffMap
	}{
		"no-diff": {
			newCountries: baseCountries,
			diff:         map[string]Diff{},
		},
		"full-diff": {
			newCountries: newCountries,
			diff: DiffMap{
				"Colombia": Diff{Type: "addition", Base: &colombia},
				"Comoros":  Diff{Type: "deletion", Base: &comoros},
				"Dominica": Diff{Type: "addition", Base: &dominica},
				"France":   Diff{Type: "update", Base: &france, New: &newFrance},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			diff := GetDiff(baseCountries, tc.newCountries)

			if !reflect.DeepEqual(tc.diff, diff) {
				t.Errorf("wrong diff\n got %v\n want %v", diff, tc.diff)
			}
		})
	}
}

func TestSplitMinors(t *testing.T) {
	approved := "Approved"
	precluded := "Precluded"

	a := Diff{Type: "addition", Base: &Country{}}
	b := Diff{Type: "deletion", Base: &Country{}}
	c := Diff{
		Type: "update",
		Base: &Country{
			BtStatus: &approved,
		},
		New: &Country{
			BtStatus: &precluded,
		},
	}
	d := Diff{
		Type: "update",
		Base: &Country{
			BtStatus: &approved,
		},
		New: &Country{
			BtStatus: &approved,
		},
	}

	diffs := DiffMap{
		"A": a,
		"B": b,
		"C": c,
		"D": d,
	}

	expectedMinors := DiffMap{
		"D": d,
	}
	expectedMajors := DiffMap{
		"A": a,
		"B": b,
		"C": c,
	}

	minors, majors := diffs.SplitMinors()
	if !reflect.DeepEqual(minors, expectedMinors) {
		t.Errorf("minors should be \n %v\n got %v \n", minors, expectedMinors)
	}
	if !reflect.DeepEqual(majors, expectedMajors) {
		t.Errorf("majors should be \n %v\n got %v \n", majors, expectedMajors)
	}
}
