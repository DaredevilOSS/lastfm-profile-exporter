package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"lastfm-profile-exporter/internal/types"
	"log"
	"net/http"
	"os"
	"strconv"
)

const BaseURL = "https://ws.audioscrobbler.com/2.0/"
const MaxRetries = 3

var ApiKey = os.Getenv("API_KEY")

func hasMorePages(currentPage string, totalPages string) bool {
	return *parseInt(currentPage) < *parseInt(totalPages)
}

func attemptGet(itemType ItemType, username string, page int) (*types.ResponseHolder, bool, error) {
	pageNum := strconv.Itoa(page)
	url := fmt.Sprintf("%s?method=user.%s&api_key=%s&user=%s&page=%s&limit=1000&format=json",
		BaseURL, itemType.Endpoint, ApiKey, username, pageNum)
	response, err := http.Get(url)
	if err != nil {
		return nil, true, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, true, err
	}

	var responseHolder types.ResponseHolder
	var more bool

	switch itemType {
	case TopArtists:
		var data types.GetTopArtistsResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, true, err
		}
		more = hasMorePages(data.TopArtists.Attr.Page, data.TopArtists.Attr.TotalPages)
		responseHolder.TopArtistsResponse = &data
	case TopAlbums:
		var data types.GetTopAlbumsResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, true, err
		}
		more = hasMorePages(data.TopAlbums.Attr.Page, data.TopAlbums.Attr.TotalPages)
		responseHolder.TopAlbumsResponse = &data
	case TopTracks:
		var data types.GetTopTracksResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, true, err
		}
		more = hasMorePages(data.TopTracks.Attr.Page, data.TopTracks.Attr.TotalPages)
		responseHolder.TopTracksResponse = &data
	}
	return &responseHolder, more, nil
}

func getWithRetry(itemType ItemType, username string, page int) (*types.ResponseHolder, bool) {
	retries := 0
	for retries < MaxRetries {
		response, more, err := attemptGet(itemType, username, page)
		if err != nil {
			retries++
			log.Println(err)
			continue
		}
		return response, more
	}
	return nil, true
}

func get(itemType ItemType, username string, inputChannel chan types.ResponseHolder) {
	defer close(inputChannel)
	var response *types.ResponseHolder
	page := 0
	nextPageExists := true
	for nextPageExists {
		page++
		log.Printf("Fetching page %d of %s for %s", page, itemType.Name, username)
		response, nextPageExists = getWithRetry(itemType, username, page)
		if response != nil {
			inputChannel <- *response
		}
	}
}
