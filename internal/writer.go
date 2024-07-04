package internal

import (
	"encoding/json"
	"lastfm-profile-exporter/internal/types"
	"log"
	"strings"
)

const batchSize = 20000

func ObjectToJson(itemType ItemType, obj types.ParsedHolder) (*string, error) {
	var jsonBytes []byte
	var err error

	switch itemType {
	case TopArtists:
		jsonBytes, err = json.Marshal(obj.TopArtistObj)
	case TopAlbums:
		jsonBytes, err = json.Marshal(obj.TopAlbumObj)
	case TopTracks:
		jsonBytes, err = json.Marshal(obj.TopTrackObj)
	}

	if err != nil {
		return nil, err
	}
	stringifiedJson := string(jsonBytes)
	return &stringifiedJson, nil
}

func write(itemType ItemType, username string, inputChannel chan types.ParsedHolder, outputChannel chan string) {
	defer close(outputChannel)
	var batch []string

	for current := range inputChannel {
		line, err := ObjectToJson(itemType, current)
		if err != nil {
			log.Println(err)
			continue
		}
		batch = append(batch, *line)

		if len(batch) >= batchSize {
			filename, err := WriteToFile(username, strings.Replace(itemType.Name, " ", "_", -1), batch)
			if err != nil {
				log.Println(err)
				continue
			}
			outputChannel <- *filename
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		filename, err := WriteToFile(username, strings.Replace(itemType.Name, " ", "_", -1), batch)
		if err != nil {
			log.Println(err)
			return
		}
		outputChannel <- *filename
	}
}
