package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type apiToken struct {
	AccessToken string `json:"access_token"`
}

func Login(instance string) string {
	log.Info(fmt.Sprintf("Authorize at: %s", instance+"/oauth/authorize"+
		"?client_id="+os.Getenv("CLIENT_ID")+
		"&redirect_uri=urn:ietf:wg:oauth:2.0:oob"+
		"&response_type=code"+
		"&scope=read+write"))
	fmt.Print("Enter code: ")
	var code string
	fmt.Scanln(&code)
	client := &http.Client{}
	req, err := http.NewRequest("POST", instance+"/oauth/token", bytes.NewBuffer([]byte("client_id="+os.Getenv("CLIENT_ID")+"&client_secret="+os.Getenv("CLIENT_SECRET")+"&grant_type=authorization_code&redirect_uri=urn:ietf:wg:oauth:2.0:oob&code="+code)))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
	var app apiToken
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		panic(err)
	}
	f, err := os.Create("access_token.key")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(app.AccessToken)

	log.Info("App verified. Access token written in ./access_token.key")
	return app.AccessToken
}
