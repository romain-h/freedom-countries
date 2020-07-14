package score

import (
	"sort"
)

type Country struct {
	Name              string  `json:"country"`
	IsTerritory       bool    `json:"is_territory"`
	Score             *int    `json:"score"`
	Status            *string `json:"status"`
	PoliticalRights   *int    `json:"political_rights"`
	CivilLiberties    *int    `json:"civil_liberties"`
	ObstaclesToAccess *int    `json:"obstacle_to_access"`
	LimitsOnContent   *int    `json:"limits_on_content"`
	ViolationsOfUR    *int    `json:"violations_of_UR"`
	NetScore          *int    `json:"net_score"`
	NetStatus         *string `json:"net_status"`
	BtStatus          *string `json:"bt_status"`
}

// type Scores map[string]Country

func ToCollection(countries []Country) map[string]Country {
	collection := make(map[string]Country)
	for _, s := range countries {
		collection[s.Name] = s
	}
	return collection
}

// Add BT logic to categorise country
func Preprocess(scores map[string]Country) {
	notFree := "Not Free"
	partlyFree := "Partly Free"

	for key, score := range scores {
		status := "Approved"
		// globStatusVal := *score.Status
		// netStatusVal := *score.NetStatus
		if *score.Status == notFree || (score.NetStatus != nil && *score.NetStatus == notFree) {
			status = "Precluded"
		} else if *score.Status == partlyFree || (score.NetStatus != nil && *score.NetStatus == partlyFree) {
			status = "Case by case"
		}
		score.BtStatus = &status
		scores[key] = score
	}
}

func GetDiff(c1 map[string]Country, c2 map[string]Country) []string {
	diff := make([]string, 0)
	for key, score := range c1 {
		old, found := c2[key]
		if !found {
			diff = append(diff, score.Name)
			continue
		}
		if *score.Score != *old.Score {
			diff = append(diff, score.Name)
			continue

		}
	}
	return diff
}

func GetSortedKey(collection map[string]Country) []string {
	keys := make([]string, 0, len(collection))
	for k := range collection {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
