package score

import (
	"encoding/json"
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

type Countries map[string]Country

func ReadBuf(buf []byte) (Countries, error) {
	countries := make(Countries)
	err := json.Unmarshal(buf, &countries)
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func WriteBuf(countries Countries) ([]byte, error) {
	return json.Marshal(countries)
}

// Add BT logic to categorise country
func Preprocess(scores Countries) {
	notFree := "Not Free"
	partlyFree := "Partly Free"

	for key, score := range scores {
		status := "Approved"
		if *score.Status == notFree || (score.NetStatus != nil && *score.NetStatus == notFree) {
			status = "Precluded"
		} else if *score.Status == partlyFree || (score.NetStatus != nil && *score.NetStatus == partlyFree) {
			status = "Case by case"
		}
		score.BtStatus = &status
		scores[key] = score
	}
}
