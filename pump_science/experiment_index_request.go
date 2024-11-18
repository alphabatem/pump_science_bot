package pump_science

import "github.com/gagliardetto/solana-go"

type ExperimentIndexResponse []*Experiment

type Experiment struct {
	Id                          string           `json:"id"`
	ContractId                  string           `json:"contract_id"`
	ExperimentId                string           `json:"experiment_id"`
	Intervention                string           `json:"intervention"`
	Intervention2               string           `json:"intervention2"`
	CompoundImageUri            string           `json:"compound_image_uri"`
	MedianLifespanCompound      int              `json:"median_lifespan_compound"`
	MedianLifespanControl       int              `json:"median_lifespan_control"`
	PercentageLifespanExtension int              `json:"percentage_lifespan_extension"`
	Ticker                      string           `json:"ticker"`
	Mint                        solana.PublicKey `json:"mint"`
	StartTime                   int64            `json:"start_time"`
	EndTime                     int64            `json:"end_time"`
	ReleaseInterval             int              `json:"release_interval"`
	Status                      string           `json:"status"`
	LiveStatus                  interface{}      `json:"live_status"`
	Description                 string           `json:"description"`
	LearnMoreLink               string           `json:"learn_more_link"`
	LearnMoreLink2              string           `json:"learn_more_link2"`
	AnimalModel                 string           `json:"animal_model"`
	ExperimentData              ExperimentData   `json:"experiment_data"`
}

type ExperimentData struct {
	Food        string `json:"food"`
	Dosage      string `json:"dosage"`
	Humidity    string `json:"humidity"`
	StrainId    string `json:"strain_id"`
	InitialAge  string `json:"initial_age"`
	Temperature string `json:"temperature"`
}
