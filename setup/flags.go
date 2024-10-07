package setup

import (
    "os"
    "flag"
    "regexp"
    "strconv"
)

var (
    LogLVL              = flag.Int("l", 1, "Log lvl")
    Epg                 = flag.String("e", getEnv("IPTVTOOLKIT_EPG", ""), "EPG URI")
    EpgPathDst          = flag.String("E", regexp.MustCompile("/$").ReplaceAllString(getEnv("IPTVTOOLKIT_EPG_DST", "./files/tvguide"), ""), "Path EPGs storage")
    Playlist            = flag.String("p", getEnv("IPTVTOOLKIT_PLAYLIST", ""), "Playlist URI")
    PlaylistPathDst     = flag.String("P", regexp.MustCompile("/$").ReplaceAllString(getEnv("IPTVTOOLKIT_PLAYLIST_DST", "./files/playlist"), ""), "Path playlists storage")
    PlalistUdpxy        = flag.String("u", getEnv("IPTVTOOLKIT_PLAYLIST_UDPXY", ""), "Create a playlist with embedded udpxy")
    EmbedEPG            = flag.String("i", getEnv("IPTVTOOLKIT_PLAYLIST_EMBED_EPG", ""), "Embed a link to EPG in the playlist")
    WebPath             = flag.String("f", getEnv("IPTVTOOLKIT_WEB_PATH", "./files"), "Web Server path")
    WebPort             = flag.Int("w", aToi(getEnv("IPTVTOOLKIT_WEB_PORT", "4022")), "Web Server port")
    Health              = flag.Bool("H", false, "Start Healthcheck server")
    Udpxy               = flag.Bool("U", false, "Start UDPXY proxy")
    Schedule            = flag.Bool("S", false, "Start Schedule job")
    Web                 = flag.Bool("W", false, "Start Web Server")
    Crontab             = flag.String("c", getEnv("IPTVTOOLKIT_CRONTAB", "30 6 * * *"), "Сrontab style task schedule")
)

func getEnv(key, fallback string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        return fallback
    }
    return value
}

func aToi(value string) int {
    a, _ := strconv.Atoi(value)
    return a
}

func Initgo() {
    flag.Parse()
}
