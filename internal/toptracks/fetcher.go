package toptracks

import (
	"encoding/json"
	"fmt"
	"io"
	"lastfm-profile-exporter/internal/common"
	"log"
	"net/http"
	"strconv"
)

type Artist struct {
	Name string `json:"name"`
}

type Track struct {
	Rank      string `json:"rank"`
	Name      string `json:"name"`
	PlayCount string `json:"playcount"`
	Duration  string `json:"duration"`
	Artist    Artist `json:"artist"`
}
type Attr struct {
	User       string `json:"user"`
	TotalPages string `json:"totalPages"`
	Page       string `json:"page"`
	PerPage    string `json:"perPage"`
	Total      string `json:"total"`
}
type TopTracks struct {
	Attr   Attr    `json:"@attr"`
	Tracks []Track `json:"track"`
}
type GetTopTracksResponse struct {
	TopTracks TopTracks `json:"toptracks"`
}

func attemptGet(username string, page int) (*GetTopTracksResponse, error) {
	pageNum := strconv.Itoa(page)
	url := fmt.Sprintf("%s?method=user.gettoptracks&api_key=%s&user=%s&page=%s&limit=1000&format=json", common.BaseURL,
		common.ApiKey, username, pageNum)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data GetTopTracksResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func getWithRetry(username string, page int) (*GetTopTracksResponse, bool) {
	retries := 0
	for retries < common.MaxAPIRetries {
		current, err := attemptGet(username, page)
		if err != nil {
			retries++
			log.Println(err)
			continue
		}
		//log.Printf("Get top tracks page %s/%s\n", current.TopTracks.Attr.Page, current.TopTracks.Attr.TotalPages)
		return current, current.TopTracks.Attr.Page != current.TopTracks.Attr.TotalPages
	}
	return nil, true // skipping failed page
}

func getter(username string, inputChannel fetchChannel) {
	defer close(inputChannel)
	var response *GetTopTracksResponse
	page := 0
	nextPageExists := true
	for nextPageExists {
		page++
		response, nextPageExists = getWithRetry(username, page)
		if response != nil {
			log.Printf("Got top tracks page %s/%s\n", response.TopTracks.Attr.Page, response.TopTracks.Attr.TotalPages)
			inputChannel <- *response
		}
	}
}
