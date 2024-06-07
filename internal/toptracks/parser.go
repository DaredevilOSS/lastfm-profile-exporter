package toptracks

import "log"

func parser(inputChannel fetchChannel, outputChannel parseChannel) {
	defer close(outputChannel)
	for current := range inputChannel {
		parsed := parseResponse(current)
		log.Printf("Parse top tracks page %s/%s\n", current.TopTracks.Attr.Page, current.TopTracks.Attr.TotalPages)
		for _, track := range parsed {
			outputChannel <- track
		}
	}
}

type TopTrackObj struct {
	Name     string
	Duration int32
}

func parseDuration(duration string) int32 {
	return 0 // TODO
}

func parseResponse(response GetTopTracksResponse) []TopTrackObj {
	var topTracks []TopTrackObj
	for _, input := range response.TopTracks.Tracks {
		output := TopTrackObj{
			Name:     input.Name,
			Duration: parseDuration(input.Duration),
		}
		topTracks = append(topTracks, output)
	}
	return topTracks
}
