package validator

import "time"

type response struct {
	Status string `json:"status"`
	Data   []data `json:"data"`
}

type data struct {
	AttesterSlot   int       `json:"attesterslot"`
	CommitteeIndex int       `json:"committeeindex"`
	Epoch          int       `json:"epoch"`
	InclusionSlot  int       `json:"inclusionslot"`
	Status         int       `json:"status"`
	ValidatorIndex int       `json:"validatorindex"`
	Week           int       `json:"week"`
	WeekStart      time.Time `json:"week_start"`
	WeekEnd        time.Time `json:"week_end"`
}
