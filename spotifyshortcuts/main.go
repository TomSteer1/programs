package main

import (
	"github.com/tomsteer1/SpotifyGoWrapper"
	"flag"
)

func main() {
	fileName := flag.String("config", "spotify.conf", "The config file to use")
	playlistID := flag.String("playlist", "playlistID", "The playlist to use")
	flag.Parse()
	spotify.LoadConfig(*fileName)
	go spotify.RefreshToken()
	currentSong := spotify.GetCurrentSong()
	spotify.AddTrackToPlaylist(*playlistID,currentSong.Item.Uri)
}
