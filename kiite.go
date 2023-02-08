package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type rawKiiteSongInfo struct {
	Title     string `json:"title"`
	StartTime string `json:"start_time"`
	VideoId   string `json:"video_id"`
}

type KiiteSongInfo struct {
	Title     string
	StartTime time.Time
	VideoId   string
}

func KiiteGetNextSong() (KiiteSongInfo, error) {
	log.Info("Getting next song...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://cafe.kiite.jp/api/cafe/next_song", nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Error(fmt.Sprintf("Error: %s (%d)", err.Error(), resp.StatusCode))
		return KiiteSongInfo{}, fmt.Errorf("Error: Failed to get next song (%s)", resp.Status)
	}
	var song rawKiiteSongInfo
	if err := json.NewDecoder(resp.Body).Decode(&song); err != nil {
		panic(err)
	}
	log.Info("Next song: " + song.Title)
	log.Info("Start time: " + song.StartTime)

	startTime, err := time.Parse(time.RFC3339, song.StartTime)
	if err != nil {
		log.Error("Error: Failed to parse start time")
		return KiiteSongInfo{}, fmt.Errorf("Error: Failed to parse start time")
	}

	return KiiteSongInfo{
		Title:     song.Title,
		StartTime: startTime,
		VideoId:   song.VideoId,
	}, nil
}
