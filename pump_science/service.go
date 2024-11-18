package pump_science

import (
	"encoding/json"
	"errors"
	"github.com/alphabatem/common/context"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"time"
)

type Service struct {
	context.DefaultService

	client *http.Client

	latest      ExperimentShowResponse
	experiments []*Experiment

	lastExperimentData map[string]int

	listeners []chan *ExperimentDelta

	apiKey string
}

const SERVICE = "pump_science_svc"

func (svc Service) Id() string {
	return SERVICE
}

func (svc *Service) Start() error {
	svc.client = &http.Client{Timeout: 2 * time.Second}

	svc.apiKey = os.Getenv("PUMP_SCIENCE_KEY")

	svc.lastExperimentData = map[string]int{}
	svc.listeners = []chan *ExperimentDelta{}

	go svc.worker()

	return svc.loadExperiments()
}

func (svc *Service) Listen(ch chan *ExperimentDelta) {
	svc.listeners = append(svc.listeners, ch)
}

func (svc *Service) Experiment(experimentID string) *Experiment {
	for _, e := range svc.experiments {
		if e.ExperimentId == experimentID {
			return e
		}
	}
	return nil
}

func (svc *Service) worker() {
	for {
		err := svc.work()
		if err != nil {
			log.Error().Err(err).Msg("PumpScience::Service Worker err")
			time.Sleep(2 * time.Second)
			continue
		}
		time.Sleep(10 * time.Second)
	}
}

func (svc *Service) work() error {
	uri := "https://npcnnhpqtqjqlgwqqoic.supabase.co/rest/v1/release_flies_prod?select=*&contract_id=eq.c951f7da-ae6c-49f3-a4dd-27ec18482fb9&experiment_id=eq.7fda0b1a-a19c-4566-a05b-27d1d150ea6b"

	var respDat ExperimentShowResponse
	err := svc._call(uri, &respDat)
	if err != nil {
		return err
	}

	svc.latest = respDat

	svc.Notify()
	return nil
}

func (svc *Service) Notify() {
	for _, e := range svc.latest {
		d := e.Delta()

		//Check if new data
		if v, ok := svc.lastExperimentData[d.ExperimentID]; ok && v >= d.DeltaTime {
			continue
		}
		svc.lastExperimentData[d.ExperimentID] = d.DeltaTime

		log.Info().Str("id", d.ExperimentID).Int("timestamp", d.DeltaTime).Msg("New DATA")
		for _, listener := range svc.listeners {
			listener <- &d
		}
	}
}

func (svc *Service) loadExperiments() error {
	uri := "https://npcnnhpqtqjqlgwqqoic.supabase.co/rest/v1/experimentsprod?select=*&order=start_time.desc&status=in.%28live%2Cupcoming%29"

	var respDat ExperimentIndexResponse
	err := svc._call(uri, &respDat)
	if err != nil {
		return err
	}

	svc.experiments = respDat

	log.Info().Int("count", len(svc.experiments)).Msg("Experiments loaded")
	return nil
}

func (svc *Service) _call(uri string, respDat interface{}) error {
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("apikey", svc.apiKey)
	req.Header.Add("Authorization", svc.apiKey)

	resp, err := svc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dat, &respDat)
	if err != nil {
		return err
	}

	return nil
}
