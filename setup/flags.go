package setup

import (
    "flag"
    "os"
    "regexp"
    "strconv"
)

var (
    LogLVL              = flag.Int("l", 1, "Log lvl")
    Epg                 = flag.String("e", getEnv("IPTVTOOLKIT_EPG", ""), "EPG URI")
    EpgDir              = flag.String("E", strippath(getEnv("IPTVTOOLKIT_EPG_DIR", "./files/tvguide")), "Path EPGs storage")
    Playlist            = flag.String("p", getEnv("IPTVTOOLKIT_PLAYLIST", ""), "Playlist URI")
    PlaylistDir         = flag.String("P", strippath(getEnv("IPTVTOOLKIT_PLAYLIST_DIR", "./files/playlist")), "Path playlists storage")
    PlaylistEmbedEPG    = flag.String("b", getEnv("IPTVTOOLKIT_PLAYLIST_EMBED_EPG", ""), "UDPXY URI")
    PlaylistEmbedUdpxy  = flag.String("a", getEnv("IPTVTOOLKIT_PLAYLIST_EMBED_UDPXY", ""), "Create a playlist with embedded udpxy")
    WebFilesPath        = flag.String("f", getEnv("IPTVTOOLKIT_WEB_FILES_PATH", "/files"), "Web Server static files path")
    WebFilesDir         = flag.String("d", getEnv("IPTVTOOLKIT_WEB_FILES_DIR", "./files"), "Directory on the host for web server files")
    WebUdpxyPath        = flag.String("u", getEnv("IPTVTOOLKIT_WEB_UDXPY_PATH", "/udp/:addr"), "Web Server multicast path")
    WebPort             = flag.Int("w", aToi(getEnv("IPTVTOOLKIT_WEB_PORT", "4022")), "Web Server port")
    Health              = flag.Bool("H", false, "Start Healthcheck server")
    Udpxy               = flag.Bool("U", false, "Start UDPXY proxy")
    Schedule            = flag.Bool("S", false, "Start Schedule job")
    Files               = flag.Bool("F", false, "Start File Server")
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
    *EpgDir = regexp.MustCompile(`/$`).ReplaceAllString(*EpgDir, "")
    *PlaylistDir = regexp.MustCompile(`/$`).ReplaceAllString(*PlaylistDir, "")
}

func strippath(a string) string {
    return regexp.MustCompile(`/$`).ReplaceAllString(a, "")
}
