package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func authenticate() error {
	apiKey := os.Getenv("API_KEY")
	sharedSecret := os.Getenv("SHARED_SECRET")

	token, err := getToken(apiKey)
	if err != nil {
		return err
	}
	err = requestAuth(apiKey, *token)
	if err != nil {
		return err
	}

	session, err := getSession(*token, apiKey, sharedSecret)
	if err != nil {
		return err
	}
	println("session key: " + session.Key)
	return nil
}

type GetTokenResponse struct {
	Token string `json:"token"`
}

func getToken(apiKey string) (*string, error) {
	url := "https://ws.audioscrobbler.com/2.0/?method=auth.gettoken&api_key=" + apiKey + "&format=json"
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data GetTokenResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data.Token, nil
}

func requestAuth(apiKey string, token string) error {
	url := "http://www.last.fm/api/auth/?api_key=" + apiKey + "&token=" + token
	err := exec.Command("open", url).Start()
	if err != nil {
		return err
	}
	fmt.Println("Please authenticate in the browser and then press Enter to continue...")
	_, err = fmt.Scanln()
	if err != nil {
		return err
	} // Waits for the user to press Enter
	return nil
}

func getMd5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func getApiSignature(token string, apiKey string, sharedSecret string, method string) string {
	apiToSign := "api_key" + apiKey + "method" + method + "token" + token + sharedSecret
	log.Println("API to sign: " + apiToSign)
	return getMd5(apiToSign)
}

type LastFmSession struct {
	Name       string `json:"name"`
	Key        string `json:"key"`
	Subscriber int8   `json:"subscriber"`
}
type GetSessionResponse struct {
	Session LastFmSession `json:"session"`
}
type GetSessionBadResponse struct {
	Message string `json:"message"`
	Error   int8   `json:"error"`
}

func getSession(token string, apiKey string, sharedSecret string) (*LastFmSession, error) {
	apiSig := getApiSignature(token, apiKey, sharedSecret, "auth.getSession")
	url := "https://ws.audioscrobbler.com/2.0/?api_key=" + apiKey + "&token=" + token + "&api_sig=" + apiSig + "&format=json"
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		var data GetSessionBadResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(data.Message)
	}

	var data GetSessionResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data.Session, nil
}
