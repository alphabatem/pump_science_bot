package pump_science

type ExperimentShowResponse []*ExperimentShow

type ExperimentShow struct {
	Id             int                   `json:"id"`
	Compound       string                `json:"compound"`
	ContractId     string                `json:"contract_id"`
	ExperimentId   string                `json:"experiment_id"`
	BoxId          string                `json:"box_id"`
	ReplicateId    string                `json:"replicate_id"`
	ExperimentType string                `json:"experiment_type"`
	ColumnNumber   int                   `json:"column_number"`
	Data           []*ExperimentShowData `json:"data"`
	Video          []string              `json:"video"`
	InitialAge     int                   `json:"initial_age"`
	NumAnimals     int                   `json:"num_animals"`
}

func (r *ExperimentShow) Delta() ExperimentDelta {
	if len(r.Data) < 2 {
		return ExperimentDelta{
			ExperimentID: r.ExperimentId,
		}
	}

	cur := r.Data[len(r.Data)-1]
	last := r.Data[len(r.Data)-2]

	return ExperimentDelta{
		ExperimentID:    r.ExperimentId,
		DeltaTime:       cur.EndTime,
		FlyCount:        cur.FlyCount - last.FlyCount,
		AverageSpeed:    cur.AverageSpeed - last.AverageSpeed,
		AverageDistance: cur.AverageDistance - last.AverageDistance,
	}
}

type ExperimentShowData struct {
	EndTime         int     `json:"end_time"`
	FlyCount        int     `json:"flyCount"`
	StartTime       int     `json:"start_time"`
	InitialAge      int     `json:"initial_age"`
	AverageSpeed    float64 `json:"averageSpeed"`
	AverageDistance float64 `json:"averageDistance"`
}

type ExperimentDelta struct {
	ExperimentID    string  `json:"experimentID"`
	DeltaTime       int     `json:"deltaTime"`
	FlyCount        int     `json:"flyCount"`
	AverageSpeed    float64 `json:"averageSpeed"`
	AverageDistance float64 `json:"averageDistance"`
}
