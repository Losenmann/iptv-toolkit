package setup

import (
    "flag"
    "regexp"
)

var (
    LogLVL              = flag.Int("l", 1, "Log lvl")
    Epg                 = flag.String("e", "", "EPG URI")
    EpgPathDst          = flag.String("E", "./files/tvguide", "Path EPGs storage")
    Playlist            = flag.String("p", "", "Playlist URI")
    PlaylistPathDst     = flag.String("P", "./files/playlist", "Path playlists storage")
    WebPort             = flag.Int("w", 4023, "Web Server port")
    WebPath             = flag.String("f", "./files", "Web Server path")
    EmbedEPG            = flag.String("i", "", "Embed a link to EPG in the playlist")
    PlalistUdpxy        = flag.String("u", "", "Create a playlist with embedded udpxy")
    Health              = flag.Bool("H", false, "Start Healthcheck server")
    Udpxy               = flag.Bool("U", false, "Start UDPXY proxy")
    Schedule            = flag.Bool("S", false, "Start Schedule job")
    Web                 = flag.Bool("W", false, "Start Web Server")
    Crontab             = flag.String("c", "30 6 * * *", "Create a playlist with embedded udpxy")
)

func Initgo() {
    flag.Parse()
    *EpgPathDst = regexp.MustCompile("/$").ReplaceAllString(*EpgPathDst, "")
    *PlaylistPathDst = regexp.MustCompile("/$").ReplaceAllString(*PlaylistPathDst, "")
}
