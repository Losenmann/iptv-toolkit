package setup

import (
    "flag"
    "regexp"
)

var (
    LogLVL              = flag.Int("l", 1, "Log lvl")
    Epg                 = flag.String("e", "", "Path out ")
    EpgPathDst          = flag.String("E", "./", "Path EPGs storage")
    Playlist            = flag.String("p", "", "Playlist URL")
    PlaylistPathDst     = flag.String("P", "./", "Path playlists storage")
    InputPlaylist       = flag.String("A", "./", "Input playlist")
    InputEPG            = flag.String("B", "./", "Input EPG")
    EmbedEPG            = flag.String("i", "/home/maxim/Документы/repo/iptv-toolkit/tmp/test.xml", "Embed a link to EPG in the playlist")
    PlalistUdpxy        = flag.String("u", "", "Create a playlist with embedded udpxy")
    Health              = flag.Bool("H", false, "Start Healthcheck server")
    Udpxy               = flag.Bool("U", false, "Start UDPXY proxy")
    Schedule            = flag.Bool("S", false, "Start Schedule job")
    Crontab             = flag.String("c", "*30 6 * * *", "Create a playlist with embedded udpxy")
)


func Initgo() {
    flag.Parse()
    *EpgPathDst = regexp.MustCompile("/$").ReplaceAllString(*EpgPathDst, "")
    *PlaylistPathDst = regexp.MustCompile("/$").ReplaceAllString(*PlaylistPathDst, "")
}
