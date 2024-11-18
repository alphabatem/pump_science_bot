package services

import (
	"github.com/alphabatem/common/context"
	"github.com/alphabatem/pump_science_bot/pump_science"
	"github.com/rs/zerolog/log"
)

type BotService struct {
	context.DefaultService

	pumpSciSvc *pump_science.Service
	swapSvc    *SwapService

	deltaCh chan *pump_science.ExperimentDelta
}

const BOT_SVC = "bot_svc"

func (svc BotService) Id() string {
	return BOT_SVC
}

func (svc *BotService) Shutdown() {
	close(svc.deltaCh)
}

func (svc *BotService) Start() error {
	svc.deltaCh = make(chan *pump_science.ExperimentDelta, 100)

	svc.swapSvc = svc.Service(SWAP_SVC).(*SwapService)

	svc.pumpSciSvc = svc.Service(pump_science.SERVICE).(*pump_science.Service)
	svc.pumpSciSvc.Listen(svc.deltaCh)

	return svc.work()
}

func (svc *BotService) work() error {
	log.Info().Msg("Listening on delta changes")
	for d := range svc.deltaCh {
		err := svc.onDelta(d)
		if err != nil {
			log.Error().Err(err).Msg("BotService::onDelta error")
		}

	}
	return nil
}

func (svc *BotService) onDelta(delta *pump_science.ExperimentDelta) error {
	experiment := svc.pumpSciSvc.Experiment(delta.ExperimentID)

	if delta.FlyCount < 0 {
		log.Info().Str("id", delta.ExperimentID).Int("delta", delta.FlyCount).Str("mint", experiment.Mint.String()).Msg("Fly has died!")
		return svc.swapSvc.Sell(experiment.Mint, 1)
	}

	if delta.AverageDistance < 0 {
		log.Info().Str("id", delta.ExperimentID).Float64("delta", delta.AverageDistance).Msg("Flys Moving Less")
	} else {
		log.Info().Str("id", delta.ExperimentID).Float64("delta", delta.AverageDistance).Msg("Flys Moving More")
	}

	if delta.AverageSpeed < 0 {
		log.Info().Str("id", delta.ExperimentID).Float64("delta", delta.AverageSpeed).Msg("Flys Slowing Down")
	} else {
		log.Info().Str("id", delta.ExperimentID).Float64("delta", delta.AverageSpeed).Msg("Flys Speeding Up")
	}

	return nil
}
