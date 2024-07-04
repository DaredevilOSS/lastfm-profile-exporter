package types

type LeanArtist struct {
	Name string `json:"name"`
}
type Attr struct {
	User       string `json:"user"`
	TotalPages string `json:"totalPages"`
	Page       string `json:"page"`
	PerPage    string `json:"perPage"`
	Total      string `json:"total"`
}

type Artist struct {
	Name      string `json:"name"`
	PlayCount string `json:"playcount"`
}
type Album struct {
	Name      string     `json:"name"`
	PlayCount string     `json:"playcount"`
	Artist    LeanArtist `json:"artist"`
}
type Track struct {
	Rank      string     `json:"rank"`
	Name      string     `json:"name"`
	PlayCount string     `json:"playcount"`
	Duration  string     `json:"duration"`
	Artist    LeanArtist `json:"artist"`
}

type GetTopArtistsResponse struct {
	TopArtists struct {
		Attr    Attr     `json:"@attr"`
		Artists []Artist `json:"artist"`
	} `json:"topartists"`
}
type GetTopAlbumsResponse struct {
	TopAlbums struct {
		Attr   Attr    `json:"@attr"`
		Albums []Album `json:"album"`
	} `json:"topalbums"`
}
type GetTopTracksResponse struct {
	TopTracks struct {
		Attr   Attr    `json:"@attr"`
		Tracks []Track `json:"track"`
	} `json:"toptracks"`
}
type ResponseHolder struct {
	TopArtistsResponse *GetTopArtistsResponse
	TopAlbumsResponse  *GetTopAlbumsResponse
	TopTracksResponse  *GetTopTracksResponse
}
