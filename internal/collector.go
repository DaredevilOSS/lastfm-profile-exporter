package internal

import (
	"lastfm-profile-exporter/internal/types"
)

type ItemType struct {
	Name     string
	Endpoint string
}

var (
	TopArtists = ItemType{
		Name:     "top artists",
		Endpoint: "gettopartists",
	}
	TopAlbums = ItemType{
		Name:     "top albums",
		Endpoint: "gettopalbums",
	}
	TopTracks = ItemType{
		Name:     "top tracks",
		Endpoint: "gettoptracks",
	}
)

var itemTypes = []ItemType{
	TopArtists,
	TopAlbums,
	TopTracks,
}

func Collect(username string) {
	for _, itemType := range itemTypes {
		fetchChannel := make(chan types.ResponseHolder)
		parseChannel := make(chan types.ParsedHolder)
		fileChannel := make(chan string)

		go get(itemType, username, fetchChannel)
		go parse(itemType, fetchChannel, parseChannel)
		go write(itemType, username, parseChannel, fileChannel)
		upload(fileChannel)
	}
}
