package toptracks

import (
	"lastfm-profile-exporter/internal/common"
	"log"
)

type fetchChannel chan GetTopTracksResponse
type parseChannel chan TopTrackObj

func Export(username string) {
	log.Println("Fetching top tracks for " + username)
	fetchChannel := make(fetchChannel)
	go getter(username, fetchChannel)
	parseChannel := make(parseChannel)
	go parser(fetchChannel, parseChannel)
	fileChannel := make(common.FileChannel)
	go writer(username, parseChannel, fileChannel)
	common.Uploader(fileChannel)
}
