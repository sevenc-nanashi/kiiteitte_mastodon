package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type apiError struct {
	Message string `json:"error"`
}

type apiApp struct {
	Name string `json:"name"`
}

type apiPostStatus struct {
	Status string `json:"status"`
  Visibility string `json:"visibility"`
}

type apiStatus struct {
	Uri string `json:"uri"`
}

func MastodonVerifyApp(accessToken string, instance string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", instance+"/api/v1/apps/verify_credentials", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var error apiError

		if err := json.NewDecoder(resp.Body).Decode(&error); err != nil {
			log.Error(fmt.Sprintf("Error: %s (%d)", err.Error(), resp.StatusCode))
			os.Exit(1)
		}
		log.Error(fmt.Sprintf("Error: %s (%d)", error.Message, resp.StatusCode))
		os.Exit(1)
	}
	var app apiApp
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		panic(err)
	}
	log.Info("App verified: " + app.Name)
}

func MastodonPostSong(accessToken string, instance string, nextSong KiiteSongInfo) (bool, error) {
	client := &http.Client{}

	reqBody := apiPostStatus{
		Status: fmt.Sprintf("\u266a%s #%s #Kiite\nKiite Cafeできいてます https://cafe.kiite.jp", nextSong.Title, nextSong.VideoId),
    Visibility: "unlisted",
	}

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", instance+"/api/v1/statuses", bytes.NewBuffer(reqBodyJson))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	log.Info("Posting song...")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var error apiError

		if err := json.NewDecoder(resp.Body).Decode(&error); err != nil {
			log.Error(fmt.Sprintf("Error: %s (%d)", err.Error(), resp.StatusCode))
			return false, fmt.Errorf("Error: Failed to post song (%s)", resp.Status)
		}
		log.Error(fmt.Sprintf("Error: %s (%d)", error.Message, resp.StatusCode))
		return false, fmt.Errorf("Error: Failed to post song (%s)", resp.Status)
	}

	var status apiStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Error(fmt.Sprintf("Error: %s (%d)", err.Error(), resp.StatusCode))
		return false, fmt.Errorf("Error: Failed to get status data (%s)", resp.Status)
	}

	log.Info("Posted song: " + status.Uri)

	return true, nil
}
