package main

import (
	"github.com/alphabatem/common/context"
	"github.com/alphabatem/pump_science_bot/pump_science"
	"github.com/alphabatem/pump_science_bot/services"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	ctx, err := context.NewCtx(
		&pump_science.Service{},

		&services.SwapService{},
		&services.BotService{},
	)

	if err != nil {
		log.Fatal().Err(err)
		return
	}

	err = ctx.Run()
	if err != nil {
		log.Fatal().Err(err)
		return
	}
}
