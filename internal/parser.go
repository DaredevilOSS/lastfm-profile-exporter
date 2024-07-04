package internal

import (
	"lastfm-profile-exporter/internal/types"
	"log"
)

func parseResponse(itemType ItemType, response types.ResponseHolder) []types.ParsedHolder {
	var parsedObjects []types.ParsedHolder

	switch itemType {
	case TopArtists:
		for _, input := range response.TopArtistsResponse.TopArtists.Artists {
			var output types.ParsedHolder
			output.TopArtistObj = &types.TopArtistObj{
				ArtistName: input.Name,
				PlayCount:  *parseInt(input.PlayCount),
			}
			parsedObjects = append(parsedObjects, output)
		}
	case TopAlbums:
		for _, input := range response.TopAlbumsResponse.TopAlbums.Albums {
			var output types.ParsedHolder
			output.TopAlbumObj = &types.TopAlbumObj{
				ArtistName: input.Artist.Name,
				AlbumName:  input.Name,
				PlayCount:  *parseInt(input.PlayCount),
			}
			parsedObjects = append(parsedObjects, output)
		}
	case TopTracks:
		for _, input := range response.TopTracksResponse.TopTracks.Tracks {
			var output types.ParsedHolder
			output.TopTrackObj = &types.TopTrackObj{
				ArtistName: input.Name,
				PlayCount:  *parseInt(input.PlayCount),
			}
			parsedObjects = append(parsedObjects, output)
		}
	}
	return parsedObjects
}

func parse(itemType ItemType, inputChannel chan types.ResponseHolder, outputChannel chan types.ParsedHolder) {
	defer close(outputChannel)
	for current := range inputChannel {
		parsedResponse := parseResponse(itemType, current)
		log.Printf("Parse %s objects\n", itemType.Name)

		for _, parsedObj := range parsedResponse {
			outputChannel <- parsedObj
		}
	}
}
