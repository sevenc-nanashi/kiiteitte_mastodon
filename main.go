package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/joho/godotenv"
	colorable "github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()

	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())

	Title()

	instance, found := os.LookupEnv("INSTANCE")
	if !found {
		log.Error("INSTANCE not found! Set it in .env.")
		os.Exit(1)
	}

	accessTokenBytes, err := fs.ReadFile(os.DirFS("."), "access_token.key")
	var accessToken string
	if err != nil {
		accessToken = Login(instance)
	} else {
		accessToken = string(accessTokenBytes)
	}

	MastodonVerifyApp(accessToken, instance)

	log.Info("Starting...")

	for {
		songInfo, err := KiiteGetNextSong()
		if err != nil {
			log.Error(err)
			log.Info("Retrying in 10 seconds...")
			time.Sleep(10 * time.Second)
			continue
		}
		currentTime := time.Now()
		timeUntilStart := (songInfo.StartTime.Sub(currentTime))
		log.Info(fmt.Sprintf("Waiting for %d seconds...", int(timeUntilStart.Seconds())))
		time.Sleep(timeUntilStart)

		MastodonPostSong(accessToken, instance, songInfo)
		log.Info("Waiting for 10 seconds...")
		time.Sleep(10 * time.Second)
	}
}
