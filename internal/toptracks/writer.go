package toptracks

import (
	"encoding/json"
	"lastfm-profile-exporter/internal/common"
	"log"
)

func ObjectToJson(topTrack TopTrackObj) (*string, error) {
	marshaled, err := json.Marshal(topTrack)
	if err != nil {
		return nil, err
	}
	jsonString := string(marshaled)
	return &jsonString, nil
}

func writer(username string, inputChannel parseChannel, outputChannel common.FileChannel) {
	defer close(outputChannel)
	var batch []string
	batchSize := 100000

	for current := range inputChannel {
		line, err := ObjectToJson(current)
		if err != nil {
			log.Println(err)
			continue
		}
		batch = append(batch, *line)

		if len(batch) >= batchSize {
			outputChannel <- common.WriteToFile(username, "toptracks", batch)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		outputChannel <- common.WriteToFile(username, "toptracks", batch)
	}
}
