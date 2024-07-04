package types

type TopArtistObj struct {
	ArtistName string
	PlayCount  int
}
type TopAlbumObj struct {
	ArtistName string
	AlbumName  string
	PlayCount  int
}
type TopTrackObj struct {
	ArtistName string
	SongName   string
	Duration   int
	PlayCount  int
}
type ParsedHolder struct {
	TopArtistObj *TopArtistObj
	TopAlbumObj  *TopAlbumObj
	TopTrackObj  *TopTrackObj
}
